package main_test

import (
	"testing"
)

func TestEternity(t *testing.T) {
	t.Run("division by zero", func(t *testing.T) {
		onlyErrorTest(t, outBinPath,
			"div 10 0",
			"Computation error: error in expression 10 / 0: divide by zero",
		)
	})
}
