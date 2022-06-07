package main

import (
	"log"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"github.com/spf13/cobra"
)

// k8s config received from --kubeconfig flag
var kconfig = genericclioptions.NewConfigFlags(true)

var rootCmd = &cobra.Command{
	Use: "kubectlx",
	Short: "kubectl implemented by client-go",
	Long: `kubectl implemented by client-go`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, _ := kconfig.ToRESTConfig()
		// Get clientset for built-in resources
		clientset, _ := kubernetes.NewForConfig(conf)

		// Get k8s cluster's version
		version, _ := clientset.Discovery().ServerVersion()
		log.Println("[INFO] k8s cluster version: ", version)
	},
}

func main() {
	kconfig.AddFlags(rootCmd.PersistentFlags())
	_ = rootCmd.Execute()
}
