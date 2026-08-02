package port_test

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_ToString(t *testing.T) {
	saveDefaultCtx(t)

	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(20),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-9000000000000000),
		ToExpPos:  decimal.IntPtr(9000000000000000),
		MinE:      decimal.IntPtr(-9000000000000000),
		MaxE:      decimal.IntPtr(9000000000000000),
	})

	// --- Non-exponential toString tests ---
	t.Run("non_exponential", func(t *testing.T) {
		type tc struct{ expected string; v interface{} }
		cases := []tc{
			{"0", float64(0)}, {"0", "0"},
			{"NaN", "NaN"}, {"Infinity", "Infinity"},
			{"1", float64(1)}, {"9", float64(9)}, {"90", float64(90)}, {"90.12", 90.12},
			{"0.1", 0.1}, {"0.01", 0.01}, {"0.0123", 0.0123},
			{"111111111111111111111", "111111111111111111111"},
			{"1111111111111111111111", "1111111111111111111111"},
			{"11111111111111111111111", "11111111111111111111111"},
			{"0.00001", 0.00001}, {"0.000001", 0.000001},
			{"0", math.Copysign(0, -1)}, {"0", "-0"},
			{"-Infinity", "-Infinity"},
			{"-1", float64(-1)}, {"-9", float64(-9)}, {"-90", float64(-90)},
			{"-90.12", -90.12}, {"-0.1", -0.1}, {"-0.01", -0.01},
			{"-0.0123", -0.0123},
			{"-111111111111111111111", "-111111111111111111111"},
			{"-1111111111111111111111", "-1111111111111111111111"},
			{"-11111111111111111111111", "-11111111111111111111111"},
			{"-0.00001", -0.00001}, {"-0.000001", -0.000001},
		}
		for _, c := range cases {
			d := newDec(t, c.v)
			got := d.ToString()
			if got != c.expected {
				t.Errorf("ToString(%v) = %q, want %q", c.v, got, c.expected)
			}
		}
	})

	// --- Exponential format tests ---
	t.Run("exponential", func(t *testing.T) {
		ctx := decimal.GetDefaultContext()
		ctx.ToExpNeg = 0
		ctx.ToExpPos = 0

		type tc struct{ expected, input string }
		cases := []tc{
			{"1e-7", "0.0000001"},
			{"1.2e-7", "0.00000012"},
			{"1.23e-7", "0.000000123"},
			{"1e-8", "0.00000001"},
			{"1.2e-8", "0.000000012"},
			{"1.23e-8", "0.0000000123"},
			{"-1e-7", "-0.0000001"},
			{"-1.2e-7", "-0.00000012"},
			{"-1.23e-7", "-0.000000123"},
			{"-1e-8", "-0.00000001"},
			{"-1.2e-8", "-0.000000012"},
			{"-1.23e-8", "-0.0000000123"},
			{"5.73447902457635174479825134e+14", "573447902457635.174479825134"},
			{"1.07688e+1", "10.7688"},
			{"3.171194102379077141557759899307946350455841e+27", "3171194102379077141557759899.307946350455841"},
			{"1e+0", "1"},
			{"2.1320000803e+7", "21320000.803"},
			{"5.0878741e+4", "50878.741"},
			{"3e+0", "3"},
			{"1.56e+0", "1.56"},
			{"3.3431e+0", "3.3431"},
			{"4.89e+0", "4.89"},
			{"4.9e+1", "49"},
			{"1.3e+1", "13"},
			{"7.8e+1", "78"},
			{"5.22e+0", "5.22"},
			{"2.45e+0", "2.45"},
			{"1.23e+0", "1.23"},
			{"8.91411e+0", "8.91411"},
			{"6.9e+1", "69"},
			{"1.8e+1", "18"},
			{"5.11e+1", "51.1"},
			{"1.21e+1", "12.1"},
			{"2.44e+0", "2.44"},
			{"1e+0", "1"},
			{"5.43e+1", "54.3"},
			{"1.7972e+1", "17.972"},
		}
		for _, c := range cases {
			d, err := decimal.New(c.input)
			if err != nil {
				t.Errorf("New(%q): %v", c.input, err)
				continue
			}
			if got := d.ToString(); got != c.expected {
				t.Errorf("ToString(%q) = %q, want %q", c.input, got, c.expected)
			}
		}
	})
}
