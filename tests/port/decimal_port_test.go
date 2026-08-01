package port_test

// ============================================================================
// tests/port/decimal_port_test.go
//
// Go-idiomatic test suite for the decimal-go API.
// Focuses on Go-specific interfaces: Context structs, (Decimal, error) return
// signatures, type conversions (int64, float64), and Go error handling.
// ============================================================================

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// ---------------------------------------------------------------------------
// 1. Constructors & Error Handling
// ---------------------------------------------------------------------------
func TestGoAPI_Constructors(t *testing.T) {
	t.Run("New_ValidStrings", func(t *testing.T) {
		valid := []string{"0", "123.456", "-0.00123", "1e10", "-2.5e-4", "0x1a", "0b1010", "Infinity", "-Infinity", "NaN"}
		for _, input := range valid {
			d, err := decimal.New(input)
			if err != nil {
				t.Errorf("New(%q) unexpected error: %v", input, err)
			}
			if d == nil {
				t.Errorf("New(%q) returned nil Decimal", input)
			}
		}
	})

	t.Run("New_InvalidStrings", func(t *testing.T) {
		invalid := []string{"", "abc", "1.2.3", "++1", "0xXYZ", "---5"}
		for _, input := range invalid {
			d, err := decimal.New(input)
			if err == nil {
				t.Errorf("New(%q) expected error, got nil (Decimal: %v)", input, d)
			}
		}
	})

	t.Run("NewFromInt64", func(t *testing.T) {
		tests := []struct {
			input int64
			want  string
		}{
			{0, "0"},
			{100, "100"},
			{-100, "-100"},
			{9223372036854775807, "9223372036854775807"},
			{-9223372036854775808, "-9223372036854775808"},
		}
		for _, tt := range tests {
			d := decimal.NewFromInt64(tt.input)
			if d.String() != tt.want {
				t.Errorf("NewFromInt64(%d) = %q, want %q", tt.input, d.String(), tt.want)
			}
		}
	})

	t.Run("NewFromFloat64", func(t *testing.T) {
		tests := []struct {
			input float64
			want  string
		}{
			{0.0, "0"},
			{1.5, "1.5"},
			{-2.25, "-2.25"},
			{math.Inf(1), "Infinity"},
			{math.Inf(-1), "-Infinity"},
		}
		for _, tt := range tests {
			d, err := decimal.NewFromFloat64(tt.input)
			if err != nil {
				t.Errorf("NewFromFloat64(%g) unexpected error: %v", tt.input, err)
			}
			if d.String() != tt.want {
				t.Errorf("NewFromFloat64(%g) = %q, want %q", tt.input, d.String(), tt.want)
			}
		}

		// Check NaN separately
		nanD, err := decimal.NewFromFloat64(math.NaN())
		if err != nil || !nanD.IsNaN() {
			t.Errorf("NewFromFloat64(NaN) = %v, err=%v, want IsNaN()=true", nanD, err)
		}
	})
}

