package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_AdvancedTranscendentalSpecialValues(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input string
		call  func(*decimal.Decimal) *decimal.Decimal
		want  string
	}{
		{"NaN", (*decimal.Decimal).Atan, "NaN"},
		{"Infinity", (*decimal.Decimal).Atan, "1.5707963267948966192"},
		{"-Infinity", (*decimal.Decimal).Sinh, "-Infinity"},
		{"Infinity", (*decimal.Decimal).Cosh, "Infinity"},
		{"-Infinity", (*decimal.Decimal).Tanh, "-1"},
		{"-0", (*decimal.Decimal).Asinh, "-0"},
		{"0.9", (*decimal.Decimal).Acosh, "NaN"},
		{"1", (*decimal.Decimal).Acosh, "0"},
	} {
		ctx.Precision, ctx.Rounding = 20, decimal.RoundHalfUp
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := tt.call(x).ValueOf(); got != tt.want {
			t.Errorf("Decimal(%s) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestOriginal_AdvancedTranscendentalValues(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input     string
		precision int
		rounding  decimal.RoundingMode
		call      func(*decimal.Decimal) *decimal.Decimal
		want      string
	}{
		{"0.2", 30, decimal.RoundHalfDown, (*decimal.Decimal).Atan, "0.197395559849880758370049765195"},
		{"0.1", 7, decimal.RoundFloor, (*decimal.Decimal).Sinh, "0.1001667"},
		{"0.39", 10, decimal.RoundFloor, (*decimal.Decimal).Cosh, "1.077018834"},
		{"0.5", 8, decimal.RoundDown, (*decimal.Decimal).Tanh, "0.46211715"},
		{"14", 10, decimal.RoundHalfUp, (*decimal.Decimal).Asinh, "3.333477587"},
		{"76", 10, decimal.RoundUp, (*decimal.Decimal).Acosh, "5.023837236"},
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

func TestOriginal_Atan2(t *testing.T) {
	ctx := missingModuleContext()
	ctx.Precision, ctx.Rounding = 20, decimal.RoundHalfUp
	for _, tt := range []struct {
		y, x string
		want string
	}{
		{"0", "-0", "3.1415926535897932385"},
		{"-0", "-1", "-3.1415926535897932385"},
		{"1", "0", "1.5707963267948966192"},
		{"3", "4", "0.6435011087932843868"},
	} {
		y, err := ctx.New(tt.y)
		if err != nil {
			t.Fatal(err)
		}
		x, err := ctx.New(tt.x)
		if err != nil {
			t.Fatal(err)
		}
		if got := ctx.Atan2(y, x).ValueOf(); got != tt.want {
			t.Errorf("Atan2(%s, %s) = %s, want %s", tt.y, tt.x, got, tt.want)
		}
	}
}

func TestOriginal_Atanh(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		input     string
		precision int
		rounding  decimal.RoundingMode
		want      string
	}{
		{"NaN", 40, decimal.RoundHalfUp, "NaN"},
		{"Infinity", 40, decimal.RoundHalfUp, "NaN"},
		{"2", 40, decimal.RoundHalfUp, "NaN"},
		{"0", 40, decimal.RoundHalfUp, "0"},
		{"-0", 40, decimal.RoundHalfUp, "-0"},
		{"1", 40, decimal.RoundHalfUp, "Infinity"},
		{"-1", 40, decimal.RoundHalfUp, "-Infinity"},
		{"0.8", 10, decimal.RoundCeil, "1.098612289"},
		{"-0.6", 5, decimal.RoundHalfEven, "-0.69315"},
		{"0.7", 8, decimal.RoundCeil, "0.86730053"},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, tt.rounding
		x, err := ctx.New(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		if got := x.Atanh().ValueOf(); got != tt.want {
			t.Errorf("Decimal(%s).Atanh() = %s, want %s", tt.input, got, tt.want)
		}
	}
}
