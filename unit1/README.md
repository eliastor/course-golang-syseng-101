# Unit 1: A tour of Go
Please finish all materials on https://go.dev/tour/ or https://go.dev/tour/list except "Generics" and "Concurrency"

Note: if you stuck with module try to read https://go.dev/blog/using-go-modules.

## FAQ

### Defer, why and how?

Defer is used to delay some function call when return from function is happened. In most cases defer is used to perform cleanup and closing resources.

For an instance when we are working with files we need to close it. For example we have function that opens file and call other funciton(s) performs processing of opened file which can return errors.

Every time we are facing with return function file must be closed by calling `file.Close()`:

```go
func routine(filename string) error
file, err := os.Open(filename)
if err != nil {
	return err
}

err = Foo(file)
if err != nil {
	file.Close()
	return err
}

file.Close()
return nil
```

Such approach is not convinient and causes places for potential errors.
instead of that we can use `defer` statement:

```go
func routine(filename string) error
file, err := os.Open(filename)
if err != nil {
	return err
}

defer file.Close() //starting from this point no matter where routine returns file.Close() will be called.

err = Foo(file)
if err != nil {
	return err
}
return nil
```

Note that defer is a statement that delays function call and arguments to this call. The actual call and argument resolving will happen "during" `return` statement


### Slices and arrays. Mechanincs and difference

Both arrays and slices represent solid piece of memory.
Arrays is static structure points to fixed amount of memory, while slice is dynamic and can point to flexible amount of memory.

Arrays and slices have 2 major properties:
 - length `len()` - number of elements defined in array or slice
 - capacity `cap()` - total length of memory allocated for slice or array

Array properties are known at compilation time, while slice properties are known only at runtime.

#### Array:

```go 
nums := [5]int{1, 2, 3} // nums == [1 2 3 0 0], len()=5, cap()=5
nums[0]=0 				// nums == [0 2 3 0 0], len()=5, cap()=5
nums[3]=4 				// nums == [0 2 3 4 0], len()=5, cap()=5
```

Note: `[...]int{1, 2, 3}` will define array of len()=3 cap()=3. Three dots direct compiler to count number of elements in definition and set cap and len to this value.

#### Slice:

As we saw for array its cap and len are always the same. Slice is slightly trickier.

```go
nums := make([]int, 0,2)	// nums[], mem[0 0], cap=2, len=0
nums[0]=1					// panic, len=0, there is no first element
nums = append(nums, 1)  	// nums[1], mem[1 0], cap=2, len=1
append(nums[1:],2)			// nums[1], mem[1 2], cap=2, len=1
nums = nums[:cap(nums)]		// nums[1 2], mem[1 2], cap=2, len=2
```

```go
nums := []int{1, 2, 3} 	// nums == [1 2 3 0 0], len()=3, cap()=3
nums[0]=0 				// nums == [0 2 3 0 0], len()=3, cap()=3
nums[3]=4 				// nums == [0 2 3 4 0], len()=3, cap()=3
nums2:=make([]int,0,1) 	// nums2 == [], 		len()=0, cap()=1
nums2 = append(nums2, 9)// nums2 == [9], 		len()=1, cap()=1
nums2 = append(nums2, 8)// nums2 == [9 8], 		len()=2, cap()=2
```

As one might see we can append new elements to the end of the slice and capacity of slice will be increased if needed. For arrays we can't append elements. Slices and arrays are different types as `float` and `int` for example:

```go
append(nums, 5)			// compile error: nums is not slice
nums=[]int{0}			// compile error: nums is array and []int{0} is slice
```

We can easily make slice from array and slice:

```go
nums := [5]int{1, 2, 3}	// [1 2 3 0 0], len()=5, cap()=5
slice := nums[:] 		// [1 2 3 0 0], len()=5, cap()=5
```

Note that when we make slice from another slice or array, created slice points to the same memory as original slice or array:

```go
nums := [3]int{1, 2}
slice := nums[0:1]
// nums: [1 2 0] len()=3, cap()=3, slice [1] len()=1, cap()=3

slice[0]=10 // We changed element in slice but element 0 is also changed in nums:
// nums: [10 2 0] len()=3, cap()=3, slice [10] len()=1, cap()=3
// the reason is that element 0 of slice points to the same memory as element 0 of nums

// We can even append element to slice and it will appear in nums:
slice = append(slice, 11, 12)
// nums: [10 11 12] len()=3, cap()=3, slice [10,11,12] len()=3, cap()=3
```

What if one will continue to add elements to slice built on top of array?

```go
slice = append(slice, 13)
// nums: [10 11 12] len()=3, cap()=3, slice [10,11,12,13] len()=4, cap()=4
slice[0]=9000
// nums: [10 11 12] len()=3, cap()=3, slice [9000,11,12,13] len()=4, cap()=4
```

As one can see, when we exceed capacity of slice array, append will alocate new slice and copy all values from old slice and append values to the end of **new** slice. Values are remain but slice doesn't point to the same memory as array anymore. Changing element in slice will not affect orifinal array.

The same situation will happen for slice of slice. Just keep in mind slice internals with cap, len and underlying memory.

