#!/bin/bash
UNITN="6"
rm -f "../unit${UNITN}/exercises/e$1/main_test.go"

if [[ $1 == 0 ]]; then
    echo "no files for e0"
else
    cp "unit${UNITN}/main_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_test.go"
    cp "unit${UNITN}/main_cycle_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_cycle_test.go"

    if [[ $1 == 3 ]]; then 
        cp "unit${UNITN}/main_cycle_cancelled_e3_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_cycle_cancelled_e3_test.go"
        
        cp "unit${UNITN}/main_e3_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_e3_test.go"
        cp "unit${UNITN}/main_internal_e3_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_internal_e3_test.go"
    else
        cp "unit${UNITN}/main_cycle_default_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_cycle_default_test.go"
    fi
   
    if [[ $1 > 1 ]]; then
       
        cp "unit${UNITN}/main_timeout_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_timeout_test.go"
    fi

    if [[ $1 < 3 ]]; then
        cp "unit${UNITN}/main_process_noctx_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_process_noctx_test.go"
    else 
        cp "unit${UNITN}/main_process_ctx_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_process_ctx_test.go"
    fi
fi

cd ..

go mod init course || true
go get github.com/stretchr/testify/assert

CGO_ENABLED=0 go test "./unit${UNITN}/exercises/e$1/..."
