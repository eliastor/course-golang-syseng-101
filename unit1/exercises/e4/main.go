package main

import "golang.org/x/tour/reader"

type MyReader struct{}

func (x MyReader) Read(b []byte) (int, error) {
	var i int = 0
	var e error
	for {
		b = append(b, byte('A'))
		i += 1
	}
	return i, e
}

func main() {
	reader.Validate(MyReader{})
}
