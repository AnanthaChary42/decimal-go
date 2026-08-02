package decimal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// ConfigOptions holds optional configuration values for Config and Clone.
// Pointer fields allow distinguishing "unset" (nil) from zero values.
type ConfigOptions struct {
	Precision *int
	Rounding  *int
	ToExpNeg  *int
	ToExpPos  *int
	MinE      *int
	MaxE      *int
	Crypto    *bool
	Modulo    *int
	Defaults  *bool
}

// GetDefaultContext returns the package-level default context.
// Modifying the returned context affects all Decimals created with the default.
func GetDefaultContext() *Context {
	return defaultCtx
}

// ResetDefaultContext resets the package-level default context to factory defaults.
// This is primarily used by tests to ensure isolation.
func ResetDefaultContext() {
	*defaultCtx = *defaultConfigValues()
}

// Config applies configuration options to the package-level default context.
// Panics with DecimalError for invalid values, matching JS throw behavior.
// Returns the default context.
func Config(opts ConfigOptions) *Context {
	return applyConfig(defaultCtx, opts)
}

// Set is an alias for Config, matching JS Decimal.set.
func Set(opts ConfigOptions) *Context {
	return Config(opts)
}

// Clone creates a new independent context copying from the package-level default.
// If opts are provided, they are validated and applied to the new context.
func Clone(opts ...ConfigOptions) *Context {
	return defaultCtx.Clone(opts...)
}

// Clone creates a new independent Context with the same settings as ctx.
// If opts are provided, they are validated and applied to the new context.
func (ctx *Context) Clone(opts ...ConfigOptions) *Context {
	newCtx := &Context{
		Precision: ctx.Precision,
		Rounding:  ctx.Rounding,
		ToExpNeg:  ctx.ToExpNeg,
		ToExpPos:  ctx.ToExpPos,
		MinE:      ctx.MinE,
		MaxE:      ctx.MaxE,
		Modulo:    ctx.Modulo,
		Crypto:    ctx.Crypto,
	}

	if len(opts) > 0 {
		applyConfig(newCtx, opts[0])
	}

	return newCtx
}

// defaultConfigValues returns a Context with factory default settings.
func defaultConfigValues() *Context {
	return &Context{
		Precision: 20,
		Rounding:  RoundHalfUp,
		ToExpNeg:  -7,
		ToExpPos:  21,
		MinE:      -EXP_LIMIT,
		MaxE:      EXP_LIMIT,
		Modulo:    RoundDown,
		Crypto:    false,
	}
}

// applyConfig validates and applies configuration options to a context.
// Panics with DecimalError for invalid values.
func applyConfig(ctx *Context, opts ConfigOptions) *Context {
	// Handle defaults: true — reset to factory before applying other fields.
	if opts.Defaults != nil && *opts.Defaults {
		defaults := defaultConfigValues()
		ctx.Precision = defaults.Precision
		ctx.Rounding = defaults.Rounding
		ctx.ToExpNeg = defaults.ToExpNeg
		ctx.ToExpPos = defaults.ToExpPos
		ctx.MinE = defaults.MinE
		ctx.MaxE = defaults.MaxE
		ctx.Modulo = defaults.Modulo
		ctx.Crypto = defaults.Crypto
	}

	// precision: [1, MAX_DIGITS]
	if opts.Precision != nil {
		v := *opts.Precision
		if v < 1 || float64(v) > MAX_DIGITS {
			panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: precision: %v", v)})
		}
		ctx.Precision = v
	}

	// rounding: [0, 8]
	if opts.Rounding != nil {
		v := *opts.Rounding
		if v < 0 || v > 8 {
			panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: rounding: %v", v)})
		}
		ctx.Rounding = RoundingMode(v)
	}

	// toExpNeg: [-EXP_LIMIT, 0]
	if opts.ToExpNeg != nil {
		v := *opts.ToExpNeg
		if v > 0 || float64(v) < -EXP_LIMIT {
			panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: toExpNeg: %v", v)})
		}
		ctx.ToExpNeg = v
	}

	// toExpPos: [0, EXP_LIMIT]
	if opts.ToExpPos != nil {
		v := *opts.ToExpPos
		if v < 0 || float64(v) > EXP_LIMIT {
			panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: toExpPos: %v", v)})
		}
		ctx.ToExpPos = v
	}

	// minE: [-EXP_LIMIT, 0]
	if opts.MinE != nil {
		v := *opts.MinE
		if v > 0 || float64(v) < -EXP_LIMIT {
			panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: minE: %v", v)})
		}
		ctx.MinE = v
	}

	// maxE: [0, EXP_LIMIT]
	if opts.MaxE != nil {
		v := *opts.MaxE
		if v < 0 || float64(v) > EXP_LIMIT {
			panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: maxE: %v", v)})
		}
		ctx.MaxE = v
	}

	// modulo: [0, 9]
	if opts.Modulo != nil {
		v := *opts.Modulo
		if v < 0 || v > 9 {
			panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: modulo: %v", v)})
		}
		ctx.Modulo = RoundingMode(v)
	}

	// crypto
	if opts.Crypto != nil {
		ctx.Crypto = *opts.Crypto
	}

	return ctx
}

