package parser

import (
	"fmt"
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
func ReadMap1(mapstr string) ([]string, []string) {

	var vs []string
	var ks []string

	s := NewScanner(strings.NewReader(mapstr))

	//ldone := false

	s.Scan() //open bracket

	//i := 0

	for {
		ch := s.read()
		fmt.Println("> ", ch, " ", string(ch))	

		if ch == eof {
			break
		} else if isKeyword(ch) {
			
		}

		if isMapEnd(ch) {
			fmt.Println("end")
			break
		}
	}

	// for !ldone {

	// 	fmt.Println("ldone ",ldone)

	// 	//4 cases
	// 	//keyword => scanid
	// 	//whitespace => consume
	// 	//scanid => value
	// 	//close bracked => end

	// 	_, kk := s.Scan() //":"
	// 	fmt.Println("kk ",kk, len(kk),string(kk))
	// 	if len(kk) == 0 {
	// 		fmt.Println("!!")
	// 		ldone = true
	// 	}			
	// 	switch kk {
	// 	case ":":
	// 		fmt.Println("keyword")
	// 		_, klit := s.scanIdent()
	// 		ks = append(ks, klit)
	// 		fmt.Println(klit)
	// 		//s.scanWhitespace()

	// 	case "}":
	// 		ldone = true

	// 	case "":
	// 		ldone = true

	// 	case " ":
	// 		fmt.Println("whitespace")
	// 		s.scanWhitespace()
	// 		_, vlit := s.scanIdent()
	// 		vs = append(vs, vlit)

	// 	default:
	// 		fmt.Println("default")
	// 		_, vlit := s.scanIdent()
	// 		vs = append(vs, vlit)
	// 		s.scanWhitespace()
	// 		ldone = true

	// 	}

	// 	//debug
		
	// 	//ldone = i > 20
	// 	i++

	// 	// _, klit := s.scanIdent()
	// 	// fmt.Println(klit)
	// 	// ks = append(ks, klit)
	// 	// s.scanWhitespace()

	// 	// _, vlit := s.scanIdent()
	// 	// fmt.Println(vlit)
	// 	// vs = append(vs, vlit)

	// 	// _, xlit := s.Scan() //close
	// 	// if xlit == "}" {
	// 	// 	//fmt.Println("Close map")
	// 	// 	ldone = true
	// 	// }

	// }
	return vs, ks
}
