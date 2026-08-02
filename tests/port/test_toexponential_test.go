package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_ToExponential ports assertions from test/modules/toExponential.js
// JS: T.assertEqual(expected, new Decimal(n).toExponential(dp));
func TestOriginal_ToExponential(t *testing.T) {
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
		// Lines 22-52: positive numbers
		{"1", -1, "1e+0"},
		{"11", -1, "1.1e+1"},
		{"112", -1, "1.12e+2"},

		{"1", 0, "1e+0"},
		{"11", 0, "1e+1"},
		{"112", 0, "1e+2"},
		{"1", 1, "1.0e+0"},
		{"11", 1, "1.1e+1"},
		{"112", 1, "1.1e+2"},
		{"1", 2, "1.00e+0"},
		{"11", 2, "1.10e+1"},
		{"112", 2, "1.12e+2"},
		{"1", 3, "1.000e+0"},
		{"11", 3, "1.100e+1"},
		{"112", 3, "1.120e+2"},

		{"0.1", -1, "1e-1"},
		{"0.11", -1, "1.1e-1"},
		{"0.112", -1, "1.12e-1"},
		{"0.1", 0, "1e-1"},
		{"0.11", 0, "1e-1"},
		{"0.112", 0, "1e-1"},
		{"0.1", 1, "1.0e-1"},
		{"0.11", 1, "1.1e-1"},
		{"0.112", 1, "1.1e-1"},
		{"0.1", 2, "1.00e-1"},
		{"0.11", 2, "1.10e-1"},
		{"0.112", 2, "1.12e-1"},
		{"0.1", 3, "1.000e-1"},
		{"0.11", 3, "1.100e-1"},
		{"0.112", 3, "1.120e-1"},

		// Lines 54-60: negative numbers
		{"-1", -1, "-1e+0"},
		{"-11", -1, "-1.1e+1"},
		{"-112", -1, "-1.12e+2"},
		{"-1", 0, "-1e+0"},
		{"-11", 0, "-1e+1"},
		{"-112", 0, "-1e+2"},
		{"-1", 1, "-1.0e+0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			var got string
			if tt.dp < 0 {
				got = d.ToExponential(-1)
			} else {
				got = d.ToExponential(tt.dp)
			}
			if got != tt.expected {
				t.Errorf("Decimal(%q).ToExponential(%d) = %q, want %q", tt.input, tt.dp, got, tt.expected)
			}
		})
	}
}
