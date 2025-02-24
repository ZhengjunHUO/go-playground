package main

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

func main() {
	chartPath := "./helm-chart"
	userValuesPath := "./helm-chart/examples/minimal.yaml"

	// Load Helm chart
	chart, err := loader.Load(chartPath)
	if err != nil {
		fmt.Println("Error occurred loading helm chart:", err)
		os.Exit(1)
	}

	// Load user provided values from file
	userValues, err := chartutil.ReadValuesFile(userValuesPath)
	if err != nil {
		fmt.Println("Error occurred loading values:", err)
		os.Exit(1)
	}

	//temp := chartutil.CoalesceTables(chart.Values, defaultValues)
	coalescedValues, err := chartutil.CoalesceValues(chart, userValues)
	if err != nil {
		log.Fatal(err)
	}

	release := map[string]interface{}{
		"Name":    "obs",
		"Service": "Helm",
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

	for name, content := range rslt {
		fmt.Println("=====", name, "=====")
		fmt.Println(content)
	}
}
