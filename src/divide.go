package decimal

// divide is the core long division function matching decimal.js divide algorithm.
// Ported directly from decimal.js lines 1285-1545.
func divide(x, y *Decimal, pr int, rm RoundingMode, isMod bool, moduloMode RoundingMode) *Decimal {
	ctx := x.getContext()

	// Special cases: NaN or Infinity.
	if x.d == nil || y.d == nil {
		if x.s == 0 || y.s == 0 {
			return ctx.NewFromInt64(0) // NaN represented as s=0
		}
		if x.d == nil && y.d == nil {
			// Inf / Inf = NaN
			r := &Decimal{s: 0, ctx: ctx}
			return r
		}
		if x.d == nil {
			// Inf / y = ±Inf
			s := x.s * y.s
			r := &Decimal{s: s, d: nil, e: 0, ctx: ctx}
			return r
		}
		// x / Inf = ±0
		s := x.s * y.s
		r := &Decimal{s: s, d: []int32{0}, e: 0, ctx: ctx}
		return r
	}

	// Division by zero.
	if y.IsZero() {
		if x.IsZero() {
			// 0 / 0 = NaN
			r := &Decimal{s: 0, ctx: ctx}
			return r
		}
		// x / 0 = ±Inf
		s := x.s * y.s
		r := &Decimal{s: s, d: nil, e: 0, ctx: ctx}
		return r
	}

	// 0 / y = ±0
	if x.IsZero() {
		s := x.s * y.s
		r := &Decimal{s: s, d: []int32{0}, e: 0, ctx: ctx}
		return r
	}

	// Sign of result.
	sign := x.s * y.s

	// Prepare digit slices for division.
	xd := make([]int32, len(x.d))
	copy(xd, x.d)
	yd := make([]int32, len(y.d))
	copy(yd, y.d)

	e := x.e - y.e

	xL := len(xd)
	yL := len(yd)

	// Result digits slice.
	var qd []int32

	// Determine logBase and base.
	base32 := BASE
	logBase := LOG_BASE

	// Calculate precision required.
	sd := pr
	if sd <= 0 {
		sd = ctx.Precision
	}

	// Adjust sd to base 1e7 words.
	sd = iceil(sd, logBase) + 2

	var inexact bool

	if yL == 1 {
		// Divisor < 1e7.
		var k int32
		yd0 := yd[0]
		sd++

		for i := 0; (i < xL || k != 0) && sd > 0; i++ {
			var t int32
			if i < xL {
				t = k*base32 + xd[i]
			} else {
				t = k * base32
			}
			qd = append(qd, t/yd0)
			k = t % yd0
			sd--
		}

		more := k != 0

		// Leading zero?
		if len(qd) > 0 && qd[0] == 0 {
			qd = qd[1:]
		}

		if logBase == 1 {
			inexact = more
		} else {
			// Get number of digits of qd[0].
			digits := 1
			if len(qd) > 0 {
				for k2 := qd[0]; k2 >= 10; k2 /= 10 {
					digits++
				}
			}
			q := &Decimal{
				s:   sign,
				e:   digits + e - 1,
				d:   qd,
				ctx: ctx,
			}
			if len(q.d) == 0 {
				q.d = []int32{0}
				q.e = 0
			}
			_ = inexact
			return finalise(q, pr, rm)
		}
	} else {
		// Multi-word divisor.
		// Normalize divisor so highest digit >= base / 2.
		scale := base32 / (yd[0] + 1)
		if scale > 1 {
			yd = multiplyByInt(yd, scale, base32)
			xd = multiplyByInt(xd, scale, base32)
			xL = len(xd)
			yL = len(yd)
		}

		// Main Knuth long division algorithm loop.
		// Ensure xd is long enough.
		for len(xd) < yL+1 {
			xd = append(xd, 0)
		}

		y0 := yd[0]
		y1 := yd[1]

		for sd > 0 {
			var qHat int32
			if xd[0] == y0 {
				qHat = base32 - 1
			} else {
				qHat = (xd[0]*base32 + xd[1]) / y0
			}

			if qHat != 0 {
				// Refine qHat.
				for y1*qHat > (xd[0]*base32+xd[1]-qHat*y0)*base32+xd[2] {
					qHat--
				}
			}

			// Multiply and subtract.
			rem := multiplyAndSubtract(xd, yd, qHat, base32)
			if rem < 0 {
				// Over-subtracted: add back.
				qHat--
				addBack(xd, yd, base32)
			}

			qd = append(qd, qHat)

			// Shift xd left by 1 word.
			xd = xd[1:]
			if len(xd) < yL+1 {
				xd = append(xd, 0)
			}

			// Check if remainder is 0 and we passed input digits.
			allZero := true
			for _, v := range xd {
				if v != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				break
			}
			sd--
		}

		// Remove leading zeros from qd.
		for len(qd) > 1 && qd[0] == 0 {
			qd = qd[1:]
		}

		// Determine digits of first word.
		digits := 1
		if len(qd) > 0 {
			for k2 := qd[0]; k2 >= 10; k2 /= 10 {
				digits++
			}
		}

		q := &Decimal{
			s:   sign,
			e:   digits + e - 1,
			d:   qd,
			ctx: ctx,
		}
		if len(q.d) == 0 {
			q.d = []int32{0}
			q.e = 0
		}
		return finalise(q, pr, rm)
	}

	q := &Decimal{
		s:   sign,
		e:   e,
		d:   qd,
		ctx: ctx,
	}
	if len(q.d) == 0 {
		q.d = []int32{0}
		q.e = 0
	}
	return finalise(q, pr, rm)
}

