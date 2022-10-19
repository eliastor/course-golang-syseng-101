package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	x := 0
	y := x
	z := x
	return func() int {
		if z == 0 {
			z = 1
			return 0
		} else if y == 0 {
			z = z + y
			y = 1
		} else {
			z = x + y
			x = y
			y = z
		}
		return z
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 15; i++ {
		fmt.Println(f())
	}
}
