package port_test

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_ToNumber_PositiveZero ports the positive zero section of toNumber.js
// JS: T.assert(1 / new Decimal(n).toNumber() === Infinity) → result is +0
func TestOriginal_ToNumber_PositiveZero(t *testing.T) {
	positiveZeroInputs := []string{
		"0",
		"0.0",
		"0.000000000000",
		"0e+0",
		"0e-0",
		"1e-9000000000000000",
	}

	for _, input := range positiveZeroInputs {
		t.Run("posZero_"+input, func(t *testing.T) {
			d, err := decimal.New(input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", input, err)
			}
			got := d.ToNumber()
			// Check it is positive zero: 1/+0 === +Infinity
			if got != 0 || math.Signbit(got) {
				t.Errorf("Decimal(%q).ToNumber() = %v, want +0", input, got)
			}
		})
	}
}

// TestOriginal_ToNumber_NegativeZero ports the negative zero section of toNumber.js
// JS: T.assert(1 / new Decimal(n).toNumber() === -Infinity) → result is -0
func TestOriginal_ToNumber_NegativeZero(t *testing.T) {
	negativeZeroInputs := []string{
		"-0",
		"-0.0",
		"-0.000000000000",
		"-0e+0",
		"-0e-0",
		"-1e-9000000000000000",
	}

	for _, input := range negativeZeroInputs {
		t.Run("negZero_"+input, func(t *testing.T) {
			d, err := decimal.New(input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", input, err)
			}
			got := d.ToNumber()
			// Check it is negative zero: 1/-0 === -Infinity
			if got != 0 || !math.Signbit(got) {
				t.Errorf("Decimal(%q).ToNumber() = %v, want -0", input, got)
			}
		})
	}
}

// TestOriginal_ToNumber_Values ports the value comparison section of toNumber.js
// JS: T.assertEqual(expected, new Decimal(n).toNumber());
func TestOriginal_ToNumber_Values(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		// Lines 42-47: special values
		{"Infinity", math.Inf(1)},
		{"-Infinity", math.Inf(-1)},

		// Lines 49-53: positive 1
		{"1", 1},
		{"1.0", 1},
		{"1e+0", 1},
		{"1e-0", 1},

		// Lines 55-59: negative 1
		{"-1", -1},
		{"-1.0", -1},
		{"-1e+0", -1},
		{"-1e-0", -1},

		// Lines 61-62: general values
		{"123.456789876543", 123.456789876543},
		{"-123.456789876543", -123.456789876543},

		// Lines 64-65: small values
		{"1.1102230246251565e-16", 1.1102230246251565e-16},
		{"-1.1102230246251565e-16", -1.1102230246251565e-16},

		// Lines 67-68: MAX_SAFE_INTEGER
		{"9007199254740991", 9007199254740991},
		{"-9007199254740991", -9007199254740991},

		// Lines 70-71: float64 extremes
		{"5e-324", 5e-324},
		{"1.7976931348623157e+308", 1.7976931348623157e+308},

		// Lines 73-76: overflow/underflow
		{"9.999999e+9000000000000000", math.Inf(1)},
		{"-9.999999e+9000000000000000", math.Inf(-1)},
		{"1e-9000000000000000", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := decimal.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.ToNumber()
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Decimal(%q).ToNumber() = %v, want NaN", tt.input, got)
				}
			} else if got != tt.expected {
				t.Errorf("Decimal(%q).ToNumber() = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestOriginal_ToNumber_NaN tests that NaN input returns NaN
func TestOriginal_ToNumber_NaN(t *testing.T) {
	d, _ := decimal.New("NaN")
	got := d.ToNumber()
	if !math.IsNaN(got) {
		t.Errorf("Decimal('NaN').ToNumber() = %v, want NaN", got)
	}
}

// TestOriginal_ToNumber_NegZeroResult tests that '-1e-9000000000000000' → -0
func TestOriginal_ToNumber_NegZeroResult(t *testing.T) {
	d, _ := decimal.New("-1e-9000000000000000")
	got := d.ToNumber()
	if got != 0 || !math.Signbit(got) {
		t.Errorf("Decimal('-1e-9000000000000000').ToNumber() = %v, want -0", got)
	}
}
