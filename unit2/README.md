# Unit 2: Concurrency

Please finish all materials on https://go.dev/tour/concurrency

Concurency was kept in mind during creation of go.
There 2 major primitives that can help in most of your cases: channels and goroutines.

Channels can be used as FIFO queues (buffered) or for syncronization (unbuffered).

Gorountines can be used not only for paralelisation of some computations but for lifecycle management of data, channels and streams.

For this exercises you will need to initialise your environment.
Go to repo root and initilise module:

```sh
go mod init course
```

For the next excercises you'll need some packages. If you'll try to run exercise 0 it will fail with error:
    
    unit2/exercises/e0/main.go:10:2: no required module provides package github.com/brianvoe/gofakeit; to add it:
        go get github.com/brianvoe/gofakeit

Let's add unresolved dependency package by running `go get github.com/brianvoe/gofakeit`

More information about modules: https://go.dev/blog/using-go-modules.
Also https://go.dev/doc/code might be helpful.

Channels explained:

- https://yourbasic.org/golang/channels-explained/
- https://yourbasic.org/golang/select-explained/
- https://yourbasic.org/golang/detect-deadlock/
- https://yourbasic.org/golang/broadcast-channel/
- https://yourbasic.org/golang/stop-goroutine/
- https://yourbasic.org/golang/detect-data-races/

Workerks pattern explained:
- https://gobyexample.com/worker-pools

Wait groups:
- https://pkg.go.dev/sync#WaitGroup
- https://gobyexample.com/waitgroups

---

## E0. Bureaucracy office (FanOut, FanIn)

Let's start with sample exercise.

There is Bureaucracy office that generates meaningless documents.
For some reasons the office have to switch to electronic document workflow.

Bureaucrat generate documents.
Each document must be signed with Signature of the office by Executor.
Then the signature must be verified by Publisher and sent to nowhere.
in addition we will count how many documents we created, signed, verified and sent to nowhere.


Find [source code](exercises/e0/main.go) of this exercise.

Here some information about useful patterns that have been applied in this exercise among with several more syncronisation primitives.

### Fan-out pattern

The main idea of Fan Out Pattern is to have:
- a channel that provides a signaling semantics ( close(chan) )
- channel can be buffered, so we canget FIFO buffering of messages
- a goroutine that start multiple workers subscribed to one channel and processed messages from this channel in parallell.
- multiple workers use signaling channel to signal that the processing is done


In this [excersice](exercises/e0/main.go) `docsNew` channel have 2 `executor`s goroutines spawned by `SpawnExecutors` function. Every `executor` listens same channel and as soon as message appears in channel it will be taken from the channel by first free executor.

Original Fan-out pattern assumes that each worker has its own output channel or no outputs at all. In this exercise we connect workers output channel to one channel and receive Fan-In pattern

### Fan-in pattern

The main idea of Fan In Pattern is to have:
- a channel that provides a signaling semantics ( close(chan) )
- channel can be buffered, so we canget FIFO buffering of messages
- collect messages from multiple goroutines to one channel for further processing

