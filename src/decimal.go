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
	Precision int          // Significant digits, default 20
	Rounding  RoundingMode // Default RoundHalfUp (4)
	ToExpNeg  int          // Exponent threshold for exponential notation (neg), default -7
	ToExpPos  int          // Exponent threshold for exponential notation (pos), default 21
	MinE      int          // Minimum exponent, default -9e15
	MaxE      int          // Maximum exponent, default 9e15
	Modulo    RoundingMode // Modulo mode, default RoundFloor (3)
}

// Default Context matching decimal.js default settings.
var defaultCtx = &Context{
	Precision: 20,
	Rounding:  RoundHalfUp,
	ToExpNeg:  -7,
	ToExpPos:  21,
	MinE:      MIN_EXP,
	MaxE:      MAX_EXP,
	Modulo:    RoundFloor,
}

// Decimal represents an arbitrary-precision decimal number.
// Mirrored after decimal.js internal structure.
type Decimal struct {
	s   int      // Sign: 1 (positive), -1 (negative), 0 (NaN)
	e   int      // Exponent (base 10)
	d   []int32  // Digits array (base 1e7)
	ctx *Context // Configuration context for operations
}

// New creates a new Decimal from a string using the default context.
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
		return x, nil
	}

	// Try base conversion (hex, binary, octal) or special values (Infinity, NaN).
	err := parseOther(x, s)
	if err != nil {
		return nil, err
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
		v = -v
	} else {
		x.s = 1
	}

	if v < int64(BASE) {
		x.e = getBase10Exponent([]int32{int32(v)}, 0)
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
		x.e = 0
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
	c := &Decimal{
		s:   x.s,
		e:   x.e,
		ctx: x.ctx,
	}
	if x.d != nil {
		c.d = make([]int32, len(x.d))
		copy(c.d, x.d)
	}
	return c
}

// getContext returns the Decimal's context or defaultCtx if nil.
func (x *Decimal) getContext() *Context {
	if x.ctx != nil {
		return x.ctx
	}
	return defaultCtx
}

// Predicates matching decimal.js instance methods.

// IsFinite returns true if x is a finite number (not NaN or Infinity).
func (x *Decimal) IsFinite() bool {
	return x.d != nil
}

// IsNaN returns true if x is NaN.
func (x *Decimal) IsNaN() bool {
	return x.s == 0
}

// IsZero returns true if x is zero.
func (x *Decimal) IsZero() bool {
	return x.d != nil && len(x.d) == 1 && x.d[0] == 0
}

// IsPos returns true if x is positive (including +0).
func (x *Decimal) IsPos() bool {
	return x.s > 0
}

// IsNeg returns true if x is negative (including -0).
func (x *Decimal) IsNeg() bool {
	return x.s < 0
}

// IsInt returns true if x is an integer.
func (x *Decimal) IsInt() bool {
	return x.d != nil && ifloorDiv(x.e, LOG_BASE) > len(x.d)-2
}

// Sd returns the number of significant digits.
func (x *Decimal) Sd() int {
	if x.d == nil {
		return 0 // NaN or Infinity
	}
	return getPrecision(x.d)
}

// Dp returns the number of decimal places.
func (x *Decimal) Dp() int {
	if x.d == nil {
		return 0 // NaN or Infinity
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

	if x.e < 0 {
		return n
	}
	if n < 0 {
		return 0
	}
	return n
}

// Helper constructors.

func (ctx *Context) newFromDecimal(y *Decimal) *Decimal {
	if y.ctx == ctx {
		return y.copy()
	}
	c := y.copy()
	c.ctx = ctx
	return c
}
