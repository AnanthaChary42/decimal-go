package port_test

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_Sum(t *testing.T) {
	saveDefaultCtx(t)

	inf := math.Inf(1)
	ninf := math.Inf(-1)
	nan := math.NaN()

	// sum helper: builds []*Decimal from mixed args, calls decimal.Sum, asserts equal to expected.
	tSum := func(t *testing.T, expected *decimal.Decimal, args ...interface{}) {
		t.Helper()
		decs := make([]*decimal.Decimal, len(args))
		for i, v := range args {
			decs[i] = newDec(t, v)
		}
		result := decimal.Sum(decs...)
		assertEqualDecimal(t, expected, result, "Sum")
	}

	t.Run("sum_0", func(t *testing.T) {
		expected := newDec(t, 0)
		tSum(t, expected, "0")
		tSum(t, expected, "0", newDec(t, 0))
		tSum(t, expected, 1, 0, "-1")
		tSum(t, expected, 0, newDec(t, "-10"), 0, 0, 0, 0, 0, 10)
		tSum(t, expected, 11, -11)
		tSum(t, expected, 1, "2", newDec(t, 3), newDec(t, "4"), -10)
		tSum(t, expected, newDec(t, -10), "9", newDec(t, 0.01), 0.99)
	})

	t.Run("sum_10", func(t *testing.T) {
		expected := newDec(t, 10)
		tSum(t, expected, "10")
		tSum(t, expected, "0", newDec(t, "10"))
		tSum(t, expected, 10, 0)
		tSum(t, expected, 0, 0, 0, 0, 0, 0, 10)
		tSum(t, expected, 11, -1)
		tSum(t, expected, 1, "2", newDec(t, 3), newDec(t, "4"))
		tSum(t, expected, "9", newDec(t, 0.01), 0.99)
	})

	t.Run("sum_600", func(t *testing.T) {
		expected := newDec(t, 600)
		tSum(t, expected, 100, 200, 300)
		tSum(t, expected, "100", "200", "300")
		tSum(t, expected, newDec(t, 100), newDec(t, 200), newDec(t, 300))
		tSum(t, expected, 100, "200", newDec(t, 300))
		tSum(t, expected, 99.9, 200.05, 300.05)
	})

	t.Run("sum_NaN", func(t *testing.T) {
		expected := newDec(t, nan)
		tSum(t, expected, nan)
		tSum(t, expected, "1", nan)
		tSum(t, expected, 100, 200, nan)
		tSum(t, expected, nan, 0, "9", newDec(t, 0), 11, inf)
		tSum(t, expected, 0, newDec(t, "-Infinity"), "9", newDec(t, nan), 11)
		tSum(t, expected, 4, "-Infinity", 0, "9", newDec(t, 0), inf, 2)
	})

	t.Run("sum_Infinity", func(t *testing.T) {
		expected := newDec(t, inf)
		tSum(t, expected, inf)
		tSum(t, expected, 100, 200, "Infinity")
		tSum(t, expected, 0, newDec(t, "Infinity"), "9", newDec(t, 0), 11)
		tSum(t, expected, 0, "9", newDec(t, 0), 11, inf)
		tSum(t, expected, 4, newDec(t, inf), 0, "9", newDec(t, 0), inf, 2)
	})

	t.Run("sum_neg_Infinity", func(t *testing.T) {
		expected := newDec(t, ninf)
		tSum(t, expected, ninf)
		tSum(t, expected, 100, 200, "-Infinity")
		tSum(t, expected, 0, newDec(t, "-Infinity"), "9", newDec(t, 0), 11)
		tSum(t, expected, 0, "9", newDec(t, 0), 11, ninf)
		tSum(t, expected, 4, newDec(t, ninf), 0, "9", newDec(t, 0), ninf, 2)
	})
}
