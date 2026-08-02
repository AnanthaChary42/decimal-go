package decimal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// Abs returns a new Decimal whose value is the absolute value of x.
func (x *Decimal) Abs() *Decimal {
	r := x.copy()
	if r.IsNeg() {
		r.s = 1
	}
	return finalise(r, -999999, 0)
}

func (x *Decimal) Neg() *Decimal {
	r := x.copy()
	r.s = -r.s
	return finalise(r, -999999, 0)
}

// Plus returns a new Decimal whose value is x + y.
func (x *Decimal) Plus(y *Decimal) *Decimal {
	return x.Add(y)
}

// Add returns a new Decimal whose value is x + y.
func (x *Decimal) Add(y *Decimal) *Decimal {
	ctx := x.getContext()
	// Convert y to use x's context.
	y = x.newFromDecimal(y)

	// If either is not finite...
	if x.d == nil || y.d == nil {
		// Return NaN if either is NaN.
		if x.s == 0 || y.s == 0 {
			r := x.copy()
			r.s = 0
			r.d = nil
			return r
		}

		// Return x if y is finite and x is ±Infinity.
		// Return x if both are ±Infinity with the same sign.
		// Return NaN if both are ±Infinity with different signs.
		// Return y if x is finite and y is ±Infinity.
		if x.d == nil {
			if y.d != nil || x.s == y.s {
				return x.copy()
			}
			r := x.copy()
			r.s = 0
			r.d = nil
			return r
		}
		return y.copy()
	}

	// If signs differ...
	if x.s != y.s {
		y2 := y.copy()
		y2.s = -y2.s
		return x.Minus(y2)
	}

	xd := make([]int32, len(x.d))
	copy(xd, x.d)
	yd := y.d
	pr := ctx.Precision
	rm := ctx.Rounding

	// If either is zero...
	if xd[0] == 0 || yd[0] == 0 {
		if yd[0] == 0 {
			result := x.copy()
			if external {
				return finalise(result, pr, rm)
			}
			return result
		}
		result := y.copy()
		if external {
			return finalise(result, pr, rm)
		}
		return result
	}

	// Calculate base 1e7 exponents.
	k := ifloorDiv(x.e, LOG_BASE)
	e := ifloorDiv(y.e, LOG_BASE)

	i := k - e

	// If base 1e7 exponents differ...
	if i != 0 {
		var d []int32
		var length int
		xdSmaller := false

		if i < 0 {
			d = xd
			i = -i
			length = len(yd)
			xdSmaller = true
		} else {
			d = make([]int32, len(yd))
			copy(d, yd)
			yd = d
			e = k
			length = len(xd)
		}

		// Limit number of zeros prepended to max(ceil(pr / LOG_BASE), len) + 1.
		kk := iceil(pr, LOG_BASE)
		if kk > length {
			length = kk + 1
		} else {
			length = length + 1
		}

		if i > length {
			i = length
			d = d[:1]
		}

		// Prepend zeros to equalise exponents.
		// Reverse, push zeros, reverse back.
		reversed := make([]int32, len(d))
		for idx := 0; idx < len(d); idx++ {
			reversed[len(d)-1-idx] = d[idx]
		}
		for j := 0; j < i; j++ {
			reversed = append(reversed, 0)
		}
		d2 := make([]int32, len(reversed))
		for idx := 0; idx < len(reversed); idx++ {
			d2[len(reversed)-1-idx] = reversed[idx]
		}

		if xdSmaller {
			// d was xd
			xd = d2
		} else {
			yd = d2
		}
	}

	length := len(xd)
	i2 := len(yd)

	// If yd is longer than xd, swap xd and yd so xd points to the longer array.
	if length-i2 < 0 {
		i2 = length
		d := yd
		yd = xd
		xd = d
	}

	// Only start adding at yd.length - 1.
	var carry int32
	for i3 := i2 - 1; i3 >= 0; i3-- {
		sum := xd[i3] + yd[i3] + carry
		carry = sum / BASE
		xd[i3] = sum % BASE
	}

	if carry != 0 {
		xd = append([]int32{carry}, xd...)
		e++
	}

	// Remove trailing zeros.
	for len(xd) > 0 && xd[len(xd)-1] == 0 {
		xd = xd[:len(xd)-1]
	}

	result := &Decimal{
		d:   xd,
		e:   getBase10Exponent(xd, e),
		s:   y.s,
		ctx: ctx,
	}

	if external {
		return finalise(result, pr, rm)
	}
	return result
}

