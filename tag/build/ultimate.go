// +build premium,plus

package main

/* here the +build tag logic is AND
// +build premium
// +build plus
*/

func init() {
	functions = append(functions,
		"Ultimate Fonction 1",
		"Ultimate Fonction 2",
	)
}
