package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Times ports assertions from test/modules/times.js
// JS: T.assertEqual(expected, new Decimal(multiplicand).times(multiplier).valueOf());
func TestOriginal_Times(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 300,
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
		// Lines 18-50: special values (NaN, Infinity)
		{"1", "NaN", "NaN"},
		{"-1", "NaN", "NaN"},
		{"1", "Infinity", "Infinity"},
		{"1", "-Infinity", "-Infinity"},
		{"-1", "Infinity", "-Infinity"},
		{"-1", "-Infinity", "Infinity"},
		{"0", "NaN", "NaN"},
		{"-0", "NaN", "NaN"},
		{"0", "Infinity", "NaN"},
		{"0", "-Infinity", "NaN"},
		{"-0", "Infinity", "NaN"},
		{"-0", "-Infinity", "NaN"},
		{"NaN", "1", "NaN"},
		{"NaN", "-1", "NaN"},
		{"NaN", "0", "NaN"},
		{"NaN", "-0", "NaN"},
		{"NaN", "NaN", "NaN"},
		{"NaN", "Infinity", "NaN"},
		{"NaN", "-Infinity", "NaN"},
		{"Infinity", "1", "Infinity"},
		{"Infinity", "-1", "-Infinity"},
		{"-Infinity", "1", "-Infinity"},
		{"-Infinity", "-1", "Infinity"},
		{"Infinity", "0", "NaN"},
		{"Infinity", "-0", "NaN"},
		{"-Infinity", "0", "NaN"},
		{"-Infinity", "-0", "NaN"},
		{"Infinity", "NaN", "NaN"},
		{"-Infinity", "NaN", "NaN"},
		{"Infinity", "Infinity", "Infinity"},
		{"Infinity", "-Infinity", "-Infinity"},
		{"-Infinity", "Infinity", "-Infinity"},
		{"-Infinity", "-Infinity", "Infinity"},

		// Lines 52-60: zero multiplications & signs
		{"1", "0", "0"},
		{"1", "-0", "-0"},
		{"-1", "0", "-0"},
		{"-1", "-0", "0"},
		{"0", "1", "0"},
		{"0", "-1", "-0"},
		{"-0", "1", "-0"},
		{"-0", "-1", "0"},
		{"0", "0", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.a+"*"+tt.b, func(t *testing.T) {
			a, _ := ctx.New(tt.a)
			b, _ := ctx.New(tt.b)

			got := a.Times(b).ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Times(%q) = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}
