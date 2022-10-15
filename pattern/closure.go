package main

import (
	"fmt"
)

func incremCounter() func () int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

func main() {
	c := incremCounter()

	for i:=0; i<3; i++ {
		fmt.Println(c())
	}

	// reset counter
	c = incremCounter()
	for i:=0; i<3; i++ {
		fmt.Println(c())
	}
}
