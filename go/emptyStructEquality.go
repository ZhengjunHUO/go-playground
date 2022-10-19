package main

import (
	"fmt"
	"unsafe"
)

type empty struct {}

type notEmpty struct {
	id int
}

// 不相同
func compareEmpty() {
	a := empty{}
	b := empty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))
	fmt.Println("&a's size: ", unsafe.Sizeof(a))
	fmt.Println("&b's size: ", unsafe.Sizeof(b))
}

// 相同
func compareEmptyPrintAddr() {
	a := empty{}
	b := empty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))

	// 相比compareEmpty()增加一行printf
	fmt.Printf("%v and %v is equal ? %v\n", &a, &b, (&a == &b) )

	fmt.Println("&a's size: ", unsafe.Sizeof(a))
	fmt.Println("&b's size: ", unsafe.Sizeof(b))
}

// 整合了上述两个func，无论bool为true或false，结果都是相同
func compareEmptyWithB(printAddr bool) {
	//a := struct{}{}
	//b := struct{}{}
	a := empty{}
	b := empty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))

	// 应该是和编译器有关，false时尽管在runtime时并没有执行这一行，
	// 在调用该func时a和b还是被初始化指向了一个地址
	if printAddr {
		//fmt.Printf("%T and %T is equal ? %v\n", &a, &b, (&a == &b) )
		fmt.Printf("%v and %v is equal ? %v\n", &a, &b, (&a == &b) )
	}

	fmt.Println("&a's size: ", unsafe.Sizeof(a))
	fmt.Println("&b's size: ", unsafe.Sizeof(b))
}

// 不同(正常)
func compareNotEmpty() {
	a := notEmpty{}
	b := notEmpty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))
	fmt.Println("&a's size: ", unsafe.Sizeof(a))
	fmt.Println("&b's size: ", unsafe.Sizeof(b))
}

// 不同(正常)
func compareNotEmptyPrintAddr() {
	a := notEmpty{}
	b := notEmpty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))

	fmt.Printf("%v and %v is equal ? %v\n", &a, &b, (&a == &b) )

	fmt.Println("&a's size: ", unsafe.Sizeof(a))
	fmt.Println("&b's size: ", unsafe.Sizeof(b))
}

func main() {
	fmt.Println("=> Result of compareEmpty(): ")
	compareEmpty()
	fmt.Println("=> Result of compareEmptyPrintAddr(): ")
	compareEmptyPrintAddr()
	fmt.Println("=> Result of compareEmptyWithB(false): ")
	compareEmptyWithB(false)
	fmt.Println("=> Result of compareEmptyWithB(true): ")
	compareEmptyWithB(true)

	fmt.Println("=> Result of compareNotEmpty(): ")
	compareNotEmpty()
	fmt.Println("=> Result of compareNotEmptyPrintAddr(): ")
	compareNotEmptyPrintAddr()
}
