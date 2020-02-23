package ntwk

import (
	"log"
	"time"
)

var Timechan = make(chan string, 100)

func TakeChan() string {
	x := ""
	x = <-Timechan
	log.Println(len(Timechan))
	return x
}

func Publisher() chan string {
	tchan := make(chan string, 100)
	return tchan
}

func Pubtime(tchan chan string) {
	d := 1000 * time.Millisecond
	for _ = range time.Tick(d) {
		PublishTime(tchan)
	}
}

func Subtime(tchan chan string, name string) {
	d := 100 * time.Millisecond
	for _ = range time.Tick(d) {
		x := <-tchan
		log.Printf("subscriber [%s] %s", name, x)
	}
}

func Subout(tchan chan string, name string, outchan chan Message) {
	d := 100 * time.Millisecond
	for _ = range time.Tick(d) {
		x := <-tchan
		log.Printf("subscriber [%s] %s", name, x)

		//TODO put in outchan
		msg := EncodeMsg("TEST", x, "test")
		outchan <- msg

	}
}

func PublishTime(tchan chan string) {
	t := time.Now()
	timeFormat := "2006-01-02T15:04:05"
	tf := t.Format(timeFormat)
	go func() {
		log.Println("publish ", tf)
		tchan <- tf
		//log.Println(len(Timechan), cap(Timechan))
		//log.Println("do")
	}()

}

// func PublishLoop() {
// 	tTime := 1000 * time.Millisecond
// 	for _ = range time.Tick(tTime) {
// 		//log.Println(tt)
// 		PublishTime()
// 	}
// }
