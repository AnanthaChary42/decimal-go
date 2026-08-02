package port_test

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_IntPow(t *testing.T) {
	saveDefaultCtx(t)

	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(40),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-7),
		ToExpPos:  decimal.IntPtr(21),
		MinE:      decimal.IntPtr(-9000000000000000),
		MaxE:      decimal.IntPtr(9000000000000000),
	})

	ctx := decimal.GetDefaultContext()

	tPow := func(expected, n, exp interface{}) {
		t.Helper()
		base := newDec(t, n)
		exponent := newDec(t, exp)
		result := base.Pow(exponent)
		expStr := ""
		switch v := expected.(type) {
		case string:
			expStr = v
		case float64:
			if math.IsNaN(v) {
				expStr = "NaN"
			} else if math.IsInf(v, 1) {
				expStr = "Infinity"
			} else if math.IsInf(v, -1) {
				expStr = "-Infinity"
			} else {
				expStr = newDec(t, v).ValueOf()
			}
		case int:
			expStr = newDec(t, v).ValueOf()
		}
		if result.ValueOf() != expStr {
			t.Errorf("(%v)^(%v) = %s, want %s", n, exp, result.ValueOf(), expStr)
		}
	}

	t.Run("basic", func(t *testing.T) {
		tPow("4", 2, 2)
		tPow("2147483648", 2, 31)
		tPow("0.25", 2, -2)
		tPow("0.0625", 2, -4)
		tPow("1", 1, 100)
		tPow("0", 0, 1000)
	})

	t.Run("NaN_cases", func(t *testing.T) {
		nan := math.NaN()
		tPow("NaN", 2, nan)
		tPow("NaN", float64(0), nan)
		tPow("NaN", math.Copysign(0, -1), nan)
		tPow("NaN", math.Inf(1), nan)
		tPow("NaN", math.Inf(-1), nan)
		tPow("1", nan, float64(0))
		tPow("1", nan, math.Copysign(0, -1))
		tPow("NaN", nan, nan)
		tPow("NaN", nan, 2.2)
		tPow("NaN", nan, 1)
		tPow("NaN", nan, -1)
		tPow("NaN", nan, -2.2)
		tPow("NaN", nan, math.Inf(1))
		tPow("NaN", nan, math.Inf(-1))
	})

	t.Run("Infinity_cases", func(t *testing.T) {
		tPow("Infinity", 1.1, math.Inf(1))
		tPow("Infinity", -1.1, math.Inf(1))
		tPow("Infinity", 2, math.Inf(1))
		tPow("Infinity", -2, math.Inf(1))
		tPow("0", 2, math.Inf(-1))
		tPow("0", -2, math.Inf(-1))
		tPow("Infinity", 1.0/1.1, math.Inf(-1))
		tPow("Infinity", 1.0/-1.1, math.Inf(-1))
		tPow("Infinity", 0.5, math.Inf(-1))
		tPow("Infinity", -0.5, math.Inf(-1))
		tPow("NaN", 1, math.Inf(1))
		tPow("NaN", 1, math.Inf(-1))
		tPow("NaN", -1, math.Inf(1))
		tPow("NaN", -1, math.Inf(-1))
		tPow("0", 0.1, math.Inf(1))
		tPow("0", -0.1, math.Inf(1))
		tPow("0", 0.999, math.Inf(1))
		tPow("0", -0.999, math.Inf(1))
		tPow("Infinity", 0.1, math.Inf(-1))
		tPow("Infinity", -0.1, math.Inf(-1))
		tPow("Infinity", 0.999, math.Inf(-1))
		tPow("Infinity", -0.999, math.Inf(-1))
		tPow("Infinity", math.Inf(1), 2)
		// Original JavaScript case is (1 / Infinity)^-2, i.e. (+0)^-2.
		tPow("Infinity", float64(0), -2)
		tPow("-Infinity", math.Inf(-1), 3)
		tPow("-Infinity", math.Inf(-1), 13)
	})

	t.Run("neg_zero_pow", func(t *testing.T) {
		neg0 := math.Copysign(0, -1)
		tPow("-Infinity", neg0, -3)
		tPow("-Infinity", neg0, -13)
		tPow("-Infinity", neg0, -1)
		tPow("Infinity", neg0, -2)
		tPow("Infinity", float64(0), -2)
	})

	t.Run("NaN_neg_fractional_pow", func(t *testing.T) {
		tPow("NaN", -0.00001, 1.1)
		tPow("NaN", -0.00001, -1.1)
	})

	t.Run("large_exponents", func(t *testing.T) {
		ctx.Precision = 20
		tPow("1.9801312458591796501e+301030", 2, 1000001)
		tPow("5.0501702959901511235e-301031", 2, -1000001)
		tPow("-1.9801312458591796501e+301030", -2, 1000001)
		tPow("-5.0501702959901511235e-301031", -2, -1000001)
	})

	t.Run("precision_600", func(t *testing.T) {
		ctx.Precision = 600
		tPow("4096", "8", 4)
		tPow("-1.331", "-1.1", 3)
		tPow("5.125696", "-2.264", 2)
		tPow("1", "61818", 0)
		tPow("3.2", "3.2", 1)
		tPow("1280630.81718016", "5.8", 8)
		tPow("3965.318943552", "15.828", 3)
		tPow("16", "4", 2)
		tPow("1", "-1", 4)
		tPow("-8", "-2", 3)
		tPow("1", "-1", 8)
		tPow("16807", "7", 5)
		tPow("9", "3", 2)
		tPow("14641", "121", 2)
		tPow("390625", "-5", 8)
		tPow("64", "-2", 6)
		tPow("128", "2", 7)
		tPow("262144", "4", 9)
		tPow("16384", "4", 7)
		tPow("1", "1", 1)
		tPow("-15.625", "-2.5", 3)
		tPow("17179869184", "2", 34)
		tPow("4294967296", "2", 32)
		tPow("8589934592", "2", 33)
		tPow("-1", "-1", 41)
		tPow("1073741824", "8", 10)
		tPow("-134217728", "-8.00", 9)
		tPow("1220703125", "5.0", 13)
		tPow("131621703842267136", "-6", 22)
		tPow("22876792454961", "3", 28)
		tPow("-155568095557812224", "-14", 15)
		tPow("1", "1", 25)
	})

	t.Run("toExpNeg_toExpPos_0", func(t *testing.T) {
		ctx.ToExpNeg = 0
		ctx.ToExpPos = 0

		tPow("2e+0", 2, "1.0")
		tPow("1.6e+1", 2, "4.00000000")
		tPow("6.25e-2", 2, -4)
		tPow("-7e+0", "-7", 1)
		tPow("1e+0", "-429.32321", 0)
		tPow("8.1e+1", "-3", 4)
		tPow("1.296e+3", "-6", 4)
		tPow("2.9e+0", "2.9", 1)
		tPow("1.764e+3", "-42", 2)
		tPow("2.43e+2", "3", 5)
	})
}