Read more about slices in [go blog article](https://go.dev/blog/slices-intro)

---

#### 

## Quiz

#### Q1. How will you add the number 3 to the right side?

`numbers := []int{1, 1, 2}`

1. `numbers.append(3)`
2. `numbers.insert(3, 3)`
3. `append(numbers, 3)`
4. `numbers = append(numbers, 3)` - this doing it

#### Q2. From where is the variable fooVar accessible if it is declared outside of any functions in a file in package fooPackage located inside module fooModule

1. anywhere inside `fooPackage`, not the rest of `fooModule` - here
2. by any application that imports `fooModule`
3. from anywhere in `fooModule`
4. by other packages in `fooModule` as long as they import `fooPackage`

#### Q3. What should the idiomatic name be for an interface with a single method and the signature Serve() error

1. Servable
2. Server
3. ServeInterface
4. IServe

#### Q4. Which is **not** valid loop construct?

1. `for i,r:=0,rand.Int(); i < r%10; i++ { ... }`
2. `for { ... }`
3. `{ ... } for false` -- here is a failure
4. `for _,c := range "foobar" { ... }`

## Excercises
`E0` is for illustration how to work and submit Excercises.
Let's make 

### E0. Exercise: Loops and Functions
More: https://go.dev/tour/flowcontrol/8
let's create implementation in `unit1/exercises/e0/main.go`, wherer `unit1` is number of current unit, `e0` is number of current exercise.

If task has multiple steps to do, we assume that the last steps is final. Input values must be the same as in Exercie definition if other is not mentioned. 

At first we've been asked to implement Sqrt function with partial implementation of Newton method and send output of computation steps on each iteration of method:

```go
package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	z := 1.0
	for i:=0; i<10; i++ {
		delta := (z*z - x) / (2*z)
		z -= delta
		fmt.Println(z)
	}
	return z
}

func main() {
	fmt.Println(Sqrt(2))
}
```

Next we've asked to change loop condition to stop once the value has stopped changing (or only changes by a very small amount). We can achieve that by adding `epsilon` constant that will be compared with `delta`. 
**Note that their absolute values should be compared**. One approach is to use `Abs()` function form `math` package, or we can notice that `x` argument must be positive, let's make our own condition based on partial comparison of positive-only numbers:

```go
package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	epsilon := 0.0000000000001
	z := 1.0
	for {
		delta := (z*z - x) / (2 * z)
        // or: if math.Abs(delta) < epsilon {
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
```

If you haven't managed to solve all steps publish code for the step you succeeded and make a comment in code.
If you exercise requires to make multiple files or even packages, don't hesitate to create them in `unit1/exercises/eX/` folder

As soon as you implemented the code and place it into `unit1/exercises/e0/main.go`, make PR to this repo.

---

### E1. Fibonacci closure

More: https://go.dev/tour/moretypes/26

**Note**: Edit only `fibonacci` function definition:

```go
func fibonacci() func() int {}
```
Keep `main` function untouch. Don't add additional Prints to output. It is checked in tests.

Share your implementation `unit1/exercises/e1/main.go` in github PR.

---

### E2. Stringers
More: https://go.dev/tour/methods/18

Implement `func (ip IPAddr) String() string` method for `IPAddr` structure. 
Note that receiver variable `ip` is not pointer.

Keep `main` function untouch. Don't add additional Prints to output. It is checked in tests.

Share your implementation  `unit1/exercises/e2/main.go` in github PR.

Note: about interfaces, behaviour and who is calling String()

---

### E3. Errors
More: https://go.dev/tour/methods/20

**Note**: errors are just values and must be treated as values. In this scenario error value will be `ErrNegativeSqrt` struct that satisfy error type. Nil is valid error, pointer to struct with `Error() string` method is valid error too.

Keep `main` function untouch. Don't add additional Prints to output. It is checked in tests.

Share your implementation  `unit1/exercises/e3/main.go` in github PR.

---

### E4. rot13Reader
More: https://go.dev/tour/methods/23

Add more info about reader and the entire workflow of reaeder.

`'A'` is a byte. It is treated as 8 bit number, so it can be compared and added/substracted. 

not that A-Z and a-z letters have their own ranges and should be mapped to letter from their group: letter from A-Z must be mapped to letter from A-Z. Same for a-z.

https://www.rapidtables.com/code/text/ascii-table.html


Keep `main` function untouch. Don't add additional Prints to output. It is checked in tests.

Share your implementation  `unit1/exercises/e4/main.go` in github PR.

---

### E5. Images
More: https://go.dev/tour/methods/25

**Note**: read the docs for [golang.org/x/tour/pic](https://pkg.go.dev/golang.org/x/tour/pic) package.

Keep `pic.ShowImage(m)` in `main` function untouch. Don't add additional Prints to output. It is checked in tests. 

Share your implementation  `unit1/exercises/e5/main.go` in github PR.

Note about what to do 

---

## Summary of what you learnt

Don't worry if you haven't deep understanding why interfaces are needed. You got almost all Go syntax. Key point after this unit is to be able to read the code and tear down it to known primiteves.