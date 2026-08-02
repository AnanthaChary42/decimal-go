package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_DivToInt ports assertions from test/modules/divToInt.js
// JS: T.assertEqual(expected, new Decimal(dividend).divToInt(divisor).valueOf());
func TestOriginal_DivToInt(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 100,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  0,
		ToExpPos:  0,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		dividend string
		divisor  string
		expected string
	}{
		// Lines 18-21: basic integer division
		{"0", "1", "0e+0"},
		{"1", "3", "0e+0"},
		{"5", "2", "2e+0"},
		{"5", "3", "1e+0"},

		// Lines 22-30: precision divisions
		{"36895644873019582.63", "7.718100793996604914489691733070282", "4.780404643292329e+15"},
		{"687735627911862004719398564058596925.0620935155879", "661665211998060549450363044996.55374284667971784770628959", "1.039401e+6"},
		{"36341626756628760.847811686627543419", "5675392302633600.96", "6e+0"},
	}

	for _, tt := range tests {
		t.Run(tt.dividend+"_divToInt_"+tt.divisor, func(t *testing.T) {
			d, err := ctx.New(tt.dividend)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.dividend, err)
			}
			y, err := ctx.New(tt.divisor)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.divisor, err)
			}
			got := d.DivToInt(y).ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).DivToInt(%q).ValueOf() = %q, want %q", tt.dividend, tt.divisor, got, tt.expected)
			}
		})
	}
}
