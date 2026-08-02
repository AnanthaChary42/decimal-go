package decimal

import (
	"errors"
	"math"
	"strconv"
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
		return math.NaN()
	}
	if !x.IsFinite() {
		if x.IsNeg() {
			return math.Inf(-1)
		}
		return math.Inf(1)
	}
	if x.IsZero() {
		if x.IsNeg() {
			return math.Copysign(0, -1)
		}
		return 0
	}

	s := x.ValueOf()
	f, err := strconv.ParseFloat(s, 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return 0
	}
	return f
}

