#!/bin/bash

cp unit2/e$1_main_test.go.tpl ../unit2/exercises/e$1/main_test.go

cd ..

go mod init course || true

go get github.com/brianvoe/gofakeit
go get github.com/stretchr/testify/assert
go get golang.org/x/tour/tree

CGO_ENABLED=0 go test ./unit2/exercises/e$1/...
