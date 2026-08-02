package port_test

import (
	"strconv"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_ToDP ports assertions from test/modules/toDP.js
// JS: T.assertEqual(expected, new Decimal(n).toDP(dp, rm).valueOf());
func TestOriginal_ToDP(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfCeil, // rounding 7 in decimal.js is RoundHalfCeil
		ToExpNeg:  -9000000000000000,
		ToExpPos:  300,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		hasDP    bool
		dp       int
		rm       *decimal.RoundingMode
		expected string
	}{
		// Lines 24-35: special values (dp is omitted, meaning round to Precision sig digits)
		{"0", false, 0, nil, "0"},
		{"-1", false, 0, nil, "-1"},
		{"9.999e+9000000000000000", false, 0, nil, "9.999e+9000000000000000"},
		{"-9.999e+9000000000000000", false, 0, nil, "-9.999e+9000000000000000"},
		{"1e-9000000000000000", false, 0, nil, "1e-9000000000000000"},
		{"-1e-9000000000000000", false, 0, nil, "-1e-9000000000000000"},
		{"Infinity", false, 0, nil, "Infinity"},
		{"-Infinity", false, 0, nil, "-Infinity"},
		{"NaN", false, 0, nil, "NaN"},

		// Lines 36-56: rounding to 0 dp (dp is explicitly 0)
		{"0.5", true, 0, nil, "1"},
		{"0.7", true, 0, nil, "1"},
		{"1", true, 0, nil, "1"},
		{"1.1", true, 0, nil, "1"},
		{"1.49999", true, 0, nil, "1"},
		{"-0.5", true, 0, nil, "-0"},
		{"-0.500001", true, 0, nil, "-1"},
		{"-0.7", true, 0, nil, "-1"},
		{"-1", true, 0, nil, "-1"},
		{"-1.1", true, 0, nil, "-1"},
		{"-1.49999", true, 0, nil, "-1"},
		{"-1.5", true, 0, nil, "-1"},

		{"0.4", true, 0, nil, "0"},
		{"-0.4", true, 0, nil, "-0"},
		{"0.6", true, 0, nil, "1"},
		{"-0.6", true, 0, nil, "-1"},
		{"1.5", true, 0, nil, "2"},
		{"1.6", true, 0, nil, "2"},
		{"-1.6", true, 0, nil, "-2"},

		// Lines 58-60: zeros
		{"0", true, 0, nil, "0"},
		{"-0", true, 1, nil, "-0"},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			var got string
			if tt.hasDP {
				if tt.rm != nil {
					got = d.ToDP(tt.dp, *tt.rm).ValueOf()
				} else {
					got = d.ToDP(tt.dp).ValueOf()
				}
			} else {
				got = d.ToSD(ctx.Precision).ValueOf()
			}
			if got != tt.expected {
				t.Errorf("Decimal(%q).ToDP(%d) = %q, want %q", tt.input, tt.dp, got, tt.expected)
			}
		})
	}
}
