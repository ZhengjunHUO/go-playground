package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func concat() {
	// load principle config into global viper
	conf := "./config.yaml"
	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}

	fmt.Println("[DEBUG] viper's content: ")
	print_all(viper.GetViper())

	// load patch into a sub viper
	for _, val := range []string{"./delta1.yaml", "./delta2.yaml"} {
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

func print_all(v *viper.Viper) {
	for _, key := range v.AllKeys() {
		fmt.Printf("[%v]: %v\n", key, v.Get(key))
	}
	fmt.Println()
}
