package decimal

import (
	"math"
	"regexp"
	"strings"
)

// Regex patterns for parsing, matching decimal.js patterns.
var (
	reDecimal = regexp.MustCompile(`^(\d+(\.\d*)?|\.\d+)(e[+-]?\d+)?$`)
	reBinary  = regexp.MustCompile(`(?i)^0b([01]+(\.[01]*)?|\.[01]+)(p[+-]?\d+)?$`)
	reHex     = regexp.MustCompile(`(?i)^0x([0-9a-f]+(\.[0-9a-f]*)?|\.[0-9a-f]+)(p[+-]?\d+)?$`)
	reOctal   = regexp.MustCompile(`(?i)^0o([0-7]+(\.[0-7]*)?|\.[0-7]+)(p[+-]?\d+)?$`)
)

// isDecimalRegex tests whether a string (sign already stripped) matches a decimal number.
func isDecimalRegex(s string) bool {
	return reDecimal.MatchString(s)
}

// parseDecimalStr parses a decimal string (sign already set on x).
// Equivalent to parseDecimal in decimal.js.
func parseDecimalStr(x *Decimal, str string) {
	var e int

	// Decimal point?
	dotIdx := strings.IndexByte(str, '.')
	if dotIdx > -1 {
		str = str[:dotIdx] + str[dotIdx+1:]
	}

	// Exponential form?
	eIdx := strings.IndexAny(str, "eE")
	if eIdx > 0 {
		// Determine exponent.
		if dotIdx < 0 {
			dotIdx = eIdx
		}
		expStr := str[eIdx+1:]
		expVal := 0
		sign := 1
		if len(expStr) > 0 {
			switch expStr[0] {
			case '+':
				expStr = expStr[1:]
			case '-':
				sign = -1
				expStr = expStr[1:]
			}
		}
		for _, ch := range expStr {
			expVal = expVal*10 + int(ch-'0')
		}
		e = dotIdx + sign*expVal
		str = str[:eIdx]
	} else if dotIdx < 0 {
		// Integer.
		e = len(str)
	} else {
		e = dotIdx
	}

	// Determine leading zeros.
	i := 0
	for i < len(str) && str[i] == '0' {
		i++
	}

	// Determine trailing zeros.
	strLen := len(str)
	for strLen > i && str[strLen-1] == '0' {
		strLen--
	}

	str = str[i:strLen]

	if len(str) > 0 {
		strLen = len(str)
		x.e = e - i - 1
		x.d = make([]int32, 0)

		// Transform base.
		// e is the base 10 exponent.
		// i is where to slice str to get the first word of the digits array.
		sliceIdx := (x.e + 1) % LOG_BASE
		if x.e < 0 {
			sliceIdx += LOG_BASE
		}

		if sliceIdx < strLen {
			if sliceIdx > 0 {
				x.d = append(x.d, parseIntSlice(str[:sliceIdx]))
			}
			for strLen -= LOG_BASE; sliceIdx < strLen; {
				x.d = append(x.d, parseIntSlice(str[sliceIdx:sliceIdx+LOG_BASE]))
				sliceIdx += LOG_BASE
			}
			str = str[sliceIdx:]
			sliceIdx = LOG_BASE - len(str)
		} else {
			sliceIdx -= strLen
		}

		for sliceIdx > 0 {
			str += "0"
			sliceIdx--
		}

		x.d = append(x.d, parseIntSlice(str))

		if external {
			ctx := x.getContext()
			// Overflow?
			if x.e > ctx.MaxE {
				x.d = nil
				x.e = 0
			} else if x.e < ctx.MinE {
				// Underflow → zero.
				x.e = 0
				x.d = []int32{0}
			}
		}
	} else {
		// Zero.
		x.e = 0
		x.d = []int32{0}
	}
}

// parseOther handles non-decimal strings: hex, binary, octal, Infinity, NaN, underscore separators.
// Equivalent to parseOther in decimal.js.
func parseOther(x *Decimal, str string) error {
	// Handle underscore separators.
	if strings.Contains(str, "_") {
		// Remove underscores between digits.
		cleaned := strings.ReplaceAll(str, "_", "")
		if isDecimalRegex(cleaned) {
			parseDecimalStr(x, cleaned)
			return nil
		}
		str = cleaned
	}

	if str == "Infinity" || str == "NaN" {
		if str == "NaN" {
			x.s = 0
		}
		x.e = 0
		x.d = nil
		return nil
	}

	var base int
	if reHex.MatchString(str) {
		base = 16
		str = strings.ToLower(str)
	} else if reBinary.MatchString(str) {
		base = 2
	} else if reOctal.MatchString(str) {
		base = 8
	} else {
		return newInvalidArgError(str)
	}

	// Is there a binary exponent part?
	pIdx := strings.IndexAny(str, "pP")
	var p int
	hasP := false
	if pIdx > 0 {
		pStr := str[pIdx+1:]
		p = parseSimpleInt(pStr)
		hasP = true
		str = str[2:pIdx]
	} else {
		str = str[2:]
	}

	// Convert to integer then handle fraction.
	dotIdx := strings.IndexByte(str, '.')
	isFloat := dotIdx >= 0
	ctx := x.getContext()

	if isFloat {
		str = str[:dotIdx] + str[dotIdx+1:]
	}

	xd := convertBase(str, base, int(BASE))
	xe := len(xd) - 1

	// Remove trailing zeros.
	for xe >= 0 && xd[xe] == 0 {
		xd = xd[:xe]
		xe--
	}

	if xe < 0 {
		x2 := ctx.NewFromInt64(0)
		x2.s = x.s
		x.e = x2.e
		x.d = x2.d
		return nil
	}

	x.e = getBase10Exponent(xd, xe)
	x.d = xd
	external = false

	if isFloat {
		fracLen := len(str) - dotIdx
		divisor := intPow(ctx, ctx.NewFromInt64(int64(base)), fracLen, fracLen*2)
		external = false // intPow resets external to true; restore it for divide
		result := divide(x, divisor, len(str)*4, RoundingMode(0), false, 0)
		x.d = result.d
		x.e = result.e
		x.s = result.s
	}

	if hasP {
		var powResult *Decimal
		absP := p
		if absP < 0 {
			absP = -absP
		}
		if absP < 54 {
			powResult, _ = ctx.NewFromFloat64(math.Pow(2, float64(p)))
		} else {
			powResult = ctx.NewFromInt64(2).Pow(ctx.NewFromInt64(int64(p)))
		}
		result := x.Times(powResult)
		x.d = result.d
		x.e = result.e
	}

	external = true
	return nil
}

// parseIntSlice parses a string of digits to int32.
func parseIntSlice(s string) int32 {
	var result int32
	for _, ch := range s {
		result = result*10 + int32(ch-'0')
	}
	return result
}

// parseSimpleInt parses a simple integer string (possibly with sign).
func parseSimpleInt(s string) int {
	neg := false
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	} else if len(s) > 0 && s[0] == '+' {
		s = s[1:]
	}
	result := 0
	for _, ch := range s {
		result = result*10 + int(ch-'0')
	}
	if neg {
		return -result
	}
	return result
}
