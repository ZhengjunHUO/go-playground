package main

import (
	"fmt"
)

var functions = []string{
	"Standard Function 1",
	"Standard Function 2",
	"Standard Function 3",
}

func main() {
	for i := range functions {
		fmt.Printf("Load %s !\n", functions[i])
	}
}
