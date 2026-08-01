package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

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
		{"-123456.7891011", "123456.7891011"},
		{"99", "99"},
		{"-99", "99"},
		{"999.999", "999.999"},
		{"-999.999", "999.999"},
		{"-0.001", "0.001"},
		{"Infinity", "Infinity"},
		{"-Infinity", "Infinity"},
		{"NaN", "NaN"},
		{"-NaN", "NaN"},
		{"11.121", "11.121"},
		{"-0.023842", "0.023842"},
		{"-1.19", "1.19"},
		{"-0.00000000009622", "9.622e-11"},
		{"-0.000000000509", "5.09e-10"},
		{"-1e20", "100000000000000000000"},
		{"1e20", "100000000000000000000"},
		{"-0.000001", "0.000001"},
		{"0.000001", "0.000001"},
		{"-1.234567e+10", "12345670000"},
		{"-999999999999", "999999999999"},
		{"-0.000000000001", "1e-12"},
		{"-1234.5678", "1234.5678"},
		{"-987654321", "987654321"},
		{"-0.0000001", "1e-7"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := decimal.New(tt.input)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.input, err)
			}
			got := d.Abs().String()
			if got != tt.expected {
				t.Errorf("New(%q).Abs().String() = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestOriginal_Cmp(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
		isNaN    bool
	}{
		{"1", "0", 1, false},
		{"1", "-0", 1, false},
		{"-1", "0", -1, false},
		{"-1", "-0", -1, false},
		{"0", "1", -1, false},
		{"0", "-1", 1, false},
		{"-0", "1", -1, false},
		{"-0", "-1", 1, false},
		{"0", "0", 0, false},
		{"0", "-0", 0, false},
		{"-0", "0", 0, false},
		{"-0", "-0", 0, false},
		{"0", "0.1", -1, false},
		{"0", "-0.1", 1, false},
		{"-0", "0.1", -1, false},
		{"-0", "-0.1", 1, false},
		{"0.1", "0", 1, false},
		{"0.1", "-0", 1, false},
		{"-0.1", "0", -1, false},
		{"-0.1", "-0", -1, false},
		{"NaN", "1", 0, true},
		{"NaN", "-1", 0, true},
		{"NaN", "0", 0, true},
		{"NaN", "NaN", 0, true},
		{"1", "NaN", 0, true},
		{"Infinity", "1", 1, false},
		{"Infinity", "-1", 1, false},
		{"-Infinity", "1", -1, false},
		{"-Infinity", "-1", -1, false},
		{"Infinity", "Infinity", 0, false},
		{"-Infinity", "-Infinity", 0, false},
		{"Infinity", "-Infinity", 1, false},
		{"-Infinity", "Infinity", -1, false},
		{"100", "200", -1, false},
		{"200", "100", 1, false},
		{"100", "100", 0, false},
		{"-100", "-200", 1, false},
		{"-200", "-100", -1, false},
	}

	for _, tt := range tests {
		name := tt.a + "_vs_" + tt.b
		t.Run(name, func(t *testing.T) {
			da, err := decimal.New(tt.a)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.a, err)
			}
			db, err := decimal.New(tt.b)
			if err != nil {
				t.Fatalf("New(%q) error: %v", tt.b, err)
			}
			got, ok := da.Cmp(db)
			if tt.isNaN {
				if ok {
					t.Errorf("Cmp(%q, %q) returned ok=true, want false for NaN", tt.a, tt.b)
				}
			} else {
				if !ok {
					t.Errorf("Cmp(%q, %q) returned ok=false", tt.a, tt.b)
				}
				if got != tt.expected {
					t.Errorf("Cmp(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.expected)
				}
			}
		})
	}
}

func TestOriginal_Plus(t *testing.T) {
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
		{"123.456", "789.101", "912.557"},
		{"-123.456", "-789.101", "-912.557"},
		{"0.000001", "0.000002", "0.000003"},
		{"1e10", "1e10", "20000000000"},
		{"1e-5", "2e-5", "0.00003"},
		{"-1e-5", "1e-5", "0"},
		{"1.23e5", "4.56e5", "579000"},
		{"-1.23e5", "-4.56e5", "-579000"},
		{"9999999", "1", "10000000"},
		{"0.9999999", "0.0000001", "1"},
		{"-0.9999999", "-0.0000001", "-1"},
		{"0", "-0", "0"},
		{"-0", "0", "0"},
		{"-0", "-0", "0"},
		{"Infinity", "1", "Infinity"},
		{"1", "Infinity", "Infinity"},
		{"-Infinity", "-1", "-Infinity"},
		{"-Infinity", "Infinity", "NaN"},
		{"Infinity", "-Infinity", "NaN"},
		{"NaN", "1", "NaN"},
		{"1", "NaN", "NaN"},
		{"NaN", "NaN", "NaN"},
	}

	for _, tt := range tests {
		name := tt.a + "+" + tt.b
		t.Run(name, func(t *testing.T) {
			da, _ := decimal.New(tt.a)
			db, _ := decimal.New(tt.b)
			res := da.Plus(db)
			if res.String() != tt.expected {
				t.Errorf("%s + %s = %q, want %q", tt.a, tt.b, res.String(), tt.expected)
			}
		})
	}
}

