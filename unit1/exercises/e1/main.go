package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	start_num := 0
	next_num := 1
	return func() int {
		sum := start_num + next_num
		start_num = next_num
		next_num = sum
		return sum
	}
}

func fibonacci1() func() int {
	prev := 0
	current := 1
	return func() int {
		result := prev + current
		prev = current
		current = result

		return result
	}

}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
