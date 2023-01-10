package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mul(t *testing.T) {
	r, err := mul(9999999999, 999999999)
	assert.Error(t, err)
	assert.Zero(t, r)
	assert.ErrorIs(t, err, ErrIntOverflow)
}

func Test_Div(t *testing.T) {
	r, err := div(10, 0)
	assert.Error(t, err)
	assert.Zero(t, r)
	assert.ErrorIs(t, err, ErrDivideByZero)
}

func Test_Pow(t *testing.T) {
	r, err := pow(999, 999)
	assert.Error(t, err)
	assert.Zero(t, r)
	assert.ErrorIs(t, err, ErrIntOverflow)
}
