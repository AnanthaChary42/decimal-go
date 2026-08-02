package decimal

import (
	"fmt"
	"math"
	"math/big"
)

// decimalArgument converts the small set of Go values that correspond to a
// JavaScript Decimal constructor argument. Keeping this conversion in one
// place lets the methods below retain the optional/overloaded shape of their
// decimal.js counterparts without forcing callers to pre-construct values.
func decimalArgument(ctx *Context, value any) (*Decimal, error) {
	switch v := value.(type) {
	case *Decimal:
		if v == nil {
			return nil, newInvalidArgError(value)
		}
		return (&Decimal{ctx: ctx}).newFromDecimal(v), nil
	case string:
		return ctx.New(v)
	case int:
		return ctx.NewFromInt64(int64(v)), nil
	case int8:
		return ctx.NewFromInt64(int64(v)), nil
	case int16:
		return ctx.NewFromInt64(int64(v)), nil
	case int32:
		return ctx.NewFromInt64(int64(v)), nil
	case int64:
		return ctx.NewFromInt64(v), nil
	case uint:
		return ctx.New(fmt.Sprint(v))
	case uint8:
		return ctx.NewFromInt64(int64(v)), nil
	case uint16:
		return ctx.NewFromInt64(int64(v)), nil
	case uint32:
		return ctx.New(fmt.Sprint(v))
	case uint64:
		return ctx.New(fmt.Sprint(v))
	case float32:
		return ctx.NewFromFloat64(float64(v))
	case float64:
		return ctx.NewFromFloat64(v)
	default:
		return nil, newInvalidArgError(value)
	}
}

func mustDecimalArgument(ctx *Context, value any) *Decimal {
	x, err := decimalArgument(ctx, value)
	if err != nil {
		panic(err)
	}
	return x
}

func roundingArgument(value any) RoundingMode {
	var n int
	switch v := value.(type) {
	case RoundingMode:
		n = int(v)
	case int:
		n = v
	case int8:
		n = int(v)
	case int16:
		n = int(v)
	case int32:
		n = int(v)
	case int64:
		if v > math.MaxInt || v < math.MinInt {
			panic(newInvalidArgError(value))
		}
		n = int(v)
	default:
		panic(newInvalidArgError(value))
	}
	if n < int(RoundUp) || n > int(RoundHalfFloor) {
		panic(newInvalidArgError(value))
	}
	return RoundingMode(n)
}

// Hypot returns sqrt(a*a + b*b + ...), using the package default context.
// It is the Go equivalent of Decimal.hypot(...). Values may be Decimals,
// strings, integer values, or float values accepted by decimalArgument.
func Hypot(values ...any) *Decimal {
	return defaultCtx.Hypot(values...)
}

// Hypot returns sqrt(a*a + b*b + ...) using ctx. It deliberately postpones
// rounding until the square root, as Decimal.hypot does.
func (ctx *Context) Hypot(values ...any) *Decimal {
	total := ctx.NewFromInt64(0)
	previousExternal := external
	external = false
	defer func() { external = previousExternal }()

	for _, value := range values {
		n := mustDecimalArgument(ctx, value)
		if !n.IsFinite() {
			if !n.IsNaN() {
				external = previousExternal
				return ctx.infinity(1)
			}
			total = n
			continue
		}
		if total.IsFinite() {
			total = total.Plus(n.Times(n))
		}
	}

	external = previousExternal
	return total.Sqrt()
}

func (ctx *Context) infinity(sign int8) *Decimal {
	return &Decimal{s: sign, d: nil, e: 0, ctx: ctx}
}

// ToNearest returns x rounded to the nearest multiple of y. With no y (or a
// nil y), the multiple is one. The optional rounding argument uses the same
// integer constants as decimal.js (0 through 8).
func (x *Decimal) ToNearest(args ...any) *Decimal {
	ctx := x.getContext()
	y := ctx.NewFromInt64(1)
	rm := ctx.Rounding
	if len(args) > 2 {
		panic(newInvalidArgError(args))
	}
	if len(args) > 0 && args[0] != nil {
		y = mustDecimalArgument(ctx, args[0])
		if len(args) == 2 {
			rm = roundingArgument(args[1])
		}
	}

	// The special-value ordering exactly follows Decimal#toNearest.
	if !x.IsFinite() {
		if y.IsNaN() {
			return y
		}
		return x.copy()
	}
	if !y.IsFinite() {
		if y.IsNaN() {
			return y
		}
		y.s = x.s
		return y
	}
	if y.IsZero() {
		y.s = x.s
		return y
	}

	q := roundedIntegerQuotient(x, y, rm)
	previousExternal := external
	external = false
	qDecimal, err := ctx.New(q.String())
	if err != nil {
		external = previousExternal
		panic(err)
	}
	result := y.Times(qDecimal)
	external = previousExternal
	// Multiplication by an integer is an exact decimal operation. finalise with
	// the no-rounding sentinel performs only the JS overflow/underflow checks.
	result.s = x.s
	return finalise(result, -999999, RoundDown)
}

