package main

import (
	"fmt"
)

type K8sCNI int

const (
	CILIUM K8sCNI = iota
	FLANNEL
	WEAVE
	CALICO 
)

func (c K8sCNI) String() string {
    return [...]string{"Cilium", "Flannel", "Weave", "Calico"}[c]
}

func showCniNumber(c K8sCNI) {
	fmt.Println(c)
}

func main() {
	/*
	i := 3
	showCniNumber(K8sCNI(i))
	*/

	showCniNumber(CILIUM)
	showCniNumber(FLANNEL)
	showCniNumber(WEAVE)
	showCniNumber(CALICO)
}
