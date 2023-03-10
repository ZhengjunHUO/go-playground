package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	patches []string
	mainConfig string
	envConfig string
)

var ConfCmd = &cobra.Command{
	Use:   "conf",
	Short: "generate config yaml",
	Run: func(cmd *cobra.Command, args []string) {
		envViper := viper.New()
		envViper.SetConfigFile(envConfig)
		if err := envViper.ReadInConfig(); err != nil {
			fmt.Println("Error reading env config file: ", err)
		}

		for _, v := range envViper.AllKeys() {
			fmt.Printf("[Debug] env: %v\n", v)
			name := envViper.GetString(v)
			fmt.Printf("[Debug]     %v\n", name)
			value, isSet := os.LookupEnv(name)
			if !isSet {
				fmt.Printf("[Debug]     Env %v not set, skip\n", name)
			} else {
				fmt.Printf("[Debug]     [%v: %v]\n", name, value)
				fmt.Printf("[Debug]     [%v: %v]\n", v, value)
				viper.Set(v, value)
			}
		}

		fmt.Println("[DEBUG] main config is: ", mainConfig)
		fmt.Println("[DEBUG] delta(s): ")
		for _, v := range patches {
			fmt.Println("  - ", v)
		}

		viper.SetConfigFile(mainConfig)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Error reading config file: ", err)
		}

		fmt.Println("[DEBUG] viper's content: ")
		print_all(viper.GetViper())

		// load patch into a sub viper
		for _, val := range patches {
			delta := viper.New()
			delta.SetConfigFile(val)
			if err := delta.ReadInConfig(); err != nil {
				fmt.Printf("Error reading %v : %v\n", val, err)
			}

			fmt.Printf("[DEBUG] %v's content: \n", val)
			print_all(delta)

			// merge the patch to the principle viper
			if err := viper.MergeConfigMap(delta.AllSettings()); err != nil {
				fmt.Printf("Error merge %v to viper: %v\n", val, err)
			}

			fmt.Printf("[DEBUG] viper's content after merge %v: \n", val)
			print_all(viper.GetViper())
		}
	},
}

func init() {
	ConfCmd.PersistentFlags().StringVarP(&mainConfig, "config", "c", "", "")
	ConfCmd.PersistentFlags().StringVarP(&envConfig, "env", "e", "", "")
	ConfCmd.PersistentFlags().StringSliceVarP(&patches, "override", "o", []string{}, "")
	if err := ConfCmd.MarkPersistentFlagRequired("config"); err != nil {
		fmt.Println(err)
	}
}

func print_all(v *viper.Viper) {
	for _, key := range v.AllKeys() {
		fmt.Printf("[%v]: %v\n", key, v.Get(key))
	}
	fmt.Println()
}
