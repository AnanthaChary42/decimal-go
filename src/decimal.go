package decimal

import (
	"math"
	"strconv"
	"strings"
)

// RoundingMode specifies how rounding is performed.
type RoundingMode int

const (
	RoundUp        RoundingMode = 0 // Away from zero.
	RoundDown      RoundingMode = 1 // Towards zero.
	RoundCeil      RoundingMode = 2 // Towards +Infinity.
	RoundFloor     RoundingMode = 3 // Towards -Infinity.
	RoundHalfUp    RoundingMode = 4 // Towards nearest neighbour. If equidistant, up.
	RoundHalfDown  RoundingMode = 5 // Towards nearest neighbour. If equidistant, down.
	RoundHalfEven  RoundingMode = 6 // Towards nearest neighbour. If equidistant, towards even.
	RoundHalfCeil  RoundingMode = 7 // Towards nearest neighbour. If equidistant, towards +Infinity.
	RoundHalfFloor RoundingMode = 8 // Towards nearest neighbour. If equidistant, towards -Infinity.
	Euclid         RoundingMode = 9 // Euclidian division.
)

// Context holds configuration for decimal operations.
// Equivalent to the Decimal constructor's configuration in decimal.js.
type Context struct {
	Precision int          // Maximum number of significant digits. Default: 20.
	Rounding  RoundingMode // Rounding mode. Default: RoundHalfUp (4).
	ToExpNeg  int          // Exponent at/below which toString uses exponential notation. Default: -7.
	ToExpPos  int          // Exponent at/above which toString uses exponential notation. Default: 21.
	MinE      int          // Minimum exponent (underflow to zero). Default: -EXP_LIMIT.
	MaxE      int          // Maximum exponent (overflow to Infinity). Default: EXP_LIMIT.
	Modulo    RoundingMode // Modulo mode. Default: RoundDown (1).
}

// DefaultContext returns a new Context with default settings matching decimal.js defaults.
func DefaultContext() *Context {
	return &Context{
		Precision: 20,
		Rounding:  RoundHalfUp,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -EXP_LIMIT,
		MaxE:      EXP_LIMIT,
		Modulo:    RoundDown,
	}
}

// Decimal represents an arbitrary-precision decimal number.
// Internal representation matches decimal.js:
//   - d: digit array in base 1e7 (nil for NaN/Infinity)
//   - e: base-10 exponent of the most significant digit
//   - s: sign (-1 or 1; 0 for NaN)
type Decimal struct {
	d   []int32  // digit array, base 1e7; nil = NaN or Infinity
	e   int      // exponent
	s   int8     // sign: 1, -1, or 0 (NaN)
	ctx *Context // configuration context
}

// external controls whether finalise should check for overflow/underflow.
// Matches the `external` variable in decimal.js. Some internal operations
// temporarily set this to false.
var external = true

// Package-level default context used when no context is specified.
var defaultCtx = DefaultContext()

// ---- Constructors ----

// New creates a new Decimal from a string using the default context.
// Accepts decimal ("123.456"), exponential ("1.23e+5"), hex ("0xff"),
// binary ("0b101"), octal ("0o77"), "Infinity", "NaN", and underscore separators.
func New(s string) (*Decimal, error) {
	return defaultCtx.New(s)
}

// New creates a new Decimal from a string using this context.
func (ctx *Context) New(s string) (*Decimal, error) {
	x := &Decimal{ctx: ctx}
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, newInvalidArgError(s)
	}

	// Handle sign.
	if s[0] == '-' {
		x.s = -1
		s = s[1:]
	} else {
		if s[0] == '+' {
			s = s[1:]
		}
		x.s = 1
	}

	if isDecimalRegex(s) {
		parseDecimalStr(x, s)
	} else {
		err := parseOther(x, s)
		if err != nil {
			return nil, err
		}
	}

	return x, nil
}

// NewFromInt64 creates a Decimal from an int64 using the default context.
func NewFromInt64(v int64) *Decimal {
	return defaultCtx.NewFromInt64(v)
}

// NewFromInt64 creates a Decimal from an int64 using this context.
func (ctx *Context) NewFromInt64(v int64) *Decimal {
	x := &Decimal{ctx: ctx}

	if v == 0 {
		x.s = 1
		x.e = 0
		x.d = []int32{0}
		return x
	}

	if v < 0 {
		x.s = -1
		if v == math.MinInt64 {
			// -MinInt64 overflows in two's complement; handle via string parsing
			parseDecimalStr(x, "9223372036854775808")
			return x
		}
		v = -v
	} else {
		x.s = 1
	}

	// Fast path for small integers (< 1e7).
	if v < int64(BASE) {
		e := 0
		for i := v; i >= 10; i /= 10 {
			e++
		}
		x.e = e
		x.d = []int32{int32(v)}
		return x
	}

	// Convert to string and parse.
	parseDecimalStr(x, strconv.FormatInt(v, 10))
	return x
}

