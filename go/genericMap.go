package main

import (
	"fmt"
)

func GetKeysFromMap[K comparable, V any](dict map[K]V) []K{
	rslt := make([]K, 0, len(dict))
	for k := range dict {
		rslt = append(rslt, k)
	}

	return rslt
}

func main() {
	dict1 := map[string]int{
		"fufu": 8,
		"foo": 3,
		"bar": 5,
	}
	dict2 := map[int]bool{
		1: false,
		2: false,
		7: true,
	}
	fmt.Println(GetKeysFromMap(dict1))
	fmt.Println(GetKeysFromMap(dict2))
}
