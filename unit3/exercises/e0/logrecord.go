package main

import (
	"encoding"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Logrecord struct {
	IP         string `json:"ip"`
	Username   string `json:"user"`
	Timestamp  uint64 `json:"time"`
	HTTPMethod string `json:"method"`
	URIPath    string `json:"path"`
	Size       uint   `json:"size"`
	HTTPCode   uint   `json:"code"`
}

const (
	apacheDatetimeFormat = "02/01/2006:15:04:05 -0700"
)

var (
	_ encoding.TextUnmarshaler = &Logrecord{}
)

var (
	errMalformed error = errors.New("malformed text")
)

func (r *Logrecord) UnmarshalText(text []byte) error {
	parts := strings.SplitN(string(text), " - ", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: there must be exactly one \" - \" separator", errMalformed)
	}
	lefts := strings.SplitN(parts[0], " ", 2)
	if len(lefts) != 2 {
		return fmt.Errorf("%w: part before \" - \" must contain two fields separated by space", errMalformed)
	}
	r.IP = lefts[0]
	r.Username = lefts[1]

	right := parts[1]
	if right[0] != '[' {
		return fmt.Errorf("%w:  part after \" - \" must starts from \"[\"", errMalformed)
	}
	if right[26] != ']' {
		return fmt.Errorf("%w:  part after \" - \" must have date and time in square brackets", errMalformed)
	}
	t, err := time.Parse(apacheDatetimeFormat, string(right[1:26]))
	if err != nil {
		return fmt.Errorf("%w: couldn't parse date: %v", errMalformed, err)
	}
	r.Timestamp = uint64(t.Unix())
	right = right[27:]

	if right[1] != '"' {
		return fmt.Errorf("%w:  part after \" - \" must contain URI path in double quotes", errMalformed)
	}
	right = right[2:]

	left, right, found := strings.Cut(right, "\"")
	if !found {
		return fmt.Errorf("%w:  part after \" - \" must contain URI path in double quotes", errMalformed)
	}
	lefts = strings.SplitN(left, " ", 3)
	if len(lefts) != 3 {
		return fmt.Errorf("%w:  part after \" - \" must three fields separated by space", errMalformed)
	}
	r.HTTPMethod = lefts[0]
	r.URIPath = lefts[1]

	right = right[1:]
	rights := strings.SplitN(right, " ", 2)
	if len(rights) != 2 {
		return fmt.Errorf("%w: log line must ends with two fields separated by space after HTTP version", errMalformed)
	}
	code, err := strconv.Atoi(rights[0])
	if err != nil {
		return fmt.Errorf("%w: can't convert string to http code: %v", errMalformed, err)
	}
	r.HTTPCode = uint(code)
	size, err := strconv.Atoi(rights[1])
	if err != nil {
		return fmt.Errorf("%w: can't convert string to size: %v", errMalformed, err)
	}
	r.Size = uint(size)

	return nil
}
