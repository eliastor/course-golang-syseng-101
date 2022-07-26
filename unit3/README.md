# Unit 3: First project. Streams

## Materials

For this unit you need your module "course" to be initialized as in unit2.
Note that your module must be named course.

Install VSCode go plugin: https://marketplace.visualstudio.com/items?itemName=golang.go.
Open "Command Palette" in VSCode (View -> Command Palette) and paste "go.tools.install" in it. Install all suggested tools and go.

your course repository must be like this:

```
.
├── README.md
├── go.mod
├── go.sum
├── unit1
├── unit2
├── unit3
│   ├── README.md
│   └── exercises
│       ├── e0
|       |   ├── fake.log
|       |   ├── fake.log.gz
|       |   ├── generator.go
|       |   ├── logrecord.go
│       │   └── main.go
│       ├── e1
│       ├── e2
...
```

Open course repository in VSCode IDE so you can see file tree structure mentioned above.

In this unit we'll use library "github.com/adamliesko/fakelog/generator":

```sh
go get "github.com/adamliesko/fakelog/generator"
```

More information about modules: https://go.dev/blog/using-go-modules
Also https://go.dev/doc/code might be helpful

This unit covers two most popular concepts: interfaces and IO operations.

### Interfaces

Refresh interfaces syntax and key concepts:
https://gobyexample.com/interfaces and https://go.dev/tour/methods/9

Interface is set of function definitions. Specific type implements interface when type has all functions defined in interfaces with same definition.

Let's imagine that you application requires some Key-Value in-memory storage for strings, for example cache. You've implemented it using `map[string]string`:

```go
type Cache struct{
    data map[string]string
    mu sync.Mutex
}
func  (c *Cache) Get(key string) string {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.data[key]
}
func  (c *Cache) Put(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key]=value
}
func NewCache() *Cache {
    cache := &Cache {
        data: make(map[string]string)
    }
    return cache
}
```

All your code uses cache as follows:

```go
cache := NewCache()
...
func GetFoo(cache *Cache) string{
    value := cache.Get("foo") // try to get value from cache
    if bar == "" {
        // get or generate it from somewhere else
        // value = GetFooFromDatabase()
        cache.Put("foo", value)
    }
    return value
}
...
_ = GetFoo(cache)
```

In code above you are using Cache struct with some method defined. Function `NewCache()` prepares new instance of cache and returns it. Such generator function are named **Factory method** or **Factory Pattern**. Factory is widely used concept in Go.

Suddenly you realized that your cache must be persistent across restarts, but you want to save in-memory one for stateless workloads or development purposes.

So all our code (`GetFoo(cache)` in this example) must be able to work with two types of caches.

Interfaces will help you, let's define interface for Generic cache:

```go
type Cacher interface {
    Get(key string) string
    Put(key, value string)
}
```

Rewrite your code:

```go
...
func GetFoo(cache Cacher) string{ // The only thing you changed is definition of function
    // now GetFoo function accept interface Cacher instead of specific Cache struct
    // Rest of your code is being unmodified because you've created interface with same methods you had.
    value := cache.Get("foo") // try to get value from cache
    if bar == "" {
        // get or generate it from somewhere else
        // value = GetFooFromDatabase()
        cache.Put("foo", value)
    }
    return value
}
...
_ = GetFoo(cache)
```


Let's define 2 types of storages:

1. in-memory: 

```go
type inMemoryStorage struct {
    data map[string]string
    mu sync.Mutex
}

func NewInMemStore() Cacher {
    store := &inMemoryStorage {
        data: make(map[string]string)
    }
    return store
}

func  (s *inMemoryStorage) Get(key string) string {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.data[key]
}
func  (s *inMemoryStorage) Put(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key]=value
}
```

2. file-backed:

```go
type fileBackedStorage string

func NewFileStorage(rootpath string) Cacher {
    store := fileBackedStorage(rootpath)
    return &store
}

func  (s *fileBackedStorage) Get(key string) string {
    filename := filepath.Join(s, key)
    data, _ := os.ReadFile(filename)
	return string(data)
}
func  (s *fileBackedStorage) Put(key, value string) {
    filename := filepath.Join(s, key)
    _ = os.WriteFile(filepath, []byte(value), 0660)
}
```

Every of these storages satisfy `Cache` interface and can be passed to every function that expects `Cache` interface.

```go
cache := NewInMemStore()
// or cache := NewFileStorage("/var/lib/app/cache")
_ = GetFoo(cache) // it will work no matter which Store was created and passed
```

As soon as `NewInMemStore()` and `NewFileStorage()` return object that satisfy `Cacher` interface so bot of them can be passed to everywhere `Cacher` interface is expected.

`GetFoo` isn't aware of underlying cache strategy (in-memory or file-backed), the only thing it's needed is `Cacher` interface. Implementation details are hided behind the interface in specific implementation of these interface.

