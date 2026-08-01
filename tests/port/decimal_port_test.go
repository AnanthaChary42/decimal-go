package port_test

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestGoAPI_Predicates(t *testing.T) {
	zero, _ := decimal.New("0")
	one, _ := decimal.New("1")
	negOne, _ := decimal.New("-1")
	inf, _ := decimal.New("Infinity")
	nan, _ := decimal.New("NaN")

	if !zero.IsZero() {
		t.Error("0 should be zero")
	}
	if one.IsZero() {
		t.Error("1 should not be zero")
	}
	if !nan.IsNaN() {
		t.Error("NaN should be NaN")
	}
	if !negOne.IsNeg() {
		t.Error("-1 should be negative")
	}
	if inf.IsNeg() {
		t.Error("Infinity should not be negative")
	}
}

func TestGoAPI_Comparisons(t *testing.T) {
	a, _ := decimal.New("10")
	b, _ := decimal.New("20")

	if !a.Lt(b) {
		t.Error("10 should be less than 20")
	}
	if !b.Gt(a) {
		t.Error("20 should be greater than 10")
	}
	if a.Eq(b) {
		t.Error("10 should not equal 20")
	}
}

func TestGoAPI_Constructors(t *testing.T) {
	d1 := decimal.NewFromInt64(100)
	if d1.String() != "100" {
		t.Errorf("NewFromInt64(100) = %q, want 100", d1.String())
	}

	d2 := decimal.NewFromInt64(math.MinInt64)
	if d2.String() != "-9223372036854775808" {
		t.Errorf("NewFromInt64(MinInt64) = %q, want -9223372036854775808", d2.String())
	}

	d3 := decimal.NewFromInt64(math.MaxInt64)
	if d3.String() != "9223372036854775807" {
		t.Errorf("NewFromInt64(MaxInt64) = %q, want 9223372036854775807", d3.String())
	}

	d4 := decimal.NewFromInt64(0)
	if d4.String() != "0" {
		t.Errorf("NewFromInt64(0) = %q, want 0", d4.String())
	}
}

func TestGoAPI_Context(t *testing.T) {
	ctx := decimal.DefaultContext()
	ctx.Precision = 5

	d1, _ := ctx.New("1.23456789")
	d2, _ := ctx.New("1")

	res := d1.Plus(d2)
	if res.String() != "2.2346" {
		t.Errorf("Precision 5 addition = %q, want 2.2346", res.String())
	}
}

func TestGoAPI_Formatting(t *testing.T) {
	d, _ := decimal.New("123.456")
	if d.String() != "123.456" {
		t.Errorf("String() = %q, want 123.456", d.String())
	}
}

func TestGoAPI_Arithmetic(t *testing.T) {
	a, _ := decimal.New("15")
	b, _ := decimal.New("3")

	if a.Add(b).String() != "18" {
		t.Errorf("15 + 3 = %q, want 18", a.Add(b).String())
	}
	if a.Sub(b).String() != "12" {
		t.Errorf("15 - 3 = %q, want 12", a.Sub(b).String())
	}
	if a.Mul(b).String() != "45" {
		t.Errorf("15 * 3 = %q, want 45", a.Mul(b).String())
	}
	if a.Div(b).String() != "5" {
		t.Errorf("15 / 3 = %q, want 5", a.Div(b).String())
	}
}

func TestGoAPI_Rounding(t *testing.T) {
	d, _ := decimal.New("2.7")
	if d.Ceil().String() != "3" {
		t.Errorf("Ceil(2.7) = %q, want 3", d.Ceil().String())
	}
	if d.Floor().String() != "2" {
		t.Errorf("Floor(2.7) = %q, want 2", d.Floor().String())
	}
	if d.Trunc().String() != "2" {
		t.Errorf("Trunc(2.7) = %q, want 2", d.Trunc().String())
	}
}
