package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var serviceExportGVR = schema.GroupVersionResource{
	Group:    "multicluster.x-k8s.io",
	Version:  "v1alpha1",
	Resource: "serviceexports",
}

func main() {
	var (
		kubeconfig string
		namespace  string
		name       string
	)

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "path to kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	}
	flag.StringVar(&namespace, "namespace", "default", "namespace of the ServiceExport")
	flag.StringVar(&name, "name", "", "name of the service to export (required)")
	flag.Parse()

	if name == "" {
		fmt.Fprintln(os.Stderr, "error: --name is required")
		flag.Usage()
		os.Exit(1)
	}

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	dynClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating dynamic client: %v\n", err)
		os.Exit(1)
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "multicluster.x-k8s.io/v1alpha1",
			"kind":       "ServiceExport",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
		},
	}

	result, err := dynClient.Resource(serviceExportGVR).
		Namespace(namespace).
		Apply(context.Background(), name, obj, metav1.ApplyOptions{
			FieldManager: "serviceexport-deployer",
			Force:        true,
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error applying ServiceExport: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ServiceExport %s/%s applied (resourceVersion: %s)\n",
		result.GetNamespace(), result.GetName(), result.GetResourceVersion())
}
