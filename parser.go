package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
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
	fmt.Println("scan whitespace")
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

func verifyTx(txmap string, sighex string, pubhex string) {
	s1 := crypto.SignatureFromHex(sighex)
	p1 := crypto.PubKeyFromHex(pubhex)
	//verified := crypto.VerifyMessageSign(s1, keypair, message)
	verified := crypto.VerifyMessageSignPub(s1, p1, txmap)
	fmt.Println(verified)
}

func stringWrap(s string) string {
	return "\"" + s + "\""
}

func createsigmap(pubkey_string string, txsighex string) string {
	v := []string{stringWrap(pubkey_string), stringWrap(txsighex)}
	k := []string{"SenderPubkey", "Signature"}
	m := makeMap(v, k)
	return m
}

func verifySigmap(sigmap string, txmap string) bool {
	var txsig block.TxSig
	edn.Unmarshal([]byte(sigmap), &txsig)
	//fmt.Println("txsig.Signature ", txsig.Signature)
	s1 := crypto.SignatureFromHex(txsig.Signature)
	p1 := crypto.PubKeyFromHex(txsig.SenderPubkey)
	verified := crypto.VerifyMessageSignPub(s1, p1, txmap)
	fmt.Println("verified => ", verified)
	return verified
}

//create the vector from tx and sig data
func txVector(simpletx string, sigmap string) string {
	vs := []string{":STX", simpletx, sigmap}
	return makeVector(vs)
}

//[:type {tx} {sig}]
//extract the components
func ScanScript(inputVector string) (string, string) {

	s := NewScanner(strings.NewReader(inputVector))
	ftok, _ := s.scanFirstKey()
	//s.scanWhitespace()
	s.Scan()

	//simple tx. first element contains tx, 2nd the signature data
	if ftok == SIMPLETX {
		//s.scanWhitespace()
		_, txmap := s.scanMap()
		fmt.Println("tx content => ", txmap)

		// var tx block.Tx
		// edn.Unmarshal([]byte(txmap), &tx)

		s.scanWhitespace()
		//log.Println("sender ", tx.Sender)

		_, sigmap := s.scanMap()
		//fmt.Println("sigmap => ", sigmap)
		//remove leading whitespace between vector elements
		//sigmap = sigmap[1:]

		//verifySigmap(sigmap, txmap)
		return sigmap, txmap

	}
	return "", ""
}

// func scanMapString() ([]string, []string) {
// 	//scan for open map
// 	//scan :keyword
// 	//scan keyword id
// 	//scan value

// }

//verify signature
//independent of balance check
func VerifyTxScriptSig(v string) bool {

	sigmap, txmap := ScanScript(v)
	valid := verifySigmap(sigmap, txmap)
	return valid
}

func CreateSimpleTxContent(sender string, receiver string, amount int) string {
	v := []string{stringWrap(sender), stringWrap(receiver), strconv.Itoa(amount)}
	k := []string{"Sender", "Receiver", "amount"}
	txmap := makeMap(v, k)
	return txmap
}

func SignMap(keypair crypto.Keypair, msg string) string {
	txsig := crypto.SignMsgHash(keypair, msg)
	txsighex := hex.EncodeToString(txsig.Serialize())
	pubkey_string := crypto.PubKeyToHex(keypair.PubKey)

	v := []string{pubkey_string, txsighex}
	k := []string{"SenderPubkey", "Signature"}
	sigmap := makeMap(v, k)
	return sigmap
}

func parseexample() {

	//`{:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42}`
	simpletx := CreateSimpleTxContent("Pa033f6528cc1", "P7ba453f23337", 42)

	keypair := crypto.PairFromSecret("test")
	sigmap := SignMap(keypair, simpletx)
	fmt.Println(sigmap)

	verifySigmap(sigmap, simpletx)

	fmt.Println("tx ", simpletx)
	fmt.Println("sigmap ", sigmap)

	v := txVector(simpletx, sigmap)
	fmt.Println("tx vector ", v)

	v = `[:STX {:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42} {:SenderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :Signature "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"}]`

	valid := VerifyTxScriptSig(v)
	fmt.Println(valid)

	//verification parser

}