// Minus returns a new Decimal whose value is x - y.
func (x *Decimal) Minus(y *Decimal) *Decimal {
	return x.Sub(y)
}

// Sub returns a new Decimal whose value is x - y.
func (x *Decimal) Sub(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = x.newFromDecimal(y)

	// If either is not finite...
	if x.d == nil || y.d == nil {
		if x.s == 0 || y.s == 0 {
			r := x.copy()
			r.s = 0
			r.d = nil
			return r
		}

		if x.d != nil {
			// x is finite, y is ±Infinity → return -y
			r := y.copy()
			r.s = -r.s
			return r
		}

		if y.d != nil || x.s != y.s {
			return x.copy()
		}

		// Both ±Infinity with same sign → NaN
		r := x.copy()
		r.s = 0
		r.d = nil
		return r
	}

	// If signs differ...
	if x.s != y.s {
		y2 := y.copy()
		y2.s = -y2.s
		return x.Plus(y2)
	}

	xd := make([]int32, len(x.d))
	copy(xd, x.d)
	yd := y.d
	pr := ctx.Precision
	rm := ctx.Rounding

	// If either is zero...
	if xd[0] == 0 || yd[0] == 0 {
		if yd[0] != 0 {
			r := y.copy()
			r.s = -r.s
			if external {
				return finalise(r, pr, rm)
			}
			return r
		}

		if xd[0] != 0 {
			r := x.copy()
			if external {
				return finalise(r, pr, rm)
			}
			return r
		}

		// Both zero.
		r := &Decimal{ctx: ctx}
		if rm == RoundFloor {
			r.s = -1
		} else {
			r.s = 1
		}
		r.e = 0
		r.d = []int32{0}
		return r
	}

	// x and y are finite, non-zero numbers with the same sign.
	e := ifloorDiv(y.e, LOG_BASE)
	xe := ifloorDiv(x.e, LOG_BASE)

	k := xe - e

	// If base 1e7 exponents differ...
	xLTy := false
	if k != 0 {
		if k < 0 {
			xLTy = true
			k = -k
		}

		var d []int32
		var length int
		if xLTy {
			d = make([]int32, len(xd))
			copy(d, xd)
			length = len(yd)
		} else {
			d = make([]int32, len(yd))
			copy(d, yd)
			e = xe
			length = len(xd)
		}

		i := maxInt(iceil(pr, LOG_BASE), length) + 2
		if k > i {
			k = i
			d = d[:1]
		}

		// Prepend zeros.
		reversed := make([]int32, len(d))
		for idx := 0; idx < len(d); idx++ {
			reversed[len(d)-1-idx] = d[idx]
		}
		for j := 0; j < k; j++ {
			reversed = append(reversed, 0)
		}
		d2 := make([]int32, len(reversed))
		for idx := 0; idx < len(reversed); idx++ {
			d2[len(reversed)-1-idx] = reversed[idx]
		}

		if xLTy {
			xd = d2
		} else {
			yd = d2
		}
	} else {
		// Exponents equal. Compare digits.
		i := len(xd)
		length := len(yd)
		if i < length {
			xLTy = true
		}
		minLen := length
		if i < minLen {
			minLen = i
		}

		for i2 := 0; i2 < minLen; i2++ {
			if xd[i2] != yd[i2] {
				xLTy = xd[i2] < yd[i2]
				break
			}
		}
		k = 0
	}

	if xLTy {
		// Swap.
		xd, yd = yd, xd
	}

	length := len(xd)

	// Append zeros to xd if shorter.
	for i := len(yd) - length; i > 0; i-- {
		xd = append(xd, 0)
		length++
	}

	// Subtract yd from xd.
	for i := len(yd) - 1; i >= k; i-- {
		if xd[i] < yd[i] {
			j := i
			for j > 0 && xd[j-1] == 0 {
				j--
				xd[j] = BASE - 1
			}
			if j > 0 {
				xd[j-1]--
			}
			xd[i] += BASE
		}
		xd[i] -= yd[i]
	}

	// Remove trailing zeros.
	for len(xd) > 0 && xd[len(xd)-1] == 0 {
		xd = xd[:len(xd)-1]
	}

	// Remove leading zeros and adjust exponent.
	for len(xd) > 0 && xd[0] == 0 {
		xd = xd[1:]
		e--
	}

	// Zero?
	if len(xd) == 0 {
		r := &Decimal{ctx: ctx}
		if rm == RoundFloor {
			r.s = -1
		} else {
			r.s = 1
		}
		r.e = 0
		r.d = []int32{0}
		return r
	}

	s := y.s
	if xLTy {
		s = -s
	}

	result := &Decimal{
		d:   xd,
		e:   getBase10Exponent(xd, e),
		s:   s,
		ctx: ctx,
	}

	if external {
		return finalise(result, pr, rm)
	}
	return result
}

