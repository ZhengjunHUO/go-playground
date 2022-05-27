// +build premium plus

package main

// here the +build tag logic is "OR"
func init() {
	functions = append(functions, 
		"Plus Fonction 1",
		"Plus Fonction 2",
		"Plus Fonction 3",
	)
}
