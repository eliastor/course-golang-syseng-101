#!/bin/bash
UNITN="4"
rm -f "../unit${UNITN}/exercises/e$1/main_test.go"
cp "unit${UNITN}/e$1_main_test.go.tpl" "../unit${UNITN}/exercises/e$1/main_test.go"

cd ..

go mod init course || true

CGO_ENABLED=0 go test "./unit${UNITN}/exercises/e$1/..."
