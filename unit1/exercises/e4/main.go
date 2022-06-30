package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(b []byte) (n int, err error) {
	for {
		n, err := r.r.Read(b)
		for i := 0; i <= n; i++ {

			b[i] = rot13(b[i])
		}
		if err == io.EOF {
			return n, io.EOF
		}
		return n, nil
	}
}

func rot13(a byte) byte {
	if (a >= 65) && (a <= 90) {
		a = a + 13
		if a > 90 {
			a = a - 26
		}
	} else if (a >= 97) && (a <= 122) {
		a = a + 13
		if a > 122 {
			a = a - 26
		}
	}
	return a
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
