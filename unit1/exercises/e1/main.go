// a function that returns an int.
func fibonacci() func() int {
	a, b := 0, 0
	return func() int {
		switch a {
		case 0:
			a = 1
			b = 1
			return 0
		default:
			current := a
			a, b = b, a+b
			return a - current
		}
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 12; i++ {
		fmt.Println(f())
	}
}