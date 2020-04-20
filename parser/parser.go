package parser

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"strconv"
	"strings"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	"olympos.io/encoding/edn"
)

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

//SCRIPT
//CONTRACT

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

func isIdent(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_' || ch == '"'
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
	switch ch {
	case eof:
		return EOF, ""
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer
	// Non-whitespace characters and EOF will cause the loop to exit
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

// scan identifier consumes the current rune and all contiguous ident runes
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it
	var buf bytes.Buffer
	//buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer
	// Non-ident characters and EOF will cause the loop to exit
	for {
		ch := s.read()
		if ch == eof {
			break
		} else if !isIdent(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

//scan for map contents as string
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
		if tok == OPENVECTOR {
		}

		if tok == KEYWORD {
			firstKey = true

			idt, idtlit := s.scanIdent()
			if idtlit == STX {
				return SIMPLETX, idtlit
			}
			idt += 0
		}
		lit += ""
	}

	return ILLEGAL, ""
}

func (s *Scanner) scanSimpletx() (tok Token, lit string) {

	//scan first keyword
	firstKey := false

	for tok, lit := s.Scan(); !firstKey; tok, lit = s.Scan() {
		if tok == OPENMAP {
		}

		if tok == KEYWORD {
			firstKey = true

			idt, idtlit := s.scanIdent()
			if idtlit == STX {
				return SIMPLETX, idtlit
			}
			idt += 0
		}
		lit += ""
	}

	return ILLEGAL, ""
}

func (s *Scanner) ReadMap() ([]string, []string) {

	var vs []string
	var ks []string

	for {
		ch := s.read()

		if ch == eof {
			break
		} else if isMapEnd(ch) {

		} else if isMapStart(ch) {

		} else if isWhitespace(ch) {

		} else if isKeyword(ch) {
			_, idtlit := s.scanIdent()
			ks = append(ks, idtlit)
		} else {
			_, vlit := s.scanIdent()
			vlit = string(ch) + vlit
			vs = append(vs, vlit)
		}
	}

	return vs, ks
}

func ReadMapP(s string) ([]string, []string) {
	sc := NewScanner(strings.NewReader(s))
	vs, ks := sc.ReadMap()
	return vs, ks
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

func verifyTx(txmap string, sighex string, pubhex string) bool {
	s1 := crypto.SignatureFromHex(sighex)
	p1 := crypto.PubKeyFromHex(pubhex)
	verified := crypto.VerifyMessageSignPub(s1, p1, txmap)
	return verified
}

func CreateSigmap(pubkey_string string, txsighex string) string {
	v := []string{StringWrap(pubkey_string), StringWrap(txsighex)}
	k := []string{"senderPubkey", "signature"}
	m := MakeMapArr(v, k)
	return m
}

func SignMap(keypair crypto.Keypair, msg string) string {
	txsig := crypto.SignMsgHash(keypair, msg)
	txsighex := hex.EncodeToString(txsig.Serialize())
	pubkey_string := crypto.PubKeyToHex(keypair.PubKey)
	sigmap := CreateSigmap(pubkey_string, txsighex)
	return sigmap
}

func VerifySigmap(sigmap string, txmap string) bool {
	var txsig block.TxSig
	edn.Unmarshal([]byte(sigmap), &txsig)
	s1 := crypto.SignatureFromHex(txsig.Signature)
	p1 := crypto.PubKeyFromHex(txsig.SenderPubkey)
	verified := crypto.VerifyMessageSignPub(s1, p1, txmap)
	return verified
}

//create the vector from tx and sig data
func TxVector(simpletx string, sigmap string) string {
	vs := []string{MakeKeyword(STX), simpletx, sigmap}
	return MakeVector(vs)
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

		// var tx block.Tx
		// edn.Unmarshal([]byte(txmap), &tx)

		s.scanWhitespace()
		//log.Println("sender ", tx.Sender)

		_, sigmap := s.scanMap()
		//remove leading whitespace between vector elements
		//sigmap = sigmap[1:]

		//verifySigmap(sigmap, txmap)
		return sigmap, txmap

	}
	return "", ""
}

//TODO
// func MakeBlockStr(txs []string, pubk string) string {

// 	txstr := MakeVector(txs)
// 	k := []string{"txs", "pubkey"}
// 	v := []string{txstr, pubk}
// 	txmap := MakeMap(k, v)
// 	return txmap
// }

//verify signature
//independent of balance check
func VerifyTxScriptSig(v string) bool {

	sigmap, txmap := ScanScript(v)
	valid := VerifySigmap(sigmap, txmap)
	return valid
}

func CreateSimpleTxContent(sender string, receiver string, amount int) string {
	m := map[string]string{"sender": StringWrap(sender), "receiver": StringWrap(receiver), "amount": strconv.Itoa(amount)}

	//m := map[string]string{}

	txmap := MakeMap(m)
	return txmap
}
