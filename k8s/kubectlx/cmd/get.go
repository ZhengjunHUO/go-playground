package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"
)

var getCmd = &cobra.Command{
	Use: "get <resource type>",
	Short: "Print k8s resources",
	Long: `Print k8s resources`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化一个resource builder
		builder := resource.NewBuilder(kconfig)

		// 如果通过-n <namespace>指定了namespace则覆盖默认
		if kconfig.Namespace != nil && len(*kconfig.Namespace) > 0 {
			namespace = *kconfig.Namespace
		}

		// 依靠scheme初始化各种k8s资源
		// 使用namespace和从命令行参数中得到的需要查看的资源种类来缩小查找范围
		rslt, _ := builder.WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
			NamespaceParam(namespace).ResourceTypeOrNameArgs(true, args...).Do().Object()

		fmt.Println(rslt)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	kconfig.AddFlags(getCmd.PersistentFlags())
}
