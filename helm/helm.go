package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/chart"
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

	// Load user provided values from file
	userValuesPath := "./helm-chart/examples/nolookup.yaml"
	userValues, err := chartutil.ReadValuesFile(userValuesPath)
	if err != nil {
		fmt.Println("Error occurred loading values:", err)
		os.Exit(1)
	}

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
	userValuesPath := "./helm-chart/examples/nolookup.yaml"
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

	// Create client action
	client := action.NewInstall(actionConfig)
	client.ReleaseName = "obs"
	client.Namespace = namespace
	client.CreateNamespace = true

	// Load chart
	chartRef, err := client.ChartPathOptions.LocateChart(chartPath, cli.New())
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
	rslt, err := client.Run(chart, coalescedValues)
	if err != nil {
		log.Fatalf("Error occurred installing chart: %v", err)
	}

	// Print release info
	fmt.Printf("Release [%s] installed\n", rslt.Name)
}

func poc_helm_upgrade() {
	namespace := "opensee-obs-agents"
	chartPath := "./helm-chart"

	// Load user provided values from file
	userValuesPath := "./helm-chart/examples/nolookup.yaml"
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

	// Create client action
	client := action.NewUpgrade(actionConfig)
	client.Namespace = namespace

	// Load chart
	chartRef, err := client.ChartPathOptions.LocateChart(chartPath, cli.New())
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
	rslt, err := client.Run("obs", chart, coalescedValues)
	if err != nil {
		log.Fatalf("Error occurred upgrading chart: %v", err)
	}

	// Print release info
	fmt.Printf("Release %s upgraded.\n", rslt.Name)
}

func poc_helm_uninstall() {
	namespace := "opensee-obs-agents"

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = "/home/huo/.kube/config"
	}

	actionConfig, err := newActionConfig(kubeconfig, namespace)
	if err != nil {
		log.Fatalf("Error occurred initializing Helm: %v", err)
	}

	client := action.NewUninstall(actionConfig)
	client.DeletionPropagation = "foreground"

	rslt, err := client.Run("obs")
	if err != nil {
		log.Fatalf("Error occurred installing chart: %v", err)
	}

	fmt.Printf("Release %s uninstalled.\n", rslt.Release.Name)
}

func poc_helm_status() {
	namespace := "opensee-obs-agents"

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = "/home/huo/.kube/config"
	}

	actionConfig, err := newActionConfig(kubeconfig, namespace)
	if err != nil {
		log.Fatalf("Error occurred initializing Helm: %v", err)
	}

	client := action.NewStatus(actionConfig)
	client.ShowResources = true

	rslt, err := client.Run("obs")
	if err != nil {
		log.Fatalf("Error occurred calling status on chart: %v", err)
	}

	fmt.Printf("[Status]: \nNAME: %s\nLAST DEPLOYED: %v\nNAMESPACE: %s\nSTATUS: %v\nREVISION: %v\n", rslt.Name, rslt.Info.LastDeployed, rslt.Namespace, rslt.Info.Status, rslt.Version)
	//fmt.Printf("%v\n", rslt.Info.Resources)
}

type HelmEngine struct {
	namespace    string
	releaseName  string
	actionConfig *action.Configuration
	chart        *chart.Chart
	userValues   chartutil.Values
}

type HelmEngineOption func(*HelmEngine)

func withHelmNamespace(namespace string) HelmEngineOption {
	return func(h *HelmEngine) {
		h.namespace = namespace
	}
}

func withHelmReleaseName(releaseName string) HelmEngineOption {
	return func(h *HelmEngine) {
		h.releaseName = releaseName
	}
}

func withHelmActionConfig(kubeconfigPath string) HelmEngineOption {
	return func(h *HelmEngine) {
		actionConfig, err := newActionConfig(kubeconfigPath, h.namespace)
		if err != nil {
			log.Fatalf("Error occurred initializing Helm: %v\n", err)
		}

		h.actionConfig = actionConfig
	}
}

func withHelmChart(chartPath string) HelmEngineOption {
	return func(h *HelmEngine) {
		chart, err := loader.Load(chartPath)
		if err != nil {
			log.Fatalf("Error occurred loading helm chart: %v\n", err)
		}

		h.chart = chart
	}
}

func withHelmUserValuesFromFile(userValuesPath string) HelmEngineOption {
	return func(h *HelmEngine) {
		userValues, err := chartutil.ReadValuesFile(userValuesPath)
		if err != nil {
			log.Fatalf("Error occurred loading values: %v\n", err)
		}

		h.userValues = userValues
	}
}

