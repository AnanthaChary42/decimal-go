package port_test

// ============================================================================
// tests/port/test_original_test.go
//
// Test cases ported directly from the decimal.js test suite.
// Source: https://github.com/MikeMcl/decimal.js  (test/modules/*.js)
//
// NO new test cases were invented — every assertion below is a direct
// translation of a `t(...)` call from the original JavaScript tests.
// ============================================================================

import (
	"fmt"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// ---------------------------------------------------------------------------
// toString.js — lines 18-52 (non-exponential format, precision 20, toExpNeg/Pos = ±9e15)
// ---------------------------------------------------------------------------
func TestOriginal_ToString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"NaN", "NaN"},
		{"Infinity", "Infinity"},
		{"1", "1"},
		{"9", "9"},
		{"90", "90"},
		{"90.12", "90.12"},
		{"0.1", "0.1"},
		{"0.01", "0.01"},
		{"0.0123", "0.0123"},
		{"111111111111111111111", "111111111111111111111"},
		{"1111111111111111111111", "1111111111111111111111"},
		{"11111111111111111111111", "11111111111111111111111"},
		{"0.00001", "0.00001"},
		{"0.000001", "0.000001"},
		// Negatives
		{"-Infinity", "-Infinity"},
		{"-1", "-1"},
		{"-9", "-9"},
		{"-90", "-90"},
		{"-90.12", "-90.12"},
		{"-0.1", "-0.1"},
		{"-0.01", "-0.01"},
		{"-0.0123", "-0.0123"},
		{"-111111111111111111111", "-111111111111111111111"},
		{"-1111111111111111111111", "-1111111111111111111111"},
		{"-11111111111111111111111", "-11111111111111111111111"},
		{"-0.00001", "-0.00001"},
		{"-0.000001", "-0.000001"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			d, err := decimal.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.String()
			if got != tt.expected {
				t.Errorf("New(%q).String() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// abs.js — lines 18-75 (precision 20, rounding 4, toExpNeg -7, toExpPos 21)
// ---------------------------------------------------------------------------
func TestOriginal_Abs(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"-0", "0"},
		{"1", "1"},
		{"-1", "1"},
		{"0.5", "0.5"},
		{"-0.5", "0.5"},
		{"0.1", "0.1"},
		{"-0.1", "0.1"},
		{"1.1", "1.1"},
		{"-1.1", "1.1"},
		{"1.5", "1.5"},
		{"-1.5", "1.5"},
		{"-1e-5", "0.00001"},
		{"-9e9", "9000000000"},
		{"123456.7891011", "123456.7891011"},
		{"99", "99"},
		{"-99", "99"},
		{"999.999", "999.999"},
		{"-999.999", "999.999"},
		{"-0.001", "0.001"},
		{"Infinity", "Infinity"},
		{"-Infinity", "Infinity"},
		{"NaN", "NaN"},
		{"11.121", "11.121"},
		{"-0.023842", "0.023842"},
		{"-1.19", "1.19"},
		{"3838.2", "3838.2"},
		{"127.0", "127"},
		{"4.23073", "4.23073"},
		{"-2.5469", "2.5469"},
		{"-29949", "29949"},
		{"-277.10", "277.1"},
		{"53.456", "53.456"},
		{"-100564", "100564"},
		{"-12431.9", "12431.9"},
		{"-97633.7", "97633.7"},
		{"220", "220"},
		{"18.720", "18.72"},
		{"-2817", "2817"},
		{"-44535", "44535"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := decimal.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.Abs().String()
			if got != tt.expected {
				t.Errorf("Abs(%s) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// cmp.js — lines 20-79 (precision 20, rounding 4)
// ---------------------------------------------------------------------------
func TestOriginal_Cmp(t *testing.T) {
	tests := []struct {
		a, b string
		want int // 1, -1, 0; use -99 for NaN
	}{
		{"1", "0", 1},
		{"-1", "0", -1},
		{"0", "1", -1},
		{"0", "-1", 1},
		{"0", "0", 0},
		{"0.1", "0", 1},
		{"-0.1", "0", -1},
		// NaN comparisons → NaN (our Cmp returns ok=false)
		{"NaN", "1", -99},
		{"NaN", "-1", -99},
		{"NaN", "0", -99},
		{"NaN", "NaN", -99},
		{"NaN", "Infinity", -99},
		{"NaN", "-Infinity", -99},
		{"1", "NaN", -99},
		{"-1", "NaN", -99},
		{"0", "NaN", -99},
		{"Infinity", "NaN", -99},
		{"-Infinity", "NaN", -99},
		// Infinity comparisons
		{"Infinity", "1", 1},
		{"Infinity", "-1", 1},
		{"-Infinity", "1", -1},
		{"-Infinity", "-1", -1},
		{"Infinity", "0", 1},
		{"-Infinity", "0", -1},
		{"Infinity", "Infinity", 0},
		{"Infinity", "-Infinity", 1},
		{"-Infinity", "Infinity", -1},
		{"-Infinity", "-Infinity", 0},
		{"Infinity", "123", 1},
		{"3", "-Infinity", 1},
		{"1", "Infinity", -1},
		{"1", "-Infinity", 1},
		{"-1", "Infinity", -1},
		{"-1", "-Infinity", 1},
		{"0", "Infinity", -1},
		{"0", "-Infinity", 1},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s_cmp_%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, _ := decimal.New(tt.a)
			b, _ := decimal.New(tt.b)
			cmp, ok := a.Cmp(b)
			if tt.want == -99 {
				// NaN comparison: ok should be false
				if ok {
					t.Errorf("Cmp(%s, %s) ok=true, want false (NaN)", tt.a, tt.b)
				}
			} else {
				if !ok {
					t.Errorf("Cmp(%s, %s) ok=false, want true", tt.a, tt.b)
				} else if cmp != tt.want {
					t.Errorf("Cmp(%s, %s) = %d, want %d", tt.a, tt.b, cmp, tt.want)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// plus.js — lines 18-31, 49-100 (precision 200, rounding 4)
// ---------------------------------------------------------------------------
func TestOriginal_Plus(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		// plus.js lines 18-31
		{"1", "0", "1"},
		{"-1", "0", "-1"},
		{"0", "1", "1"},
		{"0", "-1", "-1"},
		{"1", "NaN", "NaN"},
		{"-1", "NaN", "NaN"},
		{"1", "Infinity", "Infinity"},
		{"1", "-Infinity", "-Infinity"},
		{"-1", "Infinity", "Infinity"},
		{"-1", "-Infinity", "-Infinity"},

		// plus.js lines 86-100 (precision 200, rounding 4 — string operands)
		{"1", "0", "1"},
		{"1", "1", "2"},
		{"1", "-45", "-44"},
		{"1", "22", "23"},
		{"1", "0144", "145"},
		{"1", "6.1915", "7.1915"},
		{"1", "-1.02", "-0.02"},
		{"1", "0.09", "1.09"},
		{"1", "-0.0001", "0.9999"},
		{"1", "8e5", "800001"},
		{"1", "9E12", "9000000000001"},
		{"1", "1e-14", "1.00000000000001"},
		{"1", "3.345E-9", "1.000000003345"},
		{"1", "-345.43e+4", "-3454299"},
		{"1", "-94.12E+0", "-93.12"},

		// plus.js line 109-113
		{"-1", "1", "0"},
		{"-0.01", "0.01", "0"},
		{"54", "-54", "0"},
		{"9.99", "-9.99", "0"},
		{"0.0000023432495704937", "-0.0000023432495704937", "0"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s+%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, _ := decimal.New(tt.a)
			b, _ := decimal.New(tt.b)
			got := a.Plus(b).String()
			if got != tt.expected {
				t.Errorf("(%s).Plus(%s) = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// minus.js — lines 18-31, 52-57, 87-113 (precision 200, rounding 4)
// ---------------------------------------------------------------------------
func TestOriginal_Minus(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		// minus.js lines 18-31
		{"1", "0", "1"},
		{"-1", "0", "-1"},
		{"1", "NaN", "NaN"},
		{"-1", "NaN", "NaN"},
		{"1", "Infinity", "-Infinity"},
		{"1", "-Infinity", "Infinity"},
		{"-1", "Infinity", "-Infinity"},
		{"-1", "-Infinity", "Infinity"},
		{"0", "1", "-1"},
		{"0", "-1", "1"},

		// minus.js lines 52-57 (rounding 4)
		{"0", "0", "0"},
		{"1", "1", "0"},
		{"-1", "-1", "0"},

		// minus.js lines 87-113
		{"1", "0", "1"},
		{"1", "1", "0"},
		{"1", "-45", "46"},
		{"1", "22", "-21"},
		{"1", "0144", "-143"},
		{"1", "6.1915", "-5.1915"},
		{"1", "-1.02", "2.02"},
		{"1", "0.09", "0.91"},
		{"1", "-0.0001", "1.0001"},
		{"1", "8e5", "-799999"},
		{"1", "9E12", "-8999999999999"},
		{"1", "1e-14", "0.99999999999999"},
		{"1", "3.345E-9", "0.999999996655"},
		{"1", "-345.43e+4", "3454301"},
		{"1", "-94.12E+0", "95.12"},
		{"0", "0", "0"},
		{"0", "0.001", "-0.001"},
		{"0", "111.1111111110000", "-111.111111111"},
		{"-1", "1", "-2"},
		{"-0.01", "0.01", "-0.02"},
		{"54", "-54", "108"},
		{"9.99", "-9.99", "19.98"},
		{"0.0000023432495704937", "-0.0000023432495704937", "0.0000046864991409874"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s-%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, _ := decimal.New(tt.a)
			b, _ := decimal.New(tt.b)
			got := a.Minus(b).String()
			if got != tt.expected {
				t.Errorf("(%s).Minus(%s) = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// times.js — lines 52-79 (precision 300, rounding 4, toExpNeg -7, toExpPos 21)
// ---------------------------------------------------------------------------
func TestOriginal_Times(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		// times.js lines 52-63 (zero * ±1)
		{"1", "0", "0"},
		{"-1", "0", "0"},
		{"0", "1", "0"},
		{"0", "-1", "0"},
		{"0", "0", "0"},

		// times.js lines 65-79
		{"1", "1", "1"},
		{"1", "-45", "-45"},
		{"1", "22", "22"},
		{"1", "0144", "144"},
		{"1", "6.1915", "6.1915"},
		{"1", "-1.02", "-1.02"},
		{"1", "0.09", "0.09"},
		{"1", "-0.0001", "-0.0001"},
		{"1", "8e5", "800000"},
		{"1", "9E12", "9000000000000"},
		{"1", "-345.43e+4", "-3454300"},
		{"1", "-94.12E+0", "-94.12"},

		// times.js NaN/Infinity lines 18-50
		{"1", "NaN", "NaN"},
		{"-1", "NaN", "NaN"},
		{"1", "Infinity", "Infinity"},
		{"1", "-Infinity", "-Infinity"},
		{"-1", "Infinity", "-Infinity"},
		{"-1", "-Infinity", "Infinity"},
		{"NaN", "1", "NaN"},
		{"NaN", "-1", "NaN"},
		{"NaN", "0", "NaN"},
		{"NaN", "NaN", "NaN"},
		{"NaN", "Infinity", "NaN"},
		{"NaN", "-Infinity", "NaN"},
		{"Infinity", "1", "Infinity"},
		{"Infinity", "-1", "-Infinity"},
		{"-Infinity", "1", "-Infinity"},
		{"-Infinity", "-1", "Infinity"},
		{"Infinity", "Infinity", "Infinity"},
		{"Infinity", "-Infinity", "-Infinity"},
		{"-Infinity", "Infinity", "-Infinity"},
		{"-Infinity", "-Infinity", "Infinity"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s*%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, _ := decimal.New(tt.a)
			b, _ := decimal.New(tt.b)
			got := a.Times(b).String()
			if got != tt.expected {
				t.Errorf("(%s).Times(%s) = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// div.js — lines 36-80 (precision 40, rounding 4, toExpNeg -7, toExpPos 21)
// ---------------------------------------------------------------------------
func TestOriginal_Div(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		// div.js lines 36-68 (special values)
		{"1", "NaN", "NaN"},
		{"-1", "NaN", "NaN"},
		{"0", "0", "NaN"},
		{"0", "NaN", "NaN"},
		{"NaN", "1", "NaN"},
		{"NaN", "-1", "NaN"},
		{"NaN", "0", "NaN"},
		{"NaN", "NaN", "NaN"},
		{"NaN", "Infinity", "NaN"},
		{"NaN", "-Infinity", "NaN"},
		{"Infinity", "1", "Infinity"},
		{"Infinity", "-1", "-Infinity"},
		{"-Infinity", "1", "-Infinity"},
		{"-Infinity", "-1", "Infinity"},
		{"Infinity", "NaN", "NaN"},
		{"-Infinity", "NaN", "NaN"},
		{"Infinity", "Infinity", "NaN"},
		{"Infinity", "-Infinity", "NaN"},
		{"-Infinity", "Infinity", "NaN"},
		{"-Infinity", "-Infinity", "NaN"},

		// div.js lines 70-80 (actual division, precision 40)
		{"1", "1", "1"},
		{"1", "-45", "-0.02222222222222222222222222222222222222222"},
		{"1", "22", "0.04545454545454545454545454545454545454545"},
		{"1", "0144", "0.006944444444444444444444444444444444444444"},
		{"1", "6.1915", "0.1615117499798110312525236210934345473633"},
		{"1", "-1.02", "-0.9803921568627450980392156862745098039216"},
		{"1", "-0.0001", "-10000"},
		{"1", "8e5", "0.00000125"},
		{"1", "1e-14", "100000000000000"},
		{"1", "-94.12E+0", "-0.01062473438164045898852528686782830429239"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s/%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, _ := decimal.New(tt.a)
			b, _ := decimal.New(tt.b)
			got := a.Div(b).String()
			if got != tt.expected {
				t.Errorf("(%s).Div(%s) = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// isFiniteEtc.js (subset) — predicates
// ---------------------------------------------------------------------------
func TestOriginal_Predicates(t *testing.T) {
	// isFinite
	for _, s := range []string{"1", "-1", "0.1", "1e999", "0"} {
		d, _ := decimal.New(s)
		if !d.IsFinite() {
			t.Errorf("IsFinite(%s) = false, want true", s)
		}
	}
	for _, s := range []string{"NaN", "Infinity", "-Infinity"} {
		d, _ := decimal.New(s)
		if d.IsFinite() {
			t.Errorf("IsFinite(%s) = true, want false", s)
		}
	}

	// isNaN
	nanD, _ := decimal.New("NaN")
	if !nanD.IsNaN() {
		t.Error("IsNaN(NaN) = false")
	}
	oneD, _ := decimal.New("1")
	if oneD.IsNaN() {
		t.Error("IsNaN(1) = true")
	}

	// isZero
	zeroD, _ := decimal.New("0")
	if !zeroD.IsZero() {
		t.Error("IsZero(0) = false")
	}
	if oneD.IsZero() {
		t.Error("IsZero(1) = true")
	}

	// isNeg / isPos
	negD, _ := decimal.New("-1")
	if !negD.IsNeg() {
		t.Error("IsNeg(-1) = false")
	}
	if oneD.IsNeg() {
		t.Error("IsNeg(1) = true")
	}
	if !oneD.IsPos() {
		t.Error("IsPos(1) = false")
	}
}

// ---------------------------------------------------------------------------
// Decimal.js (constructor) — lines 73-95 (string input, multi-word digits)
// ---------------------------------------------------------------------------
func TestOriginal_Constructor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"9", "9"},
		{"-99", "-99"},
		{"999", "999"},
		{"-9999", "-9999"},
		{"99999", "99999"},
		{"-999999", "-999999"},
		{"9999999", "9999999"},
		{"-99999999", "-99999999"},
		{"999999999", "999999999"},
		{"-9999999999", "-9999999999"},
		{"99999999999", "99999999999"},
		{"-999999999999", "-999999999999"},
		{"9999999999999", "9999999999999"},
		{"-99999999999999", "-99999999999999"},
		{"999999999999999", "999999999999999"},
		{"-9999999999999999", "-9999999999999999"},
		{"99999999999999999", "99999999999999999"},
		{"-999999999999999999", "-999999999999999999"},
		{"9999999999999999999", "9999999999999999999"},
		{"-99999999999999999999", "-99999999999999999999"},
		{"999999999999999999999", "999999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := decimal.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.String()
			if got != tt.expected {
				t.Errorf("New(%q).String() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
