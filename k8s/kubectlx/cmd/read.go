package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"
)

var readCmd = &cobra.Command{
	Use: "read <path_to_yaml/json>...",
	Short: "read in k8s object from yaml/json file(s)",
	Long: `read in k8s object from yaml/json file(s)`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		builder := resource.NewBuilder(kconfig)

		namespace := ""
		if kconfig.Namespace != nil && len(*kconfig.Namespace) > 0 {
			namespace = *kconfig.Namespace
		}
		// builder won't work correctly if specified namespace differs from the namespace in the manifest
		enforceNS := namespace != ""

		_ = builder.WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
			NamespaceParam(namespace).DefaultNamespace().
			FilenameParam(enforceNS, &resource.FilenameOptions{Filenames: args},).Do().
			Visit(func(info *resource.Info, _ error) error {
				fmt.Println("[DEBUG] Find a resource: ")
				fmt.Println(info.Object)
				return nil
			})
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	kconfig.AddFlags(readCmd.PersistentFlags())
}
