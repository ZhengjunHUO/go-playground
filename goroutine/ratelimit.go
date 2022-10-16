package main

import (
	"fmt"
	"time"
)

var (
	allowBurstNum = 3
	limitedConsumInterval = 1*time.Second
)

// 控制从chan中的读取速度为1次/秒，允许burst为可同时读取3次
func main() {
	// 1. 准备limiter
	limiterWithBurst := make(chan time.Time, allowBurstNum)
	// 初始状态: 事先填满buffer (允许burst)
	for i:=0; i<allowBurstNum; i++ {
		limiterWithBurst <- time.Now()
	}

	// 每秒尝试向其中写入一个值
	go func() {
		for t := range time.Tick(limitedConsumInterval){
			limiterWithBurst <- t
		}
	}()

	// 2. 生成任务
	tasks := make(chan int, 5)
	go func() {
		for i:=0; i<10; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	// 3. 处理任务
	for t := range tasks {
		<-limiterWithBurst
		fmt.Printf("[%v] Processing task %d\n", time.Now(), t)
	}
}
