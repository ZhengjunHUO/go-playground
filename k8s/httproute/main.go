package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

func main() {
	var (
		kubeconfig string
		namespace  string
		name       string
		action     string
		gateway    string
		backend    string
		labels     string
	)

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "path to kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	}
	flag.StringVar(&namespace, "namespace", "default", "namespace of the HTTPRoute")
	flag.StringVar(&name, "name", "", "name of the HTTPRoute (required)")
	flag.StringVar(&action, "action", "create", "action to perform: create, get, delete, patch")
	flag.StringVar(&gateway, "gateway", "", "parent gateway name (required for create)")
	flag.StringVar(&backend, "backend", "", "backend service name (required for create)")
	flag.StringVar(&labels, "labels", "", `JSON labels for patch action, e.g. '{"env":"prod"}' (required for patch)`)
	flag.Parse()

	if name == "" {
		fmt.Fprintln(os.Stderr, "error: --name is required")
		flag.Usage()
		os.Exit(1)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	client, err := gatewayclient.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating gateway client: %v\n", err)
		os.Exit(1)
	}

	routes := client.GatewayV1().HTTPRoutes(namespace)

	switch action {
	case "create":
		if gateway == "" || backend == "" {
			fmt.Fprintln(os.Stderr, "error: --gateway and --backend are required for create")
			os.Exit(1)
		}

		pathType := gatewayv1.PathMatchPathPrefix
		pathValue := "/"
		port := gatewayv1.PortNumber(8080)
		parentNs := gatewayv1.Namespace(namespace)

		httproute := &gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: gatewayv1.HTTPRouteSpec{
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{
						{
							Name:      gatewayv1.ObjectName(gateway),
							Namespace: &parentNs,
						},
					},
				},
				Rules: []gatewayv1.HTTPRouteRule{
					{
						Matches: []gatewayv1.HTTPRouteMatch{
							{
								Path: &gatewayv1.HTTPPathMatch{
									Type:  &pathType,
									Value: &pathValue,
								},
							},
						},
						BackendRefs: []gatewayv1.HTTPBackendRef{
							{
								BackendRef: gatewayv1.BackendRef{
									BackendObjectReference: gatewayv1.BackendObjectReference{
										Name: gatewayv1.ObjectName(backend),
										Port: &port,
									},
								},
							},
						},
					},
				},
			},
		}

		result, err := routes.Create(context.Background(), httproute, metav1.CreateOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating HTTPRoute: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("HTTPRoute %s/%s created (resourceVersion: %s)\n",
			result.Namespace, result.Name, result.ResourceVersion)

	case "get":
		result, err := routes.Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "HTTPRoute %s/%s not found: %v\n", namespace, name, err)
			os.Exit(1)
		}
		fmt.Printf("HTTPRoute %s/%s exists (resourceVersion: %s)\n",
			result.Namespace, result.Name, result.ResourceVersion)

	case "delete":
		err := routes.Delete(context.Background(), name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error deleting HTTPRoute: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("HTTPRoute %s/%s deleted\n", namespace, name)

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
		result, err := routes.Patch(context.Background(), name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error patching HTTPRoute: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("HTTPRoute %s/%s patched (resourceVersion: %s)\n",
			result.Namespace, result.Name, result.ResourceVersion)

	default:
		fmt.Fprintf(os.Stderr, "unknown action %q, use: create, get, delete, patch\n", action)
		os.Exit(1)
	}
}
