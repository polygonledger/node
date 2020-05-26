package xutils

import (
	"time"
)

//check if a channel is empty
func IsEmpty(c chan string, d time.Duration) bool {
	select {
	case <-c:
		//log.Println("got ", ret)
		//assert ret==nil
		return false

	case <-time.After(d):
		//log.Println("timeout")
		return true
	}
}
