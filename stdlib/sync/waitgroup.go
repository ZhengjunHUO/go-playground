package main

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
)

var wg sync.WaitGroup

func InitCluster(name string) {
	fmt.Printf("[INFO] Initializing %s ... \n", name)
	defer wg.Done()
	// simulate the init process
	time.Sleep(time.Duration(rand.Intn(10))*time.Second)

	fmt.Printf("[ OK ] Cluster %s initialized.\n", name)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	clusterNames := []string{"wang", "fufu", "huo"}
	for i := range clusterNames {
		wg.Add(1)
		go InitCluster(clusterNames[i])
	}

	wg.Wait()
	fmt.Println("[ OK ] All jobs done. Quit.")
}