NB: Go-way is to accept as much generic interfaces and return more specific implementations (structs).
Everywhere when you need behaviour instead of specific data structure, interfaces are preferable. If you take a look on go standard library you'll see that it's full of interfaces, especially in io and cryptography.


### IO and streams

At this point you're aware of most syntax and features of Go. But the language has conceptual features that are widely used in every real-world project. One of the such key concepts is IO operations and streams.

Stream is some io.Reader/io.Writer interface that allows to read/write chuncks of bytes.
Most of streams operations ends with io.Copy(r,w) under the hood.

There are lot of things shown as io.Reader/io.Writer interfaces in Go: file descriptor, standard input/output/error, http socket, unix socket, output/input of gzip compression, encryption algorithm input/output, etc...

### io.Reader

There is [io.Reader](https://pkg.go.dev/io#Reader) interface:

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

Implementation of the `Read()` should return the number of bytes read or an error if one occurred. If the source has returned all its content, Read should return io.EOF.

The behavior of a reader will depend on its implementation, however there are a few rules, from the io.Reader doc that you should be aware of when consuming directly from a reader:

1. `Read()` will read up to len(p) into p, when possible.
2. After a `Read()` call, n may be less then len(p).
3. Upon error, `Read()` may still return `n` bytes in buffer `p`. For instance, reading from a TCP socket that is abruptly closed. Depending on your use, you may choose to keep the bytes in p or retry.
4. When a `Read()` exhausts available data, a reader may return a non-zero n and `err=io.EOF`. However, depending on implementation, a reader may choose to return a non-zero `n` and `err=nil` at the end of stream. In that case, any subsequent reads must return `n=0`, `err=io.EOF`.
5. A call to `Read()` that returns `n=0` and `err=nil` does not mean `io.EOF` as the next call to `Read()` may return more data.

`Read()` is designed to be called within a loop where, with each iteration, it reads a chunk of data from the source and places it into buffer p. This loop will continue until the method returns an io.EOF error.

```go
reader := strings.NewReader("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi eleifend sapien at aliquet semper.")
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err != nil{
		    if err == io.EOF {
                fmt.Println(string(p[:n])) //should handle any remainding bytes.
                break
		    }
            fmt.Println(err)
            os.Exit(1)
		}
		fmt.Println(string(p[:n]))
	}
```

The source code above creates a 4-byte long transfer buffer p with make([]byte,4). The buffer is purposefully kept smaller then the length of the string source. This is to demonstrate how to properly stream chunks of data from a source that is larger than the buffer.

Let's implement Reader that will return stream of letters `A` with specific length:

```go
// You can easily create methods for any type, not only for structs.
// In real life code structs are used for any objects but sometimes you can find other types with methods.
// In this example we only need a counter, which holds number of returned letters `A`.
// As soon as we need to modify the counter we defined meethod receiver as pointer to aaa.
type aaa int

func (r *aaa) Read(p []byte) (int, error) {
	n := len(p)
	var err error

    // int(*r) - gets value from address r and as soon as it has type aaa, converts it to int.
	if int(*r) < len(p) {
		if int(*r) >= 0 {
			n = int(*r)
		}
		err = io.EOF
	}
    // `*r = ...` means we need to set value pointed by pointer r. Note that value must have aaa type.
    // `aaa(...)` means to convert braced value to aaa type.
	*r = aaa(int(*r) - n)
	copy(p, bytes.Repeat([]byte("A"), n))
	return n, err
}

// NewAAAReader creates aaa Reader with specified length
func NewAAAReader(len int) io.Reader {
	r := aaa(len)
	return &r
}
```

Now if we'll create `aaa` reader with length of 10 and read it by executing Read or using [io.ReadAll](https://pkg.go.dev/io#ReadAll)  we'll receive `AAAAAAAAAA` (10 A's).

Try it in [playground](https://go.dev/play/p/pRJ4xvj6RSS)

### io.Reader chaining

The standard library has many readers already implemented. It is a common practice in Go to use a reader as the source of another reader.

Let's make reader that will lower case of all letters from another reader:

```go
type lowercaser struct {
	r io.Reader
}

func (r *lowercaser) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	s := string(p)
	lowered := strings.ToLower(s)
	copy(p, []byte(lowered))
	return n, err
}

func NewLowerCaser(r io.Reader) io.Reader {
	l := &lowercaser{r}
	return l
}
```

Now we can chain our two readers:

```go
func main() {
	r := NewAAAReader(10)
	lowerer := NewLowerCaser(r)
	lowerred, err := io.ReadAll(lowerer)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(lowerred))
}
```

We used [io.ReadAll](https://cs.opensource.google/go/go/+/refs/tags/go1.18.4:src/io/io.go;l=638) which under the hood call Read iterativly and appends readed data to buffer. As soon as error occured it returns buffer and error (nil error if there was err==io.EOF).

Try it in [playground](https://go.dev/play/p/3cWjb3B86IT)

### io.Writer

There is [io.Writer](https://pkg.go.dev/io#Writer) interface:

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

The method is designed to read data from buffer p and write it to a specified target resource.
Implementation of the Write() method should return the number of bytes written or an error if any occurred.

Let's take a look on [bytes.Buffer](https://pkg.go.dev/bytes#Buffer) and how it implements `io.Writer`

```go
b := bytes.Buffer{} // a Buffer needs no initialization.
b.Write([]byte("hello world!")) //Buffer implements io.Writer so it has Write() method
fmt.Println(b.String()) // buffer has String() method that return contents of the Buffer as string.
```

### io.Copy()

Function [io.Copy()](https://pkg.go.dev/io#Copy) makes it easy to stream data from a source reader to a target writer. It abstracts out the for-loop pattern (you've seen above) and properly handle `io.EOF` and byte counts. Under the hood io.Copy read data from `io.Reader` `Read()` to buffer in a loop and writes the buffer to `io.Writer` `Write()`

Let's copy output of lowerCaser to bytes.Buffer:

```go
func main() {
	r := NewAAAReader(10)
	lowerer := NewLowerCaser(r)

    b := bytes.NewBuffer(nil)
    io.Copy(b, lowerer) // it will copy all the data from lowerer, which will read all the data from r until io.EOF, to bytes.Buffer.

	fmt.Println(b.String())
}
```

---

### Another useful stream objects:

1. to copy only specifc amount of data: [io.CopyN](https://pkg.go.dev/io#CopyN)
2. to connect Reader and Writer (for exampel for compression): [io.Pipe](https://pkg.go.dev/io#Pipe)
3. to read specifc amount of data: [io.LimitedReader](https://pkg.go.dev/io#LimitedReader)
4. to seek and read sections (for archives, files and S3 objects): [io.SectionReader](https://pkg.go.dev/io#SectionReader)
5. to make different chains for Reader and Writer (like http(s) reverse proxy): [bufio.ReadWriter](https://pkg.go.dev/bufio#ReadWriter)
6. to work with files via os.Create(), os.Open(), ... : [os.File](https://pkg.go.dev/os#File)
7. to parse simple formats: [bufio.Scanner](https://pkg.go.dev/bufio#Scanner)
8. to make a copy of stream (like tee command): [io.TeeReader](https://pkg.go.dev/io#TeeReader)


## E0. Web server logs to json converter

For purposes of this exercise let's write generator that can be used in abscense of real log file.
Our converter is reading data line by line and converting every line to json. Json must contain client IP, HTTP method, URI path, response code, response size and timestamp of the request.
If real log file is provided, the app is converting provided file.

Find [source code](exercises/e0/main.go) of this exercise.

---


#### Q1. Is following satisfy io.Reader interface?

```go
type Null struct {}
func (_ Null) Read(p []byte) (int, error) {
    return 0, io.EOF
}
```

1. yes
2. no, it doesn't have pointer receiver
3. no, pointer receiver is not defined (`(_ Null)`)
4. no, struct is empty

#### Q2. What is idiomatic way to read all data from io.Reader to memory?

1. call Read() in loop and append to bytes.Buffer{}
2. use io.ReadAll()
3. call Read() in loop and append to slice of bytes
4. use io.Copy() to bytes.Buffer{}


#### Q3. What is the most idiomatic definition of functions and methods?

1. accepts interfaces, returns interfaces
2. accepts structs, returns interfaces
3. accepts interfaces, returns structs
4. accepts structs, returns structs


----

## FAQ

TBA

## Excercises

### E1. Adding read from stdin

Extend code from exercise 0 by adding support of reading from stdin if filename is "-"

**Note**: Go has `os.Stdin` which satisfy io.Reader interface and represents standard input of application.

Don't add additional Prints to output. It is checked in tests.

Share your implementation `unit3/exercises/e1/main.go` in github PR.
Don't hesitate to copy contents of `unit3/exercises/e0/` to ``unit3/exercises/e1/` and modify necessary files.

### E2. Adding read from gzipped files

Extend code from exercise 1 by adding capability of reading from gzipped files. Note that as a result of this exercise application must be able to read from both stdin and gzipped files. Fake generator must be also present as it was in e0.

Hint 1:

```go
if strings.HasSuffix(filename, ".gz") || strings.HasSuffix(filename, ".gzip") {
    // ...
}
```

Hint 2:
[gzip example](https://pkg.go.dev/compress/gzip#example-package-WriterReader)

Don't add additional Prints to output. It is checked in tests.

Share your implementation `unit3/exercises/e2/main.go` in github PR.
Don't hesitate to copy contents of `unit3/exercises/e1/` to ``unit3/exercises/e2/` and modify necessary files.
