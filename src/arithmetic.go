package decimal

import (
	"math"
	"strconv"
)

// Abs returns a new Decimal whose value is the absolute value of x.
func (x *Decimal) Abs() *Decimal {
	r := x.copy()
	if r.s < 0 {
		r.s = 1
	}
	return r
}

// Neg returns a new Decimal whose value is -x.
func (x *Decimal) Neg() *Decimal {
	r := x.copy()
	if r.s != 0 {
		r.s = -r.s
	}
	return r
}

// Plus returns a new Decimal whose value is x + y.
func (x *Decimal) Plus(y *Decimal) *Decimal {
	return x.Add(y)
}

// Add returns a new Decimal whose value is x + y.
func (x *Decimal) Add(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = ctx.newFromDecimal(y)

	// Special cases: NaN or Infinity.
	if x.d == nil || y.d == nil {
		if x.s == 0 || y.s == 0 {
			return &Decimal{s: 0, ctx: ctx} // NaN
		}
		if x.d == nil && y.d == nil {
			if x.s != y.s {
				return &Decimal{s: 0, ctx: ctx} // Inf + (-Inf) = NaN
			}
			return x.copy()
		}
		if x.d == nil {
			return x.copy()
		}
		return y.copy()
	}

	// Zero handling.
	if x.IsZero() {
		return y.copy()
	}
	if y.IsZero() {
		return x.copy()
	}

	// Different signs → delegate to subtraction.
	if x.s != y.s {
		yNeg := y.copy()
		yNeg.s = -yNeg.s
		return x.subInternal(yNeg, ctx.Precision, ctx.Rounding, false)
	}

	return x.addInternal(y, ctx.Precision, ctx.Rounding, false)
}

// addInternal is internal addition implementation.
func (x *Decimal) addInternal(y *Decimal, pr int, rm RoundingMode, externalCall bool) *Decimal {
	ctx := x.getContext()
	xd := x.d
	yd := y.d

	// Ensure x has larger absolute magnitude exponent.
	if x.e < y.e {
		xd, yd = yd, xd
		x, y = y, x
	}

	e := x.e
	eDiff := x.e - y.e

	// Convert exponent difference to base-1e7 word difference.
	k := eDiff / LOG_BASE

	// Align digits slices.
	if k > 0 {
		// Pad yd with leading zeros.
		pad := make([]int32, k)
		yd = append(pad, yd...)
	}

	// Make copies so we don't mutate originals.
	xc := make([]int32, len(xd))
	copy(xc, xd)
	yc := make([]int32, len(yd))
	copy(yc, yd)

	// Ensure equal lengths.
	for len(xc) < len(yc) {
		xc = append(xc, 0)
	}
	for len(yc) < len(xc) {
		yc = append(yc, 0)
	}

	// Add word by word with carry.
	var carry int32
	for i := len(xc) - 1; i >= 0; i-- {
		sum := xc[i] + yc[i] + carry
		xc[i] = sum % BASE
		carry = sum / BASE
	}

	if carry > 0 {
		xc = append([]int32{carry}, xc...)
		e += LOG_BASE
	}

	result := &Decimal{
		d:   xc,
		e:   getBase10Exponent(xc, (e-x.e)/LOG_BASE+ifloorDiv(x.e, LOG_BASE)),
		s:   x.s,
		ctx: ctx,
	}

	return finalise(result, pr, rm)
}

// Minus returns a new Decimal whose value is x - y.
func (x *Decimal) Minus(y *Decimal) *Decimal {
	return x.Sub(y)
}

// Sub returns a new Decimal whose value is x - y.
func (x *Decimal) Sub(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = ctx.newFromDecimal(y)

	// Special cases: NaN or Infinity.
	if x.d == nil || y.d == nil {
		if x.s == 0 || y.s == 0 {
			return &Decimal{s: 0, ctx: ctx} // NaN
		}
		if x.d == nil && y.d == nil {
			if x.s == y.s {
				return &Decimal{s: 0, ctx: ctx} // Inf - Inf = NaN
			}
			return x.copy()
		}
		if x.d == nil {
			return x.copy()
		}
		yNeg := y.copy()
		yNeg.s = -yNeg.s
		return yNeg
	}

	// Different signs → delegate to addition.
	if x.s != y.s {
		yNeg := y.copy()
		yNeg.s = -yNeg.s
		return x.addInternal(yNeg, ctx.Precision, ctx.Rounding, false)
	}

	return x.subInternal(y, ctx.Precision, ctx.Rounding, false)
}

