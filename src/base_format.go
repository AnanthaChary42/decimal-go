package decimal

import (
	"math/big"
	"strings"
)

// ToBinary returns x in base 2. With an explicit significant-digit argument,
// it uses binary exponential notation (for example 0b1.8 is not applicable;
// binary uses forms such as 0b1.1p+3).
func (x *Decimal) ToBinary(args ...any) string {
	return toBaseString(x, 2, args...)
}

// ToOctal returns x in base 8. With an explicit significant-digit argument it
// uses a binary exponent suffix, as decimal.js does.
func (x *Decimal) ToOctal(args ...any) string {
	return toBaseString(x, 8, args...)
}

// ToHex returns x in base 16. ToHexadecimal is provided as the decimal.js
// alias for callers that use the long method name.
func (x *Decimal) ToHex(args ...any) string {
	return toBaseString(x, 16, args...)
}

func (x *Decimal) ToHexadecimal(args ...any) string {
	return x.ToHex(args...)
}

func toBaseString(x *Decimal, baseOut int, args ...any) string {
	if len(args) > 2 {
		panic(newInvalidArgError(args))
	}
	isExp := len(args) > 0 && args[0] != nil
	ctx := x.getContext()
	sd := ctx.Precision
	rm := ctx.Rounding
	if isExp {
		sd = positiveIntArgument(args[0])
		if len(args) == 2 {
			rm = roundingArgument(args[1])
		}
	}

	if !x.IsFinite() {
		return nonFiniteToString(x)
	}
	prefix := map[int]string{2: "0b", 8: "0o", 16: "0x"}[baseOut]
	if x.IsZero() {
		zero := "0"
		if isExp {
			zero += "p+0"
		}
		if x.IsNeg() {
			return "-" + prefix + zero
		}
		return prefix + zero
	}

	workingBase := baseOut
	workingDigits := sd
	if isExp && baseOut != 2 {
		workingBase = 2
		if baseOut == 8 {
			workingDigits = sd*3 - 2
		} else {
			workingDigits = sd*4 - 3
		}
	}

	digits, exponent := baseDigitsRounded(x, workingBase, workingDigits, rm)
	var body string
	if isExp {
		if baseOut != 2 {
			body = groupedBinaryMantissa(digits, baseOut)
		} else {
			body = digits[:1]
			if len(digits) > 1 {
				body += "." + digits[1:]
			}
		}
		if exponent < 0 {
			body += "p" + itoa(exponent)
		} else {
			body += "p+" + itoa(exponent)
		}
	} else {
		body = normalBaseMantissa(digits, exponent)
	}
	if x.IsNeg() {
		return "-" + prefix + body
	}
	return prefix + body
}

func positiveIntArgument(value any) int {
	var n int
	switch v := value.(type) {
	case int:
		n = v
	case int8:
		n = int(v)
	case int16:
		n = int(v)
	case int32:
		n = int(v)
	case int64:
		if v > int64(MAX_DIGITS) || v < 1 {
			panic(newInvalidArgError(value))
		}
		n = int(v)
	default:
		panic(newInvalidArgError(value))
	}
	if n < 1 || n > MAX_DIGITS {
		panic(newInvalidArgError(value))
	}
	return n
}

// baseDigitsRounded returns a normalised positive significand in base and the
// base exponent. The returned digits have no trailing zeroes and are rounded
// to sd significant digits with decimal.js rounding rules.
func baseDigitsRounded(x *Decimal, base, sd int, rm RoundingMode) (string, int) {
	numerator, denominator := decimalRational(x)
	exponent := floorLogBase(numerator, denominator, base)
	scale := sd - 1 - exponent
	scaledNumerator := new(big.Int).Set(numerator)
	scaledDenominator := new(big.Int).Set(denominator)
	if scale >= 0 {
		scaledNumerator.Mul(scaledNumerator, bigPowBase(base, scale))
	} else {
		scaledDenominator.Mul(scaledDenominator, bigPowBase(base, -scale))
	}
	q, remainder := new(big.Int), new(big.Int)
	q.QuoRem(scaledNumerator, scaledDenominator, remainder)
	if baseRoundUp(q, remainder, scaledDenominator, x.s, rm) {
		q.Add(q, big.NewInt(1))
	}
	limit := bigPowBase(base, sd)
	if q.Cmp(limit) >= 0 {
		q.Quo(q, big.NewInt(int64(base)))
		exponent++
	}
	digits := q.Text(base)
	digits = strings.TrimRight(digits, "0")
	if digits == "" {
		digits = "0"
	}
	return digits, exponent
}

