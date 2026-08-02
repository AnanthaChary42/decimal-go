package decimal

// inexact is set by divide when the result has a remainder.
// Used by toStringBinary.
var inexact bool

// divide performs division: x / y.
// pr is the precision, rm is the rounding mode.
// dp indicates whether pr is decimal places (true) or significant digits (false).
// base is used for base conversion (0 means use BASE).
//
// When pr < 0, uses context precision and rounding (equivalent to JS pr == null).
//
// This is the most complex internal function, directly ported from decimal.js.
func divide(x, y *Decimal, pr int, rm RoundingMode, dp bool, base int) *Decimal {
	ctx := x.getContext()

	sign := int8(1)
	if x.s != y.s {
		sign = -1
	}

	xd := x.d
	yd := y.d

	// Either NaN, Infinity, or 0?
	if len(xd) == 0 || len(yd) == 0 ||
		(len(xd) > 0 && xd[0] == 0) || (len(yd) > 0 && yd[0] == 0) {

		q := &Decimal{ctx: ctx}

		if x.s == 0 || y.s == 0 ||
			(len(xd) > 0 && len(yd) > 0 && xd[0] == 0 && yd[0] == 0) ||
			(len(xd) == 0 && len(yd) == 0) {
			// NaN
			q.s = 0
			q.d = nil
			q.e = 0
			return q
		}

		if xd != nil && xd[0] == 0 || yd == nil {
			// ±0 / y  or  x / ±Infinity → ±0
			q.s = sign
			q.e = 0
			q.d = []int32{0}
			return q
		}

		// x / 0 → ±Infinity  or  ±Infinity / y → ±Infinity
		q.s = sign
		q.d = nil
		q.e = 0
		return q
	}

	var logBase int
	var e int

	if base != 0 {
		logBase = 1
		e = x.e - y.e
	} else {
		base = int(BASE)
		logBase = LOG_BASE
		e = ifloorDiv(x.e, logBase) - ifloorDiv(y.e, logBase)
	}

	yL := len(yd)
	xL := len(xd)
	q := &Decimal{s: sign, ctx: ctx}
	qd := make([]int32, 0)

	// Result exponent may be one less than e.
	i := 0
	for i < yL && i < xL && yd[i] == xd[i] {
		i++
	}
	xdi := int32(0)
	if i < xL {
		xdi = xd[i]
	}
	if i < yL && yd[i] > xdi {
		e--
	}

	var sd int
	if pr < 0 {
		sd = ctx.Precision
		rm = ctx.Rounding
		pr = sd
	} else if dp {
		sd = pr + (x.e - y.e) + 1
	} else {
		sd = pr
	}

	if sd < 0 {
		roundUp := (rm == RoundUp) ||
			(rm == RoundCeil && sign > 0) ||
			(rm == RoundFloor && sign < 0)
		if roundUp {
			q.d = []int32{1}
			q.e = e
		} else {
			q.d = []int32{0}
			q.e = 0
		}
		inexact = true
		return q
	}

	// Convert precision in number of base 10 digits to base 1e7 digits.
	sd = sd/logBase + 2
	i = 0

	base64 := int64(base)
	base32 := int32(base)

	if yL == 1 {
		// Divisor < 1e7.
		k := int64(0)
		yd0 := int64(yd[0])
		sd++

		for (i < xL || k != 0) && sd > 0 {
			var t int64
			if i < xL {
				t = k*base64 + int64(xd[i])
			} else {
				t = k * base64
			}
			qd = append(qd, int32(t/yd0))
			k = t % yd0
			i++
			sd--
		}

		more := k != 0 || i < xL

		// Leading zero?
		if len(qd) > 0 && qd[0] == 0 {
			qd = qd[1:]
		}

		if logBase == 1 {
			q.e = e
			q.d = qd
			inexact = more
		} else {
			// Get number of digits of qd[0].
			digits := 1
			if len(qd) > 0 {
				for k2 := qd[0]; k2 >= 10; k2 /= 10 {
					digits++
				}
			}
			q.e = digits + e*logBase - 1
			q.d = qd
			finalise(q, boolToRM(dp, pr+q.e+1, pr), rm, more)
		}
		return q
	}

	// Divisor >= 1e7.
	// Normalise xd and yd so highest order digit of yd is >= base/2.
	k := base32 / (yd[0] + 1)

	if k > 1 {
		yd = multiplyInteger(yd, k, base32)
		xd = multiplyInteger(xd, k, base32)
		yL = len(yd)
		xL = len(xd)
	}

	xi := yL
	rem := make([]int32, yL)
	copy(rem, xd[:minInt(yL, xL)])
	remL := len(rem)

	// Add zeros to make remainder as long as divisor.
	for remL < yL {
		rem = append(rem, 0)
		remL++
	}

	yz := make([]int32, len(yd)+1)
	yz[0] = 0
	copy(yz[1:], yd)
	yd0 := yd[0]

	if len(yd) > 1 && yd[1] >= base32/2 {
		yd0++
	}

	more := false

	for {
		k = 0

		// Compare divisor and remainder.
		cmp := divCompare(yd, rem, yL, remL)

		if cmp < 0 {
			// Divisor < remainder. Calculate trial digit, k.
			rem0 := int64(rem[0])
			if yL != remL {
				rem0 = rem0*int64(base32) + int64(safeGet(rem, 1))
			}

			k = int32(rem0 / int64(yd0))

			// prod will hold the product to subtract from remainder.
			var prod []int32
			var prodL int

			if k > 1 {
				if k >= base32 {
					k = base32 - 1
				}

				// product = divisor * trial digit.
				prod = multiplyInteger(yd, k, base32)
				prodL = len(prod)
				remL = len(rem)

				// Compare product and remainder.
				cmp = divCompare(prod, rem, prodL, remL)

				// product > remainder.
				if cmp == 1 {
					k--

					// Subtract divisor from product.
					if yL < prodL {
						prod = divSubtract(prod, yz, prodL, base32)
					} else {
						prod = divSubtract(prod, yd, prodL, base32)
					}
				}
			} else {
				// cmp is -1.
				// If k is 0, there is no need to compare yd and rem again below,
				// so change cmp to 1 to avoid it.
				// If k is 1 there is a need to compare yd and rem again below.
				if k == 0 {
					cmp = 1
					k = 1
				}
				prod = make([]int32, len(yd))
				copy(prod, yd)
			}

			prodL = len(prod)
			if prodL < remL {
				prod = append([]int32{0}, prod...)
				prodL++
			}

			// Subtract product from remainder.
			rem = divSubtract(rem, prod, remL, base32)

			// If product was < previous remainder.
			if cmp == -1 {
				remL = len(rem)

				// Compare divisor and new remainder.
				cmp = divCompare(yd, rem, yL, remL)

				if cmp < 1 {
					k++

					// Subtract divisor from remainder.
					if yL < remL {
						rem = divSubtract(rem, yz, remL, base32)
					} else {
						rem = divSubtract(rem, yd, remL, base32)
					}
				}
			}

			remL = len(rem)
		} else if cmp == 0 {
			k++
			rem = []int32{0}
			remL = 1
		}
		// if cmp === 1, k will be 0

		// Add the next digit, k, to the result array.
		qd = append(qd, k)

		// Update the remainder.
		if cmp != 0 && len(rem) > 0 && rem[0] != 0 {
			xVal := int32(0)
			if xi < xL {
				xVal = xd[xi]
			}
			rem = append(rem, xVal)
			remL = len(rem)
		} else {
			if xi < xL {
				rem = []int32{xd[xi]}
			} else {
				rem = []int32{}
			}
			remL = len(rem)
		}

		xi++
		sd--

		// In JS: (xi++ < xL || rem[0] !== void 0) && sd--
		// rem[0] !== void 0: In JS, rem = [xd[xi]] where xd[xi] is undefined (out of bounds)
		// creates rem = [undefined], and rem[0] !== void 0 is false.
		// In Go, we use len(rem) > 0 to indicate the element is "defined".
		hasRemainder := len(rem) > 0
		if !(xi <= xL || hasRemainder) || sd <= 0 {
			more = len(rem) > 0 && rem[0] != 0
			break
		}
	}

	// Leading zero?
	if len(qd) > 0 && qd[0] == 0 {
		qd = qd[1:]
	}

	if logBase == 1 {
		q.e = e
		q.d = qd
		inexact = more
	} else {
		// To calculate q.e, first get the number of digits of qd[0].
		digits := 1
		if len(qd) > 0 {
			for k2 := qd[0]; k2 >= 10; k2 /= 10 {
				digits++
			}
		}
		q.e = digits + e*logBase - 1
		q.d = qd
		finalise(q, boolToRM(dp, pr+q.e+1, pr), rm, more)
	}

	return q
}

