package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func missingModuleContext() *decimal.Context {
	return &decimal.Context{
		Precision: 20,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
		Modulo:    decimal.RoundDown,
	}
}

func TestOriginal_Hypot(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		values    []any
		precision int
		rounding  decimal.RoundingMode
		want      string
	}{
		{[]any{"1", "NaN"}, 20, decimal.RoundHalfUp, "NaN"},
		{[]any{"Infinity", "NaN"}, 20, decimal.RoundHalfUp, "Infinity"},
		{[]any{"0", "-0"}, 20, decimal.RoundHalfUp, "0"},
		{[]any{"-8", "-5482493793"}, 53, decimal.RoundCeil, "5482493793.0000000058367599140481142689669647662254134"},
		{[]any{"183", "0.6"}, 22, decimal.RoundFloor, "183.0009836039139981977"},
		{[]any{"0.19", "-5288.143883"}, 42, decimal.RoundFloor, "5288.14388641329593020355143132324946951485"},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, tt.rounding
		if got := ctx.Hypot(tt.values...).ValueOf(); got != tt.want {
			t.Errorf("Hypot(%v) = %s, want %s", tt.values, got, tt.want)
		}
	}
}

func TestOriginal_ToNearest(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		x, y any
		rm   decimal.RoundingMode
		want string
	}{
		{"123.456", nil, decimal.RoundHalfUp, "123"},
		{"123.456", "0.1", decimal.RoundHalfUp, "123.5"},
		{"123.456", "-0.02", decimal.RoundHalfUp, "123.46"},
		{"-1.5", "-3", decimal.RoundDown, "-0"},
		{"-1.5", "-3", decimal.RoundCeil, "-3"},
		{"83105511539.5", 1, decimal.RoundHalfDown, "83105511539"},
		{"450", 100, decimal.RoundHalfEven, "400"},
		{"-450", 100, decimal.RoundHalfCeil, "-400"},
		{"123.456", "Infinity", decimal.RoundHalfUp, "Infinity"},
		{"-123.456", "-Infinity", decimal.RoundHalfUp, "-Infinity"},
	} {
		x := mustNew(t, ctx, tt.x)
		got := x.ToNearest(tt.y, tt.rm).ValueOf()
		if got != tt.want {
			t.Errorf("Decimal(%v).ToNearest(%v, %d) = %s, want %s", tt.x, tt.y, tt.rm, got, tt.want)
		}
	}
}

func TestOriginal_ToFraction(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		x, max, want string
	}{
		{"0.1", "", "1,10"},
		{"-0.625", "", "-5,8"},
		{"543.017930", "", "54301793,100000"},
		{"5.1582612935891", "3", "5,1"},
		{"8.14969395596340", "4682", "14645,1797"},
		{"123.45", "10", "1111,9"},
	} {
		x, err := ctx.New(tt.x)
		if err != nil {
			t.Fatal(err)
		}
		var got []*decimal.Decimal
		if tt.max == "" {
			got = x.ToFraction()
		} else {
			got = x.ToFraction(tt.max)
		}
		if len(got) != 2 || got[0].ValueOf()+","+got[1].ValueOf() != tt.want {
			t.Errorf("Decimal(%s).ToFraction(%s) = %v, want %s", tt.x, tt.max, got, tt.want)
		}
	}
}

func TestOriginal_BaseFormatting(t *testing.T) {
	ctx := missingModuleContext()
	for _, tt := range []struct {
		value, want string
		precision   int
		method      func(*decimal.Decimal, ...any) string
	}{
		{"4294967295", "0b11111111111111111111111111111111", 32, (*decimal.Decimal).ToBinary},
		{"12", "0b1100", 4, (*decimal.Decimal).ToBinary},
		{"-0.000000123456789", "-0b0.000000000000000000000010000100100011111", 17, (*decimal.Decimal).ToBinary},
		{"24", "0o30", 2, (*decimal.Decimal).ToOctal},
		{"-0.000000123456789", "-0o0.0000000204437054635763657165267673054542", 33, (*decimal.Decimal).ToOctal},
		{"8.5", "0x8.8", 2, (*decimal.Decimal).ToHex},
		{"-0.000000123456789", "-0x0.000002123e2ccefcf5e755beec5961bda3d5a3b9", 35, (*decimal.Decimal).ToHex},
	} {
		ctx.Precision, ctx.Rounding = tt.precision, decimal.RoundHalfUp
		x, err := ctx.New(tt.value)
		if err != nil {
			t.Fatal(err)
		}
		if got := tt.method(x); got != tt.want {
			t.Errorf("base conversion of %s = %s, want %s", tt.value, got, tt.want)
		}
	}

	for _, tt := range []struct {
		value, want string
		sd          int
		method      func(*decimal.Decimal, ...any) string
	}{
		{"0", "0b0p+0", 40, (*decimal.Decimal).ToBinary},
		{"0.857421875", "0b1.10110111p-1", 9, (*decimal.Decimal).ToBinary},
		{"384", "0o1.4p+8", 40, (*decimal.Decimal).ToOctal},
		{"0o1.777p-4", "0o1.777p-4", 4, (*decimal.Decimal).ToOctal},
		{"384", "0x1.8p+8", 40, (*decimal.Decimal).ToHex},
		{"3.1415926", "0x1.921fb4d12d84ap+1", 14, (*decimal.Decimal).ToHex},
	} {
		x, err := ctx.New(tt.value)
		if err != nil {
			t.Fatal(err)
		}
		if got := tt.method(x, tt.sd, decimal.RoundHalfUp); got != tt.want {
			t.Errorf("base exponent conversion of %s = %s, want %s", tt.value, got, tt.want)
		}
	}
}

func mustNew(t *testing.T, ctx *decimal.Context, value any) *decimal.Decimal {
	t.Helper()
	switch v := value.(type) {
	case string:
		x, err := ctx.New(v)
		if err != nil {
			t.Fatal(err)
		}
		return x
	case int:
		return ctx.NewFromInt64(int64(v))
	default:
		t.Fatalf("unsupported test value %T", value)
		return nil
	}
}
