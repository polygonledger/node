package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"fmt"
)

func main() {
	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new EC2 client
	ec2Svc := ec2.New(sess)

	// Call to get detailed information on each instance
	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Success", result)
		result_string, _ := json.Marshal(result)
		fmt.Println(string(result_string))

		err := ioutil.WriteFile("instances.json", []byte(string(result_string)), 0644)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(result.Reservations)
	}
}
