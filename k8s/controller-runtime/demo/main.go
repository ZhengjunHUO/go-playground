package main

import (
	"fmt"
	"os"
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func main() {
	// Prepare client
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("Create client failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Create client succeeded !")

	// Prepare label selector for control plane pods
	slct, err := labels.Parse("component in (kube-apiserver, kube-scheduler, kube-controller-manager)")
	if err != nil {
		fmt.Println("Create selecter failed: ", err)
		os.Exit(1)
	}

	// Get control plane pods in kube-system
	poList := &corev1.PodList{}
	if err = cl.List(context.Background(), poList, client.InNamespace("kube-system"), client.MatchingLabelsSelector{Selector: slct}); err != nil {
		fmt.Println("List pods under ns kube-system failed: ", err)
		os.Exit(1)
	}

	// Print out control plane pods' name
	for _, po := range poList.Items {
		fmt.Printf("Find pod[%s]\n", po.Name)
	}
}
