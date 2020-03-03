package main

import (
	"time"
)

//basic threading helper
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func doEveryX(d time.Duration, f func() <-chan string) {
	for _ = range time.Tick(d) {
		f()
	}
}
