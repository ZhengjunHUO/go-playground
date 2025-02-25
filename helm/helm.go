package main

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/kube"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

func do_helm_template(chartPath string, userValues chartutil.Values, release map[string]interface{}) map[string]string {
	// Load Helm chart
	chart, err := loader.Load(chartPath)
	if err != nil {
		fmt.Println("Error occurred loading helm chart:", err)
		os.Exit(1)
	}

	// Prepare values
	coalescedValues, err := chartutil.CoalesceValues(chart, userValues)
	if err != nil {
		log.Fatal(err)
	}

	values := map[string]interface{}{
		"Values":  coalescedValues,
		"Release": release,
	}
	//fmt.Println(values)

	// Render templates
	e := engine.Engine{}
	rslt, err := e.Render(chart, values)

	if err != nil {
		fmt.Println("Error occurred rendering chart:", err)
		os.Exit(1)
	}

	return rslt
}

func poc_helm_template() {
	chartPath := "./helm-chart"

	release := map[string]interface{}{
		"Name":    "obs",
		"Service": "Helm",
	}

	userValues := chartutil.Values{
		"environmentName": "localtest",
		"namespace":       "opensee-obs-agents",
		"image": map[string]interface{}{
			"repository": "otel/opentelemetry-collector-contrib",
			"tag":        "0.97.0",
		},
		"imagePullSecrets": []map[string]interface{}{
			map[string]interface{}{
				"name": "regcred",
			},
		},
		"centralTelemetry": map[string]interface{}{
			"endpoint": "https://foo.bar.com",
			"insecure": false,
			"headers": map[string]interface{}{
				"authorization": "Basic SECRET",
			},
			"useHttpExporter": true,
		},
		"otelEdgeServer": map[string]interface{}{
			"enabled": true,
		},
		"otelClickhouseAgent": map[string]interface{}{
			"enabled": true,
		},
		"otelKeeperAgent": map[string]interface{}{
			"enabled": true,
		},
		"otelPostgresqlAgent": map[string]interface{}{
			"enabled": true,
		},
		"otelCalculatorAgent": map[string]interface{}{
			"enabled": true,
		},
	}

	// Load user provided values from file
	/*
		userValuesPath := "./helm-chart/examples/minimal.yaml"
		userValues, err := chartutil.ReadValuesFile(userValuesPath)
		if err != nil {
			fmt.Println("Error occurred loading values:", err)
			os.Exit(1)
		}
	*/

	rslt := do_helm_template(chartPath, userValues, release)
	for name, content := range rslt {
		fmt.Println("=====", name, "=====")
		fmt.Println(content)
	}
}

func newActionConfig(kubeconfigPath string, namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)

	config := kube.GetConfig(kubeconfigPath, "", namespace)

	// Initialize Helm action configuration
	if err := actionConfig.Init(config, namespace, "secrets", log.Printf); err != nil {
		return nil, fmt.Errorf("Failed to initialize Helm action: %w", err)
	}

	return actionConfig, nil
}

func poc_helm_install() {
	namespace := "opensee-obs-agents"
	chartPath := "./helm-chart"

	// Load user provided values from file
	userValuesPath := "./helm-chart/examples/minimal.yaml"
	userValues, err := chartutil.ReadValuesFile(userValuesPath)
	if err != nil {
		fmt.Println("Error occurred loading values:", err)
		os.Exit(1)
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = "/home/huo/.kube/config"
	}

	actionConfig, err := newActionConfig(kubeconfig, namespace)
	if err != nil {
		log.Fatalf("Error occurred initializing Helm: %v", err)
	}

	// Create install action
	install := action.NewInstall(actionConfig)
	install.ReleaseName = "obs"
	install.Namespace = namespace
	install.CreateNamespace = true

	// Load chart
	chartRef, err := install.ChartPathOptions.LocateChart(chartPath, cli.New())
	if err != nil {
		log.Fatalf("Error occurred locating chart: %v", err)
	}

	chart, err := loader.Load(chartRef)
	if err != nil {
		fmt.Println("Error occurred loading helm chart:", err)
		os.Exit(1)
	}

	// Prepare values
	coalescedValues, err := chartutil.CoalesceValues(chart, userValues)
	if err != nil {
		log.Fatal(err)
	}

	// Install chart
	rslt, err := install.Run(chart, coalescedValues)
	if err != nil {
		log.Fatalf("Error occurred installing chart: %v", err)
	}

	// Print release info
	fmt.Printf("Release installed: %s\n", rslt.Name)
}

func main() {
	poc_helm_install()
}
