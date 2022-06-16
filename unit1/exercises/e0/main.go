package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	epsilon := 0.0000000000001
	z := x / 2
	for {
		delta := (z*z - x) / (2 * z)
		if delta < epsilon && delta > -epsilon {
			break
		}
		z -= delta
		fmt.Println(z)
	}
	return z
}

func main() {
	fmt.Println(Sqrt(2))
}