In this [excersice](exercises/e0/main.go) `executor`s send processed messages further to [docsSigned channel](exercises/e0/main.go#L126), so in the channel there will be all messages from all executors.

### Cancel signal propagation

Channels in go is used not only for sharind data by communication, but for locking and syncronisations. 

One of the popular mechanics is closing of channel, which triggers receive event (https://go.dev/tour/concurrency/4).

Note: when you call close(chan) in fact it will send empty message to channel with special flag that shows us that channel is closed.

For ex, we have string channel `ch`. Receiver code is like this:

```go
str, ok := <- ch
```

Once channel is closed, all messages from channel will received by receiver. The last message will be `str == ""` and `ok == false`. This special message can be read many times. This `ok` signal can help us to understand that channel has been closed.

In our exercise we want whole pipeline to be shutted down gracefully by closing channel `done`. It will signal `bureaucrat` to [return from function](exercises/e0/main.go#L27-L28). Bureaucrat was called in [goroutine](exercises/e0/main.go#L73-76) from [SpawnBureaucrat](exercises/e0/main.go#L73-76), this goroutine will close output channel so signal from closed input channel will be prapagated to output channel.

Next we have multiple executors connected to `bureaucrat` output channel `docsNew` and spawned by [SpawnExecutors](exercises/e0/main.go#L126). Each `executor` is working in its [own goroutine](exercises/e0/main.go#L99-L102).

We decided to make Fan-In pattern for workers that's why we put their outputs into one channel that is returnet by `SpawnExecutors`. How we will propagate channel signal in this case? 
We make additiona channel `fakeIn` and add `fanoutProxy` goroutine that will forward messages from `In` to `fakeIn`. As soon as `In` will be closed, both `fakeIn` and `out` channels will be closed, so signal will be propagated from input channel to workers and to output channel.

Another solution to achieve the same is to use `WaitGroup` form `sync` package ([Learn more](https://nanxiao.me/en/use-sync-waitgroup-in-golang/)):

```go
func SpawnExecutors(n int, priv ed25519.PrivateKey, in <-chan *Document) <-chan *Document {
	out := make(chan *Document)
	totals := make([]int, n)
    wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
        wg.Add(1)
		go func(i int, wg sync.WaitGroup) {
			totals[i] = executor(priv, fakeIn, out)
            // ...
			wg.Done()
		}(i,wg)
	}
    go func(wg sync.WaitGroup) {
        wg.Wait()
        close(out)
    }
}
```


---

### Mutexes

More info: https://gobyexample.com/mutexes

Generaly if you want to access map or want to modify underlying data storage for any data type from multiple places concurrently, you must use sync.Mutex for every access of the data type.

### sync.Once

If something is needed to be run only Once for multiple call, sync.Once can help you ([open in Playground](https://go.dev/play/p/rmPXf540Qof)):

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
```

### Atomics

More info: https://gobyexample.com/atomic-counters

If you are interesting on what operation we can do atomicaly feel free to take a look atomic package documentation.

In general atomic should be used for logic that allows you to operate with data that doesn't need to be synced but eliminate undefined behaviour.

Ideal example for atomic counters, swaps and access operations is collection of gauge and counter metrics from your code for further shipping to monitoring. We will cover this example deeply on of the next units.

---

### Quiz for Exercice 0

#### Q1. Why [fmt.Println output](exercises/e0/main.go#L101) in executor goroutine is not shown in program output?

1. `fmt` package doesn't work in goroutines other than main
2. `fmt.Output` must be configured to stdout of goroutine
3. goroutine is not in time to finish before program exits
4. goroutine panics because of concurrent access of `totals` slice

## FAQ

### Why accessing counter variable in spawned goroutine is bad?

```go
for i:=0; i < n; i++ {
    go func(){
        _ = i
    }()
}
```

If you are access variables define out of scope a goroutine you should keep in mind will these variables be changed outside of goroutine or not.

I above code at the moment goroutine will access `i` variable, the variable will be changed by surrounding for loop and the value of `i` variable will be unpredicted.

Right way is to do something like that:

```go
for i:=0; i < n; i++ {
    go func(i int){
        _ = i
    }(i)
}
```

We just defined `i` variable in goroutine scope and called goroutine with **value of i**. So we can guarantee that inside goroutine `i` value will be predictable.

---

### Why is it concurrent safe?

```go
ints := make([]int, n)
for i:=0; i < n; i++ {
    go func(i int){
        ints[i]++
    }(i)
}
```

According to the mutex paragraph indeed one can think that we've forgotten to use mutex to access slice elements from different concurrent goroutines.

But nothing unpredictable happens here: every goroutine will access element assigned for particular goroutine. In this code we can guarantee that i-th goroutine will modify only i-th element of predefined array.

Underlying memory container of slice will not be changed, so such approach is concurrent safe.

Opposite situation will be if we'll try to append elements to slice concurrently, in that case we must use implicit syncronisation because append() modifies (or could do) underlying memory and among that it will access undefined index concurrently.

## More quizes

#### Q2. What is an idiomatic way to pause execution of the current scope until an arbitrary number of goroutines have returned?
1. Pass an int and Mutex to each and count when they return.
2. Loop over a select statement.
3. Sleep for a safe amount of time.
4. sync.WaitGroup


#### Q3. What does a sync.Mutex block while it is locked?

1. all goroutines
2. any other call to lock that Mutex
3. any reads or writes of the variable it is locking
4. any writes to the variable it is locking

#### Q4. What is the select statement used for?

1. executing a function concurrently
1. executing a different case based on the type of a variable
2. executing a different case based on the value of a variable
3. executing a different case based on which channel returns first

#### Q5. What is a channel?

1. a global variable
2. a medium for sending values between goroutines
3. a dynamic array of values
4. a lightweight thread for concurrent programming


## Excercises

### E1. Equivalent Binary Trees

More: https://go.dev/tour/concurrency/7 and https://go.dev/tour/concurrency/8

**Note**: according to [documentation](https://cs.opensource.google/go/x/tour/+/refs/tags/v0.1.0:tree/tree.go;l=28) of tree library in the exercise we can see that structure of trees will be the same if they will have same values.
So our task is to get all values stored in tree and compare them one by one.

Keep `main` function as follow:
```go
func main() {
	fmt.Println(Same(tree.New(1), tree.New(2)))
	fmt.Println(Same(tree.New(2), tree.New(2)))
}
```

**Hint**: Feel free do add aditional helper functions. for example you might want to use recursion while tree traversal. If you stuck with traversal you use naive approach:

    1. If left subtree exists (is not nil) run traverse on it (recursion)
    2. Process the value (send it to channel)
    3. If right subtree exists (is not nil) run traverse on it (recursion)

Share your implementation `unit2/exercises/e1/main.go` in github PR, during submission remember following restrictions:

**Don't add additional Prints** to output. It is checked in tests.

**Don't change** `func Same(t1, t2 *tree.Tree) bool` and `func Walk(t *tree.Tree, ch chan int)` definitions. They are checked in tests.

### E2. Web Crawler

More:  https://go.dev/tour/concurrency/10

**Note**: please remember that concurrent access to map is not safe.

**Note**: don't hesitate to spawn a lot of goroutines.

**Note**: don't forget to syncronize workflow of goroutines with crawlers. Crawl function must wait all spawned crawler/workers, etc...

Keep `main` function untouch. Don't add additional Prints to output. It is checked in tests.

Share your implementation `unit2/exercises/e2/main.go` in github PR.