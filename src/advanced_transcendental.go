package decimal

// Atan returns the inverse tangent of x, in radians, rounded using x's
// context.  atanSeries contains the arbitrary-precision Decimal evaluation;
// this wrapper handles decimal.js's non-finite and signed-zero cases.
func (x *Decimal) Atan() *Decimal {
	ctx := x.getContext()
	if x.IsNaN() {
		return decimalNaN(ctx)
	}
	if !x.IsFinite() {
		work := trigWorkContext(ctx)
		value := piForContext(work).Div(work.NewFromInt64(2))
		value.s = x.s
		return finalise(value, ctx.Precision, ctx.Rounding)
	}
	if x.IsZero() {
		return x.copy()
	}
	work := trigWorkContext(ctx)
	return finalise(atanSeries(work, workCopy(work, x)), ctx.Precision, ctx.Rounding)
}

// InverseTangent is the long-form decimal.js alias for Atan.
func (x *Decimal) InverseTangent() *Decimal { return x.Atan() }

// Atan2 returns the angle of the point (x, y), in radians, in the range
// [-pi, pi]. It follows decimal.js/ECMAScript's signed-zero and infinity
// conventions rather than converting the operands to float64.
func (ctx *Context) Atan2(yArg, xArg any) *Decimal {
	y := mustDecimalArgument(ctx, yArg)
	x := mustDecimalArgument(ctx, xArg)
	if y.IsNaN() || x.IsNaN() {
		return decimalNaN(ctx)
	}

	work := trigWorkContext(ctx)
	pi := piForContext(work)
	if !y.IsFinite() && !x.IsFinite() {
		value := pi.Times(work.NewFromInt64(1).Div(work.NewFromInt64(4)))
		if x.IsNeg() {
			value = pi.Times(work.NewFromInt64(3).Div(work.NewFromInt64(4)))
		}
		value.s = y.s
		return finalise(value, ctx.Precision, ctx.Rounding)
	}
	if !x.IsFinite() || y.IsZero() {
		var value *Decimal
		if x.IsNeg() {
			value = pi
		} else {
			value = work.NewFromInt64(0)
		}
		value.s = y.s
		return finalise(value, ctx.Precision, ctx.Rounding)
	}
	if !y.IsFinite() || x.IsZero() {
		value := pi.Div(work.NewFromInt64(2))
		value.s = y.s
		return finalise(value, ctx.Precision, ctx.Rounding)
	}

	// Both operands are finite and non-zero. Keep the extra precision while
	// taking the ratio and adding/subtracting pi in the opposite-x quadrant.
	ratio := workCopy(work, y).Div(workCopy(work, x))
	value := atanSeries(work, ratio)
	if x.IsNeg() {
		if y.IsNeg() {
			value = value.Sub(pi)
		} else {
			value = value.Plus(pi)
		}
	}
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// Atan2 is the package-level decimal.js-style overload.
func Atan2(y, x any) *Decimal { return defaultCtx.Atan2(y, x) }

// Sinh returns the hyperbolic sine of x, rounded using x's context.
func (x *Decimal) Sinh() *Decimal {
	ctx := x.getContext()
	if x.IsNaN() {
		return decimalNaN(ctx)
	}
	if !x.IsFinite() || x.IsZero() {
		return x.copy()
	}
	work := hyperbolicWorkContext(ctx, x, false)
	value := workCopy(work, x)
	positive := value.IsPos()
	abs := value.Abs()
	if abs.e <= 2 {
		value = hyperbolicSineSeries(work, abs)
	} else {
		forward := abs.Exp()
		backward := abs.Neg().Exp()
		value = forward.Sub(backward).Div(work.NewFromInt64(2))
	}
	if !positive {
		value.s = -value.s
	}
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// HyperbolicSine is the long-form decimal.js alias for Sinh.
func (x *Decimal) HyperbolicSine() *Decimal { return x.Sinh() }

// Cosh returns the hyperbolic cosine of x, rounded using x's context.
func (x *Decimal) Cosh() *Decimal {
	ctx := x.getContext()
	if x.IsNaN() {
		return decimalNaN(ctx)
	}
	if !x.IsFinite() {
		return ctx.infinity(1)
	}
	if x.IsZero() {
		return ctx.NewFromInt64(1)
	}
	work := hyperbolicWorkContext(ctx, x, false)
	abs := workCopy(work, x).Abs()
	var value *Decimal
	if abs.e <= 2 {
		value = hyperbolicCosineSeries(work, abs)
	} else {
		value = abs.Exp().Plus(abs.Neg().Exp()).Div(work.NewFromInt64(2))
	}
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// HyperbolicCosine is the long-form decimal.js alias for Cosh.
func (x *Decimal) HyperbolicCosine() *Decimal { return x.Cosh() }

// Tanh returns the hyperbolic tangent of x, rounded using x's context.
func (x *Decimal) Tanh() *Decimal {
	ctx := x.getContext()
	if x.IsNaN() {
		return decimalNaN(ctx)
	}
	if !x.IsFinite() {
		return ctx.NewFromInt64(int64(x.s))
	}
	if x.IsZero() {
		return x.copy()
	}
	work := hyperbolicWorkContext(ctx, x, false)
	value := workCopy(work, x)
	positive := value.IsPos()
	if value.Abs().e <= 2 {
		value = hyperbolicSineSeries(work, value.Abs()).Div(hyperbolicCosineSeries(work, value.Abs()))
		if !positive {
			value.s = -value.s
		}
		return finalise(value, ctx.Precision, ctx.Rounding)
	}
	one := work.NewFromInt64(1)
	decay := value.Abs().Times(work.NewFromInt64(2)).Neg().Exp()
	value = one.Sub(decay).Div(one.Plus(decay))
	if !positive {
		value.s = -value.s
	}
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// HyperbolicTangent is the long-form decimal.js alias for Tanh.
func (x *Decimal) HyperbolicTangent() *Decimal { return x.Tanh() }

// Asinh returns the inverse hyperbolic sine of x, rounded using x's context.
func (x *Decimal) Asinh() *Decimal {
	ctx := x.getContext()
	if x.IsNaN() || !x.IsFinite() || x.IsZero() {
		return x.copy()
	}
	work := hyperbolicWorkContext(ctx, x, true)
	if workCopy(work, x).Abs().e < 0 {
		return finalise(asinhSeries(work, workCopy(work, x)), ctx.Precision, ctx.Rounding)
	}
	previousExternal := external
	external = false
	value := workCopy(work, x)
	value = value.Times(value).Plus(work.NewFromInt64(1)).Sqrt().Plus(value)
	external = previousExternal
	return finalise(value.Ln(), ctx.Precision, ctx.Rounding)
}

// InverseHyperbolicSine is the long-form decimal.js alias for Asinh.
func (x *Decimal) InverseHyperbolicSine() *Decimal { return x.Asinh() }

// Acosh returns the inverse hyperbolic cosine of x, rounded using x's context.
func (x *Decimal) Acosh() *Decimal {
	ctx := x.getContext()
	if x.IsNaN() || x.IsNeg() || x.IsZero() {
		return decimalNaN(ctx)
	}
	one := ctx.NewFromInt64(1)
	if x.Eq(one) {
		return ctx.NewFromInt64(0)
	}
	if !x.IsFinite() {
		return x.copy()
	}
	work := hyperbolicWorkContext(ctx, x, true)
	previousExternal := external
	external = false
	value := workCopy(work, x)
	value = value.Times(value).Sub(work.NewFromInt64(1)).Sqrt().Plus(value)
	external = previousExternal
	return finalise(value.Ln(), ctx.Precision, ctx.Rounding)
}

// InverseHyperbolicCosine is the long-form decimal.js alias for Acosh.
func (x *Decimal) InverseHyperbolicCosine() *Decimal { return x.Acosh() }

// Atanh returns the inverse hyperbolic tangent of x, rounded using x's
// context. Its domain is [-1, 1].
func (x *Decimal) Atanh() *Decimal {
	ctx := x.getContext()
	if !x.IsFinite() {
		return decimalNaN(ctx)
	}
	if x.IsZero() {
		return x.copy()
	}
	one := ctx.NewFromInt64(1)
	cmp, _ := x.Abs().Cmp(one)
	if cmp > 0 {
		return decimalNaN(ctx)
	}
	if cmp == 0 {
		return ctx.infinity(x.s)
	}

	// When x is sufficiently small, all terms after the first are beyond the
	// configured precision. Returning x directly also avoids losing a tiny x
	// when (1+x)/(1-x) is formed at finite precision.
	if x.e < 0 {
		xsd := x.Sd()
		threshold := xsd
		if ctx.Precision > threshold {
			threshold = ctx.Precision
		}
		if threshold < 2*(-x.e)-1 {
			return finalise(x.copy(), ctx.Precision, ctx.Rounding, true)
		}
	}

	guard := x.Sd() + 12
	if negExponent := -x.e; negExponent > guard {
		guard = negExponent + 12
	}
	work := ctx.Clone(ConfigOptions{Precision: intPtr(ctx.Precision + guard), Rounding: intPtr(int(RoundDown))})
	value := workCopy(work, x)
	oneWork := work.NewFromInt64(1)
	ratio := oneWork.Plus(value).Div(oneWork.Sub(value))
	value = ratio.Ln().Div(work.NewFromInt64(2))
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// InverseHyperbolicTangent is the long-form decimal.js alias for Atanh.
func (x *Decimal) InverseHyperbolicTangent() *Decimal { return x.Atanh() }

// Package-level functions mirror decimal.js's static methods.
func Atan(value any) *Decimal  { return defaultCtx.Atan(value) }
func Sinh(value any) *Decimal  { return defaultCtx.Sinh(value) }
func Cosh(value any) *Decimal  { return defaultCtx.Cosh(value) }
func Tanh(value any) *Decimal  { return defaultCtx.Tanh(value) }
func Asinh(value any) *Decimal { return defaultCtx.Asinh(value) }
func Acosh(value any) *Decimal { return defaultCtx.Acosh(value) }
func Atanh(value any) *Decimal { return defaultCtx.Atanh(value) }

func (ctx *Context) Atan(value any) *Decimal  { return mustDecimalArgument(ctx, value).Atan() }
func (ctx *Context) Sinh(value any) *Decimal  { return mustDecimalArgument(ctx, value).Sinh() }
func (ctx *Context) Cosh(value any) *Decimal  { return mustDecimalArgument(ctx, value).Cosh() }
func (ctx *Context) Tanh(value any) *Decimal  { return mustDecimalArgument(ctx, value).Tanh() }
func (ctx *Context) Asinh(value any) *Decimal { return mustDecimalArgument(ctx, value).Asinh() }
func (ctx *Context) Acosh(value any) *Decimal { return mustDecimalArgument(ctx, value).Acosh() }
func (ctx *Context) Atanh(value any) *Decimal { return mustDecimalArgument(ctx, value).Atanh() }

func hyperbolicWorkContext(ctx *Context, x *Decimal, cancellation bool) *Context {
	guard := 32
	if x.e > 0 {
		if cancellation {
			guard += 2 * x.e
		} else {
			guard += x.e
		}
	}
	if x.Sd() > guard {
		guard = x.Sd() + 12
	}
	return ctx.Clone(ConfigOptions{Precision: intPtr(ctx.Precision + guard), Rounding: intPtr(int(RoundDown))})
}

func hyperbolicSineSeries(ctx *Context, x *Decimal) *Decimal {
	if x.IsZero() {
		return x
	}
	x2 := x.Times(x)
	term := workCopy(ctx, x)
	sum := workCopy(ctx, x)
	previous := sum
	for n := int64(3); n < 100000; n += 2 {
		term = term.Times(x2).Div(ctx.NewFromInt64(n * (n - 1)))
		sum = previous.Plus(term)
		if sum.Eq(previous) || term.IsZero() {
			break
		}
		previous = sum
	}
	return sum
}

func hyperbolicCosineSeries(ctx *Context, x *Decimal) *Decimal {
	if x.IsZero() {
		return ctx.NewFromInt64(1)
	}
	x2 := x.Times(x)
	term := ctx.NewFromInt64(1)
	sum := workCopy(ctx, term)
	previous := sum
	for n := int64(2); n < 100000; n += 2 {
		term = term.Times(x2).Div(ctx.NewFromInt64(n * (n - 1)))
		sum = previous.Plus(term)
		if sum.Eq(previous) || term.IsZero() {
			break
		}
		previous = sum
	}
	return sum
}

func asinhSeries(ctx *Context, x *Decimal) *Decimal {
	if x.IsZero() {
		return x
	}
	x2 := x.Times(x)
	term := workCopy(ctx, x)
	sum := workCopy(ctx, x)
	previous := sum
	for k := int64(1); k < 100000; k++ {
		numerator := -(2*k - 1) * (2*k - 1)
		denominator := 2 * k * (2*k + 1)
		term = term.Times(x2).Times(ctx.NewFromInt64(numerator)).Div(ctx.NewFromInt64(denominator))
		sum = previous.Plus(term)
		if sum.Eq(previous) || term.IsZero() {
			break
		}
		previous = sum
	}
	return sum
}