// roundedIntegerQuotient returns round(x/y) under a decimal.js rounding mode.
// It avoids the context precision intentionally: toNearest is independent of
// Decimal.precision.
func roundedIntegerQuotient(x, y *Decimal, rm RoundingMode) *big.Int {
	xn, xd := decimalRational(x)
	yn, yd := decimalRational(y)
	// Work with magnitudes; the sign only selects directed/tie rounding.
	numerator := new(big.Int).Mul(xn, yd)
	denominator := new(big.Int).Mul(xd, yn)
	q, remainder := new(big.Int), new(big.Int)
	q.QuoRem(numerator, denominator, remainder)
	if remainder.Sign() == 0 {
		if x.s != y.s {
			q.Neg(q)
		}
		return q
	}

	sign := x.s * y.s
	increment := false
	switch rm {
	case RoundUp:
		increment = true
	case RoundCeil:
		increment = sign > 0
	case RoundFloor:
		increment = sign < 0
	case RoundHalfUp, RoundHalfDown, RoundHalfEven, RoundHalfCeil, RoundHalfFloor:
		comparison := new(big.Int).Lsh(remainder, 1).Cmp(denominator)
		if comparison > 0 {
			increment = true
		} else if comparison == 0 {
			switch rm {
			case RoundHalfUp:
				increment = true
			case RoundHalfEven:
				increment = q.Bit(0) == 1
			case RoundHalfCeil:
				increment = sign > 0
			case RoundHalfFloor:
				increment = sign < 0
			}
		}
	}
	if increment {
		q.Add(q, big.NewInt(1))
	}
	if sign < 0 {
		q.Neg(q)
	}
	return q
}

// ToFraction returns a numerator/denominator pair. For a finite value it is
// always a two-element slice; non-finite values are returned as a one-element
// slice containing that value, mirroring decimal.js where toFraction returns
// the Decimal itself for NaN and infinities.
func (x *Decimal) ToFraction(maxArgs ...any) []*Decimal {
	ctx := x.getContext()
	if len(maxArgs) > 1 {
		panic(newInvalidArgError(maxArgs))
	}
	if !x.IsFinite() {
		return []*Decimal{x.copy()}
	}

	numerator, denominator := decimalRational(x)
	limit := new(big.Int).Set(denominator)
	if len(maxArgs) == 1 && maxArgs[0] != nil {
		max := mustDecimalArgument(ctx, maxArgs[0])
		if !max.IsFinite() || !max.IsInt() || max.IsNeg() || max.IsZero() {
			panic(newInvalidArgError(maxArgs[0]))
		}
		maxNumerator, maxExactDenominator := decimalRational(max)
		// max is an integer, so its exact denominator is one.
		if maxExactDenominator.Cmp(big.NewInt(1)) != 0 {
			panic(newInvalidArgError(maxArgs[0]))
		}
		limit.Set(maxNumerator)
		if limit.Cmp(denominator) > 0 {
			limit.Set(denominator)
		}
	}

	// Exact values already fit the limit. This also handles zero and preserves
	// decimal.js's lowest possible denominator.
	if denominator.Cmp(limit) <= 0 {
		return fractionDecimals(ctx, numerator, denominator, x.s)
	}

	// Continued fractions, including the final semiconvergent, reproduce the
	// algorithm used by decimal.js. In a tie it selects the previous convergent.
	n := new(big.Int).Set(numerator)
	d := new(big.Int).Set(denominator)
	p0, q0 := big.NewInt(0), big.NewInt(1)
	p1, q1 := big.NewInt(1), big.NewInt(0)
	for {
		a, remainder := new(big.Int), new(big.Int)
		a.QuoRem(n, d, remainder)
		p2 := new(big.Int).Add(new(big.Int).Mul(a, p1), p0)
		q2 := new(big.Int).Add(new(big.Int).Mul(a, q1), q0)
		if q2.Cmp(limit) > 0 {
			break
		}
		p0, p1 = p1, p2
		q0, q1 = q1, q2
		if remainder.Sign() == 0 {
			return fractionDecimals(ctx, p1, q1, x.s)
		}
		n, d = d, remainder
	}

	k := new(big.Int).Quo(new(big.Int).Sub(limit, q0), q1)
	semiP := new(big.Int).Add(new(big.Int).Mul(k, p1), p0)
	semiQ := new(big.Int).Add(new(big.Int).Mul(k, q1), q0)
	currentDistance := fractionDistance(numerator, denominator, p1, q1)
	semiDistance := fractionDistance(numerator, denominator, semiP, semiQ)
	if currentDistance.Cmp(semiDistance) <= 0 {
		return fractionDecimals(ctx, p1, q1, x.s)
	}
	return fractionDecimals(ctx, semiP, semiQ, x.s)
}

func fractionDistance(n, d, p, q *big.Int) *big.Rat {
	diff := new(big.Int).Sub(new(big.Int).Mul(n, q), new(big.Int).Mul(p, d))
	diff.Abs(diff)
	return new(big.Rat).SetFrac(diff, new(big.Int).Mul(d, q))
}

func fractionDecimals(ctx *Context, numerator, denominator *big.Int, sign int8) []*Decimal {
	n, err := ctx.New(numerator.String())
	if err != nil {
		panic(err)
	}
	if sign < 0 {
		n.s = -1
	}
	d, err := ctx.New(denominator.String())
	if err != nil {
		panic(err)
	}
	return []*Decimal{n, d}
}

// decimalRational returns the positive exact rational representation of a
// finite Decimal. Its callers handle signs separately.
func decimalRational(x *Decimal) (*big.Int, *big.Int) {
	coefficientText := digitsToStringExact(x.d)
	coefficient := new(big.Int)
	coefficient.SetString(coefficientText, 10)
	power := x.e - len(coefficientText) + 1
	var denominator *big.Int
	if power >= 0 {
		coefficient.Mul(coefficient, bigPow10(power))
		denominator = big.NewInt(1)
	} else {
		denominator = bigPow10(-power)
	}
	if coefficient.Sign() != 0 {
		gcd := new(big.Int).GCD(nil, nil, coefficient, denominator)
		coefficient.Quo(coefficient, gcd)
		denominator.Quo(denominator, gcd)
	}
	return coefficient, denominator
}
