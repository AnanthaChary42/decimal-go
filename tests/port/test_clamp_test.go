package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Clamp ports all assertions from test/modules/clamp.js
func TestOriginal_Clamp(t *testing.T) {
	tests := []struct {
		x, min, max, expected string
	}{
		// t('-0', '0', '0', '-0')
		{"-0", "0", "0", "-0"},
		// t('-0', '-0', '0', '-0')
		{"-0", "-0", "0", "-0"},
		// t('-0', '0', '-0', '-0')
		{"-0", "0", "-0", "-0"},
		// t('-0', '-0', '-0', '-0')
		{"-0", "-0", "-0", "-0"},

		// t('0', '0', '0', '0')
		{"0", "0", "0", "0"},
		// t('0', '-0', '0', '0')
		{"0", "-0", "0", "0"},
		// t('0', '0', '-0', '0')
		{"0", "0", "-0", "0"},
		// t('0', '-0', '-0', '0')
		{"0", "-0", "-0", "0"},

		// t(0, 0, 1, '0')
		{"0", "0", "1", "0"},
		// t(-1, 0, 1, '0')
		{"-1", "0", "1", "0"},
		// t(-2, 0, 1, '0')
		{"-2", "0", "1", "0"},
		// t(1, 0, 1, '1')
		{"1", "0", "1", "1"},
		// t(2, 0, 1, '1')
		{"2", "0", "1", "1"},

		// t(1, 1, 1, '1')
		{"1", "1", "1", "1"},
		// t(-1, 1, 1, '1')
		{"-1", "1", "1", "1"},
		// t(-1, -1, 1, '-1')
		{"-1", "-1", "1", "-1"},
		// t(2, 1, 2, '2')
		{"2", "1", "2", "2"},
		// t(3, 1, 2, '2')
		{"3", "1", "2", "2"},
		// t(1, 0, 1, '1')
		{"1", "0", "1", "1"},
		// t(2, 0, 1, '1')
		{"2", "0", "1", "1"},

		// t(Infinity, 0, 1, '1')
		{"Infinity", "0", "1", "1"},
		// t(0, -Infinity, 0, '0')
		{"0", "-Infinity", "0", "0"},
		// t(-Infinity, 0, 1, '0')
		{"-Infinity", "0", "1", "0"},
		// t(-Infinity, -Infinity, Infinity, '-Infinity')
		{"-Infinity", "-Infinity", "Infinity", "-Infinity"},
		// t(Infinity, -Infinity, Infinity, 'Infinity')
		{"Infinity", "-Infinity", "Infinity", "Infinity"},
		// t(0, 1, Infinity, '1')
		{"0", "1", "Infinity", "1"},

		// t(0, NaN, 1, 'NaN')
		{"0", "NaN", "1", "NaN"},
		// t(0, 0, NaN, 'NaN')
		{"0", "0", "NaN", "NaN"},
		// t(NaN, 0, 1, 'NaN')
		{"NaN", "0", "1", "NaN"},
	}

	for _, tt := range tests {
		name := tt.x + "_clamp_" + tt.min + "_" + tt.max
		t.Run(name, func(t *testing.T) {
			x, _ := decimal.New(tt.x)
			min, _ := decimal.New(tt.min)
			max, _ := decimal.New(tt.max)
			got := x.Clamp(min, max).ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Clamp(%q, %q).ValueOf() = %q, want %q",
					tt.x, tt.min, tt.max, got, tt.expected)
			}
		})
	}
}
