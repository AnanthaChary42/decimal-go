package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Round ports assertions from test/modules/round.js
// JS: Decimal.rounding = rm; T.assertEqual(expected, new Decimal(n).round().valueOf());
func TestOriginal_Round(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		rm       decimal.RoundingMode
		expected string
	}{
		// Lines 19-32: basic values
		{"Infinity", decimal.RoundHalfUp, "Infinity"},
		{"-Infinity", decimal.RoundHalfUp, "-Infinity"},
		{"NaN", decimal.RoundHalfUp, "NaN"},
		{"0", decimal.RoundHalfUp, "0"},
		{"-0", decimal.RoundHalfUp, "-0"},
		{"1", decimal.RoundHalfUp, "1"},
		{"-1", decimal.RoundHalfUp, "-1"},
		{"0.1", decimal.RoundHalfUp, "0"},
		{"-0.1", decimal.RoundHalfUp, "-0"},
		{"1e-9000000000000000", decimal.RoundHalfUp, "0"},
		{"-1e-9000000000000000", decimal.RoundHalfUp, "-0"},
		{"9.999e+9000000000000000", decimal.RoundHalfUp, "9.999e+9000000000000000"},
		{"-9.999e+9000000000000000", decimal.RoundHalfUp, "-9.999e+9000000000000000"},

		// Lines 34-60: various numbers and rounding modes
		{"913844451463.2968716291911429329", decimal.RoundHalfUp, "913844451463"},
		{"28923471241.98034587766743552", decimal.RoundDown, "28923471241"},
		{"4521981137129744.67989385131579785422314397", decimal.RoundUp, "4521981137129745"},
		{"16604074202686787.6", decimal.RoundUp, "16604074202686788"},
		{"21333681.9", decimal.RoundUp, "21333682"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			// Replace context rounding mode
			customCtx := *ctx
			customCtx.Rounding = tt.rm
			d, _ = customCtx.New(tt.input)

			got := d.Round().ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Round() [rm=%d] = %q, want %q", tt.input, tt.rm, got, tt.expected)
			}
		})
	}
}
