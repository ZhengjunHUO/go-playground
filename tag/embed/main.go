package main

import (
	"fmt"
	_ "embed"
)

//go:embed hello.txt
var salutation []byte

func main() {
	fmt.Print(string(salutation))
}
