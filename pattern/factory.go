package main

import (
	"fmt"
)

// public interface
type Cluster interface {
	GetInfo()
}

// private struct, implements interface
type k8sCluster struct {
	name		string
	cni		string
	size		int
}

func (kc k8sCluster) GetInfo() {
	fmt.Printf("k8s cluster [%s] of %d nodes is using %s as CNI.\n", kc.name, kc.size, kc.cni)
}

// factory function, create new instance
func NewK8sCluster(name, cni string, size int) Cluster {
	return &k8sCluster{
		name: name,
		cni: cni,
		size: size,
	}
}

func main() {
	c := NewK8sCluster("huo", "cilium", 10)
	c.GetInfo()
}
