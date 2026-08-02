package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_ValueOf ports all assertions from test/modules/valueOf.js
// JS: T.assertEqual(expected, new Decimal(n).valueOf());
// Config section 1: toExpNeg: -9e15, toExpPos: 9e15 (very wide — no exponential)
func TestOriginal_ValueOf(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Lines 18-32: positive values
		{"0", "0"},
		{"NaN", "NaN"},
		{"Infinity", "Infinity"},
		{"1", "1"},
		{"9", "9"},
		{"90", "90"},
		{"90.12", "90.12"},
		{"0.1", "0.1"},
		{"0.01", "0.01"},
		{"0.0123", "0.0123"},
		{"111111111111111111111", "111111111111111111111"},
		{"0.00001", "0.00001"},

		// Lines 34-46: negative values
		{"-0", "-0"},
		{"-Infinity", "-Infinity"},
		{"-1", "-1"},
		{"-9", "-9"},
		{"-90", "-90"},
		{"-90.12", "-90.12"},
		{"-0.1", "-0.1"},
		{"-0.01", "-0.01"},
		{"-0.0123", "-0.0123"},
		{"-111111111111111111111", "-111111111111111111111"},
		{"-0.00001", "-0.00001"},
	}

	// Use context with very wide toExpNeg/toExpPos
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	for _, tt := range tests {
		t.Run("default_"+tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).ValueOf() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestOriginal_ValueOf_ExpFormat ports the exponential format section of valueOf.js
// JS: Decimal.toExpNeg = Decimal.toExpPos = 0;
func TestOriginal_ValueOf_ExpFormat(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  0,
		ToExpPos:  0,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		expected string
	}{
		// Lines 51-56: small floats in exp notation
		{"0.0000001", "1e-7"},
		{"0.000000123", "1.23e-7"},
		{"0.000000012", "1.2e-8"},
		{"-0.0000001", "-1e-7"},
		{"-0.000000123", "-1.23e-7"},
		{"-0.000000012", "-1.2e-8"},

		// Lines 58-63: large numbers in exp notation
		{"573447902457635.174479825134", "5.73447902457635174479825134e+14"},
		{"10.7688", "1.07688e+1"},
		{"3171194102379077141557759899.307946350455841", "3.171194102379077141557759899307946350455841e+27"},
		{"49243534668981911776986533197425948906.34579", "4.924353466898191177698653319742594890634579e+37"},
		{"6855582439265693973.28633907445409866949445343654692955", "6.85558243926569397328633907445409866949445343654692955e+18"},
		{"1", "1e+0"},
	}

	for _, tt := range tests {
		t.Run("exp_"+tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).ValueOf() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
