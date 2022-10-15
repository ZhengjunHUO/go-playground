package main

import (
	"fmt"
	"time"
)

type task struct {
	id	int
	value	*int
}

var one int = 1

func main() {
	tasks := make(chan task, 2)
	canQuit := make(chan bool)

	go func() {
		// 消费端，较慢地读取chan中的内容
		for {
			time.Sleep(3*time.Second)
			t, hasNext := <-tasks

			if hasNext {
				fmt.Printf("Recv task: %+v\n", t)
			}else{
				// 通道已经被发送方关闭，会额外读到一个内容的空值
				fmt.Printf("All tasks received, the extra t's value is: %+v\n", t)
				canQuit <- true
				return
			}
		}
	}()

	// 生产端，快速填满chan的buffer后block直到消费端消费内容
	for i:=0; i<5; i++ {
		tasks <- task{i, &one}
		fmt.Println("Sent task: ", i)
	}
	// 生产端先行关闭通道，但消费端仍旧可以读取其buffer中内容
	// 有效内容读取完成后会额外读到一个(内容的空值，false)
	// 另消费方自行关闭chan会在发送方产生panic: send on closed channel错误
	close(tasks)
	fmt.Println("All tasks sent")
	<-canQuit
}
