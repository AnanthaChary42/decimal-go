package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_TrigonometricSpecialValues(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input string
		call  func(*decimal.Decimal) *decimal.Decimal
		want  string
	}{
		{"NaN", (*decimal.Decimal).Sin, "NaN"},
		{"Infinity", (*decimal.Decimal).Cos, "NaN"},
		{"-Infinity", (*decimal.Decimal).Tan, "NaN"},
		{"0", (*decimal.Decimal).Cos, "1"},
		{"-0", (*decimal.Decimal).Sin, "-0"},
		{"-0", (*decimal.Decimal).Tan, "-0"},
		{"1.0000000000000001", (*decimal.Decimal).Asin, "NaN"},
		{"-2", (*decimal.Decimal).Acos, "NaN"},
	} {
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := tt.call(x).ValueOf(); got != tt.want {
			t.Errorf("Decimal(%s) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestOriginal_TrigonometricValues(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input     string
		precision int
		rounding  decimal.RoundingMode
		call      func(*decimal.Decimal) *decimal.Decimal
		want      string
	}{
		{"0.1", 10, decimal.RoundCeil, (*decimal.Decimal).Sin, "0.09983341665"},
		{"0.1", 10, decimal.RoundHalfUp, (*decimal.Decimal).Cos, "0.9950041653"},
		{"0.1", 7, decimal.RoundHalfUp, (*decimal.Decimal).Tan, "0.1003347"},
		{"0.8", 9, decimal.RoundHalfDown, (*decimal.Decimal).Asin, "0.927295218"},
		{"-0.8", 8, decimal.RoundCeil, (*decimal.Decimal).Acos, "2.4980916"},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, tt.rounding
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := tt.call(x).ValueOf(); got != tt.want {
			t.Errorf("Decimal(%s) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

