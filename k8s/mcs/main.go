package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
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
		action     string
		labels     string
	)

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "path to kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	}
	flag.StringVar(&namespace, "namespace", "default", "namespace of the ServiceExport")
	flag.StringVar(&name, "name", "", "name of the service to export (required)")
	flag.StringVar(&action, "action", "apply", "action to perform: apply, delete, patch")
	flag.StringVar(&labels, "labels", "", `JSON labels for patch action, e.g. '{"env":"prod"}'`)
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

	res := dynClient.Resource(serviceExportGVR).Namespace(namespace)

	switch action {
	case "apply":
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
		result, err := res.Apply(context.Background(), name, obj, metav1.ApplyOptions{
			FieldManager: "serviceexport-deployer",
			Force:        true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error applying ServiceExport: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("ServiceExport %s/%s applied (resourceVersion: %s)\n",
			result.GetNamespace(), result.GetName(), result.GetResourceVersion())

	case "get":
		result, err := res.Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "ServiceExport %s/%s not found: %v\n", namespace, name, err)
			os.Exit(1)
		}
		fmt.Printf("ServiceExport %s/%s exists (resourceVersion: %s)\n",
			result.GetNamespace(), result.GetName(), result.GetResourceVersion())

	case "delete":
		err := res.Delete(context.Background(), name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error deleting ServiceExport: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("ServiceExport %s/%s deleted\n", namespace, name)

	case "patch":
		if labels == "" {
			fmt.Fprintln(os.Stderr, "error: --labels is required for patch action")
			os.Exit(1)
		}
		var labelMap map[string]string
		if err := json.Unmarshal([]byte(labels), &labelMap); err != nil {
			fmt.Fprintf(os.Stderr, "error parsing --labels JSON: %v\n", err)
			os.Exit(1)
		}
		patch := map[string]interface{}{
			"metadata": map[string]interface{}{
				"labels": labelMap,
			},
		}
		patchBytes, _ := json.Marshal(patch)
		result, err := res.Patch(context.Background(), name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error patching ServiceExport: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("ServiceExport %s/%s patched (resourceVersion: %s)\n",
			result.GetNamespace(), result.GetName(), result.GetResourceVersion())

	default:
		fmt.Fprintf(os.Stderr, "unknown action %q, use: apply, get, delete, patch\n", action)
		os.Exit(1)
	}
}
