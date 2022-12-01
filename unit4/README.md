# Unit 4: Web server

For this unit your environemnt must be initialized as in previous units.

## Materials

Extessive material can be found on [official go site](https://go.dev/doc/articles/wiki/). Take a look, try to reproduce it, **but don't dig it too much**.

[Go standard library](https://pkg.go.dev/net/http) offers great capabilities to build http servers.

## Servers, multiplexers, routers

Most of the examples with implementation ща Go http servers with standard library can be reduced to code:

```go

import (
    "net/http"
)

func newRootHandler() {
    // Here we define function and return it
    // note that we return function and under the hood it is used like this:
    //      handler := newRootHandler()
    // handler variable is address of function which we can call:
    //      handler(w, r)

    return func(w http.ResponseWriter, r *http.Request) {
	    log.Println("request /")
	    io.WriteString(w, "Some data to be returned\n")
    }
}

func main() {
    http.HandleFunc("/", newRootHandler())
    // more handlers

	err := http.ListenAndServe(":3333", nil) // application execution will stop here until server will be stopped
}
```

Note that in example above `http` is a package.
By calling HandleFunc you add handlers for particular paths in URL, for example if there is **handler1** `/api/` it handles all paths started from `/api`: `/api/v1`, `/api/v1/foo`, `/api/` and `/api`. If additional **handler2** is added to `/api/v1/test` path, this partical path will be handled by handle2 while all other `/api` paths will be handled by **hanlder1**:

```shell
/ - default handler
├── /api/ - handler1
|    └── /v1 - handler1
|    |   ├── /foo - handler1
|    |   |   └── /* - handler1 # any subpath of /api/v1/foo
|    |   ├── /test - handler2
|    |   |   └── /* - handler1 # any subpath of /api/v1/test
|    |   └── /* - handler1 # any subpath of /api/v1/
|    └── /* - handler1 # any other subpath of /api/
└── /* - default handler # any other subpath of /
```

Note that for example `/api/v1/test/foo` will be managed by handler1 instead of handler2 because hanlder2 manages `/api/v1/test` only. If you want handler2 to manage all subpaths of `/api/v1/test` you should define path with trailing `/`: `/api/v1/test/`

Such tree of handlers is managed by things commonly called as **http multiplexers** or **routers**.

Tearing down `http` package you can find that there is (DefaultServeMux)[]. Everytime you specify nil as handler in [`http.ListenAndServe(":3333", nil)`](https://pkg.go.dev/net/http#ListenAndServe) http library will use DefaultServeMux defined in http package. In most production application it's better and more maintainable to create own instance of [**multiplexer**](https://pkg.go.dev/net/http#ServeMux) of **mux** in short and configure it:

```go
mux := http.NewServeMux()
mux.HandleFunc("/", newRootHandler())
...
```

Non-default mux allows you to spinup multiple different servers listen different port for different purposes, for example separate servers on different ports/addresses for web, api and metrics:

```go
webmux := http.NewServeMux()
webmux.HandleFunc("/", newWebHandler())

apimux := http.NewServeMux()
apimux.HandleFunc("/api/v1/", newApiV1Handler())
apimux.HandleFunc("/api/v2/", newApiV2Handler())

metricsmux := http.NewServeMux()
metricsmux.HandleFunc("/metrics", newPrometheusMetricsHandler())

wg := sync.WaitGroup{}

wg.Add(1)
go func(){
    defer wg.Done
    log.Println(http.ListenAndServe(":80", webmux)) // goroutine execution will stop here until server will be stopped
}

wg.Add(1)
go func(){
    defer wg.Done
    log.Println(http.ListenAndServe(":8080", apimux)) // goroutine execution will stop here until server will be stopped
}

wg.Add(1)
go func(){
    defer wg.Done
    log.Println(http.ListenAndServe(":3000", metricsmux)) // goroutine execution will stop here until server will be stopped
}

wg.Wait() // application execution will stop here until aat least one server works
```

---

When we call [`http.ListenAndServe()`](https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/net/http/server.go;l=3253) it creates [`http.Server`](https://pkg.go.dev/net/http#Server) in-the-fly and calls `ListenAndServe()` fucntion of created `http.Server`. 

By creating new [`http.Server`](https://pkg.go.dev/net/http#Server) instance you can tune http server that can be also be useful in produciton code, for example if you need custom TLS or timeouts configuration.

Server listens the socket, accepts incoming connections and pass handling of the connection to mux handler.

---

### E0. Dumb web server

To build exercise, being in root folder of the repo you can run:

```bash
go build ./unit4/exercises/e0
```

This command will build local folder with all `.go` files in it and place result application to `e0` file in current (repo root folder).

If you want to specify name of path of the file:

```bash
go build -o ./unit4/e0 ./unit4/exercises/e0
```

for oxercies of these and next unit it is handy to build and run in one command:

```bash
go run ./unit4/exercises/e0
```

It will build and run code in `./unit4/exercises/e0`:

```bash
./unit4/exercises/e0
```


Find [source code](exercises/e0/main.go) of this exercise.

---

## FAQ

TBA

---

## Quiz

#### Q1. What handler will be called for "/a/b/c?foo=bar"  for following mux code?

```
	mux.HandleFunc("/", defaultHandler)
	mux.HandleFunc("/a/", aHandler)
	mux.HandleFunc("/b", bHandler)
	mux.HandleFunc("/c", cHandler)
	mux.HandleFunc("/a/b", abHandler)
```

1. aHandler
2. bHandler
3. cHandler
4. abHandler


## Excercises

### E1. Echo

Extend code from exercise 0: Add handler to "/echo" path which for "POST" requests will read request body and echo them back.

**Note**: Test verifies the output of your program by running it and sending generated request to your http server and verifies a reply. Server must listens port 8080. You may organize code as your own.

Don't add additional Prints to output. It is checked in tests.

Share your implementation `unit4/exercises/e1/main.go` in github PR.
Don't hesitate to copy contents of `unit4/exercises/e0/` to `unit4/exercises/e1/` and modify necessary files or add new ones.

**Hint**

```
		switch r.Method { // just for illustration let's read whole body
		case http.MethodPost, http.MethodPut:
			io.Copy(io.Discard, r.Body) // read request Body to null
		}
```

### E2. Ungzipper

Extend code from exercise 1: Add handler to "/ungzip" path which for "POST" requests will read request body, ungzip the input and send it back. In case of error, server must return HTTP Status Code 400 (Bad Request)

**Note**: Test verifies the output of your program by running it and sending generated request to your http server and verifies a reply. Server must listens port 8080. You may organize code as your own.

**Note**: To test your implementation you can run 
`echo "AAAAAAAAAA=" | gzip | curl --data-binary @- -v http://localhost:8080/ungzip`. It must return something like:

```
*   Trying ::1:8080...
* Connected to localhost (::1) port 8080 (#0)
> POST /ungzip HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.74.0
> Accept: */*
> Content-Length: 25
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 25 out of 25 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 27 Oct 2022 06:05:26 GMT
< Content-Length: 12
< Content-Type: text/plain; charset=utf-8
< 
AAAAAAAAAA=
* Connection #0 to host localhost left intact
```

Key point here is that you receive ungzipped "AAAAAAAAAA=" string.

For malformed request (for example, for non gzipped request body: `echo "AAAAAAAAAA=" | curl --data-binary @- -v http://localhost:8080/ungzip`), you must receive:

```
*   Trying ::1:8080...
* Connected to localhost (::1) port 8080 (#0)
> POST /ungzip HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.74.0
> Accept: */*
> Content-Length: 12
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 12 out of 12 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 400 Bad Request
< Date: Thu, 27 Oct 2022 06:07:47 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

Key point here is to get `HTTP/1.1 400 Bad Request` response.

Don't add additional Prints to output. It is checked in tests.

Share your implementation `unit4/exercises/e1/main.go` in github PR.
Don't hesitate to copy contents of `unit4/exercises/e0/` to `unit4/exercises/e1/` and modify necessary files or add new ones.

**Hints**

 - [ResponseWriter.WriteHeader(statusCode int)](https://pkg.go.dev/net/http#ResponseWriter)
 - `http.StatusBadRequest` - constant defined in http package for HTTP status code 400
 - [gzip.Reader](https://pkg.go.dev/compress/gzip#Reader)

---

## Additional materials on web servers