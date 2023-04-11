package main

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NamespaceExists(clientset *kubernetes.Clientset, namespace string) (bool, error) {
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func main() {
	// Use the path to the Kubernetes config file
	config, err := clientcmd.BuildConfigFromFlags("", "/home/huo/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// Create a new Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Check if namespace "my-namespace" exists
	exists, err := NamespaceExists(clientset, "default")
	if err != nil {
		panic(err.Error())
	}

	if exists {
		fmt.Println("Namespace exists!")
	} else {
		fmt.Println("Namespace does not exist!")
	}
}