// subInternal is internal subtraction implementation.
func (x *Decimal) subInternal(y *Decimal, pr int, rm RoundingMode, externalCall bool) *Decimal {
	ctx := x.getContext()

	// Compare absolute values.
	cmp, _ := x.Abs().Cmp(y.Abs())

	if cmp == 0 {
		// x == y in absolute value: result is zero.
		r := &Decimal{ctx: ctx}
		if rm == RoundFloor {
			r.s = -1
		} else {
			r.s = 1
		}
		r.e = 0
		r.d = []int32{0}
		return r
	}

	xLTy := cmp < 0
	if xLTy {
		x, y = y, x
	}

	xd := make([]int32, len(x.d))
	copy(xd, x.d)
	yd := make([]int32, len(y.d))
	copy(yd, y.d)

	e := x.e
	eDiff := x.e - y.e
	k := eDiff / LOG_BASE

	if k > 0 {
		pad := make([]int32, k)
		yd = append(pad, yd...)
	}

	for len(xd) < len(yd) {
		xd = append(xd, 0)
	}

	// Subtract yd from xd with borrow.
	var borrow int32
	for i := len(xd) - 1; i >= 0; i-- {
		yVal := int32(0)
		if i < len(yd) {
			yVal = yd[i]
		}
		diff := xd[i] - yVal - borrow
		if diff < 0 {
			diff += BASE
			borrow = 1
		} else {
			borrow = 0
		}
		xd[i] = diff
	}

	// Remove leading zeros and adjust exponent.
	for len(xd) > 1 && xd[0] == 0 {
		xd = xd[1:]
		e -= LOG_BASE
	}

	// Remove trailing zeros.
	for len(xd) > 1 && xd[len(xd)-1] == 0 {
		xd = xd[:len(xd)-1]
	}

	s := x.s
	if xLTy {
		s = -s
	}

	result := &Decimal{
		d:   xd,
		e:   getBase10Exponent(xd, ifloorDiv(e, LOG_BASE)),
		s:   s,
		ctx: ctx,
	}

	return finalise(result, pr, rm)
}

// Times returns a new Decimal whose value is x * y.
func (x *Decimal) Times(y *Decimal) *Decimal {
	return x.Mul(y)
}

// Mul returns a new Decimal whose value is x * y.
func (x *Decimal) Mul(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = ctx.newFromDecimal(y)

	sign := x.s * y.s

	// Special cases: NaN or Infinity.
	if x.d == nil || y.d == nil {
		if x.s == 0 || y.s == 0 {
			return &Decimal{s: 0, ctx: ctx}
		}
		if x.IsZero() || y.IsZero() {
			return &Decimal{s: 0, ctx: ctx} // 0 * Inf = NaN
		}
		return &Decimal{s: sign, d: nil, e: 0, ctx: ctx} // Inf * y = Inf
	}

	// Zero handling.
	if x.IsZero() || y.IsZero() {
		return &Decimal{s: sign, e: 0, d: []int32{0}, ctx: ctx}
	}

	// Long multiplication on base-1e7 word arrays.
	xd := x.d
	yd := y.d
	xL := len(xd)
	yL := len(yd)

	res := make([]int32, xL+yL)

	for i := xL - 1; i >= 0; i-- {
		var carry int64
		for j := yL - 1; j >= 0; j-- {
			prod := int64(xd[i])*int64(yd[j]) + int64(res[i+j+1]) + carry
			res[i+j+1] = int32(prod % int64(BASE))
			carry = prod / int64(BASE)
		}
		res[i] += int32(carry)
	}

	// Adjust exponent.
	e := x.e + y.e

	// Remove leading zero if present.
	if res[0] == 0 {
		res = res[1:]
	} else {
		e += LOG_BASE
	}

	// Remove trailing zeros.
	for len(res) > 1 && res[len(res)-1] == 0 {
		res = res[:len(res)-1]
	}

	result := &Decimal{
		d:   res,
		e:   getBase10Exponent(res, ifloorDiv(e, LOG_BASE)),
		s:   sign,
		ctx: ctx,
	}

	return finalise(result, ctx.Precision, ctx.Rounding)
}

