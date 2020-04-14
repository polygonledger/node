package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/polygonledger/node/block"
	"olympos.io/encoding/edn"
)

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

var eof = rune(0)

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

	//fmt.Println("letter ", string(ch))

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

// scanIdent consumes the current rune and all contiguous ident runes.
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

func (s *Scanner) scanMap() (tok Token, lit string) {
	fmt.Println("scanmap")
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		ch := s.read()
		fmt.Println(string(ch))
		if isMapEnd(ch) {
			fmt.Println("> break ", string(ch))
			break
		} else if isMapStart(ch) {

		} else {
			fmt.Println("write ", ch)
			_, _ = buf.WriteRune(ch)
		}
		//  else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
		// 	fmt.Println("> ", string(ch))
		// 	s.unread()
		// 	break
		// }

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

	msg := `{:TxType :simple,
			  :Sender "abc",
			  :Receiver "xyz",
		      :amount 42,
			  :nonce 1}`

	var tx block.Tx
	edn.Unmarshal([]byte(msg), &tx)
	log.Println(tx)

	//inputs := "{:stx {:body {:Sender abc :Receiver xyz :amount 42} :sig {:SenderPubkey sdfasa: Signature afwfswf}}"
	inputs := "[:STX {:Sender abc :Receiver xyz :amount 42} {:SenderPubkey sdfasa: Signature afwfswf}]"
	s := NewScanner(strings.NewReader(inputs))

	ftok, flit := s.scanFirstKey()
	fmt.Println("first => ", ftok, flit)

	if ftok == SIMPLETX {
		_, firstmap := s.scanMap()
		fmt.Println("? ", firstmap)

		rest := s.scanRest()
		fmt.Println("rest ", rest)

		var tx block.Tx
		edn.Unmarshal([]byte(msg), &tx)
		log.Println(tx.Signature)
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
