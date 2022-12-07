package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var results []string

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
		fmt.Fprint(w, string(body))
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
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

	addr := ":8080"

	fmt.Println("Starting server on " + addr)

	http.ListenAndServe(addr, mux)
}