func baseRoundUp(q, remainder, denominator *big.Int, sign int8, rm RoundingMode) bool {
	if remainder.Sign() == 0 {
		return false
	}
	switch rm {
	case RoundUp:
		return true
	case RoundCeil:
		return sign > 0
	case RoundFloor:
		return sign < 0
	case RoundHalfUp, RoundHalfDown, RoundHalfEven, RoundHalfCeil, RoundHalfFloor:
		comparison := new(big.Int).Lsh(remainder, 1).Cmp(denominator)
		if comparison > 0 {
			return true
		}
		if comparison < 0 {
			return false
		}
		switch rm {
		case RoundHalfUp:
			return true
		case RoundHalfEven:
			return q.Bit(0) == 1
		case RoundHalfCeil:
			return sign > 0
		case RoundHalfFloor:
			return sign < 0
		}
	}
	return false
}

func floorLogBase(numerator, denominator *big.Int, base int) int {
	// A bit-length estimate makes the correction loop constant-time for the
	// bases used here, while exact comparisons keep the result mathematical.
	log2Base := map[int]int{2: 1, 8: 3, 16: 4}[base]
	exponent := (numerator.BitLen() - denominator.BitLen()) / log2Base
	if compareToBasePower(numerator, denominator, base, exponent) < 0 {
		for compareToBasePower(numerator, denominator, base, exponent) < 0 {
			exponent--
		}
	} else {
		for compareToBasePower(numerator, denominator, base, exponent+1) >= 0 {
			exponent++
		}
	}
	return exponent
}

// compareToBasePower compares numerator/denominator with base^exponent.
func compareToBasePower(numerator, denominator *big.Int, base, exponent int) int {
	if exponent >= 0 {
		return numerator.Cmp(new(big.Int).Mul(denominator, bigPowBase(base, exponent)))
	}
	return new(big.Int).Mul(numerator, bigPowBase(base, -exponent)).Cmp(denominator)
}

func bigPowBase(base, exponent int) *big.Int {
	if exponent == 0 {
		return big.NewInt(1)
	}
	return new(big.Int).Exp(big.NewInt(int64(base)), big.NewInt(int64(exponent)), nil)
}

func normalBaseMantissa(digits string, exponent int) string {
	if exponent < 0 {
		return "0." + strings.Repeat("0", -exponent-1) + digits
	}
	position := exponent + 1
	if position >= len(digits) {
		return digits + strings.Repeat("0", position-len(digits))
	}
	return digits[:position] + "." + digits[position:]
}

func groupedBinaryMantissa(binaryDigits string, baseOut int) string {
	if len(binaryDigits) == 1 {
		return "1"
	}
	groupSize := 3
	if baseOut == 16 {
		groupSize = 4
	}
	fraction := binaryDigits[1:]
	if remainder := len(fraction) % groupSize; remainder != 0 {
		fraction += strings.Repeat("0", groupSize-remainder)
	}
	var out strings.Builder
	out.WriteString("1.")
	for i := 0; i < len(fraction); i += groupSize {
		value := 0
		for _, ch := range fraction[i : i+groupSize] {
			value = value*2 + int(ch-'0')
		}
		out.WriteByte(NUMERALS[value])
	}
	return strings.TrimRight(out.String(), "0")
}
