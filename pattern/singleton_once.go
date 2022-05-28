package main

import (
	"fmt"
	"sync"
)

type k8s struct {
	Name	string
	Size	int
}

// 一个private的单例
var (
	singleK8s 	*k8s
	once		sync.Once
)

// 返回单例，如果不存在创建它，暴露给外部调用
func GetK8s() *k8s {
	// 如不存在则先创建
	if singleK8s == nil {
		once.Do(func () {
			fmt.Println("[INFO] Creating k8s singleton ...")
			singleK8s = &k8s{"Global", 16}
			fmt.Println("[INFO] Done !")
		})
	}else{
		fmt.Println("[INFO] Return exist k8s.")
	}

	// 返回单例
	return singleK8s
}

func main() {
	for i := 0; i < 10; i++ {
		go GetK8s()
	}

	fmt.Scanln()
}
