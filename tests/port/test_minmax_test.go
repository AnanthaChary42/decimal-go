package port_test

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_MinAndMax(t *testing.T) {
	saveDefaultCtx(t)

	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(20),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-7),
		ToExpPos:  decimal.IntPtr(21),
		MinE:      decimal.IntPtr(-9000000000000000),
		MaxE:      decimal.IntPtr(9000000000000000),
	})

	type minMaxCase struct {
		minExpected, maxExpected interface{}
		arr                     []interface{}
	}

	tMinMax := func(t *testing.T, c minMaxCase) {
		t.Helper()
		args := make([]*decimal.Decimal, len(c.arr))
		for i, v := range c.arr {
			args[i] = newDec(t, v)
		}

		maxResult := decimal.Max(args...)
		minResult := decimal.Min(args...)

		expectedMax := newDec(t, c.maxExpected)
		expectedMin := newDec(t, c.minExpected)

		if !maxResult.Eq(expectedMax) && !(maxResult.IsNaN() && expectedMax.IsNaN()) {
			t.Errorf("Max: got %s, want %s", maxResult.ValueOf(), expectedMax.ValueOf())
		}
		if !minResult.Eq(expectedMin) && !(minResult.IsNaN() && expectedMin.IsNaN()) {
			t.Errorf("Min: got %s, want %s", minResult.ValueOf(), expectedMin.ValueOf())
		}
	}

	nan := math.NaN()
	inf := math.Inf(1)
	ninf := math.Inf(-1)
	neg0 := math.Copysign(0, -1)

	t.Run("NaN_cases", func(t *testing.T) {
		cases := []minMaxCase{
			{nan, nan, []interface{}{nan}},
			{nan, nan, []interface{}{-2, 0, -1, nan}},
			{nan, nan, []interface{}{-2, nan, 0, -1}},
			{nan, nan, []interface{}{nan, -2, 0, -1}},
			{nan, nan, []interface{}{-2, 0, -1, newDec(t, nan)}},
			{nan, nan, []interface{}{inf, -2, "NaN", 0, -1, ninf}},
			{nan, nan, []interface{}{"NaN", inf, -2, 0, -1, ninf}},
			{nan, nan, []interface{}{inf, -2, nan, 0, -1, ninf}},
		}
		for _, c := range cases {
			tMinMax(t, c)
		}
	})

	t.Run("zero_and_negzero", func(t *testing.T) {
		cases := []minMaxCase{
			{0, 0, []interface{}{0, 0, 0}},
			{neg0, 0, []interface{}{neg0, 0, 0}},
			{neg0, 0, []interface{}{0, neg0, 0}},
			{neg0, 0, []interface{}{0, 0, neg0}},
			{neg0, 1, []interface{}{1, 0, neg0}},
			{-2, 0, []interface{}{0, -1, neg0, -2}},
			{-2, inf, []interface{}{-2, -1, neg0, 0, inf}},
			{ninf, 0, []interface{}{-2, 0, -1, ninf}},
			{ninf, inf, []interface{}{ninf, -2, 0, -1, inf}},
			{ninf, inf, []interface{}{inf, -2, 0, -1, ninf}},
			{ninf, inf, []interface{}{ninf, -2, 0, newDec(t, inf)}},
		}
		for _, c := range cases {
			tMinMax(t, c)
		}
	})

	t.Run("sorted_arrays", func(t *testing.T) {
		cases := []minMaxCase{
			{-2, 0, []interface{}{-2, 0, -1}},
			{-2, 0, []interface{}{-2, -1, neg0, 0}},
			{-2, 0, []interface{}{0, -2, -1}},
			{-2, 0, []interface{}{0, -1, -2}},
			{-2, 0, []interface{}{-1, -2, 0}},
			{-2, 0, []interface{}{-1, 0, neg0, -2}},
			{-1, 1, []interface{}{-1, 0, 1}},
			{-1, 1, []interface{}{-1, 1, 0}},
			{-1, 1, []interface{}{0, -1, 1}},
			{-1, 1, []interface{}{0, 1, -1}},
			{-1, 1, []interface{}{1, -1, 0}},
			{-1, 1, []interface{}{1, 0, -1}},
			{0, 2, []interface{}{0, 1, 2}},
			{0, 2, []interface{}{0, 2, 1}},
			{0, 2, []interface{}{1, 0, 2}},
			{0, 2, []interface{}{1, 2, 0}},
			{0, 2, []interface{}{2, 1, 0}},
			{0, 2, []interface{}{2, 0, 1}},
		}
		for _, c := range cases {
			tMinMax(t, c)
		}
	})

	t.Run("mixed_types", func(t *testing.T) {
		cases := []minMaxCase{
			{-1, 1, []interface{}{"-1", 0, newDec(t, 1)}},
			{-1, 1, []interface{}{"-1", newDec(t, 1)}},
			{-1, 1, []interface{}{0, "-1", newDec(t, 1)}},
			{0, 1, []interface{}{0, newDec(t, 1)}},
			{1, 1, []interface{}{newDec(t, 1)}},
			{-1, -1, []interface{}{newDec(t, -1)}},
		}
		for _, c := range cases {
			tMinMax(t, c)
		}
	})

	t.Run("decimal_string_values", func(t *testing.T) {
		cases := []minMaxCase{
			{0.0009999, 0.0010001, []interface{}{0.001, 0.0009999, 0.0010001}},
			{-0.0010001, -0.0009999, []interface{}{-0.001, -0.0009999, -0.0010001}},
			{-0.000001, "999.001", []interface{}{2, neg0, "1e-9000000000000000", "324.32423423", -0.000001, "999.001", 10}},
			{"-9.99999e+9000000000000000", inf, []interface{}{10, "-9.99999e+9000000000000000", newDec(t, inf), "9.99999e+9000000000000000", 0}},
			{-3, 3, []interface{}{1, "2", 3, "-1", -2, "-3"}},
		}
		for _, c := range cases {
			tMinMax(t, c)
		}
	})
}
