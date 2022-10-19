package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	filename := "fake.log"

	log.SetFlags(0) // Don't show any additional information while printing to application log (stderr)

	var scannerInput io.Reader

	switch len(os.Args) {
	case 2:
		// if filename specified open the file
		filename = os.Args[1]
		// os.Open returns file handler which satisfy io.Reader interface so we can read stream of bytes from file.
		flog, err := os.Open(filename)
		if err != nil {
			log.Println(err)
			return
		}
		defer flog.Close()

		scannerInput = flog
	case 1:
		generator := NewFakeLogGenerator()
		defer generator.Close()

		scannerInput = generator
	default:
		log.Println("zero or one argument allowed")
		os.Exit(1)
	}

	// linescanner allows us to scan input stream of bytes from flog and split the stream to lines: https://pkg.go.dev/bufio#Scanner
	// as soon as flog satisfy io.Reader we can use it as argument for NewScanner

	linescanner := bufio.NewScanner(scannerInput)
	linescanner.Split(bufio.ScanLines)

	enc := json.NewEncoder(os.Stdout)

	for i := 0; linescanner.Scan() && (i < 10000); i++ {
		text := linescanner.Bytes() // if you need string, use linescanner.Text()
		rec := &Logrecord{}
		err := rec.UnmarshalText(text)
		if err != nil {
			log.Println("unable to parse line:", err)
			continue
		}

		enc.Encode(rec)
	}

}
