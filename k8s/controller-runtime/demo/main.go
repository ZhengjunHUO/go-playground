package main

import (
	"fmt"
	"os"
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1 "k8s.io/api/core/v1"
)

func main() {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("Create client failed !")
		os.Exit(1)
	}
	fmt.Println("Create client succeeded !")

	poList := &corev1.PodList{}
	if err = cl.List(context.Background(), poList, client.InNamespace("kube-system")); err != nil {
		fmt.Println("List pods under ns kube-system failed !")
		os.Exit(1)
	}
	fmt.Println("List pods under ns kube-system succeeded !")
}
