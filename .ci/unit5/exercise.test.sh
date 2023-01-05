#!/bin/bash
UNITN="5"
rm -f "../unit${UNITN}/exercises/e$1/main_test.go"

if [[ $1 == 0 ]]; then
    cp "unit${UNITN}/e$1_main_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_test.go"
else
    cp "unit${UNITN}/main_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_test.go"
    cp "unit${UNITN}/main_internal_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_internal_test.go"

    if [[ $1 < 3 ]]; then
    cp "unit${UNITN}/main_divzero_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_divzero_test.go"
    else 
    cp "unit${UNITN}/main_divzero_eternity_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_divzero_eternity_test.go.tpl"
    fi


    if [[ $1 > 1 ]]; then
    cp "unit${UNITN}/main_expressionerror_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_expressionerror_test.go"
    fi
fi

cd ..

go mod init course || true
go get github.com/stretchr/testify/assert

CGO_ENABLED=0 go test "./unit${UNITN}/exercises/e$1/..."
