#!/bin/bash
rm -f ../unit3/exercises/e$1/main_test.go
cp unit3/e$1_main_test.go.tpl ../unit3/exercises/e$1/main_test.go

cd ..

go mod init course || true

go get github.com/brianvoe/gofakeit
go get github.com/stretchr/testify/assert
go get golang.org/x/tour/tree
go get github.com/adamliesko/fakelog

CGO_ENABLED=0 go test ./unit3/exercises/e$1/...
