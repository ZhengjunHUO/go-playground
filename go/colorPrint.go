package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	fmt.Println("\033[1;31;42mBold red text with green background\033[0m")

	suffix := "\033[0m"
	for i := 0; i < 10; i++ {
		prefix := fmt.Sprintf("\033[38;5;%vm", r.Intn(256) + 1)
		fmt.Printf("%sRustacean here !%s\n", prefix, suffix)
	}
}
