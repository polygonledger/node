package main

// import (
// 	"encoding/json"
// 	"fmt"
// )

// type Message struct {
// 	Type string `json:type`
// 	Data json.RawMessage
// 	//Timestamp string `json:timestamp`
// }

// type Account struct {
// 	ABC string `json:abc`
// }

// type Balance struct {
// 	DEF string `json:def`
// }

// var JSONEventData = []byte(`{
//    "type": "account",
//    "data": {
//         "abc": "asfsdaf"
//    }}`)

// var JSONNewsData = []byte(`{
//    "type": "balance",
//    "data": {
//         "def": "xxxafew20"
//    }}`)

// func show(m Message) {

// 	switch m.Type {
// 	case "account":
// 		var e Account
// 		if err := json.Unmarshal([]byte(m.Data), &e); err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Printf("JSON type is [%v] and [%v] \n", m.Type, e.ABC)

// 	case "balance":
// 		var n Balance
// 		if err := json.Unmarshal([]byte(m.Data), &n); err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Printf("JSON type is [%v] and [%v] \n", m.Type, n.DEF)

// 	default:
// 		fmt.Println("unable to unmarshal JSON data or differentiate the type")
// 	}
// }

// func main() {

// 	var m Message

// 	// unmarshal JSONEventData to message struct
// 	if err := json.Unmarshal(JSONEventData, &m); err != nil {
// 		fmt.Println(err)
// 	}
// 	show(m)

// 	if err := json.Unmarshal(JSONNewsData, &m); err != nil {
// 		fmt.Println(err)
// 	}
// 	show(m)

// 	var m2 Message
// 	m2.Type = "balance"
// 	var b Balance
// 	b.DEF = "test"
// 	bjson, _ := json.Marshal(b)
// 	fmt.Println(bjson)
// 	m2.Data = bjson

// 	msgjson, _ := json.Marshal(m2)
// 	fmt.Println(string(msgjson))

// }
