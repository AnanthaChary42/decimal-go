package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Cbrt ports assertions from test/modules/cbrt.js
// JS: T.assertEqual(expected, new Decimal(n).cbrt().valueOf());
func TestOriginal_Cbrt(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 40,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		expected string
	}{
		// Lines 20-28: special & exact values
		{"NaN", "NaN"},
		{"0", "0"},
		{"-0", "-0"},
		{"Infinity", "Infinity"},
		{"-Infinity", "-Infinity"},
		{"1", "1"},
		{"-1", "-1"},
		{"8", "2"},
		{"-8", "-2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.Cbrt().ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Cbrt().ValueOf() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
