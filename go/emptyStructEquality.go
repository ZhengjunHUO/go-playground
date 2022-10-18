package main

import (
	"fmt"
)

type test struct {
	//id int
}

func main() {
	//a := struct{}{}
	//b := struct{}{}
	a := test{}
	b := test{}
	fmt.Println(&a == &b)
	fmt.Printf("is equal ? %v\n", (&a == &b))
	//fmt.Printf("%T and %T is equal ? %v\n", &a, &b, (&a == &b) )
	//fmt.Printf("%v and %v is equal ? %v\n", &a, &b, (&a == &b) )
	//fmt.Printf("%p == %p: %v\n", &a, &b, &a == &b)
}
