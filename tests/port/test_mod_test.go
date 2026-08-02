package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Mod ports assertions from test/modules/mod.js
// JS: T.assertEqual(expected, new Decimal(a).mod(b).valueOf());
func TestOriginal_Mod(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 400,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		a, b     string
		expected string
	}{
		// Lines 19-24: Division by zero / NaN
		{"1", "0", "NaN"},
		{"1", "-0", "NaN"},
		{"-1", "0", "NaN"},
		{"-1", "-0", "NaN"},
		{"1", "NaN", "NaN"},
		{"-1", "NaN", "NaN"},

		// Lines 25-28: Infinity divisors
		{"1", "Infinity", "1"},
		{"1", "-Infinity", "1"},
		{"-1", "Infinity", "-1"},
		{"-1", "-Infinity", "-1"},

		// Lines 29-32: Zero dividends
		{"0", "1", "0"},
		{"0", "-1", "0"},
		{"-0", "1", "-0"},
		{"-0", "-1", "-0"},

		// Lines 33-38: Zero mod Zero
		{"0", "0", "NaN"},
		{"0", "-0", "NaN"},
		{"-0", "0", "NaN"},
		{"-0", "-0", "NaN"},
		{"0", "NaN", "NaN"},
		{"-0", "NaN", "NaN"},

		// Lines 39-42: Zero mod Infinity
		{"0", "Infinity", "0"},
		{"0", "-Infinity", "0"},
		{"-0", "Infinity", "-0"},
		{"-0", "-Infinity", "-0"},

		// Lines 43-60: NaN / Infinity dividends
		{"NaN", "1", "NaN"},
		{"NaN", "-1", "NaN"},
		{"Infinity", "1", "NaN"},
		{"-Infinity", "1", "NaN"},
		{"Infinity", "Infinity", "NaN"},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_mod_"+tt.b, func(t *testing.T) {
			a, _ := ctx.New(tt.a)
			b, _ := ctx.New(tt.b)
			got := a.Mod(b).ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Mod(%q).ValueOf() = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}
