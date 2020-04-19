package parser

//see also
//https://github.com/candid82/joker
//https://github.com/go-edn/

//basic edn utils
//create structs from strings, no nesting

func StringWrap(s string) string {
	return "\"" + s + "\""
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

func MakeMap(els []string, keys []string) string {
	vs := `{`
	for i, s := range els {
		k := keys[i]
		vs += ":" + k + " " + s
		if i < len(els)-1 {
			vs += " "
		}
	}
	vs += `}`
	return vs
}

// func scanMapString() ([]string, []string) {
// 	//scan for open map
// 	//scan :keyword
// 	//scan keyword id
// 	//scan value

// }
