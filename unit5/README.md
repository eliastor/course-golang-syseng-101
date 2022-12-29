# Unit 5: Errors

For this unit your environment must be initialized as in previous units.

## Materials

Excessive material can be found on [GoByExample](https://gobyexample.com/errors) and [official go site](https://go.dev/blog/error-handling-and-go) amd [standard library package](https://pkg.go.dev/errors)

## Errors are just values

The error type in Go is implemented as the following interface:
```go
type error interface {
    Error() string
}
```

In most programming languages one can find exception and try/catch blocks/methods to handle errors. In Go there is no such mechanism. Instead of that all the errors are treated as additional values that must be returned from methods as last return values and satisfy error interface.

Main idea is that you have two options how to deal with errors: escalate it to calling function and let them to deal with the error, or deal by your own and make some logic on this error.

Creating of errors is easy, you can create generic error with [errors.New()](https://pkg.go.dev/errors#New) from `errors` package: 

```go
func div(a, b int) (int, error) {
    if b == 0 {
        return 0,  errors.New("divide by zero")
    }
    return a/b, nil
}
```

Errors can be returned as nil, and in fact, it’s the default, or zero value of on error in Go.

This is important since checking if err != nil is the idiomatic way to determine if an error was encountered (replacing the try/catch statements you may be familiar with in other programming languages).

Errors are returned as the last argument in a function. Hence in example above, we return an int and an error, in that order.

When we do return an error, the other arguments returned by the function are typically returned as their default “zero” value. A user of a function may expect that if a non-nil error is returned, then the other arguments returned are not relevant.

Error messages are usually written in lower-case and don’t end in punctuation. Exceptions can be made though, for example when including a proper noun, a function name that begins with a capital letter, etc.

For convenience you can use [fmt.Errorf()](https://pkg.go.dev/fmt#Errorf) to construct more verbose error message:

```go
func copyFile(src, dst string) (error) {
    if src == "/dev/null" {
        return fmt.Errorf("can't copy from /dev/null to %s", dst)
    }
    // ...
}
```

### Defining common errors

There is common technique to create common errors that are used to check the meaning of the error.

In example above with div function let's imagine that we are making library that works with numbers and division to zero can be found in multiple places in the code or we need to figure out is error represents "division to zero" error or not.

For such cases common errors are created in package scope:

```go
package main

import (
    "errors"
    "fmt"    
)

var ErrDivideByZero =  errors.New("divide by zero")

func div(a, b int) (int, error) {
    if b == 0 {
        // now we can utilise the error we defined above in package scope
        return 0, ErrDivideByZero
    }
    return a/b, nil
}
```

Errors package offers [errors.Is()](https://pkg.go.dev/errors#Is) function to check if errors are equal:

```go
package main

import (
    "errors"
    "fmt"    
)

// ErrDivideByZero is used by inside functions (div in this example) to signal about computational error with division by zero
var ErrDivideByZero =  errors.New("divide by zero")

func div(a, b int) (int, error) {
    if b == 0 {
        // now we can utilise the error we defined above in package scope
        return 0, ErrDivideByZero //this called escalation to calling function
    }
    return a/b, nil
}

func main() {
    result, err := div(1, 0)
    if err != nil {
        if errors.Is(err, ErrDivideByZero) {
                // 
                fmt.Println("eternity")
            } else {
                // for all other errors we can't work with, let's escalate them. As far as it's main function, we can only escalate it to end-user via output:
                fmt.Println("unexpected div error:", err)
            }
        return
    }
    fmt.Println(result)
}
```

In example above we **escalated** error from `div` function to calling function `main` and in main we use go-idiomatic way to check error and work with error:

```go
func inner() (string, error)
_, err := somefunc()
if err !=nil {
    // work on error or escalate it
}
```

### Customizing errors

Many use cases can be handled with technique above, however there can be cases when you might need more functionality: additional data fields, dynamic values in error message and so on.

For achieving that we can create our own custom error type. Also by using [`errors.As`](https://pkg.go.dev/errors#As) it is possible to check and convert (cast) generic error to our custom error:

```go
package main

import (
	"errors"
	"fmt"
)

type ComputationError struct {
	Args []interface{}
	Msg  string
}

func (e *ComputationError) Error() string {
	return fmt.Sprintf("computational error (args: %v): %s", e.Args, e.Msg)
}

func div(a, b int) (int, error) {
	if b == 0 {
		return 0, &ComputationError{
			Args: []interface{}{a, b},
			Msg:  fmt.Sprintf("cannot divide %d by zero", a),
		}
	}
	return a / b, nil
}

func main() {
	result, err := div(1, 0)
	if err != nil {
		var compErr *ComputationError
		if errors.As(err, &compErr) {
			fmt.Printf("Computation error: %v\n", compErr)
		} else {
			fmt.Println("unexpected div error:", err)
		}
		return
	}
	fmt.Println(result)
}
```

> Note: also you can take a look on [go blog](https://go.dev/blog/go1.13-errors) to get more information about errors customizations for changing behavior for Is and As functions. 

### Wrapping errors

In real world scenarios you can find that errors can happen at different stages and abstraction levels of program and the actual error can be generated in one function call (in example above it is `div()`) and handled on another (in example above it is `main()`).

When we have much more function calls between error production and handling it's better to get more context what was exactly done to get clear sight of program flow while error was been produced, for achieving that we can use wrapping. 

Let's slightly modify example our second example:

```go
package main

import (
	"errors"
	"fmt"
)

var ErrDivideByZero =  errors.New("divide by zero")
var ErrIntOverflow =  errors.New("integer overflow")

func mul(a, b int) (int, error) {
    r := a * b
    if r / a != b {
         // here we wrapped error ErrIntOverflow with additional message "error in expression..." and also provided expression.
        return 0, fmt.Errorf("error in expression %d * %d: %w", a, b, ErrIntOverflow)
    }
    return r, nil
}

func div(a, b int) (int, error) {
	if b == 0 {
        // here we wrapped error ErrDivideByZero with additional message "error in expression..." and also provided expression.
		return 0, fmt.Errorf("error in expression %d / %d: %w", a, b, ErrDivideByZero)
	}
	return a / b, nil
}

type Operation struct {
    F func(a, b int) (int, error)
    A int
    B int
}

func main() {
    ops := []Operation{
        {mul, 5, 0},
        {div, 5, 0},
    }

    for _, op := range ops {
        result, err := op.F(op.A, op.B)
	    if err != nil {
		   log.Fatalln("Computation error:", err)
	    }
        fmt.Println(result)
    }
}
```

Now we have all the context to get clues where error is happened and can determine the expression causing the error.

Moreover, as soon as we wrapped `ErrDivideByZero` error with all the messages we still can compare generic error message, we received during execution of operation function F, with `ErrDivideByZero` error using `errors.Is` function: for example, we can use `errors.Is` to detect division by zero and print "eternity" for division by zero operation instead of the error: 

```go
...
    result, err := op.F(op.A, op.B)
	if err != nil {
        if errors.Is(err, ErrDivideByZero) {
            fmt.Println("eternity")
        } else {
            fmt.Println("unexpected div error:", err)
        }
	    log.Fatalln("Computation error:", err)
	return
    }
...
```

Wrapping errors allow your to store errors you wrapped using custom error types with `Wrap()` method defined or for errors created using `fmt.Errorf()`: `%w` in pattern string will inform Errorf to wrap provided error with surrounding message.

Generally, it’s a good idea to wrap an error with at least the function’s name, every time you “escalate” it up - i.e. every time you receive the error from a function and want to continue returning it back up the function chain.

---

### E0. Math script processor

To build exercise, being in root folder of the repo you can run:

```bash
go build ./unit5/exercises/e0
```

This command will build local folder with all `.go` files in it and place result application to `e0` file in current (repo root folder).

If you want to specify name of path of the file:

```bash
go build -o ./unit5/e0 ./unit5/exercises/e0
cat unit5/exercises/e0/instructions.txt | ./unit5/e0
```

for exercise of these and next unit it is handy to build and run in one command:

```bash
cat unit5/exercises/e0/instructions.txt | go run ./unit5/exercises/e0/
```

Note that the program waits for stdin with instruction like. Sample of instructions format is provided in `unit5/exercises/e0/instructions.txt`

Find [source code](exercises/e0/main.go) of this exercise.

---

## FAQ

TBA

---

## Quiz

#### Q1. TBA

## Excercises

### E1. pow operation for Math script processor

Extend code from exercise 0: Add more function `pow` that will count first argument in a power of the second one, for example:

```
"pow 2 3" = 8
"pow 2 4" = 16
"pow 3 2" = 9
"pow 3 3" = 27
```

Wrap error with integer overflow with additional data as it was done with `mul()`. error message is up to you, only existence of arguments, results and error message are tested.

**Note**: Test verifies the output of your program by running it and generating input file.

**Note**: Don't forget to check integer overflow and return appropriate error message. Such case is verified in tests. 

Don't add additional Prints to output. It is checked in tests.

Don't change definitions of `ErrDivideByZero` and `ErrIntOverflow` errors, they are checked in tests.

Share your implementation `unit5/exercises/e1/main.go` in github PR.
Don't hesitate to copy contents of `unit5/exercises/e0/` to `unit5/exercises/e1/` and modify necessary files or add new ones.

**Hint**

Remember that `a in power of b` it's just multiplying `a` to itself `b` times. Any number in zero order is 1.

**Hint**

To check overflow you can make inverse operation and divide calculated result by `a` and you should get `1`

### E2. Enhance errors

In previous exercise you can note that there are multiple statement that are almost identical:

- `fmt.Errorf("error in expression %d * %d: %w", a, b, ErrIntOverflow)`

- `fmt.Errorf("error in expression %d / %d: %w", a, b, ErrDivideByZero)`

- same for pow

We should eliminate duplicates of logic and constants in our applications, so let's create custom Error type that must:

- be defined as ExpressionError
- must wrap underlying error (in example above they are `ErrIntOverflow` or `ErrDivideByZero`)
- must print error message in same format as examples above

Don't add additional Prints to output. It is checked in tests.

Don't change definitions of `ErrDivideByZero` and `ErrIntOverflow` errors, they are checked in tests.

Share your implementation `unit5/exercises/e2/main.go` in github PR.
Don't hesitate to copy contents of `unit5/exercises/e1/` to `unit5/exercises/e2/` and modify necessary files or add new ones.

**Hint**

You can add `sign` field to the struct of ExpressionError and fill it with sign of the operation and further in String method you can return it as part of error message as to get the same message as done manually with fmt.Errorf().

**Hint**

To not filling ExpressionError struct manually feel free to make function NewExpressionError(...) that will help to fill structure from arguments and return complete ExpressionError

### E3. errors.Is

In previous exercise we created `ExpressionError` struct. When we got error in main, we have no clue is what kind of error is wrapped by `ExpressionError`.

Let's change behavior of the program: if the error wrapped by `ExpressionError` is `ErrDivideByZero`, then just return `eternity` in stdout as normal result of operation and continue the work.

Don't add additional Prints to output. It is checked in tests.

Don't change definitions of `ErrDivideByZero` and `ErrIntOverflow` errors, they are checked in tests.

Share your implementation `unit5/exercises/e3/main.go` in github PR.
Don't hesitate to copy contents of `unit5/exercises/e2/` to `unit5/exercises/e3/` and modify necessary files or add new ones.

---
