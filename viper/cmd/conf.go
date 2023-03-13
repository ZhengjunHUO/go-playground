package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"

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
		// (1) read in principle config to viper
		fmt.Println("[DEBUG] main config is: ", mainConfig)
		viper.SetConfigFile(mainConfig)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Error reading config file: ", err)
		}

		// debug output
		fmt.Println("[DEBUG] viper's content: ")
		print_all(viper.GetViper())

		fmt.Println("[DEBUG] delta(s): ")
		for _, v := range patches {
			fmt.Println("  - ", v)
		}

		// (2) override values in viper with patches
		load_patches()

		// (3) read in config file, get env from os recursively
		if len(envConfig) != 0 {
			load_from_env_var(envConfig)
		}

		// debug output
		fmt.Println("[DEBUG] viper's final content: ")
		print_all(viper.GetViper())
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
		val := v.Get(key)
		if reflect.ValueOf(val).Kind() == reflect.String && strings.HasPrefix(val.(string), "$") {
			envVar := val.(string)[1:]
			fmt.Printf("Find an env var: %v\n", envVar)
			load_env_var_to_viper(key, envVar)
		}
		fmt.Printf("[%v]: %v(%v)\n", key, val, reflect.TypeOf(val))
	}
	fmt.Println()
}

func load_patches() {
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
}

func load_from_env_var(conf string) {
	envViper := viper.New()
	envViper.SetConfigFile(conf)
	if err := envViper.ReadInConfig(); err != nil {
		fmt.Println("Error reading env config file: ", err)
	}

	for _, v := range envViper.AllKeys() {
		if v == "include" {
			l := envViper.GetStringSlice(v)
			fmt.Printf("[Debug] include: %v\n", l)
			for _, f := range l {
				load_from_env_var(f)
			}
			continue
		}

		fmt.Printf("[Debug] env: %v\n", v)
		name := envViper.GetString(v)
		fmt.Printf("[Debug]     %v\n", name)
		load_env_var_to_viper(v, name)
	}
}

func load_env_var_to_viper(keyname, envname string) {
	value, isSet := os.LookupEnv(envname)
	if !isSet {
		fmt.Printf("[Debug]     Env %v not set, skip\n", envname)
	} else {
		fmt.Printf("[Debug]     Env var [%v: %v]\n", envname, value)
		fmt.Printf("[Debug]     Load [%v: %v]\n", keyname, value)
		viper.Set(keyname, value)
	}
}
