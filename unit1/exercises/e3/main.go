package main

import (
	"fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprint("cannot Sqrt negative number: -2\n")
}

func Sqrt(x float64) (float64, error) {
	var e ErrNegativeSqrt
	epsilon := 0.0000000000001
	z := x / 2
	i := 1
	if x < 0 {
		return 0, e
	}
	for {
		delta := (z*z - x) / (2 * z)
		if delta < epsilon && delta > -epsilon {
			break
		}
		z -= delta
		//fmt.Println(i," : ", z)
		i += 1
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
