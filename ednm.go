package main

//basic edn create
//create structs from strings, no nesting

func makeVector(vectorels []string) string {
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

func makeMap(els []string, keys []string) string {
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