func TestOriginal_Minus(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"3", "2", "1"},
		{"0", "0", "0"},
		{"1", "1", "0"},
		{"0.3", "0.1", "0.2"},
		{"100", "200", "-100"},
		{"-5", "-3", "-2"},
		{"123.456", "23.456", "100"},
		{"-123.456", "-23.456", "-100"},
		{"0.003", "0.001", "0.002"},
		{"1e10", "5e9", "5000000000"},
		{"1e-5", "5e-6", "0.000005"},
		{"10000000", "1", "9999999"},
		{"1", "0.0000001", "0.9999999"},
		{"-1", "-0.0000001", "-0.9999999"},
		{"0", "-0", "0"},
		{"-0", "0", "0"},
		{"-0", "-0", "0"},
		{"Infinity", "1", "Infinity"},
		{"1", "Infinity", "-Infinity"},
		{"-Infinity", "-1", "-Infinity"},
		{"Infinity", "Infinity", "NaN"},
		{"-Infinity", "-Infinity", "NaN"},
		{"NaN", "1", "NaN"},
		{"1", "NaN", "NaN"},
		{"NaN", "NaN", "NaN"},
		{"5000", "2000", "3000"},
		{"0.5", "0.25", "0.25"},
		{"-0.5", "-0.25", "-0.25"},
		{"10", "0.1", "9.9"},
		{"0.1", "10", "-9.9"},
		{"999", "999", "0"},
		{"-999", "-999", "0"},
	}

	for _, tt := range tests {
		name := tt.a + "-" + tt.b
		t.Run(name, func(t *testing.T) {
			da, _ := decimal.New(tt.a)
			db, _ := decimal.New(tt.b)
			res := da.Minus(db)
			if res.String() != tt.expected {
				t.Errorf("%s - %s = %q, want %q", tt.a, tt.b, res.String(), tt.expected)
			}
		})
	}
}

func TestOriginal_Times(t *testing.T) {
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
		{"12.34", "5.67", "69.9678"},
		{"-12.34", "5.67", "-69.9678"},
		{"-12.34", "-5.67", "69.9678"},
		{"0.001", "0.001", "0.000001"},
		{"1e5", "1e5", "10000000000"},
		{"1e-5", "1e-5", "1e-10"},
		{"9999999", "2", "19999998"},
		{"0", "0", "0"},
		{"-0", "0", "0"},
		{"-0", "-0", "0"},
		{"Infinity", "2", "Infinity"},
		{"Infinity", "-2", "-Infinity"},
		{"-Infinity", "-2", "Infinity"},
		{"Infinity", "0", "NaN"},
		{"-Infinity", "0", "NaN"},
		{"NaN", "2", "NaN"},
		{"2", "NaN", "NaN"},
		{"NaN", "NaN", "NaN"},
		{"1.5", "2", "3"},
		{"2.5", "4", "10"},
		{"-2.5", "4", "-10"},
		{"-2.5", "-4", "10"},
		{"0.5", "0.5", "0.25"},
		{"0.2", "0.2", "0.04"},
		{"0.05", "0.05", "0.0025"},
		{"1000", "0.001", "1"},
		{"1000000", "0.000001", "1"},
		{"0.1", "100", "10"},
		{"0.01", "100", "1"},
	}

	for _, tt := range tests {
		name := tt.a + "*" + tt.b
		t.Run(name, func(t *testing.T) {
			da, _ := decimal.New(tt.a)
			db, _ := decimal.New(tt.b)
			res := da.Times(db)
			if res.String() != tt.expected {
				t.Errorf("%s * %s = %q, want %q", tt.a, tt.b, res.String(), tt.expected)
			}
		})
	}
}

func TestOriginal_Div(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"6", "3", "2"},
		{"1", "2", "0.5"},
		{"100", "10", "10"},
		{"-6", "3", "-2"},
		{"-6", "-3", "2"},
		{"10", "3", "3.3333333333333333333"},
		{"1", "3", "0.33333333333333333333"},
		{"2", "3", "0.66666666666666666667"},
		{"100", "4", "25"},
		{"1", "8", "0.125"},
		{"1", "10", "0.1"},
		{"1", "100", "0.01"},
		{"1", "1000", "0.001"},
		{"5", "2", "2.5"},
		{"-5", "2", "-2.5"},
		{"-5", "-2", "2.5"},
		{"0", "5", "0"},
		{"-0", "5", "0"},
		{"5", "0", "Infinity"},
		{"-5", "0", "-Infinity"},
		{"0", "0", "NaN"},
		{"Infinity", "5", "Infinity"},
		{"-Infinity", "5", "-Infinity"},
		{"Infinity", "-5", "-Infinity"},
		{"-Infinity", "-5", "Infinity"},
		{"Infinity", "Infinity", "NaN"},
		{"-Infinity", "Infinity", "NaN"},
		{"NaN", "5", "NaN"},
		{"5", "NaN", "NaN"},
		{"NaN", "NaN", "NaN"},
	}

	for _, tt := range tests {
		name := tt.a + "/" + tt.b
		t.Run(name, func(t *testing.T) {
			da, _ := decimal.New(tt.a)
			db, _ := decimal.New(tt.b)
			res := da.Div(db)
			if res.String() != tt.expected {
				t.Errorf("%s / %s = %q, want %q", tt.a, tt.b, res.String(), tt.expected)
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
