package decimal

import (
	"fmt"
	"testing"
)

// TestNewDecimalBasic tests basic parsing of decimal strings.
func TestNewDecimalBasic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"1", "1"},
		{"-1", "-1"},
		{"123", "123"},
		{"-123", "-123"},
		{"123.456", "123.456"},
		{"-123.456", "-123.456"},
		{"0.1", "0.1"},
		{"0.0001", "0.0001"},
		{"1e5", "100000"},
		{"1e-5", "0.00001"},
		{"1.23e+5", "123000"},
		{"1.23e-5", "0.0000123"},
		{"9999999", "9999999"},
		{"10000000", "10000000"},
		{"99999999999999", "99999999999999"},
		{"1e20", "100000000000000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) returned error: %v", tt.input, err)
			}
			got := d.String()
			if got != tt.expected {
				t.Errorf("New(%q).String() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestNewFromInt64 tests construction from int64 values.
func TestNewFromInt64(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{100, "100"},
		{-100, "-100"},
		{9999999, "9999999"},
		{10000000, "10000000"},
		{123456789012345, "123456789012345"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%d", tt.input)
		t.Run(name, func(t *testing.T) {
			d := NewFromInt64(tt.input)
			got := d.String()
			if got != tt.expected {
				t.Errorf("NewFromInt64(%d).String() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestPredicates tests IsFinite, IsNaN, IsZero, etc.
func TestPredicates(t *testing.T) {
	zero, _ := New("0")
	one, _ := New("1")
	negOne, _ := New("-1")
	inf, _ := New("Infinity")
	negInfD, _ := New("-Infinity")
	nan, _ := New("NaN")

	// IsFinite
	if !zero.IsFinite() {
		t.Error("0 should be finite")
	}
	if !one.IsFinite() {
		t.Error("1 should be finite")
	}
	if inf.IsFinite() {
		t.Error("Infinity should not be finite")
	}
	if nan.IsFinite() {
		t.Error("NaN should not be finite")
	}

	// IsZero
	if !zero.IsZero() {
		t.Error("0 should be zero")
	}
	if one.IsZero() {
		t.Error("1 should not be zero")
	}

	// IsNaN
	if !nan.IsNaN() {
		t.Error("NaN should be NaN")
	}
	if one.IsNaN() {
		t.Error("1 should not be NaN")
	}

	// IsNeg
	if !negOne.IsNeg() {
		t.Error("-1 should be negative")
	}
	if one.IsNeg() {
		t.Error("1 should not be negative")
	}
	if !negInfD.IsNeg() {
		t.Error("-Infinity should be negative")
	}
}

// TestComparison tests Cmp, Eq, Gt, Lt etc.
func TestComparison(t *testing.T) {
	a, _ := New("1")
	b, _ := New("2")
	c, _ := New("1")

	cmp, ok := a.Cmp(b)
	if !ok || cmp != -1 {
		t.Errorf("1.Cmp(2) = %d, %v; want -1, true", cmp, ok)
	}

	cmp, ok = b.Cmp(a)
	if !ok || cmp != 1 {
		t.Errorf("2.Cmp(1) = %d, %v; want 1, true", cmp, ok)
	}

	cmp, ok = a.Cmp(c)
	if !ok || cmp != 0 {
		t.Errorf("1.Cmp(1) = %d, %v; want 0, true", cmp, ok)
	}

	if !a.Eq(c) {
		t.Error("1.Eq(1) should be true")
	}
	if a.Eq(b) {
		t.Error("1.Eq(2) should be false")
	}
	if !a.Lt(b) {
		t.Error("1.Lt(2) should be true")
	}
	if !b.Gt(a) {
		t.Error("2.Gt(1) should be true")
	}
}

// TestAddition tests Plus/Add operations.
func TestAddition(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"1", "2", "3"},
		{"0", "0", "0"},
		{"0.1", "0.2", "0.3"},
		{"1", "-1", "0"},
		{"999999999999", "1", "1000000000000"},
		{"-5", "3", "-2"},
		{"100", "200", "300"},
		{"0.001", "0.002", "0.003"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s+%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, err := New(tt.a)
			if err != nil {
				t.Fatal(err)
			}
			b, err := New(tt.b)
			if err != nil {
				t.Fatal(err)
			}
			result := a.Plus(b)
			got := result.String()
			if got != tt.expected {
				t.Errorf("%s + %s = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// TestSubtraction tests Minus/Sub operations.
func TestSubtraction(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"3", "2", "1"},
		{"0", "0", "0"},
		{"1", "1", "0"},
		{"0.3", "0.1", "0.2"},
		{"100", "200", "-100"},
		{"-5", "-3", "-2"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s-%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, err := New(tt.a)
			if err != nil {
				t.Fatal(err)
			}
			b, err := New(tt.b)
			if err != nil {
				t.Fatal(err)
			}
			result := a.Minus(b)
			got := result.String()
			if got != tt.expected {
				t.Errorf("%s - %s = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// TestMultiplication tests Times/Mul operations.
func TestMultiplication(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"2", "3", "6"},
		{"0", "1000", "0"},
		{"1", "1", "1"},
		{"-2", "3", "-6"},
		{"-2", "-3", "6"},
		{"100", "100", "10000"},
		{"0.1", "0.1", "0.01"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s*%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, err := New(tt.a)
			if err != nil {
				t.Fatal(err)
			}
			b, err := New(tt.b)
			if err != nil {
				t.Fatal(err)
			}
			result := a.Times(b)
			got := result.String()
			if got != tt.expected {
				t.Errorf("%s * %s = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// TestDivision tests Div operations.
func TestDivision(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"6", "3", "2"},
		{"1", "2", "0.5"},
		{"100", "10", "10"},
		{"-6", "3", "-2"},
		{"-6", "-3", "2"},
		{"10", "3", "3.3333333333333333333"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s/%s", tt.a, tt.b)
		t.Run(name, func(t *testing.T) {
			a, err := New(tt.a)
			if err != nil {
				t.Fatal(err)
			}
			b, err := New(tt.b)
			if err != nil {
				t.Fatal(err)
			}
			result := a.Div(b)
			got := result.String()
			if got != tt.expected {
				t.Errorf("%s / %s = %q, want %q", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// TestSpecialValues tests NaN and Infinity handling.
func TestSpecialValues(t *testing.T) {
	inf, _ := New("Infinity")
	negInfD, _ := New("-Infinity")
	nan, _ := New("NaN")

	if inf.String() != "Infinity" {
		t.Errorf("Infinity.String() = %q", inf.String())
	}
	if negInfD.String() != "-Infinity" {
		t.Errorf("-Infinity.String() = %q", negInfD.String())
	}
	if nan.String() != "NaN" {
		t.Errorf("NaN.String() = %q", nan.String())
	}
}

// TestRounding tests Ceil, Floor, Round, Trunc.
func TestRounding(t *testing.T) {
	tests := []struct {
		input                          string
		ceil, floor, round, trunc      string
	}{
		{"1.5", "2", "1", "2", "1"},
		{"-1.5", "-1", "-2", "-2", "-1"},
		{"2.0", "2", "2", "2", "2"},
		{"0.7", "1", "0", "1", "0"},
		{"-0.7", "0", "-1", "-1", "0"},
	}

	for _, tt := range tests {
		t.Run("Ceil_"+tt.input, func(t *testing.T) {
			d, _ := New(tt.input)
			got := d.Ceil().String()
			if got != tt.ceil {
				t.Errorf("Ceil(%s) = %q, want %q", tt.input, got, tt.ceil)
			}
		})
		t.Run("Floor_"+tt.input, func(t *testing.T) {
			d, _ := New(tt.input)
			got := d.Floor().String()
			if got != tt.floor {
				t.Errorf("Floor(%s) = %q, want %q", tt.input, got, tt.floor)
			}
		})
		t.Run("Trunc_"+tt.input, func(t *testing.T) {
			d, _ := New(tt.input)
			got := d.Trunc().String()
			if got != tt.trunc {
				t.Errorf("Trunc(%s) = %q, want %q", tt.input, got, tt.trunc)
			}
		})
	}
}
