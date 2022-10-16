package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	cost      = 1*time.Second
	workerNum = 3
	taskNum   = 10
)

type worker struct {
	id	int

	task	<-chan int
	result	chan<- int
	taskwg	*sync.WaitGroup

	quit	<-chan bool
	quitwg	*sync.WaitGroup
}

func (w *worker) run() {
	for {
		select {
		case t := <-w.task:
			fmt.Printf("[worker #%d] Recv task %v, processing ...\n", w.id, t)
			time.Sleep(cost)
			fmt.Printf("[worker #%d] Task done.\n", w.id)
			w.result <- t*10
			w.taskwg.Done()
		case <-w.quit:
			fmt.Printf("[worker #%d] Terminated.\n", w.id)
			w.quitwg.Done()
			return
		}
	}
}

func main() {
	task := make(chan int)
	result := make(chan int)
	done := make(chan bool)
	var taskwg sync.WaitGroup
	var quitwg sync.WaitGroup

	defer close(result)
	defer close(task)

	// Pop up workers
	workerPool := make([]worker, workerNum)
	workerQuit := make([]chan bool, workerNum)
	for i := 0; i < workerNum; i++ {
		workerQuit[i] = make(chan bool)
		workerPool[i] = worker{
			id:     i,
			task:   task,
			result: result,
			taskwg: &taskwg,
			quit:	workerQuit[i],
			quitwg: &quitwg,
		}
		go workerPool[i].run()
	}

	// Print result calculated by workers
	go func(result <-chan int, done <-chan bool) {
		for {
			select {
			case rslt := <-result:
				fmt.Println("Get result: ", rslt)
			case <-done:
				fmt.Println("Stop watching on result.")
				return
			}
		}
	}(result, done)

	// Dispatch tasks to running workers
	for i := 0; i < taskNum; i++ {
		taskwg.Add(1)
		task <- i
	}
	taskwg.Wait()

	// Close workers after all tasks done
	quitwg.Add(workerNum)
	for i := range workerQuit {
		workerQuit[i] <- true
		close(workerQuit[i])
	}
	quitwg.Wait()

	done <- true
	time.Sleep(1*time.Second)
	fmt.Println("All worker terminated. Quit program.")
}
