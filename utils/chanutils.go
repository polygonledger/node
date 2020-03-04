package utils

import "time"

//basic threading helper
func DoEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func DoEveryX(d time.Duration, f func() <-chan string) {
	for _ = range time.Tick(d) {
		f()
	}
}
