package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_Clone(t *testing.T) {
	saveDefaultCtx(t)

	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(10),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-7),
		ToExpPos:  decimal.IntPtr(21),
		MinE:      decimal.IntPtr(-9000000000000000),
		MaxE:      decimal.IntPtr(9000000000000000),
	})

	D1 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(1)})
	D2 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(2)})
	D3 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(3)})
	D4 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(4)})
	D5 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(5)})
	D6 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(6)})
	D7 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(7)})
	D8 := decimal.Clone()
	D8.Precision = 8
	D9 := decimal.Clone(decimal.ConfigOptions{Precision: decimal.IntPtr(9)})

	// Clone creates a distinct context.
	ctx := decimal.GetDefaultContext()
	if D9 == ctx {
		t.Error("D9 should not be the same pointer as Decimal default context")
	}

	t.Run("division_precision", func(t *testing.T) {
		type divCase struct {
			ctx      *decimal.Context
			expected string
		}
		cases := []divCase{
			{D1, "2"}, {D2, "1.7"}, {D3, "1.67"}, {D4, "1.667"},
			{D5, "1.6667"}, {D6, "1.66667"}, {D7, "1.666667"},
			{D8, "1.6666667"}, {D9, "1.66666667"},
		}
		for _, c := range cases {
			x, _ := c.ctx.New("5")
			y, _ := c.ctx.New("3")
			result := x.Div(y)
			exp, _ := c.ctx.New(c.expected)
			if !result.Eq(exp) {
				t.Errorf("5/3 at precision %d: got %s, want %s", c.ctx.Precision, result.ValueOf(), c.expected)
			}
		}
		// Default context (precision 10)
		x, _ := decimal.New("5")
		y, _ := decimal.New("3")
		result := x.Div(y)
		exp, _ := decimal.New("1.666666667")
		if !result.Eq(exp) {
			t.Errorf("5/3 at default precision 10: got %s, want 1.666666667", result.ValueOf())
		}
	})

	t.Run("precision_readback", func(t *testing.T) {
		if ctx.Precision != 10 { t.Errorf("Decimal.precision = %d", ctx.Precision) }
		if D9.Precision != 9 { t.Errorf("D9.precision = %d", D9.Precision) }
		if D8.Precision != 8 { t.Errorf("D8.precision = %d", D8.Precision) }
		if D7.Precision != 7 { t.Errorf("D7.precision = %d", D7.Precision) }
		if D6.Precision != 6 { t.Errorf("D6.precision = %d", D6.Precision) }
		if D5.Precision != 5 { t.Errorf("D5.precision = %d", D5.Precision) }
		if D4.Precision != 4 { t.Errorf("D4.precision = %d", D4.Precision) }
		if D3.Precision != 3 { t.Errorf("D3.precision = %d", D3.Precision) }
		if D2.Precision != 2 { t.Errorf("D2.precision = %d", D2.Precision) }
		if D1.Precision != 1 { t.Errorf("D1.precision = %d", D1.Precision) }
	})

	t.Run("cross_context_comparison", func(t *testing.T) {
		x, _ := decimal.New("9.99")
		y, _ := D5.New("9.99")
		if !x.Eq(y) { t.Error("Decimal(9.99) != D5(9.99)") }

		y2, _ := D3.New("-9.99")
		if x.Eq(y2) { t.Error("Decimal(9.99) == D3(-9.99)") }
	})

	// --- defaults: true ---
	t.Run("clone_defaults_true", func(t *testing.T) {
		decimal.Config(decimal.ConfigOptions{
			Precision: decimal.IntPtr(100),
			Rounding:  decimal.IntPtr(2),
			ToExpNeg:  decimal.IntPtr(-100),
			ToExpPos:  decimal.IntPtr(200),
			Defaults:  decimal.BoolPtr(true),
		})
		if ctx.Precision != 100 { t.Errorf("Precision: %d", ctx.Precision) }
		if int(ctx.Rounding) != 2 { t.Errorf("Rounding: %d", ctx.Rounding) }
		if ctx.ToExpNeg != -100 { t.Errorf("ToExpNeg: %d", ctx.ToExpNeg) }
		if ctx.ToExpPos != 200 { t.Errorf("ToExpPos: %d", ctx.ToExpPos) }

		cloneD1 := decimal.Clone(decimal.ConfigOptions{Defaults: decimal.BoolPtr(true)})
		if cloneD1.Precision != 20 { t.Errorf("D1 Precision: %d", cloneD1.Precision) }
		if int(cloneD1.Rounding) != 4 { t.Errorf("D1 Rounding: %d", cloneD1.Rounding) }
		if cloneD1.ToExpNeg != -7 { t.Errorf("D1 ToExpNeg: %d", cloneD1.ToExpNeg) }
		if cloneD1.ToExpPos != 21 { t.Errorf("D1 ToExpPos: %d", cloneD1.ToExpPos) }

		cloneD2 := decimal.Clone(decimal.ConfigOptions{
			Defaults: decimal.BoolPtr(true),
			Rounding: decimal.IntPtr(5),
		})
		if cloneD2.Precision != 20 { t.Errorf("D2 Precision: %d", cloneD2.Precision) }
		if int(cloneD2.Rounding) != 5 { t.Errorf("D2 Rounding: %d", cloneD2.Rounding) }
		if cloneD2.ToExpNeg != -7 { t.Errorf("D2 ToExpNeg: %d", cloneD2.ToExpNeg) }
		if cloneD2.ToExpPos != 21 { t.Errorf("D2 ToExpPos: %d", cloneD2.ToExpPos) }

		cloneD3 := decimal.Clone(decimal.ConfigOptions{Defaults: decimal.BoolPtr(false)})
		if int(cloneD3.Rounding) != 2 { t.Errorf("D3 Rounding: %d (expected 2, inherited from parent)", cloneD3.Rounding) }
	})
}
