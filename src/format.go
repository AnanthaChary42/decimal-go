package decimal

import (
	"math"
	"strings"
)

// finiteToString converts a finite Decimal to a string.
// If isExp is true, uses exponential notation.
// Equivalent to finiteToString in decimal.js.
func finiteToString(x *Decimal, isExp bool) string {
	ctx := x.getContext()

	if x.d == nil {
		if x.s == 0 {
			return "NaN"
		}
		if x.s < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}

	str := digitsToStringExact(x.d)

	if isExp {
		// Exponential notation.
		var sb strings.Builder
		if x.s < 0 {
			sb.WriteByte('-')
		}
		sb.WriteByte(str[0])
		if len(str) > 1 {
			sb.WriteByte('.')
			sb.WriteString(str[1:])
		}
		sb.WriteByte('e')
		e := x.e
		if e >= 0 {
			sb.WriteByte('+')
		}
		sb.WriteString(itoa(e))
		return sb.String()
	}

	// Normal notation.
	e := x.e
	eIdx := e + 1
	strL := len(str)

	var sb strings.Builder
	if x.s < 0 {
		sb.WriteByte('-')
	}

	if e < 0 {
		// Fraction less than 1.
		sb.WriteString("0.")
		k := -eIdx
		if k > 0 {
			sb.WriteString(getZeroString(k))
		}
		sb.WriteString(str)
	} else if eIdx < strL {
		// Decimal point in the middle.
		sb.WriteString(str[:eIdx])
		sb.WriteByte('.')
		sb.WriteString(str[eIdx:])
	} else {
		// Integer with optional trailing zeros.
		sb.WriteString(str)
		k := eIdx - strL
		if k > 0 {
			sb.WriteString(getZeroString(k))
		}
	}

	_ = ctx
	return sb.String()
}

// String returns the string representation of x using default formatting.
// Equivalent to Decimal.prototype.toString() in decimal.js.
func (x *Decimal) String() string {
	ctx := x.getContext()

	if x.d == nil {
		if x.s == 0 {
			return "NaN"
		}
		if x.s < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}

	// Determine whether to use exponential notation.
	e := x.e
	isExp := e <= ctx.ToExpNeg || e >= ctx.ToExpPos
	return finiteToString(x, isExp)
}

// ValueOf returns the string value, matching JS valueOf().
func (x *Decimal) ValueOf() string {
	return x.String()
}

// ToFixed returns a string representing x in normal (fixed-point) notation with dp decimal places.
func (x *Decimal) ToFixed(dp int, rm ...RoundingMode) (string, error) {
	if err := checkInt32(dp, 0, MAX_PRECISION); err != nil {
		return "", err
	}

	ctx := x.getContext()
	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	r := finalise(x.copy(), dp+x.e+1, rounding)

	if r.d == nil {
		return r.String(), nil
	}

	str := finiteToString(r, false)

	// Append trailing zeros if needed.
	dotIdx := strings.IndexByte(str, '.')
	if dotIdx < 0 {
		if dp > 0 {
			str += "." + getZeroString(dp)
		}
	} else {
		actualDp := len(str) - dotIdx - 1
		if actualDp < dp {
			str += getZeroString(dp - actualDp)
		}
	}

	return str, nil
}

// ToExponential returns a string representing x in exponential notation with dp decimal places.
func (x *Decimal) ToExponential(dp int, rm ...RoundingMode) (string, error) {
	ctx := x.getContext()
	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	r := finalise(x.copy(), dp+1, rounding)

	if r.d == nil {
		return r.String(), nil
	}

	str := finiteToString(r, true)

	// Format to exact dp decimal places.
	if dp >= 0 {
		eIdx := strings.IndexAny(str, "eE")
		head := str[:eIdx]
		tail := str[eIdx:]

		dotIdx := strings.IndexByte(head, '.')
		if dotIdx < 0 {
			if dp > 0 {
				head += "." + getZeroString(dp)
			}
		} else {
			actualDp := len(head) - dotIdx - 1
			if actualDp < dp {
				head += getZeroString(dp - actualDp)
			}
		}
		str = head + tail
	}

	return str, nil
}

// ToPrecision returns a string representing x to sd significant digits.
func (x *Decimal) ToPrecision(sd int, rm ...RoundingMode) (string, error) {
	if err := checkInt32(sd, 1, MAX_PRECISION); err != nil {
		return "", err
	}

	ctx := x.getContext()
	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	r := finalise(x.copy(), sd, rounding)

	if r.d == nil {
		return r.String(), nil
	}

	e := r.e
	isExp := e < ctx.ToExpNeg || e >= ctx.ToExpPos
	return finiteToString(r, isExp), nil
}

// Float64 returns the float64 representation of x and whether the conversion was exact.
func (x *Decimal) Float64() (float64, bool) {
	if x.d == nil {
		if x.s == 0 {
			return math.NaN(), false
		}
		if x.s < 0 {
			return negInf(), false
		}
		return posInf(), false
	}

	s := x.String()
	var f float64
	exact, _ := parseFloat(s, &f)
	return f, exact
}

// Helper: float64 conversion helper.
func parseFloat(s string, out *float64) (bool, error) {
	// Parse simple float string.
	neg := false
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}

	var val float64
	eIdx := strings.IndexAny(s, "eE")
	if eIdx >= 0 {
		// Exponential.
		baseStr := s[:eIdx]
		expVal := parseSimpleInt(s[eIdx+1:])
		bVal := 0.0
		parseFloat(baseStr, &bVal)
		val = bVal * math.Pow(10, float64(expVal))
	} else {
		// Normal.
		dotIdx := strings.IndexByte(s, '.')
		if dotIdx >= 0 {
			intPart := s[:dotIdx]
			fracPart := s[dotIdx+1:]

			var intVal float64
			for _, ch := range intPart {
				intVal = intVal*10 + float64(ch-'0')
			}

			var fracVal float64
			pow := 0.1
			for _, ch := range fracPart {
				fracVal += float64(ch-'0') * pow
				pow /= 10
			}

			val = intVal + fracVal
		} else {
			for _, ch := range s {
				val = val*10 + float64(ch-'0')
			}
		}
	}

	if neg {
		val = -val
	}
	*out = val
	return true, nil
}

func posInf() float64 {
	return math.Inf(1)
}

func negInf() float64 {
	return math.Inf(-1)
}