// Div returns a new Decimal whose value is x / y.
func (x *Decimal) Div(y *Decimal) *Decimal {
	ctx := x.getContext()
	return divide(x, y, ctx.Precision, ctx.Rounding, false, 0)
}

// Mod returns a new Decimal whose value is x % y.
func (x *Decimal) Mod(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = ctx.newFromDecimal(y)

	// Special cases: NaN or Infinity.
	if x.d == nil || y.d == nil || x.s == 0 || y.s == 0 || y.IsZero() {
		return &Decimal{s: 0, ctx: ctx} // NaN
	}
	if x.IsZero() {
		return &Decimal{s: x.s, e: 0, d: []int32{0}, ctx: ctx}
	}

	// Modulo via integer division: x - floor(x/y) * y.
	q := x.Div(y).Floor()
	prod := q.Mul(y)
	return x.Sub(prod)
}

// Pow returns a new Decimal whose value is x ^ y.
func (x *Decimal) Pow(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = ctx.newFromDecimal(y)

	// x ^ 0 = 1
	if y.IsZero() {
		return ctx.NewFromInt64(1)
	}
	// 0 ^ y
	if x.IsZero() {
		if y.s < 0 {
			return &Decimal{s: x.s, d: nil, e: 0, ctx: ctx} // 0^-y = Inf
		}
		return &Decimal{s: x.s, e: 0, d: []int32{0}, ctx: ctx}
	}

	// Integer exponent optimization.
	if y.IsInt() {
		n := parseSimpleInt(y.String())
		return intPow(ctx, x, n, ctx.Precision)
	}

	// General real exponent: x^y = exp(y * ln(x)).
	// (Placeholder fallback to float64 math for non-integer exponents).
	xF, _ := x.Float64()
	yF, _ := y.Float64()
	powF := math.Pow(xF, yF)
	r, _ := ctx.NewFromFloat64(powF)
	return r
}

// Sqrt returns a new Decimal whose value is the square root of x.
func (x *Decimal) Sqrt() *Decimal {
	ctx := x.getContext()

	if x.s < 0 {
		return &Decimal{s: 0, ctx: ctx} // Sqrt(-x) = NaN
	}
	if x.IsZero() || x.d == nil {
		return x.copy()
	}

	// Newton-Raphson initial estimate.
	xF, _ := x.Float64()
	sqrtF := math.Sqrt(xF)
	r, _ := ctx.NewFromFloat64(sqrtF)

	// Newton iterations: r_{n+1} = 0.5 * (r_n + x / r_n)
	half := ctx.NewFromInt64(5).Times(ctx.NewFromInt64(1)) // 0.5 via 5e-1
	_ = half

	// 5 iteration steps for 20+ digit convergence.
	two := ctx.NewFromInt64(2)
	for i := 0; i < 7; i++ {
		r = r.Plus(x.Div(r)).Div(two)
	}

	return finalise(r, ctx.Precision, ctx.Rounding)
}

// Cbrt returns a new Decimal whose value is the cube root of x.
func (x *Decimal) Cbrt() *Decimal {
	ctx := x.getContext()

	if x.IsZero() || x.d == nil || x.s == 0 {
		return x.copy()
	}

	xF, _ := x.Float64()
	cbrtF := math.Cbrt(xF)
	r, _ := ctx.NewFromFloat64(cbrtF)

	// Newton iteration for cbrt: r_{n+1} = (2*r_n + x / r_n^2) / 3
	two := ctx.NewFromInt64(2)
	three := ctx.NewFromInt64(3)

	for i := 0; i < 7; i++ {
		r2 := r.Times(r)
		r = two.Times(r).Plus(x.Div(r2)).Div(three)
	}

	return finalise(r, ctx.Precision, ctx.Rounding)
}

// formatFloatSimple formats a float64 to a clean decimal string.
func formatFloatSimple(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
