package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Div ports assertions from test/modules/div.js
// JS: T.assertEqual(expected, new Decimal(dividend).div(divisor).valueOf());
func TestOriginal_Div(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 40,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		dividend, divisor string
		expected          string
	}{
		// Lines 36-60: special values (division by zero, Infinity, NaN)
		{"1", "0", "Infinity"},
		{"1", "-0", "-Infinity"},
		{"-1", "0", "-Infinity"},
		{"-1", "-0", "Infinity"},
		{"1", "NaN", "NaN"},
		{"-1", "NaN", "NaN"},
		{"0", "0", "NaN"},
		{"0", "-0", "NaN"},
		{"-0", "0", "NaN"},
		{"-0", "-0", "NaN"},
		{"0", "NaN", "NaN"},
		{"-0", "NaN", "NaN"},
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
		{"Infinity", "0", "Infinity"},
		{"Infinity", "-0", "-Infinity"},
	}

	for _, tt := range tests {
		t.Run(tt.dividend+"/"+tt.divisor, func(t *testing.T) {
			a, _ := ctx.New(tt.dividend)
			b, _ := ctx.New(tt.divisor)

			got := a.Div(b).ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Div(%q) = %q, want %q", tt.dividend, tt.divisor, got, tt.expected)
			}
		})
	}
}

// TestOriginal_Div_ZeroSigns tests section 1 of div.js checking negative zero result flags
func TestOriginal_Div_ZeroSigns(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 40,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		dividend, divisor string
		isNegZero         bool
	}{
		// Lines 19-30
		{"1", "Infinity", false},
		{"1", "-Infinity", true},
		{"-1", "Infinity", true},
		{"-1", "-Infinity", false},
		{"0", "1", false},
		{"0", "-1", true},
		{"-0", "1", true},
		{"-0", "-1", false},
		{"0", "Infinity", false},
		{"0", "-Infinity", true},
		{"-0", "Infinity", true},
		{"-0", "-Infinity", false},
	}

	for _, tt := range tests {
		t.Run(tt.dividend+"/"+tt.divisor, func(t *testing.T) {
			a, _ := ctx.New(tt.dividend)
			b, _ := ctx.New(tt.divisor)

			q := a.Div(b)
			got := q.IsZero() && q.IsNeg()
			if got != tt.isNegZero {
				t.Errorf("Decimal(%q).Div(%q) isNegZero = %v, want %v", tt.dividend, tt.divisor, got, tt.isNegZero)
			}
		})
	}
}
