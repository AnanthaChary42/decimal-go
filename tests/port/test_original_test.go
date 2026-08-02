package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_Abs(t *testing.T) {
	// Config: precision: 20, rounding: 4, toExpNeg: -7, toExpPos: 21
	// JS uses valueOf() for output comparison
	tests := []struct {
		input    string
		expected string
	}{
		// Lines 18-32: basic values
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

		// Lines 34-45: misc values including Decimal constructor inputs
		{"-1e-5", "0.00001"},
		{"-9e9", "9000000000"},
		{"123456.7891011", "123456.7891011"},
		{"-123456.7891011", "123456.7891011"},
		{"99", "99"},
		{"-99", "99"},
		{"999.999", "999.999"},
		{"-999.999", "999.999"},
		{"-1", "1"},      // t('1', new Decimal(-1))
		{"-1", "1"},      // t('1', new Decimal('-1'))
		{"0.001", "0.001"},  // t('0.001', new Decimal(0.001))
		{"-0.001", "0.001"}, // t('0.001', new Decimal('-0.001'))

		// Lines 47-54: special values
		{"Infinity", "Infinity"},
		{"-Infinity", "Infinity"},
		{"NaN", "NaN"},
		{"-NaN", "NaN"},

		// Lines 56-75: more varied values
		{"11.121", "11.121"},
		{"-0.023842", "0.023842"},
		{"-1.19", "1.19"},
		{"-0.00000000009622", "9.622e-11"},
		{"-0.000000000509", "5.09e-10"},
		{"3838.2", "3838.2"},
		{"127.0", "127"},
		{"4.23073", "4.23073"},
		{"-2.5469", "2.5469"},
		{"-29949", "29949"},
		{"-277.10", "277.1"},
		{"-0.00000000000000497898", "4.97898e-15"},
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
			got := d.Abs().ValueOf()
			if got != tt.expected {
				t.Errorf("New(%q).Abs().ValueOf() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestOriginal_Abs_ExpFormat ports the exponential format section of abs.js
// JS: Decimal.toExpNeg = Decimal.toExpPos = 0;
func TestOriginal_Abs_ExpFormat(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  0,
		ToExpPos:  0,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		input    string
		expected string
	}{
		// Lines 79-115: exponential format tests
		{"-5.2452468128e+1", "5.2452468128e+1"},
		{"1.41525905257189365008396e+16", "1.41525905257189365008396e+16"},
		{"2.743068083928e+11", "2.743068083928e+11"},
		{"-1.52993064722314247378724599e+26", "1.52993064722314247378724599e+26"},
		{"3.7205576746e+10", "3.7205576746e+10"},
		{"-2.663e-10", "2.663e-10"},
		{"-1.26574209965030360615518e+17", "1.26574209965030360615518e+17"},
		{"1.052e+3", "1.052e+3"},
		{"-4.452945872502e+6", "4.452945872502e+6"},
		{"2.95732460816619226e+13", "2.95732460816619226e+13"},
		{"-1.1923100194288654481424e+18", "1.1923100194288654481424e+18"},
		{"8.99315449050893705e+6", "8.99315449050893705e+6"},
		{"5.200726538434486963e+8", "5.200726538434486963e+8"},
		{"1.182618278949368566264898065e+18", "1.182618278949368566264898065e+18"},
		{"-3.815873266712e-20", "3.815873266712e-20"},
		{"-1.316675370382742615e+6", "1.316675370382742615e+6"},
		{"-2.1032502e+6", "2.1032502e+6"},
		{"1.8e+1", "1.8e+1"},
		{"1.033525906631680944018544811261e-13", "1.033525906631680944018544811261e-13"},
		{"-1.102361746443461856816e+14", "1.102361746443461856816e+14"},
		{"8.595358491143959e+1", "8.595358491143959e+1"},
		{"1.226806049797304683867e-18", "1.226806049797304683867e-18"},
		{"-5e+0", "5e+0"},
		{"-1.091168788407093537887970016e+15", "1.091168788407093537887970016e+15"},
		{"3.87166413612272027e+12", "3.87166413612272027e+12"},
		{"1.411514e+5", "1.411514e+5"},
		{"1.0053454672509859631996e+22", "1.0053454672509859631996e+22"},
		{"6.9265714e+0", "6.9265714e+0"},
		{"1.04627709e+4", "1.04627709e+4"},
		{"2.285650225267766689304972e+5", "2.285650225267766689304972e+5"},
		{"4.5790517211306242e+7", "4.5790517211306242e+7"},
		{"-3.0033340092338313923473428e+16", "3.0033340092338313923473428e+16"},
		{"-2.83879929283797623e+1", "2.83879929283797623e+1"},
		{"4.5266377717178121183759377414e-5", "4.5266377717178121183759377414e-5"},
		{"-5.3781e+4", "5.3781e+4"},
		{"-6.722035208213298413522819127e-18", "6.722035208213298413522819127e-18"},
		{"-3.02865707828281230987116e+23", "3.02865707828281230987116e+23"},

		// Lines 117-124: extreme exponents
		{"1e-9000000000000000", "1e-9000000000000000"},
		{"-1e-9000000000000000", "1e-9000000000000000"},
		{"-9.9e-9000000000000001", "0e+0"},
		{"9.999999e+9000000000000000", "9.999999e+9000000000000000"},
		{"-9.999999e+9000000000000000", "9.999999e+9000000000000000"},
		{"1E9000000000000001", "Infinity"},
		{"-1e+9000000000000001", "Infinity"},
		{"-5.5879983320336874473209567979e+287894365", "5.5879983320336874473209567979e+287894365"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := ctx.New(tt.input)
			if err != nil {
				t.Fatalf("ctx.New(%q) error: %v", tt.input, err)
			}
			got := d.Abs().ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Abs().ValueOf() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}


func TestOriginal_Constructor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"1", "1"},
		{"-1", "-1"},
		{"123.456", "123.456"},
		{"-123.456", "-123.456"},
		{"0.000123", "0.000123"},
		{"-0.000123", "-0.000123"},
		{"1e5", "100000"},
		{"1e-5", "0.00001"},
		{"-1e5", "-100000"},
		{"-1e-5", "-0.00001"},
		{"1.2345e+6", "1234500"},
		{"-1.2345e-6", "-0.0000012345"},
		{"Infinity", "Infinity"},
		{"-Infinity", "-Infinity"},
		{"NaN", "NaN"},
		{"0x1a", "26"},
		{"0b101", "5"},
		{"0o77", "63"},
		{"999999999999", "999999999999"},
		{"0.000000000001", "1e-12"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := decimal.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			if d.String() != tt.expected {
				t.Errorf("New(%q).String() = %q, want %q", tt.input, d.String(), tt.expected)
			}
		})
	}
}

func TestOriginal_ToString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"1", "1"},
		{"-1", "-1"},
		{"123.456", "123.456"},
		{"-123.456", "-123.456"},
		{"0.000001", "0.000001"},
		{"1000000", "1000000"},
		{"1e20", "100000000000000000000"},
		{"-1e20", "-100000000000000000000"},
		{"1e-6", "0.000001"},
		{"-1e-6", "-0.000001"},
		{"Infinity", "Infinity"},
		{"-Infinity", "-Infinity"},
		{"NaN", "NaN"},
		{"0.5", "0.5"},
		{"-0.5", "-0.5"},
		{"99.99", "99.99"},
		{"-99.99", "-99.99"},
		{"0.0000001", "1e-7"},
		{"-0.0000001", "-1e-7"},
		{"1234567890", "1234567890"},
		{"-1234567890", "-1234567890"},
		{"0.123456789", "0.123456789"},
		{"-0.123456789", "-0.123456789"},
		{"100", "100"},
		{"-100", "-100"},
		{"0.01", "0.01"},
		{"-0.01", "-0.01"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := decimal.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			if d.String() != tt.expected {
				t.Errorf("New(%q).String() = %q, want %q", tt.input, d.String(), tt.expected)
			}
		})
	}
}

func TestOriginal_Predicates(t *testing.T) {
	zero, _ := decimal.New("0")
	one, _ := decimal.New("1")
	negOne, _ := decimal.New("-1")
	inf, _ := decimal.New("Infinity")
	nan, _ := decimal.New("NaN")

	if !zero.IsFinite() || !one.IsFinite() || !negOne.IsFinite() {
		t.Error("Finite check failed")
	}
	if inf.IsFinite() || nan.IsFinite() {
		t.Error("Non-finite check failed")
	}
}

func TestOriginal_Rounding(t *testing.T) {
	d1, _ := decimal.New("1.5")
	if d1.Ceil().String() != "2" || d1.Floor().String() != "1" {
		t.Error("Rounding 1.5 failed")
	}
	d2, _ := decimal.New("-1.5")
	if d2.Ceil().String() != "-1" || d2.Floor().String() != "-2" {
		t.Error("Rounding -1.5 failed")
	}
}
