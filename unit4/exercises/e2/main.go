package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// func ungzip(w http.ResponseWriter, r *http.Request) {

// 	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
// 		reader, err := gzip.NewReader(r.Body)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		io.Copy(w, reader)

// 	} else {
// 		http.Error(w, "", http.StatusBadRequest)
// 	}

// }

// func ungzip(w http.ResponseWriter, r *http.Request) {

// 	switch r.Method {
// 	case "POST":
// 		testBytes, _ := ioutil.ReadAll(r.Body)
// 		r := bytes.NewReader(testBytes)
// 		if testBytes[0] == 31 && testBytes[1] == 139 {
// 			reader, err := gzip.NewReader(r)
// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 			io.Copy(w, reader)
// 		} else {
// 			http.Error(w, "", http.StatusBadRequest)
// 		}

// 	default:
// 		http.Error(w, "", http.StatusBadRequest)
// 	}

// }
func ungzip(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		testBytes, _ := ioutil.ReadAll(r.Body)
		r := bytes.NewReader(testBytes)
		reader, err := gzip.NewReader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		io.Copy(w, reader)

	default:
		http.Error(w, "", http.StatusBadRequest)
	}

}

func echo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		r := bytes.NewReader(body)
		io.Copy(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
	}
}

func main() {
	mux := http.NewServeMux() // Creating new mux to manage handlers for different paths.

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Handler for all paths
		response := bytes.NewBuffer(nil)

		response.WriteString(r.Method) // GET, POST and other HTTP methods
		response.WriteString(" ")
		response.WriteString("URI: ")
		response.WriteString(r.RequestURI) // Request URI
		response.WriteString(" ")
		response.WriteString("handling with / handler")
		response.WriteByte('\n')

		response.WriteTo(w) // or io.Copy(w, response)
		// Note that w of type http.ResponseWriter implements io.Writer. it can be used with any code supports io.Writer
	})

	mux.HandleFunc("/echo", echo)
	mux.HandleFunc("/ungzip", ungzip)
	addr := ":8080"

	fmt.Println("Starting server on " + addr)

	http.ListenAndServe(addr, mux)
}
