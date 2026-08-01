package decimal

// Cmp compares x and y.
// Returns:
//   -1 if x < y
//    0 if x == y
//    1 if x > y
//   (0, false) if x or y is NaN
//
// Equivalent to comparedTo in decimal.js.
func (x *Decimal) Cmp(y *Decimal) (int, bool) {
	ctx := x.getContext()
	y = ctx.newFromDecimal(y)

	// NaN check.
	if x.s == 0 || y.s == 0 {
		return 0, false
	}

	// Sign comparison.
	if x.s != y.s {
		if x.IsZero() && y.IsZero() {
			return 0, true
		}
		if x.s > y.s {
			return 1, true
		}
		return -1, true
	}

	// Infinities.
	if x.d == nil || y.d == nil {
		if x.d == nil && y.d == nil {
			return 0, true
		}
		if x.d == nil {
			return x.s, true
		}
		return -y.s, true
	}

	// Exponent comparison.
	if x.e != y.e {
		if (x.e > y.e) == (x.s > 0) {
			return 1, true
		}
		return -1, true
	}

	// Digits comparison word by word.
	kX := len(x.d)
	kY := len(y.d)
	minL := kX
	if kY < minL {
		minL = kY
	}

	for i := 0; i < minL; i++ {
		if x.d[i] != y.d[i] {
			if (x.d[i] > y.d[i]) == (x.s > 0) {
				return 1, true
			}
			return -1, true
		}
	}

	// If words match so far, longer array is larger in absolute value.
	if kX != kY {
		if (kX > kY) == (x.s > 0) {
			return 1, true
		}
		return -1, true
	}

	return 0, true
}

// Eq returns true if x == y.
func (x *Decimal) Eq(y *Decimal) bool {
	cmp, ok := x.Cmp(y)
	return ok && cmp == 0
}

// Gt returns true if x > y.
func (x *Decimal) Gt(y *Decimal) bool {
	cmp, ok := x.Cmp(y)
	return ok && cmp > 0
}

// Gte returns true if x >= y.
func (x *Decimal) Gte(y *Decimal) bool {
	cmp, ok := x.Cmp(y)
	return ok && cmp >= 0
}

// Lt returns true if x < y.
func (x *Decimal) Lt(y *Decimal) bool {
	cmp, ok := x.Cmp(y)
	return ok && cmp < 0
}

// Lte returns true if x <= y.
func (x *Decimal) Lte(y *Decimal) bool {
	cmp, ok := x.Cmp(y)
	return ok && cmp <= 0
}

// Equals is an alias for Eq.
func (x *Decimal) Equals(y *Decimal) bool {
	return x.Eq(y)
}

// ComparedTo returns int matching decimal.js comparedTo (-1, 0, 1, or NaN handling).
// Returns 0 if NaN for JS API compatibility.
func (x *Decimal) ComparedTo(y *Decimal) int {
	cmp, _ := x.Cmp(y)
	return cmp
}

// Clamp returns a new Decimal clamped to [min, max].
func (x *Decimal) Clamp(min, max *Decimal) *Decimal {
	if x.Lt(min) {
		return min.copy()
	}
	if x.Gt(max) {
		return max.copy()
	}
	return x.copy()
}
