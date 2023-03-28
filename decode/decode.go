package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type MapKeyValue struct {
        ValueName string `yaml:"value_name"`
        KeyName   string `yaml:"key_name"`
}

type MapKeyValueMs struct {
        ValueName string `mapstructure:"value_name"`
        KeyName   string `mapstructure:"key_name"`
}

type Hero struct {
	Name string `yaml:"full_name"`
	Age int `yaml:"age"`
	Emails []struct {
		Address	string `yaml:"addr_name"`
		Extra	string `yaml:"extra_info"`
	} `yaml:"emails"`
	ConfigMapKeyValue []MapKeyValue `yaml:"config_map_key_value"`
}

type HeroMs struct {
	Name string `mapstructure:"full_name"`
	Age int `mapstructure:"age"`
	Emails []struct {
		Address	string `mapstructure:"addr_name"`
		Extra	string `mapstructure:"extra_info"`
	} `mapstructure:"emails"`
	ConfigMapKeyValue []MapKeyValueMs `mapstructure:"config_map_key_value"`
}

func main() {
	conf := "./input.yaml"
	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}

	fmt.Println("[DEBUG] viper's content: ")
	print_all(viper.GetViper())

	// (1) viper map[string]interface{} => go struct
	//var hero HeroYaml
	var hero HeroMs
	if err := mapstructure.Decode(viper.GetStringMap("inventory.hero1"), &hero); err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Result: %v\n", hero)

	// (2) go struct => viper map[string]interface{}
	var vp map[string]interface{}
	if err := mapstructure.Decode(hero, &vp); err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		os.Exit(1)
	}

	content, err := yaml.Marshal(vp)
	if err != nil {
		fmt.Printf("Error marshalling: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result: %v\n", string(content))
}

func print_all(v *viper.Viper) {
	for _, key := range v.AllKeys() {
		fmt.Printf("[%v]: %v\n", key, v.Get(key))
	}
	fmt.Println()
}