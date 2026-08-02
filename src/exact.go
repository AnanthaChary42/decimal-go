package decimal

import "math/big"

// divideFiniteExact calculates a finite, non-zero decimal quotient using an
// integer significand. The decimal exponent is tracked separately, so extreme
// exponents never require allocating strings of zeros.
//
// The word-array division routine is still retained for base conversion. For
// ordinary decimal division this path avoids losing precision at base-1e7 word
// boundaries before finalise applies the requested decimal rounding mode.
func divideFiniteExact(x, y *Decimal, sign int8, sd int, rm RoundingMode, dp bool, pr int) *Decimal {
	ctx := x.getContext()
	xText := digitsToStringExact(x.d)
	yText := digitsToStringExact(y.d)

	var xn, yn big.Int
	xn.SetString(xText, 10)
	yn.SetString(yText, 10)

	xDigits := len(xText)
	yDigits := len(yText)

	// Find floor(log10(xn / yn)) without expanding either decimal exponent.
	ratioExponent := xDigits - yDigits
	var left, right big.Int
	left.Set(&xn)
	right.Set(&yn)
	if xDigits < yDigits {
		left.Mul(&left, bigPow10(yDigits-xDigits))
	} else if xDigits > yDigits {
		right.Mul(&right, bigPow10(xDigits-yDigits))
	}
	if left.Cmp(&right) < 0 {
		ratioExponent--
	}

	// x = xn * 10^(x.e - xDigits + 1), and similarly for y.
	// Therefore this is the base-10 exponent of the first quotient digit.
	resultExponent := ratioExponent + (x.e - xDigits + 1) - (y.e - yDigits + 1)

	// Retain one extra significant digit and the remainder flag. finalise uses
	// both to implement all nine decimal.js rounding modes correctly.
	keptDigits := sd + 1
	scale := keptDigits - 1 - ratioExponent
	var numerator, denominator big.Int
	numerator.Set(&xn)
	denominator.Set(&yn)
	if scale >= 0 {
		numerator.Mul(&numerator, bigPow10(scale))
	} else {
		denominator.Mul(&denominator, bigPow10(-scale))
	}

	var quotient, remainder big.Int
	quotient.QuoRem(&numerator, &denominator, &remainder)

	// The coefficient-array word boundaries depend on e modulo LOG_BASE.
	// Constructing the integer coefficient and then assigning q.e would violate
	// that invariant for most exponents. Parse it in exponential form instead
	// so the parser performs the required base-1e7 alignment.
	quotientText := quotient.String()
	quoted := quotientText[:1]
	if len(quotientText) > 1 {
		quoted += "." + quotientText[1:]
	}
	quoted += "e" + itoa(resultExponent)
	q, err := ctx.New(quoted)
	if err != nil {
		// quotient is constructed from a positive integer, so this is unreachable.
		panic(err)
	}
	q.s = sign
	q.e = resultExponent
	// The original division algorithm performs a second finalisation when the
	// requested precision is decimal places rather than significant digits.
	// DivToInt and Mod use this path with pr == 0; omitting it leaves a
	// fractional quotient (for example, 1 divToInt 3 becomes 0.3).
	precision := sd
	if dp {
		precision = pr + q.e + 1
	}
	return finalise(q, precision, rm, remainder.Sign() != 0)
}

func bigPow10(n int) *big.Int {
	if n == 0 {
		return big.NewInt(1)
	}
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
}
