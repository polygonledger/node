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

// func Pubtime(tchan chan string) {
// 	d := 1000 * time.Millisecond
// 	for _ = range time.Tick(d) {
// 		PublishTime(tchan)
// 	}
// }

// func Subtime(tchan chan string, name string) {
// 	d := 100 * time.Millisecond
// 	for _ = range time.Tick(d) {
// 		x := <-tchan
// 		log.Printf("subscriber [%s] %s", name, x)
// 	}
// }

func Subout(tchan chan string, name string, outchan chan Message) {
	d := 100 * time.Millisecond
	for _ = range time.Tick(d) {
		x := <-tchan
		log.Printf("subscriber [%s] %s", name, x)

		//TODO put in outchan
		msg := EncodeMsg("TEST", x, "test")
		log.Println("republish ", msg)
		outchan <- msg

	}
}

// func PublishLoop() {
// 	tTime := 1000 * time.Millisecond
// 	for _ = range time.Tick(tTime) {
// 		//log.Println(tt)
// 		PublishTime()
// 	}
// }
