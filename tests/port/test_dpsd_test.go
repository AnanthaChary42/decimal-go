package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_DpSd ports all assertions from test/modules/dpSd.js
// JS function: t(n, dp, sd, zs)
//   T.assertEqual(dp, new Decimal(n).dp());
//   T.assertEqual(sd, new Decimal(n).sd(zs));
func TestOriginal_DpSd(t *testing.T) {
	// Config: precision: 20, rounding: 4, toExpNeg: -7, toExpPos: 21

	type testCase struct {
		input      string
		expectedDp int
		expectedSd int
		dpIsNaN    bool // true if dp should be NaN (for NaN/Inf inputs)
		sdIsNaN    bool // true if sd should be NaN (for NaN/Inf inputs)
		zs         *bool // nil = not provided, true/false = value
	}

	boolPtr := func(v bool) *bool { return &v }

	tests := []testCase{
		// t(0, 0, 1)
		{"0", 0, 1, false, false, nil},
		// t(-0, 0, 1)
		{"-0", 0, 1, false, false, nil},
		// t(NaN, NaN, NaN)
		{"NaN", 0, 0, true, true, nil},
		// t(Infinity, NaN, NaN)
		{"Infinity", 0, 0, true, true, nil},
		// t(-Infinity, NaN, NaN)
		{"-Infinity", 0, 0, true, true, nil},
		// t(1, 0, 1)
		{"1", 0, 1, false, false, nil},
		// t(-1, 0, 1)
		{"-1", 0, 1, false, false, nil},

		// t(100, 0, 1)
		{"100", 0, 1, false, false, nil},
		// t(100, 0, 1, 0)  — zs = false (0 is falsy)
		{"100", 0, 1, false, false, boolPtr(false)},
		// t(100, 0, 1, false)
		{"100", 0, 1, false, false, boolPtr(false)},
		// t(100, 0, 3, 1)  — zs = true (1 is truthy)
		{"100", 0, 3, false, false, boolPtr(true)},
		// t(100, 0, 3, true)
		{"100", 0, 3, false, false, boolPtr(true)},

		// t('0.0012345689', 10, 8)
		{"0.0012345689", 10, 8, false, false, nil},
		// t('0.0012345689', 10, 8, 0)
		{"0.0012345689", 10, 8, false, false, boolPtr(false)},
		// t('0.0012345689', 10, 8, false)
		{"0.0012345689", 10, 8, false, false, boolPtr(false)},
		// t('0.0012345689', 10, 8, 1)
		{"0.0012345689", 10, 8, false, false, boolPtr(true)},
		// t('0.0012345689', 10, 8, true)
		{"0.0012345689", 10, 8, false, false, boolPtr(true)},

		// t('987654321000000.0012345689000001', 16, 31, 0)
		{"987654321000000.0012345689000001", 16, 31, false, false, boolPtr(false)},
		// t('987654321000000.0012345689000001', 16, 31, 1)
		{"987654321000000.0012345689000001", 16, 31, false, false, boolPtr(true)},

		// t('1e+123', 0, 1)
		{"1e+123", 0, 1, false, false, nil},
		// t('1e+123', 0, 124, 1)
		{"1e+123", 0, 124, false, false, boolPtr(true)},
		// t('1e-123', 123, 1)
		{"1e-123", 123, 1, false, false, nil},
		// t('1e-123', 123, 1, 1)
		{"1e-123", 123, 1, false, false, boolPtr(true)},

		// t('9.9999e+9000000000000000', 0, 5, false)
		{"9.9999e+9000000000000000", 0, 5, false, false, boolPtr(false)},
		// t('9.9999e+9000000000000000', 0, 9000000000000001, true)
		{"9.9999e+9000000000000000", 0, 9000000000000001, false, false, boolPtr(true)},
		// t('-9.9999e+9000000000000000', 0, 5, false)
		{"-9.9999e+9000000000000000", 0, 5, false, false, boolPtr(false)},
		// t('-9.9999e+9000000000000000', 0, 9000000000000001, true)
		{"-9.9999e+9000000000000000", 0, 9000000000000001, false, false, boolPtr(true)},

		// t('1e-9000000000000000', 9e15, 1, false)
		{"1e-9000000000000000", 9000000000000000, 1, false, false, boolPtr(false)},
		// t('1e-9000000000000000', 9e15, 1, true)
		{"1e-9000000000000000", 9000000000000000, 1, false, false, boolPtr(true)},
		// t('-1e-9000000000000000', 9e15, 1, false)
		{"-1e-9000000000000000", 9000000000000000, 1, false, false, boolPtr(false)},
		// t('-1e-9000000000000000', 9e15, 1, true)
		{"-1e-9000000000000000", 9000000000000000, 1, false, false, boolPtr(true)},

		// t('55325252050000000000000000000000.000000004534500000001', 21, 53)
		{"55325252050000000000000000000000.000000004534500000001", 21, 53, false, false, nil},
	}

	for i, tt := range tests {
		d, err := decimal.New(tt.input)
		if err != nil {
			t.Fatalf("test %d: New(%q) error: %v", i, tt.input, err)
		}

		// Check dp
		if tt.dpIsNaN {
			gotDp := d.Dp()
			if gotDp != 0 {
				t.Errorf("test %d: Decimal(%q).Dp() = %d, want 0 (NaN equivalent)", i, tt.input, gotDp)
			}
		} else {
			gotDp := d.Dp()
			if gotDp != tt.expectedDp {
				t.Errorf("test %d: Decimal(%q).Dp() = %d, want %d", i, tt.input, gotDp, tt.expectedDp)
			}
		}

		// Check sd
		if tt.sdIsNaN {
			gotSd := d.Sd()
			if gotSd != 0 {
				t.Errorf("test %d: Decimal(%q).Sd() = %d, want 0 (NaN equivalent)", i, tt.input, gotSd)
			}
		} else {
			var gotSd int
			if tt.zs != nil {
				gotSd = d.Sd(*tt.zs)
			} else {
				gotSd = d.Sd()
			}
			if gotSd != tt.expectedSd {
				t.Errorf("test %d: Decimal(%q).Sd(%v) = %d, want %d", i, tt.input, tt.zs, gotSd, tt.expectedSd)
			}
		}
	}
}
