package main

import (
	"fmt"
	"github.com/ZhengjunHUO/playground/go-build/app"
)

var version = "undefined"

func main() {
	fmt.Println("version:", version)
	fmt.Println("app.Author:", app.Author)
	fmt.Println("app.Date:", app.Date)
}