// ---------------------------------------------------------------------------
// 2. Custom Context & Rounding Configurations
// ---------------------------------------------------------------------------
func TestGoAPI_Context(t *testing.T) {
	t.Run("Precision_Truncation", func(t *testing.T) {
		ctx := &decimal.Context{
			Precision: 5,
			Rounding:  decimal.RoundHalfUp,
			ToExpNeg:  -7,
			ToExpPos:  21,
			MinE:      -9e15,
			MaxE:      9e15,
		}
		d, err := ctx.New("12.3456789")
		if err != nil {
			t.Fatal(err)
		}
		// Precision 5: 12.346
		if d.String() != "12.346" {
			t.Errorf("Custom context precision 5: got %q, want %q", d.String(), "12.346")
		}
	})

	t.Run("RoundingModes", func(t *testing.T) {
		modes := []struct {
			mode decimal.RoundingMode
			want string
		}{
			{decimal.RoundUp, "1.3"},
			{decimal.RoundDown, "1.2"},
			{decimal.RoundCeil, "1.3"},
			{decimal.RoundFloor, "1.2"},
			{decimal.RoundHalfUp, "1.3"},
			{decimal.RoundHalfDown, "1.2"},
			{decimal.RoundHalfEven, "1.2"},
		}
		for _, tt := range modes {
			ctx := &decimal.Context{
				Precision: 2,
				Rounding:  tt.mode,
				ToExpNeg:  -7,
				ToExpPos:  21,
				MinE:      -9e15,
				MaxE:      9e15,
			}
			d, _ := ctx.New("1.25")
			if d.String() != tt.want {
				t.Errorf("Rounding mode %v on 1.25: got %q, want %q", tt.mode, d.String(), tt.want)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// 3. Predicates & Property Inspection
// ---------------------------------------------------------------------------
func TestGoAPI_Predicates(t *testing.T) {
	zero, _ := decimal.New("0")
	pos, _ := decimal.New("42.5")
	neg, _ := decimal.New("-10.1")
	inf, _ := decimal.New("Infinity")
	nan, _ := decimal.New("NaN")

	t.Run("IsZero", func(t *testing.T) {
		if !zero.IsZero() {
			t.Error("0.IsZero() want true")
		}
		if pos.IsZero() {
			t.Error("42.5.IsZero() want false")
		}
	})

	t.Run("IsPos_IsNeg", func(t *testing.T) {
		if !pos.IsPos() || pos.IsNeg() {
			t.Error("42.5 should be Pos and not Neg")
		}
		if !neg.IsNeg() || neg.IsPos() {
			t.Error("-10.1 should be Neg and not Pos")
		}
	})

	t.Run("IsFinite_IsNaN", func(t *testing.T) {
		if !pos.IsFinite() || !zero.IsFinite() {
			t.Error("finite numbers return false for IsFinite")
		}
		if inf.IsFinite() || nan.IsFinite() {
			t.Error("Infinity/NaN returned true for IsFinite")
		}
		if !nan.IsNaN() || pos.IsNaN() {
			t.Error("IsNaN failed for NaN or non-NaN")
		}
	})

	t.Run("IsInt", func(t *testing.T) {
		intD, _ := decimal.New("100")
		floatD, _ := decimal.New("100.5")
		if !intD.IsInt() {
			t.Error("100.IsInt() want true")
		}
		if floatD.IsInt() {
			t.Error("100.5.IsInt() want false")
		}
	})

	t.Run("Sd_Dp", func(t *testing.T) {
		d, _ := decimal.New("123.456")
		if d.Sd() != 6 {
			t.Errorf("123.456.Sd() = %d, want 6", d.Sd())
		}
		if d.Dp() != 3 {
			t.Errorf("123.456.Dp() = %d, want 3", d.Dp())
		}
	})
}

// ---------------------------------------------------------------------------
// 4. Comparison Operators & Helpers
// ---------------------------------------------------------------------------
func TestGoAPI_Comparisons(t *testing.T) {
	a, _ := decimal.New("10")
	b, _ := decimal.New("20")
	c, _ := decimal.New("10")

	t.Run("Cmp", func(t *testing.T) {
		cmp, ok := a.Cmp(b)
		if !ok || cmp != -1 {
			t.Errorf("10.Cmp(20) = (%d, %v), want (-1, true)", cmp, ok)
		}
		cmp, ok = b.Cmp(a)
		if !ok || cmp != 1 {
			t.Errorf("20.Cmp(10) = (%d, %v), want (1, true)", cmp, ok)
		}
		cmp, ok = a.Cmp(c)
		if !ok || cmp != 0 {
			t.Errorf("10.Cmp(10) = (%d, %v), want (0, true)", cmp, ok)
		}
	})

	t.Run("BooleanComparisons", func(t *testing.T) {
		if !a.Eq(c) || !a.Equals(c) {
			t.Error("10 == 10 failed")
		}
		if !a.Lt(b) || !a.Lte(b) || !a.Lte(c) {
			t.Error("Lt/Lte check failed")
		}
		if !b.Gt(a) || !b.Gte(a) || !b.Gte(c) {
			t.Error("Gt/Gte check failed")
		}
	})

	t.Run("Clamp", func(t *testing.T) {
		min, _ := decimal.New("5")
		max, _ := decimal.New("15")
		valLow, _ := decimal.New("2")
		valHigh, _ := decimal.New("25")
		valMid, _ := decimal.New("10")

		if !valLow.Clamp(min, max).Eq(min) {
			t.Errorf("Clamp(2, 5, 15) want 5")
		}
		if !valHigh.Clamp(min, max).Eq(max) {
			t.Errorf("Clamp(25, 5, 15) want 15")
		}
		if !valMid.Clamp(min, max).Eq(valMid) {
			t.Errorf("Clamp(10, 5, 15) want 10")
		}
	})
}

// ---------------------------------------------------------------------------
// 5. Formatting & Export
// ---------------------------------------------------------------------------
func TestGoAPI_Formatting(t *testing.T) {
	d, _ := decimal.New("12345.6789")

	t.Run("ToFixed", func(t *testing.T) {
		str, err := d.ToFixed(2)
		if err != nil || str != "12345.68" {
			t.Errorf("ToFixed(2) = %q, %v, want 12345.68", str, err)
		}
	})

	t.Run("ToExponential", func(t *testing.T) {
		str, err := d.ToExponential(2)
		if err != nil || str != "1.23e+4" {
			t.Errorf("ToExponential(2) = %q, %v, want 1.23e+4", str, err)
		}
	})

	t.Run("ToPrecision", func(t *testing.T) {
		str, err := d.ToPrecision(4)
		if err != nil || str != "12350" {
			t.Errorf("ToPrecision(4) = %q, %v, want 12350", str, err)
		}
	})

	t.Run("Float64", func(t *testing.T) {
		f, exact := d.Float64()
		if !exact || f != 12345.6789 {
			t.Errorf("Float64() = %g, exact=%v, want 12345.6789", f, exact)
		}
	})
}

// ---------------------------------------------------------------------------
// 6. Arithmetic Operations
// ---------------------------------------------------------------------------
func TestGoAPI_Arithmetic(t *testing.T) {
	a, _ := decimal.New("12.5")
	b, _ := decimal.New("2.5")

	t.Run("Plus_Minus_Times_Div", func(t *testing.T) {
		if a.Plus(b).String() != "15" {
			t.Errorf("12.5 + 2.5 = %q, want 15", a.Plus(b).String())
		}
		if a.Minus(b).String() != "10" {
			t.Errorf("12.5 - 2.5 = %q, want 10", a.Minus(b).String())
		}
		if a.Times(b).String() != "31.25" {
			t.Errorf("12.5 * 2.5 = %q, want 31.25", a.Times(b).String())
		}
		if a.Div(b).String() != "5" {
			t.Errorf("12.5 / 2.5 = %q, want 5", a.Div(b).String())
		}
	})

	t.Run("Abs_Neg", func(t *testing.T) {
		negA := a.Neg()
		if negA.String() != "-12.5" {
			t.Errorf("Neg(12.5) = %q, want -12.5", negA.String())
		}
		if negA.Abs().String() != "12.5" {
			t.Errorf("Abs(-12.5) = %q, want 12.5", negA.Abs().String())
		}
	})

	t.Run("Sqrt_Cbrt", func(t *testing.T) {
		nine, _ := decimal.New("9")
		if nine.Sqrt().String() != "3" {
			t.Errorf("Sqrt(9) = %q, want 3", nine.Sqrt().String())
		}
		twentySeven, _ := decimal.New("27")
		if twentySeven.Cbrt().String() != "3" {
			t.Errorf("Cbrt(27) = %q, want 3", twentySeven.Cbrt().String())
		}
	})
}

// ---------------------------------------------------------------------------
// 7. Rounding Methods
// ---------------------------------------------------------------------------
func TestGoAPI_Rounding(t *testing.T) {
	val, _ := decimal.New("3.75")
	negVal, _ := decimal.New("-3.75")

	t.Run("Ceil_Floor_Round_Trunc", func(t *testing.T) {
		if val.Ceil().String() != "4" {
			t.Errorf("Ceil(3.75) = %q, want 4", val.Ceil().String())
		}
		if val.Floor().String() != "3" {
			t.Errorf("Floor(3.75) = %q, want 3", val.Floor().String())
		}
		if val.Trunc().String() != "3" {
			t.Errorf("Trunc(3.75) = %q, want 3", val.Trunc().String())
		}
		if negVal.Ceil().String() != "-3" {
			t.Errorf("Ceil(-3.75) = %q, want -3", negVal.Ceil().String())
		}
		if negVal.Floor().String() != "-4" {
			t.Errorf("Floor(-3.75) = %q, want -4", negVal.Floor().String())
		}
	})

	t.Run("ToDP_ToSD", func(t *testing.T) {
		if val.ToDP(1).String() != "3.8" {
			t.Errorf("3.75.ToDP(1) = %q, want 3.8", val.ToDP(1).String())
		}
		if val.ToSD(2).String() != "3.8" {
			t.Errorf("3.75.ToSD(2) = %q, want 3.8", val.ToSD(2).String())
		}
	})
}
