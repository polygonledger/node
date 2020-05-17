package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/ec2"
)

func check(e error) {

	if e != nil {
		panic(e)
	}
}

func main() {

	// aws_user := os.Getenv("AWS_ACCESS_KEY_ID")
	// aws_pass := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// creds := aws.Creds(aws_user, aws_pass, "")
	// client := ec2.New(creds, "us-west-1", nil)

	client := ec2.New(session.New())

	// Only grab instances that are running or just started
	filters := []ec2.Filter{
		ec2.Filter{
			aws.String("instance-state-name"),
			[]string{"running", "pending"},
		},
	}
	request := ec2.DescribeInstancesRequest{Filters: filters}
	result, err := client.DescribeInstances(&request)
	check(err)

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println(*instance.InstanceID)
		}
	}
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/ec2"
// )

// func main() {
// 	ec2svc := ec2.New(session.New())
// 	params := &ec2.DescribeInstancesInput{
// 		Filters: []*ec2.Filter{
// 			{
// 				Name:   aws.String("tag:Environment"),
// 				Values: []*string{aws.String("prod")},
// 			},
// 			{
// 				Name:   aws.String("instance-state-name"),
// 				Values: []*string{aws.String("running"), aws.String("pending")},
// 			},
// 		},
// 	}
// 	resp, err := ec2svc.DescribeInstances(params)
// 	if err != nil {
// 		fmt.Println("there was an error listing instances in", err.Error())
// 		log.Fatal(err.Error())
// 	}

// 	for idx, res := range resp.Reservations {
// 		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
// 		for _, inst := range resp.Reservations[idx].Instances {
// 			fmt.Println("    - Instance ID: ", *inst.InstanceId)
// 		}
// 	}
// }
