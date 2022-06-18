package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawl(t *testing.T) {

}

func Test_main(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outCh := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outCh <- buf.String()
	}()

	Crawl("https://golang.org/", 4, fetcher)
	w.Close()
	os.Stdout = old
	output := <-outCh
	entries := []string{
		`found: https://golang.org/ "The Go Programming Language"`,
		`found: https://golang.org/pkg/ "Packages"`,
		`found: https://golang.org/pkg/os/ "Package os"`,
		`not found: https://golang.org/cmd/`,
		`found: https://golang.org/pkg/fmt/ "Package fmt"`,
	}

	for _, entry := range entries {
		assert.Contains(t, output, entry)
		output = strings.Replace(output, entry, "", 1)
		assert.NotContains(t, output, entry, "There must be no duplicates")
	}

}
