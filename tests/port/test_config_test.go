package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_Config(t *testing.T) {
	saveDefaultCtx(t)

	ctx := decimal.GetDefaultContext()

	// Decimal.config({}) should return without error.
	decimal.Config(decimal.ConfigOptions{})

	// Set up initial config.
	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(20),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-7),
		ToExpPos:  decimal.IntPtr(21),
		MinE:      decimal.IntPtr(-9000000000000000),
		MaxE:      decimal.IntPtr(9000000000000000),
		Crypto:    decimal.BoolPtr(false),
		Modulo:    decimal.IntPtr(1),
	})

	t.Run("initial_config_readback", func(t *testing.T) {
		if ctx.Precision != 20 { t.Errorf("Precision: %d", ctx.Precision) }
		if int(ctx.Rounding) != 4 { t.Errorf("Rounding: %d", ctx.Rounding) }
		if ctx.ToExpNeg != -7 { t.Errorf("ToExpNeg: %d", ctx.ToExpNeg) }
		if ctx.ToExpPos != 21 { t.Errorf("ToExpPos: %d", ctx.ToExpPos) }
		if ctx.MinE != -9000000000000000 { t.Errorf("MinE: %d", ctx.MinE) }
		if ctx.MaxE != 9000000000000000 { t.Errorf("MaxE: %d", ctx.MaxE) }
		if ctx.Crypto != false { t.Errorf("Crypto: %v", ctx.Crypto) }
		if int(ctx.Modulo) != 1 { t.Errorf("Modulo: %d", ctx.Modulo) }
	})

	// Second config.
	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(40),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-1000),
		ToExpPos:  decimal.IntPtr(1000),
		MinE:      decimal.IntPtr(-1000000000),
		MaxE:      decimal.IntPtr(1000000000),
		Modulo:    decimal.IntPtr(4),
	})

	t.Run("second_config_readback", func(t *testing.T) {
		if ctx.Precision != 40 { t.Errorf("Precision: %d", ctx.Precision) }
		if int(ctx.Rounding) != 4 { t.Errorf("Rounding: %d", ctx.Rounding) }
		if ctx.ToExpNeg != -1000 { t.Errorf("ToExpNeg: %d", ctx.ToExpNeg) }
		if ctx.ToExpPos != 1000 { t.Errorf("ToExpPos: %d", ctx.ToExpPos) }
		if ctx.MinE != -1000000000 { t.Errorf("MinE: %d", ctx.MinE) }
		if ctx.MaxE != 1000000000 { t.Errorf("MaxE: %d", ctx.MaxE) }
		if int(ctx.Modulo) != 4 { t.Errorf("Modulo: %d", ctx.Modulo) }
	})

	decimal.Config(decimal.ConfigOptions{
		ToExpNeg: decimal.IntPtr(-7),
		ToExpPos: decimal.IntPtr(21),
		MinE:     decimal.IntPtr(-324),
		MaxE:     decimal.IntPtr(308),
	})
	t.Run("third_config_readback", func(t *testing.T) {
		if ctx.ToExpNeg != -7 { t.Errorf("ToExpNeg: %d", ctx.ToExpNeg) }
		if ctx.ToExpPos != 21 { t.Errorf("ToExpPos: %d", ctx.ToExpPos) }
		if ctx.MinE != -324 { t.Errorf("MinE: %d", ctx.MinE) }
		if ctx.MaxE != 308 { t.Errorf("MaxE: %d", ctx.MaxE) }
	})

	// --- precision ---
	t.Run("precision_valid", func(t *testing.T) {
		for _, v := range []int{1, 20, 300000, 400000000, 1000000000} {
			decimal.Config(decimal.ConfigOptions{Precision: decimal.IntPtr(v)})
			if ctx.Precision != v {
				t.Errorf("expected precision %d, got %d", v, ctx.Precision)
			}
		}
	})

	t.Run("precision_invalid", func(t *testing.T) {
		invalids := []int{0, 1000000001}
		for _, v := range invalids {
			assertPanics(t, func() {
				decimal.Config(decimal.ConfigOptions{Precision: decimal.IntPtr(v)})
			}, "precision")
		}
	})

	// Undefined (nil) should leave precision unchanged.
	t.Run("precision_nil_unchanged", func(t *testing.T) {
		decimal.Config(decimal.ConfigOptions{Precision: decimal.IntPtr(1000000000)})
		decimal.Config(decimal.ConfigOptions{}) // nil = undefined
		if ctx.Precision != 1000000000 {
			t.Errorf("expected precision 1000000000, got %d", ctx.Precision)
		}
	})

	// --- rounding ---
	t.Run("rounding_valid", func(t *testing.T) {
		for v := 0; v <= 8; v++ {
			decimal.Config(decimal.ConfigOptions{Rounding: decimal.IntPtr(v)})
			if int(ctx.Rounding) != v {
				t.Errorf("expected rounding %d, got %d", v, ctx.Rounding)
			}
		}
	})

	t.Run("rounding_invalid", func(t *testing.T) {
		for _, v := range []int{-1, 9, 11} {
			assertPanics(t, func() {
				decimal.Config(decimal.ConfigOptions{Rounding: decimal.IntPtr(v)})
			}, "rounding")
		}
	})

	// --- toExpNeg ---
	t.Run("toExpNeg_valid", func(t *testing.T) {
		for _, v := range []int{0, -1, -999, -5675367, -98770170790791, -9000000000000000} {
			decimal.Config(decimal.ConfigOptions{ToExpNeg: decimal.IntPtr(v)})
			if ctx.ToExpNeg != v {
				t.Errorf("expected toExpNeg %d, got %d", v, ctx.ToExpNeg)
			}
		}
	})

	t.Run("toExpNeg_invalid", func(t *testing.T) {
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{ToExpNeg: decimal.IntPtr(-9000000000000001)})
		}, "toExpNeg too low")
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{ToExpNeg: decimal.IntPtr(1)})
		}, "toExpNeg positive")
	})

	// --- toExpPos ---
	t.Run("toExpPos_valid", func(t *testing.T) {
		for _, v := range []int{0, 1, 999, 5675367, 98770170790791, 9000000000000000} {
			decimal.Config(decimal.ConfigOptions{ToExpPos: decimal.IntPtr(v)})
			if ctx.ToExpPos != v {
				t.Errorf("expected toExpPos %d, got %d", v, ctx.ToExpPos)
			}
		}
	})

	t.Run("toExpPos_invalid", func(t *testing.T) {
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{ToExpPos: decimal.IntPtr(9000000000000001)})
		}, "toExpPos too high")
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{ToExpPos: decimal.IntPtr(-1)})
		}, "toExpPos negative")
	})

	// --- maxE ---
	t.Run("maxE_valid", func(t *testing.T) {
		for _, v := range []int{0, 1, 999, 5675367, 98770170790791, 9000000000000000} {
			decimal.Config(decimal.ConfigOptions{MaxE: decimal.IntPtr(v)})
			if ctx.MaxE != v {
				t.Errorf("expected maxE %d, got %d", v, ctx.MaxE)
			}
		}
	})

	t.Run("maxE_invalid", func(t *testing.T) {
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{MaxE: decimal.IntPtr(9000000000000001)})
		}, "maxE too high")
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{MaxE: decimal.IntPtr(-1)})
		}, "maxE negative")
	})

	// --- minE ---
	t.Run("minE_valid", func(t *testing.T) {
		for _, v := range []int{0, -1, -999, -5675367, -98770170790791, -9000000000000000} {
			decimal.Config(decimal.ConfigOptions{MinE: decimal.IntPtr(v)})
			if ctx.MinE != v {
				t.Errorf("expected minE %d, got %d", v, ctx.MinE)
			}
		}
	})

	t.Run("minE_invalid", func(t *testing.T) {
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{MinE: decimal.IntPtr(-9000000000000001)})
		}, "minE too low")
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{MinE: decimal.IntPtr(1)})
		}, "minE positive")
	})

	// --- crypto ---
	t.Run("crypto_valid", func(t *testing.T) {
		decimal.Config(decimal.ConfigOptions{Crypto: decimal.BoolPtr(false)})
		if ctx.Crypto != false { t.Errorf("Crypto: %v", ctx.Crypto) }
	})

	// --- modulo ---
	t.Run("modulo_valid", func(t *testing.T) {
		for v := 0; v <= 9; v++ {
			decimal.Config(decimal.ConfigOptions{Modulo: decimal.IntPtr(v)})
			if int(ctx.Modulo) != v {
				t.Errorf("expected modulo %d, got %d", v, ctx.Modulo)
			}
		}
	})

	t.Run("modulo_invalid", func(t *testing.T) {
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{Modulo: decimal.IntPtr(-1)})
		}, "modulo: -1")
		assertPanics(t, func() {
			decimal.Config(decimal.ConfigOptions{Modulo: decimal.IntPtr(10)})
		}, "modulo: 10")
	})

	// --- defaults ---
	t.Run("defaults_reset", func(t *testing.T) {
		decimal.Config(decimal.ConfigOptions{
			Precision: decimal.IntPtr(100),
			Rounding:  decimal.IntPtr(2),
			ToExpNeg:  decimal.IntPtr(-100),
			ToExpPos:  decimal.IntPtr(200),
		})
		if ctx.Precision != 100 { t.Errorf("Precision: %d", ctx.Precision) }

		decimal.Config(decimal.ConfigOptions{Defaults: decimal.BoolPtr(true)})
		if ctx.Precision != 20 { t.Errorf("after reset Precision: %d", ctx.Precision) }
		if int(ctx.Rounding) != 4 { t.Errorf("after reset Rounding: %d", ctx.Rounding) }
		if ctx.ToExpNeg != -7 { t.Errorf("after reset ToExpNeg: %d", ctx.ToExpNeg) }
		if ctx.ToExpPos != 21 { t.Errorf("after reset ToExpPos: %d", ctx.ToExpPos) }
	})

	t.Run("defaults_with_override", func(t *testing.T) {
		ctx.Rounding = 3
		decimal.Config(decimal.ConfigOptions{
			Precision: decimal.IntPtr(50),
			Defaults:  decimal.BoolPtr(true),
		})
		if ctx.Precision != 50 { t.Errorf("Precision: %d", ctx.Precision) }
		if int(ctx.Rounding) != 4 { t.Errorf("Rounding: %d", ctx.Rounding) }
	})
}
