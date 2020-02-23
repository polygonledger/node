package ntwk

import (
	"log"
	"time"
)

var Timechan = make(chan string)

func TakeChan() string {
	x := ""
	x = <-Timechan
	return x
}

func PublishTime() {
	t := time.Now()
	timeFormat := "2006-01-02T15:04:05"
	tf := t.Format(timeFormat)
	go func() {
		Timechan <- tf
		log.Println("do")
	}()

}

func PublishLoop() {
	hTime := 2000 * time.Millisecond
	for _ = range time.Tick(hTime) {
		//log.Println(x)
		PublishTime()
	}
}

func SubLoop() {
	hTime := 2000 * time.Millisecond
	for _ = range time.Tick(hTime) {
		//log.Println(x)
		PublishTime()
	}
}
