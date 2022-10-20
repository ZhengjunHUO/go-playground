package main

import (
	"fmt"
)

type list[V any] struct {
	head, tail *listNode[V]
}

type listNode[V any] struct {
        val V
        next *listNode[V]
}

func (l *list[V]) Add(v V) {
	if l.tail == nil {
		l.head = &listNode[V]{v, nil}
		l.tail = l.head
		return
	}

	l.tail.next = &listNode[V]{v, nil}
	l.tail = l.tail.next
}

func newList[V any](l []V) *list[V] {
	n := len(l)
	if n == 0 {
		return nil
	}

	rslt := &list[V]{}
	for i := 0; i < n; i++ {
		rslt.Add(l[i])
	}

	return rslt
}

func (l *list[V]) printAll() {
	if l.head != nil {
	        curr := l.head
	        for {
			fmt.Printf("%v", curr.val)

	                if curr.next == nil {
				fmt.Println()
	                        break
	                }
	                curr = curr.next
			fmt.Printf(", ")
	        }
	}
}

func main() {
	l := newList([]int{3,2,0,-4})
	l.printAll()
}