// boolToRM selects between two values based on a bool.
func boolToRM(dp bool, a, b int) int {
	if dp {
		return a
	}
	return b
}

// multiplyInteger multiplies digit array x by integer k in given base.
// Returns a new slice.
func multiplyInteger(x []int32, k int32, base int32) []int32 {
	result := make([]int32, len(x))
	copy(result, x)

	var carry int64
	for i := len(result) - 1; i >= 0; i-- {
		temp := int64(result[i])*int64(k) + carry
		result[i] = int32(temp % int64(base))
		carry = temp / int64(base)
	}

	if carry > 0 {
		result = append([]int32{int32(carry)}, result...)
	}

	return result
}

// divCompare compares two digit arrays a and b with lengths aL and bL.
// Returns 1 if a > b, -1 if a < b, 0 if equal.
func divCompare(a, b []int32, aL, bL int) int {
	if aL != bL {
		if aL > bL {
			return 1
		}
		return -1
	}

	for i := 0; i < aL; i++ {
		av := int32(0)
		if i < len(a) {
			av = a[i]
		}
		bv := int32(0)
		if i < len(b) {
			bv = b[i]
		}
		if av != bv {
			if av > bv {
				return 1
			}
			return -1
		}
	}

	return 0
}

// divSubtract subtracts b from a in-place, in the given base.
// Returns the result slice (with leading zeros removed).
// In JS, subtract mutates the array and uses shift() for leading zeros.
// In Go, we must return the new slice since re-slicing doesn't affect the caller.
func divSubtract(a, b []int32, aL int, base int32) []int32 {
	i := 0
	for aL2 := aL - 1; aL2 >= 0; aL2-- {
		bv := int32(0)
		if aL2 < len(b) {
			bv = b[aL2]
		}
		a[aL2] -= int32(i)
		if a[aL2] < bv {
			i = 1
		} else {
			i = 0
		}
		a[aL2] = int32(i)*base + a[aL2] - bv
	}

	// Remove leading zeros.
	for len(a) > 1 && a[0] == 0 {
		a = a[1:]
	}

	return a
}

// safeGet safely gets element at index i, returning 0 if out of bounds.
func safeGet(arr []int32, i int) int32 {
	if i >= 0 && i < len(arr) {
		return arr[i]
	}
	return 0
}
