package main

import (
	"fmt"
)

type kvpair struct {
	key string
	val string
}

type byValueMap map[string]kvpair
type byPointerMap map[string]*kvpair

func main() {
	// #1 by value
	test1 := byValueMap {
		"foo": { "foo/hello", "fooval" },
		"bar": { "bar/hello", "barval" },
	}

	value1 := kvpair{ "fufu/hello", "fufuval" }
	// copy the entier kvpair struct value1 into the map !
	test1["foo"] = value1
	// will not impact the map
	value1.key = "modified"

	fmt.Println(value1)
	for k, v := range test1 {
		fmt.Printf("%#v: %#v\n", k, v)
        }

	// #2 by pointer
	test2 := byPointerMap {
		"foo": &kvpair{ "foo/hello", "fooval" },
		"bar": &kvpair{ "bar/hello", "barval" },
	}

	value2 := &kvpair{ "fufu/hello", "fufuval" }
	test2["foo"] = value2
	// naturally this will impact the map
	value2.key = "modified"

	fmt.Println(*value2)
	for k, v := range test2 {
		fmt.Printf("%#v: %#v\n", k, *v)
        }
}
