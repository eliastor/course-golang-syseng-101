package main

import (
	"bytes"
	"fmt"
	"net/http"
)

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

	addr := ":8080"

	fmt.Println("Starting server on " + addr)

	http.ListenAndServe(addr, mux)
}
