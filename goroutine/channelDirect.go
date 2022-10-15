package main

import (
	"fmt"
	"time"
)

const (
	JOB_DONE = "Mission complete !"
	STILL_WORKING = "In progress ..."
)

func phoneCall(hotline chan<- string, msg string, delay time.Duration) {
	time.Sleep(delay)
	hotline <- msg
}

func main() {
	hotline := make(chan string)
	canRest := make(chan bool)

	go func(hotline <-chan string, canRest chan<- bool) {
		for {
			msg, hasNext := <-hotline
			if !hasNext {
				fmt.Println("Hotline has been closed, go to rest.")
				canRest <- true
				return
			}
			if msg == JOB_DONE {
				fmt.Println("Recv msg. Time to get some rest ~")
				canRest <- true
				return
			}
			fmt.Println("Recv msg: ", msg)
		}
	}(hotline, canRest)

	go phoneCall(hotline, STILL_WORKING, 1*time.Second)
	go phoneCall(hotline, JOB_DONE, 3*time.Second)

	<-canRest
}
