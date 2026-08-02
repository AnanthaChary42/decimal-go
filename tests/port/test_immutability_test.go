package port_test

import (
	"math/rand"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_Immutability ports the immutability.js test (Option B: subset of methods that
// exist in the Go port). Verifies that every operation does not mutate its operands.
func TestOriginal_Immutability(t *testing.T) {
	saveDefaultCtx(t)

	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(20),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-7),
		ToExpPos:  decimal.IntPtr(21),
		MinE:      decimal.IntPtr(-9000000000000000),
		MaxE:      decimal.IntPtr(9000000000000000),
	})

	ctx := decimal.GetDefaultContext()

	rng := rand.New(rand.NewSource(42))

	// Integer [0, 9e15), with each possible number of digits, [1, 16], equally likely.
	randInt := func() int64 {
		return rng.Int63n(9000000000000000) / int64Pow10(rng.Intn(16))
	}

	v := []interface{}{
		float64(0),
		"NaN",
		"Infinity",
		"-Infinity",
		0.5,
		-0.5,
		float64(1),
		float64(-1),
	}
	// Add random values like JS version.
	rx := ctx.Random()
	v = append(v, rx)
	v = append(v, rx.Neg())
	ri := randInt()
	v = append(v, float64(ri))
	v = append(v, float64(-ri))
	rx2 := ctx.Random().Add(newDec(t, float64(randInt())))
	v = append(v, rx2)
	v = append(v, rx2.Neg())

	for i, val := range v {
		a := newDec(t, val)
		aa := newDec(t, val)
		k := float64(rng.Intn(10)) / 10.0

		var b *decimal.Decimal
		if k == 0.5 {
			b = newDec(t, a)
		} else if k < 0.5 {
			b = a.Add(ctx.Random().Add(newDec(t, float64(randInt()))))
		} else {
			b = a.Sub(ctx.Random().Add(newDec(t, float64(randInt()))))
		}
		bb := newDec(t, b)
		n := rng.Intn(20) + 1
		_ = n

		label := func(op string) string {
			return op + " iter=" + string(rune('0'+i))
		}

		// Abs
		x := a.Abs()
		assertEqualDecimal(t, a, aa, label("abs-a"))
		_ = x

		// Neg
		x = a.Neg()
		assertEqualDecimal(t, a, aa, label("neg-a"))
		_ = x

		// Ceil
		x = a.Ceil()
		assertEqualDecimal(t, a, aa, label("ceil-a"))

		// Floor
		x = a.Floor()
		assertEqualDecimal(t, a, aa, label("floor-a"))

		// Round
		x = a.Round()
		assertEqualDecimal(t, a, aa, label("round-a"))

		// Trunc
		x = a.Trunc()
		assertEqualDecimal(t, a, aa, label("trunc-a"))

		// Add
		x = a.Add(b)
		assertEqualDecimal(t, a, aa, label("add-a"))
		assertEqualDecimal(t, b, bb, label("add-b"))
		_ = x

		// Sub
		x = a.Sub(b)
		assertEqualDecimal(t, a, aa, label("sub-a"))
		assertEqualDecimal(t, b, bb, label("sub-b"))

		// Mul
		x = a.Mul(b)
		assertEqualDecimal(t, a, aa, label("mul-a"))
		assertEqualDecimal(t, b, bb, label("mul-b"))

		// Div
		x = a.Div(b)
		assertEqualDecimal(t, a, aa, label("div-a"))
		assertEqualDecimal(t, b, bb, label("div-b"))

		// DivToInt
		x = a.DivToInt(b)
		assertEqualDecimal(t, a, aa, label("divToInt-a"))
		assertEqualDecimal(t, b, bb, label("divToInt-b"))

		// Mod
		x = a.Mod(b)
		assertEqualDecimal(t, a, aa, label("mod-a"))
		assertEqualDecimal(t, b, bb, label("mod-b"))

		// Pow
		x = a.Pow(b)
		assertEqualDecimal(t, a, aa, label("pow-a"))
		assertEqualDecimal(t, b, bb, label("pow-b"))

		// Sqrt
		x = a.Sqrt()
		assertEqualDecimal(t, a, aa, label("sqrt-a"))

		// Cmp
		_, _ = a.Cmp(b)
		assertEqualDecimal(t, a, aa, label("cmp-a"))
		assertEqualDecimal(t, b, bb, label("cmp-b"))

		// Eq
		_ = a.Eq(b)
		assertEqualDecimal(t, a, aa, label("eq-a"))
		assertEqualDecimal(t, b, bb, label("eq-b"))

		// Gt
		_ = a.Gt(b)
		assertEqualDecimal(t, a, aa, label("gt-a"))
		assertEqualDecimal(t, b, bb, label("gt-b"))

		// Gte
		_ = a.Gte(b)
		assertEqualDecimal(t, a, aa, label("gte-a"))
		assertEqualDecimal(t, b, bb, label("gte-b"))

		// Lt
		_ = a.Lt(b)
		assertEqualDecimal(t, a, aa, label("lt-a"))
		assertEqualDecimal(t, b, bb, label("lt-b"))

		// Lte
		_ = a.Lte(b)
		assertEqualDecimal(t, a, aa, label("lte-a"))
		assertEqualDecimal(t, b, bb, label("lte-b"))

		// IsFinite, IsNaN, IsZero, IsNeg, IsPos, IsInt
		a.IsFinite()
		assertEqualDecimal(t, a, aa, label("isFinite"))
		a.IsNaN()
		assertEqualDecimal(t, a, aa, label("isNaN"))
		a.IsZero()
		assertEqualDecimal(t, a, aa, label("isZero"))
		a.IsNeg()
		assertEqualDecimal(t, a, aa, label("isNeg"))
		a.IsPos()
		assertEqualDecimal(t, a, aa, label("isPos"))
		a.IsInt()
		assertEqualDecimal(t, a, aa, label("isInt"))

		// Sd
		_ = a.SD()
		assertEqualDecimal(t, a, aa, label("sd"))

		// Dp
		_ = a.Dp()
		assertEqualDecimal(t, a, aa, label("dp"))

		// ValueOf
		a.ValueOf()
		assertEqualDecimal(t, a, aa, label("valueOf"))

		// String
		a.String()
		assertEqualDecimal(t, a, aa, label("string"))

		// ToNumber
		a.ToNumber()
		assertEqualDecimal(t, a, aa, label("toNumber"))

		// Sign
		decimal.Sign(a)
		assertEqualDecimal(t, a, aa, label("sign"))

		// Min, Max
		decimal.Min(a, b)
		assertEqualDecimal(t, a, aa, label("min-a"))
		assertEqualDecimal(t, b, bb, label("min-b"))

		decimal.Max(a, b)
		assertEqualDecimal(t, a, aa, label("max-a"))
		assertEqualDecimal(t, b, bb, label("max-b"))
	}
}

func int64Pow10(n int) int64 {
	r := int64(1)
	for i := 0; i < n; i++ {
		r *= 10
	}
	return r
}