// ---- Static Functions ----

// Sign returns a Decimal whose value is the sign of x:
//
//	1 if x > 0, -1 if x < 0, 0 if x == 0, NaN if x is NaN.
//
// Unlike the Go sign() test helper, this returns a Decimal matching JS Decimal.sign().
func Sign(x *Decimal) *Decimal {
	if x.IsNaN() {
		return &Decimal{s: 0, d: nil, e: 0, ctx: x.getContext()}
	}
	if x.IsZero() {
		// Preserve sign of zero: +0 or -0
		r := x.copy()
		return r
	}
	if x.IsNeg() {
		return &Decimal{s: -1, d: []int32{1}, e: 0, ctx: x.getContext()}
	}
	return &Decimal{s: 1, d: []int32{1}, e: 0, ctx: x.getContext()}
}

// Min returns the minimum of the given Decimals.
// If any argument is NaN, returns NaN.
// For equal zeros, prefers negative zero (matches JS Decimal.min signed-zero semantics).
func Min(args ...*Decimal) *Decimal {
	if len(args) == 0 {
		return nil
	}

	result := args[0]
	for i := 1; i < len(args); i++ {
		arg := args[i]
		// NaN propagation.
		if arg.IsNaN() {
			return arg
		}
		if result.IsNaN() {
			return result
		}
		cmp, _ := arg.Cmp(result)
		if cmp < 0 {
			result = arg
		} else if cmp == 0 && arg.IsZero() && arg.IsNeg() && result.IsZero() && !result.IsNeg() {
			// For equal zeros, prefer negative zero for min.
			result = arg
		}
	}
	return result
}

// Max returns the maximum of the given Decimals.
// If any argument is NaN, returns NaN.
// For equal zeros, prefers positive zero (matches JS Decimal.max signed-zero semantics).
func Max(args ...*Decimal) *Decimal {
	if len(args) == 0 {
		return nil
	}

	result := args[0]
	for i := 1; i < len(args); i++ {
		arg := args[i]
		// NaN propagation.
		if arg.IsNaN() {
			return arg
		}
		if result.IsNaN() {
			return result
		}
		cmp, _ := arg.Cmp(result)
		if cmp > 0 {
			result = arg
		} else if cmp == 0 && result.IsZero() && result.IsNeg() && arg.IsZero() && !arg.IsNeg() {
			// For equal zeros, prefer positive zero for max.
			result = arg
		}
	}
	return result
}

// Sum returns a new Decimal whose value is the sum of the arguments.
// Uses the default context.
func Sum(args ...*Decimal) *Decimal {
	if len(args) == 0 {
		return defaultCtx.NewFromInt64(0)
	}
	result := args[0].copy()
	for i := 1; i < len(args); i++ {
		result = result.Add(args[i])
	}
	return result
}

// Random returns a random Decimal in [0, 1) with at most sd significant digits.
// If sd is not provided, uses the context's Precision.
func (ctx *Context) Random(sd ...int) *Decimal {
	significantDigits := ctx.Precision
	if len(sd) > 0 {
		significantDigits = sd[0]
	}

	if significantDigits < 1 || float64(significantDigits) > MAX_DIGITS {
		panic(&DecimalError{Message: fmt.Sprintf("[DecimalError] Invalid argument: %v", significantDigits)})
	}

	// Generate a random integer in [0, 10^sd).
	limit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(significantDigits)), nil)
	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		panic(err)
	}

	if n.Sign() == 0 {
		return ctx.NewFromInt64(0)
	}

	nStr := n.String()

	// Pad to sd digits to form the fractional part.
	padded := strings.Repeat("0", significantDigits-len(nStr)) + nStr

	// Remove trailing zeros for canonical representation.
	trimmed := strings.TrimRight(padded, "0")
	if len(trimmed) == 0 {
		return ctx.NewFromInt64(0)
	}

	r, _ := ctx.New("0." + trimmed)
	return r
}

// ---- Helper constructors for pointer fields in ConfigOptions ----

// IntPtr returns a pointer to the given int. Convenience for ConfigOptions.
func IntPtr(v int) *int {
	return &v
}

// BoolPtr returns a pointer to the given bool. Convenience for ConfigOptions.
func BoolPtr(v bool) *bool {
	return &v
}
