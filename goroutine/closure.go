package main

import (
	"fmt"
)

func main() {
	done := make(chan bool)

	// strs := []int{1,2,3}
	strs := []string{"huo", "foo", "bar"}

	for _, s := range strs {
		/* (1) not work as expected, s shared by all the goroutines, s may be modified when the goroutine execute
		go func() {
			fmt.Println(s)
			done <- true
		}()
		*/

		/* (2) works, the value s is passed to func 
		go func(str string) {
			fmt.Println(str)
			done <- true
		}(s)
		*/

		/* (3) works, create a new s */
		s := s
		go func() {
			fmt.Println(s)
			done <- true
		}()
	}

	for _ = range strs {
		<-done
	}
}
