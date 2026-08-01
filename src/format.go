package decimal

import (
	"math"
	"strings"
)

// finiteToString converts a finite Decimal to a string.
// If isExp is true, uses exponential notation.
// sd is the number of digits to pad to (0 means no padding).
func finiteToString(x *Decimal, isExp bool, sd int) string {
	if !x.IsFinite() {
		return nonFiniteToString(x)
	}

	e := x.e
	str := digitsToStringExact(x.d)
	strLen := len(str)

	if isExp {
		if sd > 0 {
			k := sd - strLen
			if k > 0 {
				if strLen > 1 {
					str = string(str[0]) + "." + str[1:] + getZeroString(k)
				} else {
					str = string(str[0]) + "." + getZeroString(k)
				}
			} else if strLen > 1 {
				str = string(str[0]) + "." + str[1:]
			}
		} else if strLen > 1 {
			str = string(str[0]) + "." + str[1:]
		}

		if e < 0 {
			str = str + "e" + itoa(e)
		} else {
			str = str + "e+" + itoa(e)
		}
	} else if e < 0 {
		str = "0." + getZeroString(-e-1) + str
		if sd > 0 {
			k := sd - strLen
			if k > 0 {
				str += getZeroString(k)
			}
		}
	} else if e >= strLen {
		str += getZeroString(e + 1 - strLen)
		if sd > 0 {
			k := sd - e - 1
			if k > 0 {
				str = str + "." + getZeroString(k)
			}
		}
	} else {
		k := e + 1
		if k < strLen {
			str = str[:k] + "." + str[k:]
		}
		if sd > 0 {
			k2 := sd - strLen
			if k2 > 0 {
				if e+1 == strLen {
					str += "."
				}
				str += getZeroString(k2)
			}
		}
	}

	return str
}

// nonFiniteToString returns the string for NaN or ±Infinity.
func nonFiniteToString(x *Decimal) string {
	if x.s == 0 {
		return "NaN"
	}
	if x.s < 0 {
		return "-Infinity"
	}
	return "Infinity"
}

// String returns the string representation of x.
// Implements fmt.Stringer.
// Uses exponential notation when e <= toExpNeg or e >= toExpPos.
func (x *Decimal) String() string {
	if !x.IsFinite() {
		return nonFiniteToString(x)
	}

	ctx := x.getContext()
	str := finiteToString(x, x.e <= ctx.ToExpNeg || x.e >= ctx.ToExpPos, 0)

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ValueOf returns a string representation where negative zero includes the minus sign.
// Equivalent to valueOf/toJSON in decimal.js.
func (x *Decimal) ValueOf() string {
	if !x.IsFinite() {
		return nonFiniteToString(x)
	}

	ctx := x.getContext()
	str := finiteToString(x, x.e <= ctx.ToExpNeg || x.e >= ctx.ToExpPos, 0)

	if x.IsNeg() {
		return "-" + str
	}
	return str
}

// ToFixed returns a string representing x in fixed-point notation,
// rounded to dp decimal places using the specified rounding mode.
func (x *Decimal) ToFixed(dp int, rm ...RoundingMode) string {
	if !x.IsFinite() {
		return nonFiniteToString(x)
	}

	ctx := x.getContext()
	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	var str string
	if dp < 0 {
		str = finiteToString(x, false, 0)
	} else {
		y := x.copy()
		finalise(y, dp+y.e+1, rounding)
		str = finiteToString(y, false, dp+y.e+1)
	}

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToExponential returns a string in exponential notation,
// rounded to dp decimal places using the specified rounding mode.
func (x *Decimal) ToExponential(dp int, rm ...RoundingMode) string {
	if !x.IsFinite() {
		s := nonFiniteToString(x)
		if x.IsNeg() && !x.IsZero() {
			return "-" + s
		}
		return s
	}

	ctx := x.getContext()
	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	var str string
	if dp < 0 {
		str = finiteToString(x, true, 0)
	} else {
		y := x.copy()
		finalise(y, dp+1, rounding)
		str = finiteToString(y, true, dp+1)
	}

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToPrecision returns a string rounded to sd significant digits,
// using exponential notation if sd is less than the integer part length.
func (x *Decimal) ToPrecision(sd int, rm ...RoundingMode) string {
	if !x.IsFinite() {
		s := nonFiniteToString(x)
		if x.IsNeg() && !x.IsZero() {
			return "-" + s
		}
		return s
	}

	ctx := x.getContext()
	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	var str string
	if sd <= 0 {
		str = finiteToString(x, x.e <= ctx.ToExpNeg || x.e >= ctx.ToExpPos, 0)
	} else {
		y := x.copy()
		finalise(y, sd, rounding)
		str = finiteToString(y, sd <= y.e || y.e <= ctx.ToExpNeg, sd)
	}

	if x.IsNeg() && !x.IsZero() {
		return "-" + str
	}
	return str
}

// ToNumber returns the float64 representation of x.
func (x *Decimal) ToNumber() float64 {
	if x.IsNaN() {
		return 0 // Go doesn't have NaN as easily usable
	}
	s := x.String()
	f := 0.0
	// Simple conversion via string.
	// Use strings to number conversion.
	_, _ = parseFloat(s, &f)
	return f
}

// parseFloat is a simple string to float64 converter.
func parseFloat(s string, f *float64) (bool, error) {
	s = strings.TrimSpace(s)
	val := 0.0
	neg := false
	i := 0

	if i < len(s) && s[i] == '-' {
		neg = true
		i++
	} else if i < len(s) && s[i] == '+' {
		i++
	}

	if s == "Infinity" || s == "+Infinity" {
		*f = posInf()
		return true, nil
	}
	if s == "-Infinity" {
		*f = negInf()
		return true, nil
	}

	// Integer part.
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		val = val*10 + float64(s[i]-'0')
		i++
	}

	// Fraction part.
	if i < len(s) && s[i] == '.' {
		i++
		frac := 0.1
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			val += float64(s[i]-'0') * frac
			frac /= 10
			i++
		}
	}

	// Exponent part.
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		expNeg := false
		if i < len(s) && s[i] == '-' {
			expNeg = true
			i++
		} else if i < len(s) && s[i] == '+' {
			i++
		}
		exp := 0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			exp = exp*10 + int(s[i]-'0')
			i++
		}
		if expNeg {
			exp = -exp
		}
		p := 1.0
		if exp >= 0 {
			for j := 0; j < exp; j++ {
				p *= 10
			}
			val *= p
		} else {
			for j := 0; j < -exp; j++ {
				p *= 10
			}
			val /= p
		}
	}

	if neg {
		val = -val
	}
	*f = val
	return true, nil
}

func posInf() float64 {
	return math.Inf(1)
}

func negInf() float64 {
	return math.Inf(-1)
}