// Helper: multiply slice by single int32.
func multiplyByInt(a []int32, k, base int32) []int32 {
	res := make([]int32, len(a)+1)
	var carry int32
	for i := len(a) - 1; i >= 0; i-- {
		t := a[i]*k + carry
		res[i+1] = t % base
		carry = t / base
	}
	res[0] = carry
	if res[0] == 0 {
		res = res[1:]
	}
	return res
}

// Helper: multiply y by qHat and subtract from x in-place.
func multiplyAndSubtract(x, y []int32, qHat, base int32) int32 {
	var borrow int32
	yL := len(y)
	for i := yL - 1; i >= 0; i-- {
		prod := y[i]*qHat + borrow
		xDigit := x[i+1] - prod%base
		borrow = prod / base
		if xDigit < 0 {
			xDigit += base
			borrow++
		}
		x[i+1] = xDigit
	}
	x[0] -= borrow
	return x[0]
}

// Helper: add y back to x if qHat was 1 too large.
func addBack(x, y []int32, base int32) {
	var carry int32
	yL := len(y)
	for i := yL - 1; i >= 0; i-- {
		sum := x[i+1] + y[i] + carry
		x[i+1] = sum % base
		carry = sum / base
	}
	x[0] += carry
}

// intPow calculates x^n using binary exponentiation to precision pr.
func intPow(ctx *Context, x *Decimal, n int, pr int) *Decimal {
	if n == 0 {
		return ctx.NewFromInt64(1)
	}
	if n < 0 {
		one := ctx.NewFromInt64(1)
		pow := intPow(ctx, x, -n, pr)
		return divide(one, pow, pr, ctx.Rounding, false, 0)
	}

	// Binary exponentiation.
	y := ctx.NewFromInt64(1)
	curr := x.copy()

	for n > 0 {
		if n&1 == 1 {
			y = finalise(y.Times(curr), pr, ctx.Rounding)
		}
		n >>= 1
		if n > 0 {
			curr = finalise(curr.Times(curr), pr, ctx.Rounding)
		}
	}
	return y
}

// divSubtract subtracts b from a in-place, in the given base.
func divSubtract(a, b []int32, aL int, base int32) {
	var borrow int32
	for aL2 := aL - 1; aL2 >= 0; aL2-- {
		bv := int32(0)
		if aL2 < len(b) {
			bv = b[aL2]
		}
		a[aL2] -= borrow
		if a[aL2] < bv {
			borrow = 1
		} else {
			borrow = 0
		}
		a[aL2] = borrow*base + a[aL2] - bv
	}

	// Remove leading zeros.
	for len(a) > 1 && a[0] == 0 {
		a = a[1:]
	}
}
