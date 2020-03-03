package main

import (
	"log"
	"time"
)

//check if a channel is empty
func isEmpty(c chan string, d time.Duration) bool {
	select {
	case ret := <-c:
		log.Println("got ", ret)
		return false

	case <-time.After(d):
		//log.Println("timeout")
		return true
	}
}
