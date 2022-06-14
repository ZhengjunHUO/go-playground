package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Get kubernetes cluster's version",
	Long: `Get kubernetes cluster's version`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get k8s cluster's version
		version, _ := clientset.Discovery().ServerVersion()
		fmt.Println("[INFO] k8s cluster version: ", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	kconfig.AddFlags(versionCmd.PersistentFlags())
}
