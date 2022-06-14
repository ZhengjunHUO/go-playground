package cmd

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"github.com/spf13/cobra"
)


// k8s config received from --kubeconfig flag
var (
	kconfig = genericclioptions.NewConfigFlags(true)
	clientset *kubernetes.Clientset
	namespace = "default"
)

var rootCmd = &cobra.Command{
	Use: "kubectlx",
	Short: "kubectl implemented by client-go",
	Long: `kubectl implemented by client-go`,
}

func Execute() {
	kconfig.AddFlags(rootCmd.PersistentFlags())
        _ = rootCmd.Execute()
}

func init() {
	// Get clientset for built-in resources
	conf, _ := kconfig.ToRESTConfig()
	clientset, _ = kubernetes.NewForConfig(conf)
}
