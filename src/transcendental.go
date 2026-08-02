package decimal

import (
	"math/big"
	"strings"
)

// Ln returns the natural logarithm of x, rounded using x's context.
func (x *Decimal) Ln() *Decimal {
	ctx := x.getContext()
	if x.IsZero() {
		return ctx.infinity(-1)
	}
	if x.IsNaN() || x.IsNeg() {
		return decimalNaN(ctx)
	}
	if !x.IsFinite() {
		return ctx.infinity(1)
	}
	if x.e == 0 && len(x.d) == 1 && x.d[0] == 1 {
		return ctx.NewFromInt64(0)
	}

	bits, digits := transcendentalPrecision(ctx.Precision)
	value := decimalLnBigFloat(x, bits)
	return decimalFromBigFloat(ctx, value, 0, digits, ctx.Rounding)
}

// NaturalLogarithm is the long-form decimal.js alias for Ln.
func (x *Decimal) NaturalLogarithm() *Decimal {
	return x.Ln()
}

// Exp returns e raised to x, rounded using x's context.
func (x *Decimal) Exp() *Decimal {
	ctx := x.getContext()
	if x.IsNaN() {
		return decimalNaN(ctx)
	}
	if !x.IsFinite() {
		if x.IsNeg() {
			return ctx.NewFromInt64(0)
		}
		return ctx.infinity(1)
	}
	if x.IsZero() {
		return ctx.NewFromInt64(1)
	}

	// Inputs with a decimal exponent outside this range cannot affect the
	// result's configured exponent, except on the very small side where 1+x is
	// sufficient to preserve directed rounding. Avoid converting 10^9e15 to a
	// binary float merely to determine this.
	if x.e > 16 {
		if x.IsNeg() {
			return ctx.NewFromInt64(0)
		}
		return ctx.infinity(1)
	}
	if x.e < -1000 {
		one := ctx.NewFromInt64(1)
		return one.Plus(x)
	}

	bits, digits := transcendentalPrecision(ctx.Precision)
	xf, ok := decimalBigFloat(x, bits)
	if !ok {
		return decimalNaN(ctx)
	}
	// For modest arguments the direct series retains tiny terms such as x²/2,
	// which are observable at high decimal precision. The decimal-exponent
	// decomposition below is reserved for genuinely large arguments.
	if x.e <= 1 {
		return decimalFromBigFloat(ctx, bigExp(xf, bits), 0, digits, ctx.Rounding)
	}
	ln10 := bigLn10(bits)
	quotient := bigQuo(xf, ln10, bits)
	exponent, fractional := splitFloor(quotient, bits)

	// exp(x) = 10^exponent * exp(fractional * ln(10)), where the latter
	// mantissa is in [1, 10). This handles x values near 2e16 without asking
	// math/big to materialise an exponent of e^x.
	if !exponent.IsInt64() {
		if exponent.Sign() < 0 {
			return ctx.NewFromInt64(0)
		}
		return ctx.infinity(1)
	}
	decimalExponent := int(exponent.Int64())
	if decimalExponent > ctx.MaxE {
		return ctx.infinity(1)
	}
	if decimalExponent < ctx.MinE {
		return ctx.NewFromInt64(0)
	}
	mantissa := bigExp(bigMul(fractional, ln10, bits), bits)
	return decimalFromBigFloat(ctx, mantissa, decimalExponent, digits, ctx.Rounding)
}

// NaturalExponential is the long-form decimal.js alias for Exp.
func (x *Decimal) NaturalExponential() *Decimal {
	return x.Exp()
}

// Log returns log base b of x. Omitting b returns log base 10, matching
// Decimal#log. A nil base is treated like an omitted base.
func (x *Decimal) Log(base ...any) *Decimal {
	ctx := x.getContext()
	if len(base) > 1 {
		panic(newInvalidArgError(base))
	}
	logBase := ctx.NewFromInt64(10)
	if len(base) == 1 && base[0] != nil {
		logBase = mustDecimalArgument(ctx, base[0])
	}
	if logBase.IsNaN() || !logBase.IsFinite() || logBase.IsNeg() || logBase.IsZero() ||
		(logBase.e == 0 && len(logBase.d) == 1 && logBase.d[0] == 1) {
		return decimalNaN(ctx)
	}
	if x.IsZero() {
		return ctx.infinity(-1)
	}
	if x.IsNaN() || x.IsNeg() {
		return decimalNaN(ctx)
	}
	if !x.IsFinite() {
		return ctx.infinity(1)
	}
	if x.e == 0 && len(x.d) == 1 && x.d[0] == 1 {
		return ctx.NewFromInt64(0)
	}

	bits, digits := transcendentalPrecision(ctx.Precision)
	numerator := decimalLnBigFloat(x, bits)
	denominator := decimalLnBigFloat(logBase, bits)
	value := bigQuo(numerator, denominator, bits)
	return decimalFromBigFloat(ctx, value, 0, digits, ctx.Rounding)
}

