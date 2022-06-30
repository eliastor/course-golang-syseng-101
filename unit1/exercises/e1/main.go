package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	start_num := 0
	next_num := 1
	i := 0

	return func() int {
		if i == 0 {
			i++
			return 0
		} else {
			sum := start_num + next_num
			start_num = next_num
			next_num = sum
			return start_num
		}

	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