// Times returns a new Decimal whose value is x * y.
func (x *Decimal) Times(y *Decimal) *Decimal {
	return x.Mul(y)
}

// Mul returns a new Decimal whose value is x * y.
func (x *Decimal) Mul(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = x.newFromDecimal(y)

	sign := x.s * y.s

	xd := x.d
	yd := y.d

	// If either is NaN, ±Infinity or ±0...
	if len(xd) == 0 || xd[0] == 0 || len(yd) == 0 || yd[0] == 0 {

		r := &Decimal{ctx: ctx}

		if sign == 0 || (len(xd) > 0 && xd[0] == 0 && len(yd) == 0) ||
			(len(yd) > 0 && yd[0] == 0 && len(xd) == 0) {
			// NaN
			r.s = 0
			r.d = nil
			r.e = 0
			return r
		}

		if xd == nil || yd == nil {
			// ±Infinity
			r.s = sign
			r.d = nil
			r.e = 0
			return r
		}

		// ±0
		r.s = sign
		r.e = 0
		r.d = []int32{0}
		return r
	}

	e := ifloorDiv(x.e, LOG_BASE) + ifloorDiv(y.e, LOG_BASE)
	xdL := len(xd)
	ydL := len(yd)

	// Ensure xd points to the longer array.
	if xdL < ydL {
		xd, yd = yd, xd
		xdL, ydL = ydL, xdL
	}

	// Initialise the result array with zeros.
	rL := xdL + ydL
	r := make([]int32, rL)

	// Multiply!
	for i := ydL - 1; i >= 0; i-- {
		var carry int32
		for k := xdL + i; k > i; k-- {
			t := int64(r[k]) + int64(yd[i])*int64(xd[k-i-1]) + int64(carry)
			r[k] = int32(t % int64(BASE))
			carry = int32(t / int64(BASE))
		}
		r[i] = (r[i] + carry) % BASE
	}

	// Remove trailing zeros.
	for rL > 0 && r[rL-1] == 0 {
		rL--
	}
	r = r[:rL]

	if carry := r[0]; carry != 0 {
		// There was a carry from the shift position, adjust.
	}

	// Check if the first element is the carry or the actual first digit.
	if len(r) > 0 && r[0] == 0 {
		r = r[1:]
		// don't increment e
	} else {
		e++
	}

	result := &Decimal{
		d:   r,
		e:   getBase10Exponent(r, e),
		s:   sign,
		ctx: ctx,
	}

	if external {
		return finalise(result, ctx.Precision, ctx.Rounding)
	}
	return result
}

// Div returns a new Decimal whose value is x / y.
func (x *Decimal) Div(y *Decimal) *Decimal {
	return divide(x, x.newFromDecimal(y), -1, RoundingMode(0), false, 0)
}

// DivToInt returns a new Decimal whose integer value is x / y.
func (x *Decimal) DivToInt(y *Decimal) *Decimal {
	ctx := x.getContext()
	result := divide(x, x.newFromDecimal(y), 0, RoundDown, true, 0)
	return finalise(result, ctx.Precision, ctx.Rounding)
}

// Mod returns a new Decimal whose value is x % y.
func (x *Decimal) Mod(y *Decimal) *Decimal {
	ctx := x.getContext()
	y = x.newFromDecimal(y)

	// Return NaN if x is ±Infinity or NaN, or y is NaN or ±0.
	if x.d == nil || y.s == 0 || (y.d != nil && y.d[0] == 0) {
		r := x.copy()
		r.s = 0
		r.d = nil
		return r
	}

	// Return x if y is ±Infinity or x is ±0.
	if y.d == nil || (x.d != nil && x.d[0] == 0) {
		return finalise(x.copy(), ctx.Precision, ctx.Rounding)
	}

	external = false

	var q *Decimal
	if ctx.Modulo == Euclid {
		q = divide(x, y.Abs(), 0, RoundFloor, true, 0)
		q.s *= y.s
	} else {
		q = divide(x, y, 0, ctx.Modulo, true, 0)
	}

	q = q.Times(y)
	external = true

	return x.Minus(q)
}