// Logarithm is the long-form decimal.js alias for Log.
func (x *Decimal) Logarithm(base ...any) *Decimal {
	return x.Log(base...)
}

// Log2 returns the base-2 logarithm of x.
func (x *Decimal) Log2() *Decimal {
	return x.Log(2)
}

// Log10 returns the base-10 logarithm of x.
func (x *Decimal) Log10() *Decimal {
	return x.Log()
}

// Package and Context helpers mirror decimal.js's static methods while still
// allowing independent Go Context values.
func Ln(value any) *Decimal               { return defaultCtx.Ln(value) }
func Exp(value any) *Decimal              { return defaultCtx.Exp(value) }
func Log(value any, base ...any) *Decimal { return defaultCtx.Log(value, base...) }
func Log2(value any) *Decimal             { return defaultCtx.Log2(value) }
func Log10(value any) *Decimal            { return defaultCtx.Log10(value) }

func (ctx *Context) Ln(value any) *Decimal {
	return mustDecimalArgument(ctx, value).Ln()
}

func (ctx *Context) Exp(value any) *Decimal {
	return mustDecimalArgument(ctx, value).Exp()
}

func (ctx *Context) Log(value any, base ...any) *Decimal {
	return mustDecimalArgument(ctx, value).Log(base...)
}

func (ctx *Context) Log2(value any) *Decimal {
	return mustDecimalArgument(ctx, value).Log2()
}

func (ctx *Context) Log10(value any) *Decimal {
	return mustDecimalArgument(ctx, value).Log10()
}

func decimalNaN(ctx *Context) *Decimal {
	return &Decimal{s: 0, d: nil, e: 0, ctx: ctx}
}

func transcendentalPrecision(decimalDigits int) (uint, int) {
	digits := decimalDigits + 60
	if digits < 80 {
		digits = 80
	}
	return uint(digits*4 + 64), digits
}

// decimalLnBigFloat calculates ln(x) without constructing x as a binary
// float. x is split into a compact mantissa and 10^x.e, so inputs such as
// 1e+9000000000000000 remain inexpensive.
func decimalLnBigFloat(x *Decimal, bits uint) *big.Float {
	digits := digitsToStringExact(x.d)
	mantissaText := digits[:1]
	if len(digits) > 1 {
		mantissaText += "." + digits[1:]
	}
	mantissa, ok := new(big.Float).SetPrec(bits).SetMode(big.ToNearestEven).SetString(mantissaText)
	if !ok {
		return new(big.Float).SetPrec(bits)
	}
	value := bigLn(mantissa, bits)
	if x.e != 0 {
		exponent := new(big.Float).SetPrec(bits).SetMode(big.ToNearestEven).SetInt(big.NewInt(int64(x.e)))
		value = bigAdd(value, bigMul(exponent, bigLn10(bits), bits), bits)
	}
	return value
}

func bigLn10(bits uint) *big.Float {
	value, ok := new(big.Float).SetPrec(bits).SetMode(big.ToNearestEven).SetString(LN10)
	if !ok {
		panic(ErrPrecisionLimit)
	}
	return value
}

// splitFloor decomposes x into integer floor(x) and 0 <= fraction < 1.
func splitFloor(x *big.Float, bits uint) (*big.Int, *big.Float) {
	integer, _ := x.Int(nil)
	intFloat := new(big.Float).SetPrec(bits).SetMode(big.ToNearestEven).SetInt(integer)
	fraction := new(big.Float).SetPrec(bits).SetMode(big.ToNearestEven).Sub(x, intFloat)
	if fraction.Sign() < 0 {
		integer.Sub(integer, big.NewInt(1))
		fraction.Add(fraction, bigFloat(bits, 1))
	}
	return integer, fraction
}

// decimalFromBigFloat serialises a high-precision binary calculation through
// the Decimal parser, then applies the requested decimal.js rounding mode.
// exponent is a separate base-10 scale used by Exp.
func decimalFromBigFloat(ctx *Context, value *big.Float, exponent, digits int, rm RoundingMode) *Decimal {
	text := value.Text('e', digits-1)
	position := strings.LastIndexAny(text, "eE")
	if position < 0 {
		panic(newInvalidArgError(text))
	}
	parsedExponent := parseSimpleInt(text[position+1:])
	text = text[:position] + "e" + itoa(parsedExponent+exponent)
	result, err := ctx.New(text)
	if err != nil {
		panic(err)
	}
	return finalise(result, ctx.Precision, rm)
}

