package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_ToPrecision ports assertions from test/modules/toPrecision.js
// JS: T.assertEqual(expected, new Decimal(n).toPrecision(sd, rm));
func TestOriginal_ToPrecision(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -7,
		ToExpPos:  40,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		sd       int
		expected string
	}{
		// Lines 24-38: exponential format precision
		{"1.2345e+27", 1, "1e+27"},
		{"1.2345e+27", 2, "1.2e+27"},
		{"1.2345e+27", 3, "1.23e+27"},
		{"1.2345e+27", 4, "1.235e+27"},
		{"1.2345e+27", 5, "1.2345e+27"},
		{"1.2345e+27", 6, "1.23450e+27"},
		{"1.2345e+27", 7, "1.234500e+27"},

		{"-1.2345e+27", 1, "-1e+27"},
		{"-1.2345e+27", 2, "-1.2e+27"},
		{"-1.2345e+27", 3, "-1.23e+27"},
		{"-1.2345e+27", 4, "-1.235e+27"},
		{"-1.2345e+27", 5, "-1.2345e+27"},
		{"-1.2345e+27", 6, "-1.23450e+27"},
		{"-1.2345e+27", 7, "-1.234500e+27"},

		// Lines 40-56: integers
		{"7", 1, "7"},
		{"7", 2, "7.0"},
		{"7", 3, "7.00"},

		{"-7", 1, "-7"},
		{"-7", 2, "-7.0"},
		{"-7", 3, "-7.00"},

		{"91", 1, "9e+1"},
		{"91", 2, "91"},
		{"91", 3, "91.0"},
		{"91", 4, "91.00"},

		{"-91", 1, "-9e+1"},
		{"-91", 2, "-91"},
		{"-91", 3, "-91.0"},
		{"-91", 4, "-91.00"},

		// Lines 58-60: floats
		{"91.1234", 1, "9e+1"},
		{"91.1234", 2, "91"},
		{"91.1234", 3, "91.1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.ToPrecision(tt.sd)
			if got != tt.expected {
				t.Errorf("Decimal(%q).ToPrecision(%d) = %q, want %q", tt.input, tt.sd, got, tt.expected)
			}
		})
	}
}
