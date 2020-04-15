package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	"olympos.io/encoding/edn"
)

// --- work in progress ---
// generic tx and messages parser
// inspired by Bitcoin and Clojure
// mixture of edn and script

// transactions are typed
// [:txtype {:content as map} {:signature data}]
// [:multisig <txcontent> <sig1 sig2>]
//multiplexing

//{:simple ...}
//{:script [...]}
//{:contract [...]}
//....

// ETH opcodes
// 0x30	ADDRESS	Get address of currently executing account	-	2
// 0x31	BALANCE	Get balance of the given account
// 0x32	ORIGIN	Get execution origination address	-	2
// 0x33	CALLER	Get caller address	-	2
// 0x34	CALLVALUE	Get deposited value by the instruction/transaction responsible for this execution	-	2
// 0x35	CALLDATALOAD	Get input data of current environment	-	3
// 0x36	CALLDATASIZE	Get size of input data in current environment	-	2*
// 0x37	CALLDATACOPY	Copy input data in current environment to memory	-	3
// 0x38	CODESIZE	Get size of code running in current environment	-	2
// 0x39	CODECOPY	Copy code running in current environment to memory	-

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	IDENT
	KEYWORD
	OPENMAP
	CLOSEMAP
	OPENVECTOR
	//
	MAPCONTENT
	SIMPLETX
	// Keywords
	//SELECT

)

//end of file
var eof = rune(0)

//----
//transaction types keyword
const STX = "STX"

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isKeyword(ch rune) bool {
	return ch == ':'
}

func isMapStart(ch rune) bool {
	return ch == '{'
}

func isMapEnd(ch rune) bool {
	return ch == '}'
}

func isVectorStart(ch rune) bool {
	return ch == '['
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	//fmt.Println(">> ", ch)

	if isWhitespace(ch) {
		// 	s.unread()
		// 	return s.scanWhitespace()
		return WS, ""
	} else if isKeyword(ch) {
		return KEYWORD, ":"
	} else if isVectorStart(ch) {
		return OPENVECTOR, "["
	}

	//TODO

	//else if isLetter(ch) {
	// 	s.unread()
	// 	return s.scanIdent()
	// }

	// // Otherwise read the individual character.
	// //fmt.Println("switch ", ch)
	switch ch {
	case eof:
		return EOF, ""
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scan identifier consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

//scan for map
func (s *Scanner) scanMap() (tok Token, lit string) {

	var buf bytes.Buffer
	buf.WriteRune(s.read())

	firstch := s.read()
	if !isMapStart(firstch) {
		//error
		return ILLEGAL, ""
	}

	//got open map
	buf.WriteRune(firstch)

	// Read every subsequent ident character into the buffer
	for {
		ch := s.read()
		buf.WriteRune(ch)
		if isMapEnd(ch) {
			break
		}
	}
	return MAPCONTENT, buf.String()
}

func (s *Scanner) scanFirstKey() (tok Token, lit string) {

	//scan first keyword
	firstKey := false

	for tok, lit := s.Scan(); !firstKey; tok, lit = s.Scan() {
		//fmt.Println("!! ", lit, "   ", tok)
		if tok == OPENVECTOR {
			fmt.Println("open vector")
		}

		if tok == KEYWORD {
			firstKey = true

			idt, idtlit := s.scanIdent()
			//fmt.Println(">>> ", idt)
			if idtlit == "STX" {
				//fmt.Println("simple transaction")
				return SIMPLETX, idtlit
			}
			idt += 0
		}
		lit += ""
		//fmt.Println(lit)
	}

	return ILLEGAL, ""
}

func (s *Scanner) scanSimpletx() (tok Token, lit string) {

	//scan first keyword
	firstKey := false

	for tok, lit := s.Scan(); !firstKey; tok, lit = s.Scan() {
		//fmt.Println(">>> ", lit, "   ", tok)
		if tok == OPENMAP {
			fmt.Println("open map")
		}

		if tok == KEYWORD {
			firstKey = true

			idt, idtlit := s.scanIdent()
			//fmt.Println(">>> ", idt)
			if idtlit == STX {
				//fmt.Println("simple transaction")
				return SIMPLETX, idtlit
			}
			idt += 0
		}
		lit += ""
		//fmt.Println(lit)
	}

	return ILLEGAL, ""
}

func (s *Scanner) scanRest() (rest string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func main() {
	//fmt.Println("test tx")

	keypair := crypto.PairFromSecret("test")
	pub := crypto.PubKeyToHex(keypair.PubKey)
	a := crypto.Address(pub)
	fmt.Println("from ", a)

	pubkey_string := crypto.PubKeyToHex(keypair.PubKey)
	fmt.Println(pubkey_string)

	simpletx := `{:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42}`
	signature := crypto.SignMsgHash(keypair, txmap)

	inputstring := `[:STX {:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42} {:SenderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :Signature "3044022047d9411810cd5d4feaf7fc806071bc3eb66f2d62d551f1e8b18de0ff7dbbefe80220315c5cc342dd817ba79c3612e2ad9fb560cce8f16cb1ba6ab562a5a420cddf98"}]`

	s := NewScanner(strings.NewReader(inputstring))

	ftok, flit := s.scanFirstKey()
	fmt.Println("first => ", ftok, flit)

	//simple tx. first element contains tx, 2nd the signature data
	if ftok == SIMPLETX {
		_, txmap := s.scanMap()
		fmt.Println("tx content => ", txmap)

		signature := crypto.SignMsgHash(keypair, txmap)
		sighex := hex.EncodeToString(signature.Serialize())
		fmt.Println(sighex)

		var tx block.Tx
		edn.Unmarshal([]byte(txmap), &tx)
		log.Println("sender ", tx.Sender)

		_, sigmap := s.scanMap()
		fmt.Println("sigmap => ", sigmap)

		var txsig block.TxSig
		edn.Unmarshal([]byte(sigmap), &txsig)
		fmt.Println(txsig.SenderPubkey)

		/////////////
		s1 := crypto.SignatureFromHex(txsig.Signature)
		p1 := crypto.PubKeyFromHex(txsig.SenderPubkey)
		//verified := crypto.VerifyMessageSign(s1, keypair, message)
		verified := crypto.VerifyMessageSignPub(s1, p1, txmap)
		fmt.Println(verified)

		// var tx block.Tx
		// edn.Unmarshal([]byte(msg), &tx)
		// log.Println(tx.Signature)
	}

	for tok, lit := s.Scan(); tok != EOF; tok, lit = s.Scan() {
		//fmt.Println(tok, lit)
		fmt.Println(lit, "   ", tok)
		// if tok == SELECT {
		// 	fmt.Println("S")
		// }

		tok = 0
		// if lit == "" {
		// 	fmt.Println("lit ", lit)
		// }

		//fmt.Println(tok == EOF)
	}

	// example := `{:TxType :script,
	// 			  :lock [OP_DUP]}`

}
