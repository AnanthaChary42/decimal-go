package decimal

// Cmp compares x and y.
// Returns 1 if x > y, -1 if x < y, 0 if x == y.
// Returns 0 and ok=false if either is NaN.
func (x *Decimal) Cmp(y *Decimal) (int, bool) {
	xd := x.d
	yd := y.d
	xs := x.s
	ys := y.s

	// Either NaN or ±Infinity?
	if xd == nil || yd == nil {
		if xs == 0 || ys == 0 {
			return 0, false // NaN
		}
		if xs != ys {
			return int(xs), true
		}
		// Both same sign, both non-finite.
		xIsInf := xd == nil && xs != 0
		yIsInf := yd == nil && ys != 0
		if xIsInf && yIsInf {
			return 0, true // same infinity
		}
		if xIsInf {
			// x is ±Inf, y is finite
			if xs < 0 {
				return -1, true
			}
			return 1, true
		}
		// y is ±Inf, x is finite
		if ys < 0 {
			return 1, true
		}
		return -1, true
	}

	// Either zero?
	if xd[0] == 0 || yd[0] == 0 {
		if xd[0] != 0 {
			return int(xs), true
		}
		if yd[0] != 0 {
			return -int(ys), true
		}
		return 0, true
	}

	// Signs differ?
	if xs != ys {
		return int(xs), true
	}

	// Compare exponents.
	if x.e != y.e {
		if (x.e > y.e) != (xs < 0) {
			return 1, true
		}
		return -1, true
	}

	xdL := len(xd)
	ydL := len(yd)

	// Compare digit by digit.
	j := xdL
	if ydL < j {
		j = ydL
	}
	for i := 0; i < j; i++ {
		if xd[i] != yd[i] {
			if (xd[i] > yd[i]) != (xs < 0) {
				return 1, true
			}
			return -1, true
		}
	}

	// Compare lengths.
	if xdL == ydL {
		return 0, true
	}
	if (xdL > ydL) != (xs < 0) {
		return 1, true
	}
	return -1, true
}

// CmpInt compares x and y, returning -1, 0, or 1.
// Returns 0 for NaN comparisons (matching JS behavior where NaN comparisons
// would return NaN, but we simplify to 0).
func (x *Decimal) CmpInt(y *Decimal) int {
	r, _ := x.Cmp(y)
	return r
}

// Eq returns true if x equals y.
func (x *Decimal) Eq(y *Decimal) bool {
	r, ok := x.Cmp(y)
	return ok && r == 0
}

// Gt returns true if x > y.
func (x *Decimal) Gt(y *Decimal) bool {
	r, ok := x.Cmp(y)
	return ok && r > 0
}

// Gte returns true if x >= y.
func (x *Decimal) Gte(y *Decimal) bool {
	r, ok := x.Cmp(y)
	return ok && r >= 0
}

// Lt returns true if x < y.
func (x *Decimal) Lt(y *Decimal) bool {
	r, ok := x.Cmp(y)
	return ok && r < 0
}

// Lte returns true if x <= y.
func (x *Decimal) Lte(y *Decimal) bool {
	r, ok := x.Cmp(y)
	return ok && r <= 0
}

// Clamp returns a new Decimal clamped to [min, max].
func (x *Decimal) Clamp(min, max *Decimal) *Decimal {
	_, minOk := min.Cmp(min)
	_, maxOk := max.Cmp(max)
	if !minOk || !maxOk {
		// NaN
		d := x.copy()
		d.s = 0
		d.d = nil
		return d
	}

	if min.Gt(max) {
		// Invalid range — return NaN (JS throws).
		d := x.copy()
		d.s = 0
		d.d = nil
		return d
	}

	k, _ := x.Cmp(min)
	if k < 0 {
		return min.copy()
	}
	k2, _ := x.Cmp(max)
	if k2 > 0 {
		return max.copy()
	}
	return x.copy()
}
