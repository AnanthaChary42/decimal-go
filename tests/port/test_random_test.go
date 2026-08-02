package port_test

import (
	"math/rand"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_Random(t *testing.T) {
	saveDefaultCtx(t)

	ctx := decimal.GetDefaultContext()
	maxDigits := 100
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < 996; i++ {
		sd := rng.Intn(maxDigits) + 1

		var r *decimal.Decimal
		if rng.Float64() > 0.5 {
			ctx.Precision = sd
			r = ctx.Random()
		} else {
			r = ctx.Random(sd)
		}

		// r.sd() <= sd && r >= 0 && r < 1 && r.eq(r) && r.eq(r.valueOf())
		rSD := r.SD()
		if rSD > sd {
			t.Errorf("iter %d: sd()=%d > %d", i, rSD, sd)
		}
		zero := newDec(t, 0)
		one := newDec(t, 1)
		if !r.Gte(zero) {
			t.Errorf("iter %d: r < 0: %s", i, r.ValueOf())
		}
		if !r.Lt(one) {
			t.Errorf("iter %d: r >= 1: %s", i, r.ValueOf())
		}
		if !r.Eq(r) {
			t.Errorf("iter %d: r != r", i)
		}
		// r.eq(r.valueOf())
		rStr := r.ValueOf()
		rFromStr, err := decimal.New(rStr)
		if err != nil {
			t.Fatalf("iter %d: New(%q): %v", i, rStr, err)
		}
		if !r.Eq(rFromStr) {
			t.Errorf("iter %d: r != New(r.valueOf())", i)
		}
	}

	// Invalid argument panics.
	t.Run("invalid_args", func(t *testing.T) {
		assertPanics(t, func() { ctx.Random(0) }, "Random(0)")
		assertPanics(t, func() { ctx.Random(-1) }, "Random(-1)")
		assertPanics(t, func() { ctx.Random(1000000001) }, "Random(1e9+1)")
	})
}
