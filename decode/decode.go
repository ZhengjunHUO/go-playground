package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	// Read into viper
	conf := "./input.yaml"
	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}

	fmt.Println("[DEBUG] viper's content: ")
	print_all(viper.GetViper())

	//use_mapstruct()
	use_yaml()
}

func use_yaml() {
	// (1) viper map[string]interface{} => go struct
	var hero Hero
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "yaml", Result: &hero})
	if err != nil {
		fmt.Printf("Error init decoder: %v\n", err)
		os.Exit(1)
	}

	if err := dec.Decode(viper.GetStringMap("inventory.hero1")); err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[DEBUG] hero: \n%v\n", hero)

	// (2) go struct => viper map[string]interface{}
	buf := bytes.NewBuffer([]byte{})
	enc := yaml.NewEncoder(buf)
	dummy := "dummy"

	if err := enc.Encode(map[string]interface{}{dummy: map[string]interface{}{"hero1": hero}}); err != nil {
		fmt.Printf("Error encoding: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[DEBUG] buf: \n%v\n", buf.String()[len(dummy)+2:])

	// (3) load to new viper
	vip := viper.New()
	vip.SetConfigType("yaml")
	vip.ReadConfig(buf)
	print_all(vip)
	print_all_settings(vip)
}

func use_mapstruct() {
	// (1) viper map[string]interface{} => go struct
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

func print_all_settings(v *viper.Viper) {
        content, err := yaml.Marshal(v.AllSettings())
        if err != nil {
		fmt.Printf("Error marshalling: %v\n", err)
		os.Exit(1)
        }
        br := bufio.NewReader(bytes.NewBuffer(content))

	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error read line: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%v\n", string(line))
	}
}
