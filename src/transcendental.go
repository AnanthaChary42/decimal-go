package decimal

import "math/big"

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
