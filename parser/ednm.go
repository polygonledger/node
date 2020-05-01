package parser

import (
	"sort"
	"strings"
)

//see also
//https://github.com/candid82/joker
//https://github.com/go-edn/

//basic edn utils
//create structs from strings, no nesting

func StringWrap(s string) string {
	return "\"" + s + "\""
}
func StringUnWrap(s string) string {
	return s[1 : len(s)-1]
}

func MakeKeyword(k string) string {
	return ":" + k
}

func MakeVector(vectorels []string) string {
	vs := `[`
	for i, s := range vectorels {
		vs += s
		if i < len(vectorels)-1 {
			vs += " "
		}
	}
	vs += `]`
	return vs
}

func MakeMap(m map[string]string) string {
	vs := `{`
	i := 0

	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		value := m[k]
		vs += ":" + k + " " + value
		if i < len(m)-1 {
			vs += " "
		}
		i++
	}
	vs += `}`
	return vs
}

//read map and return keys and values as strings
func ReadMap(mapstr string) ([]string, []string) {

	var vs []string
	var ks []string

	s := NewScanner(strings.NewReader(mapstr))

	//ldone := false

	//s.Scan() //open bracket

	expectvalue := false

	for {
		ch := s.read()

		if ch == eof {
			break
		} else if isKeyword(ch) {
			_, klit := s.scanIdent()
			ks = append(ks, klit)
			expectvalue = true
		} else if isWhitespace(ch) {
			if expectvalue {
				_, vlit := s.scanIdent()
				vs = append(vs, vlit)
				expectvalue = false
			}
		} else if isMapEnd(ch) {
			break
		} else if isMapStart(ch) {

		}
	}

	// 	//4 cases
	// 	//keyword => scanid
	// 	//whitespace => consume
	// 	//scanid => value
	// 	//close bracked => end

	return vs, ks
}