func withHelmUserValues(userValues chartutil.Values) HelmEngineOption {
	return func(h *HelmEngine) {
		h.userValues = userValues
	}
}

func NewHelmEngine(namespace, releaseName string, hops ...HelmEngineOption) *HelmEngine {
	helm := &HelmEngine{
		namespace:   namespace,
		releaseName: releaseName,
	}

	for _, f := range hops {
		f(helm)
	}

	return helm
}

func (h *HelmEngine) install() {
	// Create client action
	client := action.NewInstall(h.actionConfig)
	client.ReleaseName = h.releaseName
	client.Namespace = h.namespace
	client.CreateNamespace = true

	// Prepare values
	coalescedValues, err := chartutil.CoalesceValues(h.chart, h.userValues)
	if err != nil {
		log.Fatal(err)
	}

	// Install chart
	rslt, err := client.Run(h.chart, coalescedValues)
	if err != nil {
		log.Fatalf("Error occurred installing release: %v", err)
	}

	// Print release info
	fmt.Printf("Release [%s] installed\n", rslt.Name)
}

func (h *HelmEngine) upgrade() {
	// Create client action
	client := action.NewUpgrade(h.actionConfig)
	client.Namespace = h.namespace

	// Prepare values
	coalescedValues, err := chartutil.CoalesceValues(h.chart, h.userValues)
	if err != nil {
		log.Fatal(err)
	}

	// Upgrade chart
	rslt, err := client.Run(h.releaseName, h.chart, coalescedValues)
	if err != nil {
		log.Fatalf("Error occurred upgrading release: %v", err)
	}

	// Print release info
	fmt.Printf("Release %s upgraded.\n", rslt.Name)
}

func (h *HelmEngine) uninstall() {
	client := action.NewUninstall(h.actionConfig)
	client.DeletionPropagation = "foreground"

	rslt, err := client.Run(h.releaseName)
	if err != nil {
		log.Fatalf("Error occurred uninstalling release: %v", err)
	}

	fmt.Printf("Release %s uninstalled.\n", rslt.Release.Name)
}

func (h *HelmEngine) status() {
	client := action.NewStatus(h.actionConfig)
	client.ShowResources = true

	rslt, err := client.Run(h.releaseName)
	if err != nil {
		log.Fatalf("Error occurred checking status on release: %v", err)
	}

	fmt.Printf("[Status]: \nNAME: %s\nLAST DEPLOYED: %v\nNAMESPACE: %s\nSTATUS: %v\nREVISION: %v\n", rslt.Name, rslt.Info.LastDeployed, rslt.Namespace, rslt.Info.Status, rslt.Version)
}

func (h *HelmEngine) updateUserValues(userValues chartutil.Values) {
	h.userValues = userValues
}

func main() {
	userValues := chartutil.Values{
		"environmentName":      "localtest",
		"namespace":            "opensee-obs-agents",
		"disableIntrospection": true,
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
			"watchedComponents": map[string]interface{}{
				"shardsCount":   2,
				"replicasCount": 2,
			},
		},
		"otelKeeperAgent": map[string]interface{}{
			"enabled": true,
			"watchedComponents": map[string]interface{}{
				"stsCount": 3,
			},
		},
		"otelPostgresqlAgent": map[string]interface{}{
			"enabled": true,
			"watchedComponents": map[string]interface{}{
				"stsCount": 2,
			},
		},
		"otelCalculatorAgent": map[string]interface{}{
			"enabled": true,
			"watchedComponents": map[string]interface{}{
				"stsCount": 4,
			},
		},
	}

	helm := NewHelmEngine("opensee-obs-agents", "obs",
		withHelmActionConfig("/home/huo/.kube/config"),
		withHelmChart("./helm-chart"),
		withHelmUserValues(userValues))

	helm.install()
	time.Sleep(60*time.Second)
	helm.status()

	userValuesPath := "./helm-chart/examples/nolookup.yaml"
	newUserValues, err := chartutil.ReadValuesFile(userValuesPath)
	if err != nil {
		fmt.Println("Error occurred loading values:", err)
		os.Exit(1)
	}

	helm.updateUserValues(newUserValues)
	helm.upgrade()
	time.Sleep(60*time.Second)
	helm.status()
	helm.uninstall()
	//poc_helm_template()
	//poc_helm_install()
	//poc_helm_upgrade()
	//poc_helm_status()
	//poc_helm_uninstall()
}