// powPositive calculates a positive finite base raised to a finite non-integer
// power. It uses ln(x) and exp(y*ln(x)) at guard precision, matching the
// algorithm used by decimal.js without reducing the operation to float64.
func powPositive(ctx *Context, x, y *Decimal, precision int, rounding RoundingMode) *Decimal {
	decimalDigits := precision + 40
	if decimalDigits < 60 {
		decimalDigits = 60
	}
	bits := uint(decimalDigits*4 + 32)

	xf, ok := decimalBigFloat(x, bits)
	if !ok {
		return nil
	}
	yf, ok := decimalBigFloat(y, bits)
	if !ok {
		return nil
	}

	z := bigMul(bigLn(xf, bits), yf, bits)
	rf := bigExp(z, bits)
	text := rf.Text('e', decimalDigits-1)
	r, err := ctx.New(text)
	if err != nil {
		return nil
	}
	r.s = 1
	return finalise(r, precision, rounding)
}

func decimalBigFloat(x *Decimal, precision uint) (*big.Float, bool) {
	if !x.IsFinite() {
		return nil, false
	}
	text := finiteToString(x, true, 0)
	if x.s < 0 {
		text = "-" + text
	}
	f, ok := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).SetString(text)
	return f, ok
}

func bigFloat(precision uint, value int64) *big.Float {
	return new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).SetInt64(value)
}

func bigAdd(x, y *big.Float, precision uint) *big.Float {
	return new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).Add(x, y)
}

func bigMul(x, y *big.Float, precision uint) *big.Float {
	return new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).Mul(x, y)
}

func bigQuo(x, y *big.Float, precision uint) *big.Float {
	return new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).Quo(x, y)
}

func bigSmall(x, epsilon *big.Float, precision uint) bool {
	abs := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).Abs(x)
	return abs.Cmp(epsilon) <= 0
}

func bigEpsilon(precision uint) *big.Float {
	one := bigFloat(precision, 1)
	return new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).SetMantExp(one, -int(precision)+8)
}

// bigLn computes ln(x) for x > 0. Binary exponent reduction keeps the
// atanh series argument in [0, 1/3], which converges rapidly at the
// precisions used by Decimal operations.
func bigLn(x *big.Float, precision uint) *big.Float {
	mantissa := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven)
	binaryExponent := x.MantExp(mantissa)
	// x = (2*mantissa) * 2^(binaryExponent-1), where 1 <= 2*mantissa < 2.
	scaled := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).SetMantExp(mantissa, 1)
	lnScaled := bigLnUnit(scaled, precision)
	if binaryExponent == 1 {
		return lnScaled
	}
	ln2 := bigLnUnit(bigFloat(precision, 2), precision)
	factor := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).SetInt(big.NewInt(int64(binaryExponent - 1)))
	return bigAdd(lnScaled, bigMul(ln2, factor, precision), precision)
}

// bigLnUnit computes ln(x) as 2*(t + t^3/3 + t^5/5 + ...), where
// t=(x-1)/(x+1). Callers reduce x so |t| <= 1/3.
func bigLnUnit(x *big.Float, precision uint) *big.Float {
	one := bigFloat(precision, 1)
	t := bigQuo(bigAdd(x, new(big.Float).SetPrec(precision).Neg(one), precision), bigAdd(x, one, precision), precision)
	tSquared := bigMul(t, t, precision)
	term := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).Set(t)
	sum := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).Set(t)
	epsilon := bigEpsilon(precision)

	for denominator := int64(3); denominator < 100000; denominator += 2 {
		term = bigMul(term, tSquared, precision)
		addend := bigQuo(term, bigFloat(precision, denominator), precision)
		sum = bigAdd(sum, addend, precision)
		if bigSmall(addend, epsilon, precision) {
			break
		}
	}
	return bigMul(sum, bigFloat(precision, 2), precision)
}

// bigExp computes e^x. The argument is divided by a power of two before a
// Taylor expansion, then restored by repeated squaring.
func bigExp(x *big.Float, precision uint) *big.Float {
	if x.Sign() == 0 {
		return bigFloat(precision, 1)
	}
	abs := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).Abs(x)
	k := abs.MantExp(nil) + 3
	if k < 0 {
		k = 0
	}
	reduced := new(big.Float).SetPrec(precision).SetMode(big.ToNearestEven).SetMantExp(x, -k)
	term := bigFloat(precision, 1)
	sum := bigFloat(precision, 1)
	epsilon := bigEpsilon(precision)

	for n := int64(1); n < 100000; n++ {
		term = bigQuo(bigMul(term, reduced, precision), bigFloat(precision, n), precision)
		sum = bigAdd(sum, term, precision)
		if bigSmall(term, epsilon, precision) {
			break
		}
	}
	for ; k > 0; k-- {
		sum = bigMul(sum, sum, precision)
	}
	return sum
}
