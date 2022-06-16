package unit1

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExercise0(t *testing.T) {
	codepath := "../../unit1/exercises/e0/main.go"
	var callVal string
	t.Run("Syntax", newTestExercise0_Syntax(codepath, &callVal))
	t.Run("Output", newTestExercise0_Output(codepath, &callVal))
}

func newTestExercise0_Syntax(codepath string, callValue *string) func(t *testing.T) {

	return func(t *testing.T) {
		sqrtDefinedinMainScope := false
		sqrtDefined := false
		sqrtArgIsInt := false
		sqrtArgIsFloat := false
		sqrtArgStringValue := ""

		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, codepath, nil, parser.AllErrors)
		if err != nil {
			t.Fatal("can't open file for parsing:", err)
		}

		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			// t.Log(n, fmt.Sprintf("%T", n))
			funcDecl, ok := n.(*ast.FuncDecl)
			if ok {
				if funcDecl.Name.Name == "Sqrt" {
					sqrtDefined = true
					return true
				}
			}
			callExpr, ok := n.(*ast.CallExpr)
			if ok {
				ident, ok := (callExpr.Fun).(*ast.Ident)
				if ok {
					if ident.String() == "Sqrt" {
						if !assert.Equal(t, 1, len(callExpr.Args), "wrong number of arguments in Sqrt call") {
							return false
						}
						arg := callExpr.Args[0]
						basicLit, ok := arg.(*ast.BasicLit)
						if ok {
							switch basicLit.Kind {
							case token.INT:
								sqrtArgIsInt = true
								sqrtArgStringValue = basicLit.Value
							case token.FLOAT:
								sqrtArgIsFloat = true
								sqrtArgStringValue = basicLit.Value
							default:
								assert.Fail(t, "Argument for Sqrt call must be a number")
								return false
							}
						}
					}
				}

			}

			return true
		})

		for _, obj := range f.Scope.Objects {
			if obj.Kind == ast.Fun && obj.Name == "Sqrt" {
				sqrtDefinedinMainScope = true
			}
		}

		assert.True(t, sqrtDefined, "Sqrt function is not defined")
		assert.True(t, sqrtDefinedinMainScope, "Sqrt function is not defined in main package")
		assert.True(t, sqrtArgIsFloat || sqrtArgIsInt, "Sqrt argument argument in Sqrt call must be a number")
		assert.NotEmpty(t, sqrtArgStringValue, "Sqart call argument must not be empty")

		*callValue = sqrtArgStringValue
	}
}
func newTestExercise0_Output(codepath string, callValue *string) func(t *testing.T) {
	return func(t *testing.T) {
		if !assert.NotNil(t, callValue, "cellValue must not be empty") {
			assert.FailNow(t, "test error: cellValue is empty")
		}
		t.Log("Trying to build and run file:", codepath)

		cmd := exec.Command("timeout", "10", "go", "run", codepath)
		tmpDir, err := os.MkdirTemp("", "gotest-*")
		if !assert.NoError(t, err) {
			assert.FailNow(t, "test error: cannot create temporary dir")
		}
		defer os.RemoveAll(tmpDir)

		cmd.Env = []string{
			"PATH=" + os.Getenv("PATH"),
			"CGO_ENABLED=0",
			"GOCACHE=" + tmpDir,
		}

		outBuf := bytes.NewBuffer(nil)
		errBuf := bytes.NewBuffer(nil)
		cmd.Stdout = outBuf
		cmd.Stderr = errBuf

		err = cmd.Run()
		if !assert.NoError(t, err, errBuf.String()) {
			assert.FailNow(t, "test error: \"go run\" failed", err)
			return
		}

		assert.NoError(t, err, "Can't compile and run code")
		assert.Empty(t, errBuf.String(), "stderr of your code must be empty", errBuf.String())
		assert.NotEmpty(t, outBuf.String(), "stdout of your code must not be empty")

		arg, err := strconv.ParseFloat(*callValue, 64)
		if !assert.NoError(t, err) {
			assert.FailNow(t, "test error: cannot parse Sqrt argument value:", err)
		}

		epsilon := 0.0000000001

		prevVal := arg
		curVal := 0.0
		sqroot := math.Sqrt(arg)

		bfr := bufio.NewReader(outBuf)
		for {
			line, err := bfr.ReadString('\n')
			if err != nil && err != io.EOF {
				assert.FailNow(t, "test error: can't read output buffer", err)
				break
			}
			line = strings.TrimRight(line, "\r\n")
			if line == "" && err != io.EOF {
				continue
			}
			if err == io.EOF {
				break
			}
			curVal, err = strconv.ParseFloat(line, 64)
			if !assert.NoError(t, err) {
				assert.FailNow(t, "test error: cannot parse output line to number:", err)
			}
			assert.Greater(t, prevVal, curVal, "seems to be your code has error, previous iteration is less than following iteration", prevVal, curVal)
			assert.Greater(t, math.Abs(prevVal-curVal), epsilon, "seems to be your code stucked around one value", prevVal, curVal)
		}
		assert.Less(t, sqroot-curVal, epsilon, "seems to be your code has wrong computations, result value is far from real", sqroot, curVal)
	}
}
