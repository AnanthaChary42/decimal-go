package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Sqrt ports assertions from test/modules/sqrt.js
// JS: T.assertEqual(expected, new Decimal(n).sqrt().valueOf());
func TestOriginal_Sqrt(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		expected string
		sd       int
		rm       decimal.RoundingMode
	}{
		// Lines 20-29: special & exact values
		{"NaN", "NaN", 20, decimal.RoundHalfUp},
		{"4", "2", 20, decimal.RoundHalfUp},
		{"0.01", "0.1", 20, decimal.RoundHalfUp},
		{"0", "0", 20, decimal.RoundHalfUp},
		{"-0", "-0", 20, decimal.RoundHalfUp},
		{"Infinity", "Infinity", 20, decimal.RoundHalfUp},
		{"-Infinity", "NaN", 20, decimal.RoundHalfUp},
		{"-1", "NaN", 20, decimal.RoundHalfUp},
		{"-35.999", "NaN", 20, decimal.RoundHalfUp},
		{"-0.00000000000001", "NaN", 20, decimal.RoundHalfUp},

		// Lines 43-60: rounding modes test cases (precision = 2, rounding = 0 (RoundUp))
		{"101", "11", 2, decimal.RoundUp},
		{"111", "11", 2, decimal.RoundUp},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			customCtx := *ctx
			customCtx.Precision = tt.sd
			customCtx.Rounding = tt.rm
			d, err := customCtx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.Sqrt().ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Sqrt().ValueOf() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
