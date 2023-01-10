package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionError_Mul(t *testing.T) {
	_, err := mul(9999999999, 9999999999)
	assert.IsType(t, &ExpressionError{}, err)
	assert.ErrorIs(t, err, ErrIntOverflow)
}

func TestExpressionError_Div(t *testing.T) {
	_, err := div(10, 0)
	assert.IsType(t, &ExpressionError{}, err)
	assert.ErrorIs(t, err, ErrDivideByZero)
}

func TestExpressionError_Pow(t *testing.T) {
	_, err := pow(999, 999)
	assert.IsType(t, &ExpressionError{}, err)
	assert.ErrorIs(t, err, ErrIntOverflow)
}
