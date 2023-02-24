package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	isPositive := math.Signbit(x)
	if isPositive {
		return 0, ErrNegativeSqrt(x)
	} else {
		z := 1.0
		for i := 0; i < 10; i++ {
			delta := (z*z - x) / (2 * z)
			z -= delta
			//fmt.Println(z)
		}
		return z, nil
	}
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