// NewFromFloat64 creates a Decimal from a float64 using the default context.
func NewFromFloat64(v float64) (*Decimal, error) {
	return defaultCtx.NewFromFloat64(v)
}

// NewFromFloat64 creates a Decimal from a float64 using this context.
func (ctx *Context) NewFromFloat64(v float64) (*Decimal, error) {
	x := &Decimal{ctx: ctx}

	if v == 0 {
		if math.Signbit(v) {
			x.s = -1
		} else {
			x.s = 1
		}
		x.e = 0
		x.d = []int32{0}
		return x, nil
	}

	if math.IsNaN(v) {
		x.s = 0
		x.e = 0 // JS uses NaN for e, we use 0 for NaN decimals
		x.d = nil
		return x, nil
	}

	if math.IsInf(v, 0) {
		if v < 0 {
			x.s = -1
		} else {
			x.s = 1
		}
		x.d = nil
		x.e = 0
		return x, nil
	}

	if v < 0 {
		x.s = -1
		v = -v
	} else {
		x.s = 1
	}

	// Convert to string and parse.
	parseDecimalStr(x, strconv.FormatFloat(v, 'f', -1, 64))
	return x, nil
}

// copy creates a deep copy of the Decimal.
func (x *Decimal) copy() *Decimal {
	y := &Decimal{
		e:   x.e,
		s:   x.s,
		ctx: x.ctx,
	}
	if x.d != nil {
		y.d = make([]int32, len(x.d))
		copy(y.d, x.d)
	}
	return y
}

// copyWithCtx creates a copy of src using the context from x.
// This mirrors the JS pattern: `y = new x.constructor(src)`.
func (x *Decimal) newFromDecimal(src *Decimal) *Decimal {
	y := &Decimal{
		s:   src.s,
		ctx: x.ctx,
	}

	if external {
		if src.d == nil || src.e > x.ctx.MaxE {
			// Infinity.
			y.d = nil
			y.e = 0
		} else if src.e < x.ctx.MinE {
			// Zero.
			y.e = 0
			y.d = []int32{0}
		} else {
			y.e = src.e
			y.d = make([]int32, len(src.d))
			copy(y.d, src.d)
		}
	} else {
		y.e = src.e
		if src.d != nil {
			y.d = make([]int32, len(src.d))
			copy(y.d, src.d)
		}
	}

	return y
}

// ---- Predicates ----

// IsFinite returns true if x is a finite number (not NaN or Infinity).
func (x *Decimal) IsFinite() bool {
	return x.d != nil
}

// IsNaN returns true if x is NaN.
func (x *Decimal) IsNaN() bool {
	return x.s == 0
}

// IsZero returns true if x is zero (positive or negative).
func (x *Decimal) IsZero() bool {
	return x.d != nil && x.d[0] == 0
}

// IsNeg returns true if x is negative.
func (x *Decimal) IsNeg() bool {
	return x.s < 0
}

// IsPos returns true if x is positive.
func (x *Decimal) IsPos() bool {
	return x.s > 0
}

// IsInt returns true if x is an integer.
func (x *Decimal) IsInt() bool {
	return x.d != nil && ifloorDiv(x.e, LOG_BASE) > len(x.d)-2
}

// Sd returns the number of significant digits.
// If z is true and x is an integer, trailing zeros of the integer part are counted.
func (x *Decimal) Sd(z ...bool) int {
	if x.d == nil {
		return 0 // NaN
	}
	k := getPrecision(x.d)
	if len(z) > 0 && z[0] && x.e+1 > k {
		k = x.e + 1
	}
	return k
}

// Dp returns the number of decimal places.
func (x *Decimal) Dp() int {
	if x.d == nil {
		return 0 // NaN
	}
	w := len(x.d) - 1
	n := (w - ifloorDiv(x.e, LOG_BASE)) * LOG_BASE

	// Subtract the number of trailing zeros of the last word.
	lastWord := x.d[w]
	if lastWord != 0 {
		for lastWord%10 == 0 {
			lastWord /= 10
			n--
		}
	}
	if n < 0 {
		n = 0
	}
	return n
}

// getContext returns the context for this decimal, using the default if nil.
func (x *Decimal) getContext() *Context {
	if x.ctx != nil {
		return x.ctx
	}
	return defaultCtx
}
