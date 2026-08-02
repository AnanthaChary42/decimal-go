package decimal

// The trigonometric methods below intentionally use Decimal arithmetic for
// argument reduction and the series evaluation.  Converting the input to a
// float64 would lose the very property this package provides: decimal.js can
// evaluate these functions at arbitrary configured precision.

// Sin returns the sine of x, in radians, rounded using x's context.
func (x *Decimal) Sin() *Decimal {
	ctx := x.getContext()
	if !x.IsFinite() {
		return decimalNaN(ctx)
	}
	if x.IsZero() {
		return x.copy()
	}
	work, reduced, quadrant := reduceTrigArgument(ctx, x)
	value := sineSeries(work, reduced)
	if quadrant > 2 {
		value.s = -value.s
	}
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// Sine is the long-form decimal.js alias for Sin.
func (x *Decimal) Sine() *Decimal { return x.Sin() }

// Cos returns the cosine of x, in radians, rounded using x's context.
func (x *Decimal) Cos() *Decimal {
	ctx := x.getContext()
	if !x.IsFinite() {
		return decimalNaN(ctx)
	}
	if x.IsZero() {
		return ctx.NewFromInt64(1)
	}
	work, reduced, quadrant := reduceTrigArgument(ctx, x)
	value := cosineSeries(work, reduced)
	if quadrant == 2 || quadrant == 3 {
		value.s = -value.s
	}
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// Cosine is the long-form decimal.js alias for Cos.
func (x *Decimal) Cosine() *Decimal { return x.Cos() }

// Tan returns the tangent of x, in radians, rounded using x's context.
func (x *Decimal) Tan() *Decimal {
	ctx := x.getContext()
	if !x.IsFinite() {
		return decimalNaN(ctx)
	}
	if x.IsZero() {
		return x.copy()
	}
	work, reduced, quadrant := reduceTrigArgument(ctx, x)
	sine := sineSeries(work, reduced)
	cosine := cosineSeries(work, reduced)
	if quadrant > 2 {
		sine.s = -sine.s
	}
	if quadrant == 2 || quadrant == 3 {
		cosine.s = -cosine.s
	}
	value := sine.Div(cosine)
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// Tangent is the long-form decimal.js alias for Tan.
func (x *Decimal) Tangent() *Decimal { return x.Tan() }

// Asin returns the inverse sine of x, in radians, rounded using x's context.
func (x *Decimal) Asin() *Decimal {
	ctx := x.getContext()
	if x.IsZero() {
		return x.copy()
	}
	if !x.IsFinite() {
		return decimalNaN(ctx)
	}
	one := ctx.NewFromInt64(1)
	cmp, _ := x.Abs().Cmp(one)
	if cmp > 0 {
		return decimalNaN(ctx)
	}
	work, workX := inverseTrigWork(ctx, x)
	if cmp == 0 {
		value := piForContext(work).Div(work.NewFromInt64(2))
		value.s = x.s
		return finalise(value, ctx.Precision, ctx.Rounding)
	}
	value := asinWork(work, workX)
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// InverseSine is the long-form decimal.js alias for Asin.
func (x *Decimal) InverseSine() *Decimal { return x.Asin() }

// Acos returns the inverse cosine of x, in radians, rounded using x's context.
func (x *Decimal) Acos() *Decimal {
	ctx := x.getContext()
	if !x.IsFinite() {
		return decimalNaN(ctx)
	}
	one := ctx.NewFromInt64(1)
	cmp, _ := x.Abs().Cmp(one)
	if cmp > 0 {
		return decimalNaN(ctx)
	}
	if cmp == 0 {
		if x.IsNeg() {
			return finalise(piForContext(ctx), ctx.Precision, ctx.Rounding)
		}
		return ctx.NewFromInt64(0)
	}
	work, workX := inverseTrigWork(ctx, x)
	if x.IsZero() {
		return finalise(piForContext(work).Div(work.NewFromInt64(2)), ctx.Precision, ctx.Rounding)
	}
	value := piForContext(work).Div(work.NewFromInt64(2)).Sub(asinWork(work, workX))
	return finalise(value, ctx.Precision, ctx.Rounding)
}

// InverseCosine is the long-form decimal.js alias for Acos.
func (x *Decimal) InverseCosine() *Decimal { return x.Acos() }

// Package-level functions mirror decimal.js's static methods.
func Sin(value any) *Decimal  { return defaultCtx.Sin(value) }
func Cos(value any) *Decimal  { return defaultCtx.Cos(value) }
func Tan(value any) *Decimal  { return defaultCtx.Tan(value) }
func Asin(value any) *Decimal { return defaultCtx.Asin(value) }
func Acos(value any) *Decimal { return defaultCtx.Acos(value) }

func (ctx *Context) Sin(value any) *Decimal  { return mustDecimalArgument(ctx, value).Sin() }
func (ctx *Context) Cos(value any) *Decimal  { return mustDecimalArgument(ctx, value).Cos() }
func (ctx *Context) Tan(value any) *Decimal  { return mustDecimalArgument(ctx, value).Tan() }
func (ctx *Context) Asin(value any) *Decimal { return mustDecimalArgument(ctx, value).Asin() }
func (ctx *Context) Acos(value any) *Decimal { return mustDecimalArgument(ctx, value).Acos() }

// trigWorkContext provides guard digits for inverse functions and the
// Taylor series.  The original library uses a larger input-dependent guard;
// Decimal's exact coefficient representation lets a fixed guard be sufficient
// for all configured precisions while avoiding unbounded contexts for huge e.
func trigWorkContext(ctx *Context) *Context {
	guard := 32
	if ctx.Precision < 20 {
		guard = 40
	}
	return ctx.Clone(ConfigOptions{Precision: intPtr(ctx.Precision + guard), Rounding: intPtr(int(RoundDown))})
}

func trigWorkContextFor(ctx *Context, x *Decimal) *Context {
	guard := 32
	if x.e > guard {
		guard = x.e + 12
	}
	if x.Sd() > guard {
		guard = x.Sd() + 12
	}
	return ctx.Clone(ConfigOptions{Precision: intPtr(ctx.Precision + guard), Rounding: intPtr(int(RoundDown))})
}

func intPtr(v int) *int { return &v }

func piForContext(ctx *Context) *Decimal {
	value, err := ctx.New(PI)
	if err != nil {
		panic(err)
	}
	return value
}

func workCopy(ctx *Context, x *Decimal) *Decimal {
	value := x.copy()
	value.ctx = ctx
	return value
}

// reduceTrigArgument returns |r| <= pi/2 and the decimal.js quadrant (1..4).
func reduceTrigArgument(ctx *Context, x *Decimal) (*Context, *Decimal, int) {
	work := trigWorkContextFor(ctx, x)
	value := workCopy(work, x)
	isNeg := value.IsNeg()
	value.s = 1
	pi := piForContext(work)
	halfPi := pi.Div(work.NewFromInt64(2))
	if value.Lte(halfPi) {
		if isNeg {
			return work, value, 4
		}
		return work, value, 1
	}

	quotient := value.DivToInt(pi)
	if quotient.IsZero() {
		if isNeg {
			return work, pi.Sub(value), 3
		}
		return work, pi.Sub(value), 2
	}
	remainder := value.Sub(quotient.Times(pi))
	odd := !quotient.Mod(work.NewFromInt64(2)).IsZero()
	if remainder.Lte(halfPi) {
		if odd {
			if isNeg {
				return work, remainder, 2
			}
			return work, remainder, 3
		}
		if isNeg {
			return work, remainder, 4
		}
		return work, remainder, 1
	}
	if odd {
		if isNeg {
			return work, pi.Sub(remainder), 1
		}
		return work, pi.Sub(remainder), 4
	}
	if isNeg {
		return work, pi.Sub(remainder), 3
	}
	return work, pi.Sub(remainder), 2
}

func sineSeries(ctx *Context, x *Decimal) *Decimal {
	if x.IsZero() {
		return x
	}
	x2 := x.Times(x)
	term := workCopy(ctx, x)
	sum := workCopy(ctx, x)
	previous := sum
	for n := int64(3); n < 100000; n += 2 {
		term = term.Times(x2).Div(ctx.NewFromInt64(n * (n - 1)))
		if (n/2)%2 == 1 {
			sum = previous.Sub(term)
		} else {
			sum = previous.Plus(term)
		}
		if sum.Eq(previous) || term.IsZero() {
			break
		}
		previous = sum
	}
	return sum
}

func cosineSeries(ctx *Context, x *Decimal) *Decimal {
	if x.IsZero() {
		return ctx.NewFromInt64(1)
	}
	x2 := x.Times(x)
	term := ctx.NewFromInt64(1)
	sum := workCopy(ctx, term)
	previous := sum
	for n := int64(2); n < 100000; n += 2 {
		term = term.Times(x2).Div(ctx.NewFromInt64(n * (n - 1)))
		if (n/2)%2 == 1 {
			sum = previous.Sub(term)
		} else {
			sum = previous.Plus(term)
		}
		if sum.Eq(previous) || term.IsZero() {
			break
		}
		previous = sum
	}
	return sum
}

func inverseTrigWork(ctx *Context, x *Decimal) (*Context, *Decimal) {
	work := trigWorkContext(ctx)
	return work, workCopy(work, x)
}

func asinWork(ctx *Context, x *Decimal) *Decimal {
	one := ctx.NewFromInt64(1)
	x2 := x.Times(x)
	root := one.Sub(x2).Sqrt()
	ratio := x.Div(root.Plus(one))
	return atanSeries(ctx, ratio).Times(ctx.NewFromInt64(2))
}

// atanSeries is the same argument-reduction identity used by decimal.js:
// atan(x) = 2 atan(x/(1 + sqrt(1+x*x))).  Repeating it makes the alternating
// Taylor series rapidly convergent even when the original argument is large.
func atanSeries(ctx *Context, x *Decimal) *Decimal {
	if x.IsZero() {
		return x
	}
	if x.Abs().Eq(ctx.NewFromInt64(1)) {
		value := piForContext(ctx).Div(ctx.NewFromInt64(4))
		value.s = x.s
		return value
	}
	k := ctx.Precision/LOG_BASE + 2
	if k > 28 {
		k = 28
	}
	for i := 0; i < k; i++ {
		x = x.Div(x.Times(x).Plus(ctx.NewFromInt64(1)).Sqrt().Plus(ctx.NewFromInt64(1)))
	}
	x2 := x.Times(x)
	power := workCopy(ctx, x)
	term := workCopy(ctx, x)
	sum := workCopy(ctx, x)
	previous := sum
	for n := int64(3); n < 100000; n += 2 {
		power = power.Times(x2)
		term = power.Div(ctx.NewFromInt64(n))
		if ((n-1)/2)%2 == 1 {
			sum = previous.Sub(term)
		} else {
			sum = previous.Plus(term)
		}
		if sum.Eq(previous) || term.IsZero() {
			break
		}
		previous = sum
	}
	if k > 0 {
		sum = sum.Times(ctx.NewFromInt64(int64(1 << k)))
	}
	return sum
}