// Pow returns a new Decimal whose value is x raised to the power y.
func (x *Decimal) Pow(y *Decimal) *Decimal {
	ctx := x.getContext()
	yn := y.ToNumber()

	// Either ±Infinity, NaN or ±0?
	if x.d == nil || y.d == nil || x.d[0] == 0 || y.d[0] == 0 {
		xn := x.ToNumber()
		v := math.Pow(xn, yn)

		// Go's math.Pow differs from ECMAScript for ±1 raised to an
		// infinite exponent. decimal.js follows the ECMAScript result, NaN.
		if math.IsInf(yn, 0) && math.Abs(xn) == 1 {
			v = math.NaN()
		}
		r, _ := ctx.NewFromFloat64(v)
		return r
	}

	xCopy := x.copy()

	if xCopy.Eq(ctx.NewFromInt64(1)) {
		return xCopy
	}

	pr := ctx.Precision
	rm := ctx.Rounding

	if y.Eq(ctx.NewFromInt64(1)) {
		return finalise(xCopy, pr, rm)
	}

	// y exponent
	e := ifloorDiv(y.e, LOG_BASE)

	// If y is a small integer use exponentiation by squaring.
	if e >= len(y.d)-1 {
		k := yn
		if k < 0 {
			k = -k
		}
		if k <= MAX_SAFE_INTEGER {
			r := intPow(ctx, xCopy, int(k), pr)
			if y.s < 0 {
				return divide(ctx.NewFromInt64(1), r, pr, rm, false, 0)
			}
			return finalise(r, pr, rm)
		}
	}

	s := x.s

	// if x is negative
	if s < 0 {
		if e < len(y.d)-1 {
			// y is not an integer → NaN
			r := x.copy()
			r.s = 0
			r.d = nil
			return r
		}

		// Result is positive if last digit of integer y is even.
		if y.d[e]%2 == 0 {
			s = 1
		}

		// if x.eq(-1)
		if x.e == 0 && x.d[0] == 1 && len(x.d) == 1 {
			xCopy.s = s
			return xCopy
		}

		// Subsequent logarithm/exponential work is performed on |x|; restore
		// the sign after the integer-power result has been calculated.
		xCopy.s = 1
	}

	// Special case: pow(x, 0.5) == sqrt(x), pow(x, -0.5) == 1/sqrt(x)
	// This avoids the float64 fallback which overflows/underflows for extreme exponents.
	halfPos, _ := ctx.New("0.5")
	halfNeg, _ := ctx.New("-0.5")
	if y.Eq(halfPos) {
		r := xCopy.Sqrt()
		r.s = s
		return r
	}
	if y.Eq(halfNeg) {
		r := xCopy.Sqrt()
		r.s = s
		return divide(ctx.NewFromInt64(1), r, pr, rm, false, 0)
	}

	r := powPositive(ctx, xCopy, y, pr, rm)
	if r == nil {
		// Parsing finite Decimal values into big.Float cannot normally fail;
		// retain a float64 fallback only as a defensive last resort.
		result := math.Pow(x.ToNumber(), yn)
		r, _ = ctx.NewFromFloat64(result)
		if r == nil {
			return r
		}
		r = finalise(r, pr, rm)
	}
	r.s = s
	return r
}

