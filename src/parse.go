package decimal

import (
	"math"
	"regexp"
	"strings"
)

// Regex patterns for parsing, matching decimal.js patterns.
var (
	reDecimal = regexp.MustCompile(`(?i)^(\d+(\.\d*)?|\.\d+)(e[+-]?\d+)?$`)
	reBinary  = regexp.MustCompile(`(?i)^0b([01]+(\.[01]*)?|\.[01]+)(p[+-]?\d+)?$`)
	reHex     = regexp.MustCompile(`(?i)^0x([0-9a-f]+(\.[0-9a-f]*)?|\.[0-9a-f]+)(p[+-]?\d+)?$`)
	reOctal   = regexp.MustCompile(`(?i)^0o([0-7]+(\.[0-7]*)?|\.[0-7]+)(p[+-]?\d+)?$`)
)

func isDecimalRegex(s string) bool {
	return reDecimal.MatchString(s)
}

// parseDecimalStr parses a standard decimal string representation.
// Equivalent to parseDecimal in decimal.js.
func parseDecimalStr(x *Decimal, s string) {
	var e int

	// Find decimal point position and scientific notation exponent.
	dotIdx := strings.IndexByte(s, '.')
	eIdx := strings.IndexAny(s, "eE")

	if eIdx >= 0 {
		if dotIdx < 0 {
			dotIdx = eIdx
		}
		dotIdx += parseSimpleInt(s[eIdx+1:])
		s = s[:eIdx]
	} else if dotIdx < 0 {
		dotIdx = len(s)
	}

	// Remove the decimal point from the string.
	actualDot := strings.IndexByte(s, '.')
	if actualDot >= 0 {
		s = s[:actualDot] + s[actualDot+1:]
	}

	// Determine leading zeros.
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}

	// Determine trailing zeros.
	sLen := len(s)
	for sLen > i && s[sLen-1] == '0' {
		sLen--
	}
	s = s[i:sLen]

	if len(s) == 0 {
		// Zero.
		x.e = 0
		x.d = []int32{0}
		return
	}

	strLen := len(s)
	e = dotIdx - i - 1 // base-10 exponent
	x.e = e

	// Calculate first word length based on exponent alignment.
	// This is the critical step that ensures words are aligned to LOG_BASE boundaries.
	firstLen := (e + 1) % LOG_BASE
	if e < 0 {
		firstLen += LOG_BASE
	}
	if firstLen < 0 {
		firstLen += LOG_BASE
	}

	var digits []int32

	if firstLen < strLen {
		// First word (partial, only if firstLen > 0).
		if firstLen > 0 {
			digits = append(digits, parseIntSlice(s[:firstLen]))
		}

		// Middle words (full LOG_BASE chunks).
		ii := firstLen
		limit := strLen - LOG_BASE
		for ii < limit {
			digits = append(digits, parseIntSlice(s[ii:ii+LOG_BASE]))
			ii += LOG_BASE
		}

		// Last word: remaining digits padded with trailing zeros.
		s = s[ii:]
		pad := LOG_BASE - len(s)
		for pad > 0 {
			s += "0"
			pad--
		}
		digits = append(digits, parseIntSlice(s))
	} else {
		// The entire string fits in one word, pad with trailing zeros.
		pad := firstLen - strLen
		for pad > 0 {
			s += "0"
			pad--
		}
		digits = append(digits, parseIntSlice(s))
	}

	x.d = digits

	// Round if precision limit exceeded.
	ctx := x.getContext()
	finalise(x, ctx.Precision, ctx.Rounding)
}

// parseOther parses non-standard decimal strings (hex, binary, octal, NaN, Infinity).
func parseOther(x *Decimal, str string) error {
	strLower := strings.ToLower(str)

	// Infinity?
	if strLower == "infinity" || strLower == "inf" {
		x.d = nil
		x.e = 0
		return nil
	}

	// NaN?
	if strLower == "nan" {
		x.s = 0
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
