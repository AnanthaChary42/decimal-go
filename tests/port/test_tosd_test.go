package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_ToSD ports assertions from test/modules/toSD.js
// JS: T.assertEqual(expected, new Decimal(n).toSD(sd, rm).valueOf());
func TestOriginal_ToSD(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfEven, // rounding 7 is RoundHalfEven
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		sd       int
		rm       *decimal.RoundingMode
		expected string
	}{
		// Lines 24-30: default sd / special values
		{"0", 0, nil, "0"},
		{"0.5", 0, nil, "0.5"},
		{"1", 0, nil, "1"},
		{"-1", 0, nil, "-1"},
		{"Infinity", 0, nil, "Infinity"},
		{"-Infinity", 0, nil, "-Infinity"},
		{"NaN", 0, nil, "NaN"},

		// Lines 32-38: zeros with explicit sd
		{"0", 1, nil, "0"},
		{"-0", 1, nil, "-0"},

		// Lines 40-41: simple sd
		{"12.345", 0, nil, "12.345"},

		// Lines 46-60: rounding modes with precision 5
		{"123456789.12345678912346789", 5, &[]decimal.RoundingMode{decimal.RoundUp}[0], "123460000"},
		{"123456789.12345678912346789", 5, &[]decimal.RoundingMode{decimal.RoundDown}[0], "123450000"},
		{"123456789.12345678912346789", 5, &[]decimal.RoundingMode{decimal.RoundCeil}[0], "123460000"},
		{"123456789.12345678912346789", 5, &[]decimal.RoundingMode{decimal.RoundFloor}[0], "123450000"},
		{"123456789.12345678912346789", 5, &[]decimal.RoundingMode{decimal.RoundHalfUp}[0], "123460000"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			var got string
			if tt.rm != nil {
				got = d.ToSD(tt.sd, *tt.rm).ValueOf()
			} else {
				got = d.ToSD(tt.sd).ValueOf()
			}
			if got != tt.expected {
				t.Errorf("Decimal(%q).ToSD(%d) = %q, want %q", tt.input, tt.sd, got, tt.expected)
			}
		})
	}
}
