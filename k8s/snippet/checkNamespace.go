package main

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

func storageClassExists(clientset *kubernetes.Clientset, scName string) (bool, error) {
	_, err := clientset.StorageV1().StorageClasses().Get(context.Background(), scName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getServerVersion(clientset *kubernetes.Clientset) (string, error) {
	serverVersion, err := clientset.Discovery().ServerVersion()
	return serverVersion.GitVersion, err
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

	// Get server version
	serverVersion, err := getServerVersion(clientset)
	if err != nil {
		log.Fatalf("Error getting server version: %v", err)
	}
	log.Printf("Kubernetes Server Version: %s\n", serverVersion)

	// Check if namespace exists
	nsName := "default"
	exists, err := NamespaceExists(clientset, nsName)
	if err != nil {
		panic(err.Error())
	}

	if exists {
		log.Printf("Namespace %s exists!\n", nsName)
	} else {
		log.Printf("Namespace %s does not exist!\n", nsName)
	}

	// Check if storage exists
	scName := "premium-rwo"
	exists, err = storageClassExists(clientset, scName)
	if err != nil {
		panic(err.Error())
	}

	if exists {
		log.Printf("Storage class %s exists!\n", scName)
	} else {
		log.Printf("Storage class %s does not exist!\n", scName)
	}
}
