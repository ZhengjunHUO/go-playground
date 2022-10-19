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
	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))
}

// 相同
func compareEmptyPrintAddr() {
	a := empty{}
	b := empty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))

	// 相比compareEmpty()增加一行printf
	fmt.Printf("%p and %p is equal ? %v\n", &a, &b, (&a == &b) )

	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))
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
		fmt.Printf("%p and %p is equal ? %v\n", &a, &b, (&a == &b) )
	}

	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))
}

// 不同(正常)
func compareNotEmpty() {
	a := notEmpty{}
	b := notEmpty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))
	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))
}

// 不同(正常)
func compareNotEmptyPrintAddr() {
	a := notEmpty{}
	b := notEmpty{}

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))

	fmt.Printf("%p and %p is equal ? %v\n", &a, &b, (&a == &b) )

	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))
}

// 和compareEmpty()类似，不相同
func compareEmptySlice() {
	var a [100]empty
	var b [100]empty

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))

	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))
}

// 相同
func compareEmptySlicePrintAddr() {
	var a [100]empty
	var b [100]empty

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))

	fmt.Printf("%p and %p is equal ? %v\n", &a, &b, (&a == &b) )

	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))

	fmt.Printf("&a[10] and &b[10] is equal ? %v\n", (&a[0] == &b[0]))
}

// a和b已经被初始化(有指向一个长为24的slice头)所以肯定不同
// 但是slice中的空结构元素之间仍然相同
func compareEmptySliceMake() {
	a := make([]empty, 100)
	b := make([]empty, 100)

	fmt.Println(&a == &b)
	fmt.Printf("&a and &b is equal ? %v\n", (&a == &b))
	fmt.Printf("%p and %p is equal ? %v\n", &a, &b, (&a == &b) )
	fmt.Printf("&a[10] and &b[10] is equal ? %v\n", (&a[10] == &b[10]))
	fmt.Printf("&a[23] and &b[64] is equal ? %v\n", (&a[23] == &b[64]))
	fmt.Printf("&a[31] and &a[47] is equal ? %v\n", (&a[31] == &a[47]))

	fmt.Println("a's size: ", unsafe.Sizeof(a))
	fmt.Println("b's size: ", unsafe.Sizeof(b))
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

	fmt.Println("=> Result of compareEmptySlice(): ")
	compareEmptySlice()

	fmt.Println("=> Result of compareEmptySlicePrintAddr(): ")
	compareEmptySlicePrintAddr()

	fmt.Println("=> Result of compareEmptySliceMake(): ")
	compareEmptySliceMake()
}
