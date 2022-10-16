package main

import (
	"fmt"
	"time"
	"reflect"
)

var ttl = 3*time.Second

func main() {
	tasks := make(chan int)
	done := make(chan bool)
	timer := time.NewTimer(ttl)

	quitReason := map[int]string{
		0: "Channel closed by sender",
		1: "Recv quit signal",
		2: "Timeout",
	}

	cs := []reflect.SelectCase{
		reflect.SelectCase{
			Dir: reflect.SelectRecv,
			Chan: reflect.ValueOf(tasks),
		},
		reflect.SelectCase{
			Dir: reflect.SelectRecv,
			Chan: reflect.ValueOf(done),
		},
		reflect.SelectCase{
			Dir: reflect.SelectRecv,
			Chan: reflect.ValueOf(timer.C),
		},
	}

	go func() {
		tasks <- 1
		time.Sleep(2*time.Second)
		tasks <- 2
		time.Sleep(2*time.Second)

		/* case #1: send quit signal */
		done <- true

		/* case #2: close channel */
		//close(tasks)

		/* case #3: do nothing, for loop quit with timeout */
	}()

	for {
		i, val, ok := reflect.Select(cs)
		if i != 0 || !ok {
			fmt.Printf("[channel #%d] %s ! value recved: %v\n", i, quitReason[i], val)
			return
		}

		fmt.Printf("[channel #%d] Recv signal with value: %v\n", i, val)
		timer.Reset(ttl)
	}
}
