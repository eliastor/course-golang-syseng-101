package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	i := 0
	fib1 := 0
	fib2 := 1
	fib_sum := 0
	return func() int {
		if i == 0 {
			i++
			return 0
		} else {
			fib_sum = fib1 + fib2
			fib1 = fib2
			fib2 = fib_sum
			i++
			return fib1
		}

	}

}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
