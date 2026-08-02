package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_Ln(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input     string
		precision int
		rounding  decimal.RoundingMode
		want      string
	}{
		{"0", 40, decimal.RoundHalfUp, "-Infinity"},
		{"-0", 40, decimal.RoundHalfUp, "-Infinity"},
		{"1", 40, decimal.RoundHalfUp, "0"},
		{"-Infinity", 40, decimal.RoundHalfUp, "NaN"},
		{"Infinity", 40, decimal.RoundHalfUp, "Infinity"},
		{"NaN", 40, decimal.RoundHalfUp, "NaN"},
		{"2.7182818284590452353602874713526624977572", 39, decimal.RoundHalfUp, "1"},
		{"91247532.65728", 19, decimal.RoundDown, "18.3290865106890306"},
		{"130749452.494110812", 13, decimal.RoundHalfEven, "18.68879347348"},
		{"0.00000000122", 3, decimal.RoundDown, "-20.5"},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, tt.rounding
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := x.Ln().ValueOf(); got != tt.want {
			t.Errorf("Decimal(%s).Ln() = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestOriginal_Log(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input, base string
		precision   int
		rounding    decimal.RoundingMode
		want        string
	}{
		{"0", "10", 40, decimal.RoundHalfUp, "-Infinity"},
		{"-1", "10", 40, decimal.RoundHalfUp, "NaN"},
		{"Infinity", "10", 40, decimal.RoundHalfUp, "Infinity"},
		{"10", "0", 40, decimal.RoundHalfUp, "NaN"},
		{"10", "1", 40, decimal.RoundHalfUp, "NaN"},
		{"243", "9", 1, decimal.RoundHalfUp, "3"},
		{"512", "16", 2, decimal.RoundHalfDown, "2.2"},
		{"94143178827", "3486784401", 3, decimal.RoundUp, "1.15"},
		{"7625597484987", "59049", 2, decimal.RoundHalfDown, "2.7"},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, tt.rounding
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := x.Log(tt.base).ValueOf(); got != tt.want {
			t.Errorf("Decimal(%s).Log(%s) = %s, want %s", tt.input, tt.base, got, tt.want)
		}
	}
}

func TestOriginal_Log2AndLog10(t *testing.T) {
	ctx := missingModuleContext()
	ctx.ToExpNeg, ctx.ToExpPos = 0, 0
	for _, tt := range []struct {
		input     string
		precision int
		rounding  decimal.RoundingMode
		method    func(*decimal.Decimal) *decimal.Decimal
		want      string
	}{
		{"0", 40, decimal.RoundHalfUp, (*decimal.Decimal).Log2, "-Infinity"},
		{"1024", 2, decimal.RoundUp, (*decimal.Decimal).Log2, "1e+1"},
		{"7.47572e73", 2, decimal.RoundFloor, (*decimal.Decimal).Log2, "2.4e+2"},
		{"5e-2", 10, decimal.RoundCeil, (*decimal.Decimal).Log2, "-4.321928094e+0"},
		{"1", 40, decimal.RoundHalfUp, (*decimal.Decimal).Log10, "0e+0"},
		{"1000", 40, decimal.RoundHalfUp, (*decimal.Decimal).Log10, "3e+0"},
		{"1e-4", 4, decimal.RoundHalfUp, (*decimal.Decimal).Log10, "-4e+0"},
		{"7.47572e73", 2, decimal.RoundFloor, (*decimal.Decimal).Log10, "7.3e+1"},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, tt.rounding
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := tt.method(x).ValueOf(); got != tt.want {
			t.Errorf("log of %s = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestOriginal_Exp(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input     string
		precision int
		rounding  decimal.RoundingMode
		want      string
	}{
		{"0", 40, decimal.RoundHalfUp, "1"},
		{"-0", 40, decimal.RoundHalfUp, "1"},
		{"Infinity", 40, decimal.RoundHalfUp, "Infinity"},
		{"-Infinity", 40, decimal.RoundHalfUp, "0"},
		{"NaN", 40, decimal.RoundHalfUp, "NaN"},
		{"1", 40, decimal.RoundHalfUp, "2.718281828459045235360287471352662497757"},
		{"0.000000000000000004", 40, decimal.RoundCeil, "1.000000000000000004000000000000000008001"},
		{"-0.0000000000000000006", 5, decimal.RoundFloor, "0.99999"},
		{"20.72326583694641116", 20, decimal.RoundDown, "1000000000.0000000038"},
		{"-27.6310211159285483", 3, decimal.RoundDown, "0.000000000000999"},
		{"2.08E+16", 10, decimal.RoundDown, "Infinity"},
		{"-2.08E+16", 10, decimal.RoundDown, "0"},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, tt.rounding
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := x.Exp().ValueOf(); got != tt.want {
			t.Errorf("Decimal(%s).Exp() = %s, want %s", tt.input, got, tt.want)
		}
	}
}
