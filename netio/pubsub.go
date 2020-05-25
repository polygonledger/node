package netio

import "time"

/////// PUBSUB ///////

func PublishTime(ntchan Ntchan) {
	timeFormat := "2006-01-02T15:04:05"
	limiter := time.Tick(1000 * time.Millisecond)
	pubcount := 0
	//log.Println("PublishTime")

	for {
		t := time.Now()
		tf := t.Format(timeFormat)
		vlog(ntchan, "pub "+tf)
		ntchan.PUB_out <- tf
		<-limiter
		pubcount++
	}

}

//publication to writer queue. requires quit channel
func PubWriterLoop(ntchan Ntchan) {

	for {
		select {
		case msg := <-ntchan.PUB_out:
			vlog(ntchan, "sub "+msg)
			ntchan.Writer_queue <- msg
			// case <-ntchan.PUB_time_quit:
			// 	fmt.Println("stop pub")
			// 	return
			// 	// default:
			// 	// 	fmt.Println("no message received")
		}
		time.Sleep(50 * time.Millisecond)

	}

}