// Sqrt returns a new Decimal whose value is the square root of x.
func (x *Decimal) Sqrt() *Decimal {
	ctx := x.getContext()
	d := x.d
	e := x.e
	s := x.s

	// Negative/NaN/Infinity/zero?
	if s != 1 || d == nil || d[0] == 0 {
		if s == 0 || (s < 0 && (d == nil || d[0] != 0)) {
			// NaN
			r := x.copy()
			r.s = 0
			r.d = nil
			return r
		}
		if d != nil {
			return x.copy() // sqrt(0) = 0, sqrt(-0) = -0
		}
		// sqrt(Infinity) = Infinity
		r := x.copy()
		r.d = nil
		return r
	}

	external = false

	// Initial estimate.
	sFloat := math.Sqrt(x.ToNumber())

	var r *Decimal
	if sFloat == 0 || math.IsInf(sFloat, 0) {
		n := digitsToStringExact(d)
		if (len(n)+e)%2 == 0 {
			n += "0"
		}
		sFloat = math.Sqrt(parseFloatStr(n))
		e = ifloorDiv(e+1, 2) - boolToInt(e < 0 || e%2 != 0)

		if math.IsInf(sFloat, 0) {
			r, _ = ctx.New("5e" + itoa(e))
		} else {
			r, _ = ctx.New(formatFloat(sFloat, e))
		}
	} else {
		r, _ = ctx.New(formatFloatSimple(sFloat))
	}

	precision := ctx.Precision
	sd := precision + 3
	more := false
	repeating := false
	half := &Decimal{d: []int32{5000000}, e: -1, s: 1, ctx: ctx}

	// Newton-Raphson iteration. Match decimal.js's guard-digit checks: simply
	// stopping when the first precision+3 digits agree loses directed-rounding
	// information for values just above a perfect square.
	for {
		t := r
		r = t.Plus(divide(x, t, sd+2, RoundDown, false, 0)).Times(half)

		tDigits := digitsToStringExact(t.d)
		rDigits := digitsToStringExact(r.d)
		if tDigits[:minInt(sd, len(tDigits))] != rDigits[:minInt(sd, len(rDigits))] {
			continue
		}

		roundingDigits := rDigits
		if sd-3 < len(roundingDigits) {
			end := minInt(sd+1, len(roundingDigits))
			roundingDigits = roundingDigits[sd-3 : end]
		} else {
			roundingDigits = ""
		}

		// The fourth rounding digit can be one too small at a rounding boundary.
		if roundingDigits == "9999" || (!repeating && roundingDigits == "4999") {
			if !repeating {
				finalise(t, precision+1, RoundUp)
				if t.Times(t).Eq(x) {
					r = t
					break
				}
			}
			sd += 4
			repeating = true
			continue
		}

		// When the guard digits are zero (or exactly 5000), determine whether
		// an omitted non-zero tail exists before applying the caller's rounding.
		allZero := roundingDigits == "" || strings.Trim(roundingDigits, "0") == ""
		isHalfWithZeros := len(roundingDigits) > 0 && roundingDigits[0] == '5' &&
			strings.Trim(roundingDigits[1:], "0") == ""
		if allZero || isHalfWithZeros {
			finalise(r, precision+1, RoundDown)
			more = !r.Times(r).Eq(x)
		}
		break
	}

	external = true
	return finalise(r, precision, ctx.Rounding, more)
}

// Cbrt returns a new Decimal whose value is the cube root of x.
func (x *Decimal) Cbrt() *Decimal {
	ctx := x.getContext()

	if !x.IsFinite() || x.IsZero() {
		return x.copy()
	}

	external = false

	// Initial estimate.
	sFloat := math.Cbrt(x.ToNumber())
	var r *Decimal
	if sFloat == 0 || math.IsInf(sFloat, 0) {
		r, _ = ctx.New("5e" + itoa(ifloorDiv(x.e+1, 3)))
		r.s = x.s
	} else {
		r, _ = ctx.New(formatFloatSimple(sFloat))
	}

	sd := ctx.Precision + 3

	// Halley's method.
	for iter := 0; iter < 100; iter++ {
		t := r.copy()
		t3 := t.Times(t).Times(t)
		t3plusx := t3.Plus(x)
		r = divide(t3plusx.Plus(x).Times(t), t3plusx.Plus(t3), sd+2, RoundDown, false, 0)

		tStr := digitsToStringExact(t.d)[:minInt(sd, len(digitsToStringExact(t.d)))]
		rStr := digitsToStringExact(r.d)[:minInt(sd, len(digitsToStringExact(r.d)))]

		if tStr == rStr {
			break
		}
	}

	external = true
	return finalise(r, ctx.Precision, ctx.Rounding)
}

// intPow implements exponentiation by squaring.
// Used by Pow and parseOther.
func intPow(ctx *Context, x *Decimal, n int, pr int) *Decimal {
	r := ctx.NewFromInt64(1)
	k := iceil(pr, LOG_BASE) + 4

	external = false

	for {
		if n%2 != 0 {
			r = r.Times(x)
			truncateArr(&r.d, k)
		}

		n = n / 2
		if n == 0 {
			break
		}

		x = x.Times(x)
		truncateArr(&x.d, k)
	}

	external = true
	return r
}

// Helper: parse float from string.
func parseFloatStr(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return 0
	}
	return v
}

// Helper: format float with exponent.
func formatFloat(f float64, e int) string {
	s := formatFloatSimple(f)
	// Find 'e' in the result.
	for i, ch := range s {
		if ch == 'e' || ch == 'E' {
			return s[:i+1] + itoa(e)
		}
	}
	return s + "e" + itoa(e)
}

// formatFloatSimple formats a float64 to a clean decimal string.
func formatFloatSimple(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
