package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	//"os/exec"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
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
	Name   string `yaml:"full_name"`
	Age    int    `yaml:"age"`
	Emails []struct {
		Address string `yaml:"addr_name"`
		Extra   string `yaml:"extra_info"`
	} `yaml:"emails"`
	ConfigMapKeyValue []MapKeyValue `yaml:"config_map_key_value"`
}

type HeroMs struct {
	Name   string `mapstructure:"full_name"`
	Age    int    `mapstructure:"age"`
	Emails []struct {
		Address string `mapstructure:"addr_name"`
		Extra   string `mapstructure:"extra_info"`
	} `mapstructure:"emails"`
	ConfigMapKeyValue []MapKeyValueMs `mapstructure:"config_map_key_value"`
}

var dictKeys map[string]struct{}
var dictMaps map[string]struct{}

type Object struct {
	Name       string                 `yaml:"full_name,omitempty"`
	Id         int                    `yaml:"id,omitempty"`
	IsGlobal   bool                   `yaml:"is_global,omitempty"`
	Params     []string               `yaml:"params,omitempty"`
	Labels     map[string]interface{} `yaml:"labels,omitempty"`
	RefCounter *int32                 `yaml:"ref_counter,omitempty"`
	SubObject  SubObject              `yaml:"sub_object,omitempty"`
}

type SubObject struct {
	Name       string                 `yaml:"full_name,omitempty"`
	Id         int                    `yaml:"id,omitempty"`
	IsGlobal   bool                   `yaml:"is_global,omitempty"`
	Params     []string               `yaml:"params,omitempty"`
	Labels     map[string]interface{} `yaml:"labels,omitempty"`
	RefCounter *int32                 `yaml:"ref_counter,omitempty"`
}

var decodeOpt viper.DecoderConfigOption

// infer_viper_keys generates the viper keys based on the yaml tag's name in a recursive way
// if the field is a struct, it dives into the field and generates the key
// by appending the filed's name as a prefix
func infer_viper_keys(v reflect.Value, prefix string) {
	t := v.Type()
	fields := reflect.VisibleFields(t)

	for _, f := range fields {
		/*
			fmt.Println(f.Name)
			fmt.Println(f.Type)
			fmt.Println(f.Type.Kind())
			fmt.Println(strings.Split(f.Tag.Get("yaml"), ",")[0])
			fmt.Println()
		*/
		keyname := prefix + strings.Split(f.Tag.Get("yaml"), ",")[0]
		switch f.Type.Kind() {
		case reflect.Map:
			dictMaps[keyname] = struct{}{}
		case reflect.Struct:
			infer_viper_keys(v.FieldByName(f.Name), keyname+".")
		default:
			dictKeys[keyname] = struct{}{}
		}
	}
}

func main() {
	// Read into viper
	conf := "./input.yaml"
	viper.SetConfigFile(conf)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}

	//fmt.Println("[DEBUG] viper's content: ")
	//print_all(viper.GetViper())

	decodeOpt = viper.DecoderConfigOption(func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	})

	//use_mapstruct()
	use_yaml()

	/*
	dictKeys = map[string]struct{}{}
	dictMaps = map[string]struct{}{}

	v := reflect.ValueOf(Object{})
	infer_viper_keys(v, "")

	fmt.Printf("dictKeys: %#v\n", dictKeys)
	fmt.Printf("dictMaps: %#v\n", dictMaps)
	*/

	/*
	transformer := `jq -rR 'gsub("-";"+") | gsub("_";"/") | split(".") | .[1] | @base64d | fromjson | .email'`
	cmd := exec.Command("bash", "-c", transformer)
	cmd.Stdin = strings.NewReader("eyJh...")

	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
	    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	    return
	}
	fmt.Println(out.String())
	*/
}

func use_yaml() {
	// (1) viper map[string]interface{} => go struct
	var hero Hero
	/*
		dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "yaml", Result: &hero})
		if err != nil {
			fmt.Printf("Error init decoder: %v\n", err)
			os.Exit(1)
		}

		if err := dec.Decode(viper.GetStringMap("inventory.hero1")); err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			os.Exit(1)
		}
	*/
	err := viper.UnmarshalKey("inventory.hero1", &hero, decodeOpt)
	if err != nil {
		fmt.Printf("Error unmarshal key: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("[DEBUG] hero: \n%v\n", hero)

	// (2) go struct => viper map[string]interface{}
	buf, rslt := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
	enc := yaml.NewEncoder(buf)
	dummy := "dummy"

	if err := enc.Encode(map[string]interface{}{dummy: map[string]interface{}{"hero1": hero}}); err != nil {
		fmt.Printf("Error encoding: %v\n", err)
		os.Exit(1)
	}
	//fmt.Printf("[DEBUG] buf: \n%v\n", buf.String()[len(dummy)+2:])
	rslt.WriteString("inventory:\n")
	rslt.WriteString(buf.String()[len(dummy)+2:]+"\n")

	// (3) load raw data to a new viper, render it
	// and load the final result to another viper
	vip, vip_raw := viper.New(), viper.New()
	vip.SetConfigType("yaml")
	vip_raw.SetConfigType("yaml")

	// prepare the raw viper
	data := bytes.NewReader(rslt.Bytes())
	vip_raw.ReadConfig(data)

	fmt.Println("[DEBUG] vip_raw:")
	print_all_settings(vip_raw)

	// do the rendering and load the final params to another viper instance
	vip.ReadConfig(selfRender(rslt.Bytes(), vip_raw))

	//print_all(vip)
	fmt.Println("[DEBUG] vip:")
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

// print_all prints all key-value pairs in a viper instance
func print_all(v *viper.Viper) {
	for _, key := range v.AllKeys() {
		fmt.Printf("[%v]: %v\n", key, v.Get(key))
	}
	fmt.Println()
}

// print_all_settings prints the all viper's values with yaml format
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

func selfRender(data []byte, v *viper.Viper) *bytes.Buffer {
        var out bytes.Buffer
        tmpl := template.Must(template.New("selfTemplate").Parse(string(data)))
        tmpl.Delims("{{", "}}")
        err := tmpl.Execute(&out, v.AllSettings())
        if err != nil {
                panic(err)
        }
        return &out
}
