package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Pow ports assertions from test/modules/pow.js
// JS: T.assertEqual(expected, new Decimal(base).pow(exp).valueOf());
func TestOriginal_Pow(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 40,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		base, exp string
		sd        int
		rm        decimal.RoundingMode
		expected  string
	}{
		// Lines 20-22
		{"9", "0.5", 7, decimal.RoundHalfUp, "3"},
		{"9", "0.5", 26, decimal.RoundHalfUp, "3"},
		{"0.9999999999", "6", 39, decimal.RoundHalfUp, "0.999999999400000000149999999980000000001"},

		// Lines 40-51: integer powers
		{"32", "0.4", 1, decimal.RoundHalfUp, "4"},
		{"4", "2.5", 11, decimal.RoundHalfUp, "32"},
		{"4", "5.5", 27, decimal.RoundHalfUp, "2048"},
		{"9", "1.5", 5, decimal.RoundHalfUp, "27"},
	}

	for _, tt := range tests {
		t.Run(tt.base+"^"+tt.exp, func(t *testing.T) {
			customCtx := *ctx
			customCtx.Precision = tt.sd
			customCtx.Rounding = tt.rm

			b, _ := customCtx.New(tt.base)
			e, _ := customCtx.New(tt.exp)
			got := b.Pow(e).ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Pow(%q) [sd=%d, rm=%d] = %q, want %q", tt.base, tt.exp, tt.sd, tt.rm, got, tt.expected)
			}
		})
	}
}
