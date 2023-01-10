package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var ErrDivideByZero = errors.New("divide by zero")
var ErrIntOverflow = errors.New("integer overflow")
var ErrNotValidInteger = errors.New("argument is not valid integer")

type ExpressionError struct {
	A   int   `json:"var A"`
	B   int   `json:"var B"`
	Err error `json:"-"`
}

func (err *ExpressionError) Error() string {
	return fmt.Sprintf("error in expression %d * %d: %s", err.A, err.B, err.Err)
}

func NewExpressionError(a int, b int, err error) error {
	return &ExpressionError{
		A:   a,
		B:   b,
		Err: err,
	}
}

func mul(a, b int) (int, error) {
	r := a * b
	if r/a != b {
		// here we wrapped error ErrIntOverflow with additional message "error in expression..." and also provided expression.
		return 0, NewExpressionError(a, b, ErrIntOverflow) // fmt.Errorf("error in expression %d * %d: %w", a, b, ErrIntOverflow)
	}
	return r, nil
}

func div(a, b int) (int, error) {
	if b == 0 {
		// here we wrapped error ErrDivideByZero with additional message "error in expression..." and also provided expression.
		return 0, ErrDivideByZero //fmt.Errorf("error in expression %d / %d: %w", a, b, ErrDivideByZero)
	}
	return a / b, nil
}

func pow(a, b int) (int, error) {

	var list []int

	result := a
	for i := 2; i <= b; i++ {
		result *= a
		list = append(list, result)
		if len(list) > 2 {
			if list[len(list)-1] < list[len(list)-2] {
				return 0, NewExpressionError(a, b, ErrIntOverflow)
			}
		}
	}

	return result, nil
}

type Operation struct {
	F func(a, b int) (int, error)
	A int
	B int
}

func parseInstructionsStdin() ([]Operation, error) {
	if os.Stdin == nil {
		return nil, fmt.Errorf("stdin is not provided")
	}

	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	ops := []Operation{}

	for fileScanner.Scan() {
		line := fileScanner.Text()
		words := strings.Fields(line)

		a, err := strconv.Atoi(words[1])
		if err != nil {
			return nil, fmt.Errorf("wrong instruction format \"%s\", argument A: %w", line, ErrNotValidInteger)
		}

		b, err := strconv.Atoi(words[2])
		if err != nil {
			return nil, fmt.Errorf("wrong instruction format \"%s\", argument B: %w", line, ErrNotValidInteger)
		}

		var op func(a, b int) (int, error)
		switch words[0] {
		case "mul":
			op = mul
		case "div":
			op = div
		case "pow":
			op = pow
		default:
			return nil, fmt.Errorf("wrong instruction %s, on line %s", words[0], line)
		}

		ops = append(ops, Operation{
			F: op,
			A: a,
			B: b,
		})
	}
	return ops, nil
}

func main() {
	ops, err := parseInstructionsStdin()
	if err != nil {
		log.Fatalln(err)
	}

	for _, op := range ops {
		result, err := op.F(op.A, op.B)
		if err != nil {
			if errors.Is(err, ErrDivideByZero) {
				fmt.Println("eternity")
				continue
			} else {
				log.Fatalln("Computation error:", err)
			}
		}

		fmt.Println(result)
	}
}
