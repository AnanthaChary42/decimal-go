package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_ToFixed ports assertions from test/modules/toFixed.js
// JS: T.assertEqual(expected, new Decimal(n).toFixed(dp));
func TestOriginal_ToFixed(t *testing.T) {
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
		dp       int
		expected string
	}{
		// Lines 22-26
		{"99.9512986", 1, "100.0"},
		{"9.95036", 1, "10.0"},
		{"0.99", 1, "1.0"},
		{"0.09906", 2, "0.10"},
		{"0.0098034", 3, "0.010"},

		// Lines 28-30
		{"NaN", 2, "NaN"},
		{"Infinity", 2, "Infinity"},
		{"-Infinity", 2, "-Infinity"},

		// Lines 32-56
		{"1111111111111111111111", 8, "1111111111111111111111.00000000"},
		{"0.1", 1, "0.1"},
		{"0.1", 2, "0.10"},
		{"0.1", 3, "0.100"},
		{"0.01", 2, "0.01"},
		{"0.01", 3, "0.010"},
		{"0.01", 4, "0.0100"},
		{"0.001", 2, "0.00"},
		{"0.001", 3, "0.001"},
		{"0.001", 4, "0.0010"},
		{"1", 4, "1.0000"},
		{"1", 1, "1.0"},
		{"1", 0, "1"},
		{"12", 0, "12"},
		{"1.1", 0, "1"},
		{"12.1", 0, "12"},
		{"1.12", 0, "1"},
		{"12.12", 0, "12"},
		{"0.0000006", 7, "0.0000006"},
		{"0.00000006", 8, "0.00000006"},
		{"0.00000006", 9, "0.000000060"},
		{"0.00000006", 10, "0.0000000600"},
		{"0", 0, "0"},
		{"0", 1, "0.0"},
		{"0", 2, "0.00"},

		// Lines 58-60
		{"-1111111111111111111111", 8, "-1111111111111111111111.00000000"},
		{"-0.1", 1, "-0.1"},
		{"-0.10", 2, "-0.10"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.ToFixed(tt.dp)
			if got != tt.expected {
				t.Errorf("Decimal(%q).ToFixed(%d) = %q, want %q", tt.input, tt.dp, got, tt.expected)
			}
		})
	}
}
