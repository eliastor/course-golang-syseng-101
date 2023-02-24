# Unit 6: Contexts

For this unit, your environment must be initialized as in previous units.

## Materials

Excessive material can be found on [official go site](https://go.dev/blog/context). Take a look, **but don't dig it too much**.

Context package is a standard library [package](https://pkg.go.dev/context) which offers a mechanism to control cancellation and timeout of operations outside of Go program scope, for example: calling external applications, accepting and opening socket connections, working with IO and so on.

### Context

Context is an interface:

```go
// A Context carries a deadline, cancellation signal, and request-scoped values
// across API boundaries. Its methods are safe for simultaneous use by multiple
// goroutines.
type Context interface {
    // Done returns a channel that is closed when this Context is canceled
    // or times out.
    Done() <-chan struct{}

    // Err indicates why this context was canceled, after the Done channel
    // is closed.
    Err() error

    // Deadline returns the time when this Context will be canceled if any.
    Deadline() (deadline time.Time, ok bool)

    // Value returns the value associated with key or nil if none.
    Value(key interface{}) interface{}
}
```

Every time an applications hits with possibility of timeouts or interruptions in some actions contexts enter a game. A context is a small object that holds execution context, generally it is time limit and link to parent context.

In application code contexts have two parts that can be presented in code together: control and execution pieces.

Let's figure out in mechanics of contexts by exploring these two use-cases of contexts.


### Context creation and control parts

The highest possible context is one returned by [`context.Background()`](https://pkg.go.dev/context#Background) which stands for the whole application context, you should use that if you create your own contexts not linked to anything incoming (http or socket requests) to your application.
If you receive context, for example in http handler you can get context by calling[Request.Context()](https://pkg.go.dev/net/http#Request.Context). For incoming server requests, the context is canceled when the client connection is closed, the request is canceled (with HTTP/2), or when the ServeHTTP method of http server returns error:

```go
//ordinary http handler
func(w http.ResponseWriter, r *http.Request) { 
    // if you want to get data from something external that requires context for working, you must get context from current request.
    ctx := r.Context()

    // for example let's run some heavy script.
    // Following line will cancel command execution if request r is canceled
    // in most cases it can be canceled because underlying socket connection is closed.
    cmd := exec.CommandContext(ctx, "sleep", "5").Run();
    // cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // Some magic to force stop of child processes.
    // Only valid for Linux builds and architectures that use libc under the hood
	}
```

If you want to limit the execution time of your request you need to create new context inherited from one from request:

```go
func(w http.ResponseWriter, r *http.Request) { 
    // if you want to get data from something external that requires context for working, you must get context from current request.
    parentCtx := r.Context()
    ctx, cancel := context.WithTimeout(parentCtx, 1*time.Second)
    defer cancel() // Failing to call the CancelFunc leaks the child and its children until the parent is canceled or the timer fires.
    // we've just created context that will be:
    // - canceled if parentCtx is cancelled
    // - canceled after 1 second
    // ...
}
```

Examples above show you how to get context of incoming requests.
If a library or code you are using doesn't offer context capabilities but assumes that it can timeout. Then try to find another library or use [`context.TODO()`](https://pkg.go.dev/context#TODO) as placeholder. This `TODO` context is empty one which never gets canceled.

New contexts must be created using `context` package. You can create cancelable context with or without time limit.

If you don't need time limit, use [`context.WithCancel(...)`](https://pkg.go.dev/context#WithCancel).

If you need to set timeout for context, use [`context.WithTimeout(...)`](https://pkg.go.dev/context#WithTimeout). After the specified timeout context will be cancelled.

If you need to set deadline for context, use [`context.WithDeadline(...)`](https://pkg.go.dev/context#WithDeadline). After the internal time clock gets to deadline the context will be cancelled.

Pay attention to call cancel function for contexts created in your code, otherwise contexts will leak the memory.

```go
import (
	"context"
	"os/exec"
	"time"
)

func  main() {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // Failing to call the CancelFunc leaks the child and its children until the parent is canceled or the timer fires.

	if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}
}
```

----

### Support of context cancellation in your code. Execution piece.

Imagine that you process a large amount of data in goroutine (note that http handler is a separate goroutine too) or at least any routine you want to be expired or cancelable. For example we have a code of http handler that should synchronously process some files on the server. Processing time of each file is not long, but there can be a lot of files (millions for an instance):

```go
func(w http.ResponseWriter, r *http.Request) {
    files := listFiles()
    for _, file := range files { // there can be millions of files
        err := processFile(file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}
```

We want this computation to be done synchronously with our http request, so it means, that if the request was cancelled for any reason, our processing must be cancelled too. We can slightly rewrite handler and add context support:

```go
func(w http.ResponseWriter, r *http.Request) {
    files := listFiles()
    ctx := r.Context()
    for _, file := range files { // there can be millions of files
        select { // select statement blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready
            case <-ctx.Done(): //that's channel reading, it only can happen if context is cancelled.
                return
            default: // the default case in a select is run if no other case is ready
                err := processFile(file)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                }
        }
    }
}
```

What just happened is in every cycle of the loop we're checking is context still valid or not. If not we immediately stop computation and return from handler, otherwise we continue to process specific file.

For checking the context there is `Done()` function that return channel. Under the hood a context signals about its ending by closing this internal channel. As soon as it's the only message can be read from the channel, we are not interested in value of the channel (in fact it's nothing - `struct{}`).

For making the logic of checking the channel state and in case of the channel is not closed, we use following pattern:

```go
select {
    case <-ctx.Done(): 
        return
    default:
    // do normal routine
}
```

It's widely used pattern, first golang check that all cases cannot be executed: it's not possible to receive (`<-`) from not closed channel returned by `ctx.Done()`, so the code under `default:` statement will be executed.


### Note on usage

In most of the cases the work with contexts is reduced to providing context from the source (http handler for example) to any context-aware functions which are used in your code. Almost every client to database, external process execution, http and socket requests - all of them have capability to use specific context in their operation. Your responsibility as developer is to just get context from incoming event (in most cases it will be context of socket connection or http request) provide the context for these operations during your application execution.


---

### E0. The most useless but busy code

To build exercise, being in root folder of the repo you can run:

```bash
go build ./unit6/exercises/e0
```

This command will build local folder with all `.go` files in it and place result application to `e0` file in current (repo root folder).

If you want to specify name of path of the file:

```bash
go build -o ./unit6/e0 ./unit6/exercises/e0
```

for exercise of these and next unit it is handy to build and run in one command:

```bash
go run ./unit6/exercises/e0
```

It will build and run code in `./unit6/exercises/e0`:

```bash
./unit6/exercises/e0
```


Find [source code](exercises/e0/main.go) of this exercise.

---

## FAQ

TBA

---

## Exercises

### E1. http api for processing

Extend code from exercise 0: Make http server that will:
- listen 8080 port
- call cycle() and return the output for url path `/cycle/<N>`, where `<N>` is number of cycle invocations of `cycle()` function. if no `<N>` specific (for example `curl http://localhost:8080/cycle/`) assume that N is 5.
- when incoming request is cancelled, the handler must output to stdout or stderr of the server a message with following string: `incoming request %s was canceled by client`, where %s is request URI, for example: `2023/01/01 11:11:11 incoming request /cycle/ was canceled by client`.


**Note**: Test verifies response body and headers of your http server by running it and sending generated request to the server and verifying the reply. Server must listens port 8080. You may organize code as your own.

Don't add additional Prints to output. It is checked in tests.

Don't change `process()` and `cycle()` functions, they are verified in tests.

Share your implementation `unit6/exercises/e1/main.go` in github PR.
Don't hesitate to copy contents of `unit6/exercises/e0/` to `unit6/exercises/e1/` and modify necessary files or add new ones.

**Hint**

to check if error in context is present use `Err()` method of the context, find out which error is it with `errors.Is()`:
`errors.Is(err, context.Canceled)`

**Hint**

to check if error in context is present use `Err()` method of the context, find out which error is it with `errors.Is()`:
`errors.Is(err, context.Canceled)`

**Hint**

Convenient way to write string to http response is:

```go
func(w http.ResponseWriter, r *http.Request) {
		...
        io.WriteString(w, "string")
}
```

**How to test yourself**

In terminal1 run 

```bash
go run ./unit6/exercises/e1/
```

in terminal 2 run

```bash
curl -v http://localhost:8080/cycle/
```

terminal 2 will be delay for for 10 seconds and output:
```
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /cycle/ HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 09 Feb 2023 07:30:38 GMT
< Content-Length: 125
< Content-Type: text/plain; charset=utf-8
< 
iteration 0 is processed
iteration 1 is processed
iteration 2 is processed
iteration 3 is processed
iteration 4 is processed
* Connection #0 to host localhost left intact
```

in terminal 2:

```bash
curl -v http://localhost:8080/cycle/10
```

and cancel execution (Ctrl+C)

in terminal 2 similar message must appears:
```
incoming request /cycle/ was canceled by client
```



### E2. timeout of processing

Extend code from exercise 1: 

- for `/cycle/` there is url query parameter `timeout` that can setup time limit of the request handler, for ex. `curl http://localhost:8080/cycle/?timeout=11`, means that after 11 seconds `cycle()` must return value. Default value for timeout is 11
- when incoming request is cancelled due to timeout, the handler must output to stdout or stderr of the server a message with following string: `incoming request %s was canceled by client`, where %s is request URI, for example: `2023/01/01 11:11:11 incoming request /cycle/ reached timeout`.

**Note**: Test verifies response body and headers of your http server by running it and sending generated request to the server and verifying the reply. Server must listens port 8080. You may organize code as your own.

Don't add additional Prints to output. It is checked in tests.

Don't change `process()` and `cycle()` functions, they are verified in tests.

Share your implementation `unit6/exercises/e2/main.go` in github PR.
Don't hesitate to copy contents of `unit6/exercises/e1/` to `unit6/exercises/e2/` and modify necessary files or add new ones.

**Hint**

Inside handler create context with timeout from the context of the request and pass it to functions calls in this handler. Don't forget to `defer cancel()`.

**How to test yourself**

Same as in previous exercise, but by default the first test must fail with timeout in:
in terminal 1 similar message must appears:

```
incoming request /cycle/ reached timeout
```

To run normally test yourself with: 

in terminal 2:

```bash
curl -v http://localhost:8080/cycle/2?timeout=8
```

terminal 2 will be delay for for 10 seconds and output:
```
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /cycle/ HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 09 Feb 2023 07:30:38 GMT
< Content-Length: 125
< Content-Type: text/plain; charset=utf-8
< 
iteration 0 is processed
iteration 1 is processed
* Connection #0 to host localhost left intact
```


### E3. interrupt processing

Extend code from exercise 2:
Propagate context form `cycle()` to `process()` and change `process()` to be able to interrupt when context was cancelled. When `process()` it must return string `iteration %d is cancelled` where %d is a number of the iteration.

**Note**: Test verifies response body and headers of your http server by running it and sending generated request to the server and verifying the reply. Server must listens port 8080. You may organize code as your own.

**Note**: To test your implementation you can run 

**Hints**

 - [ResponseWriter.WriteHeader(statusCode int)](https://pkg.go.dev/net/http#ResponseWriter)


**How to test yourself**

Same as in previous exercises, but by default the first test (`curl -v http://localhost:8080/cycle/`) must also return response:
in terminal 2 similar message must appears:

```
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /cycle/ HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 09 Feb 2023 08:01:19 GMT
< Content-Length: 50
< Content-Type: text/plain; charset=utf-8
< 
iteration 0 is processed
iteration 1 is cancelled
* Connection #0 to host localhost left intact
```

---

## Additional materials on contexts