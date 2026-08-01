package decimal

// Ceil returns a new Decimal rounded up to a whole number.
func (x *Decimal) Ceil() *Decimal {
	return finalise(x.copy(), x.e+1, RoundCeil)
}

// Floor returns a new Decimal rounded down to a whole number.
func (x *Decimal) Floor() *Decimal {
	return finalise(x.copy(), x.e+1, RoundFloor)
}

// Round returns a new Decimal rounded to a whole number using the context's rounding mode.
func (x *Decimal) Round() *Decimal {
	ctx := x.getContext()
	return finalise(x.copy(), x.e+1, ctx.Rounding)
}

// Trunc returns a new Decimal truncated to a whole number (towards zero).
func (x *Decimal) Trunc() *Decimal {
	return finalise(x.copy(), x.e+1, RoundDown)
}

// ToDP returns a new Decimal rounded to dp decimal places.
// If rm is not provided, uses the context's rounding mode.
func (x *Decimal) ToDP(dp int, rm ...RoundingMode) *Decimal {
	r := x.copy()
	ctx := x.getContext()

	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	return finalise(r, dp+r.e+1, rounding)
}

// ToSD returns a new Decimal rounded to sd significant digits.
// If rm is not provided, uses the context's rounding mode.
func (x *Decimal) ToSD(sd int, rm ...RoundingMode) *Decimal {
	ctx := x.getContext()

	rounding := ctx.Rounding
	if len(rm) > 0 {
		rounding = rm[0]
	}

	if sd <= 0 {
		sd = ctx.Precision
	}

	return finalise(x.copy(), sd, rounding)
}
