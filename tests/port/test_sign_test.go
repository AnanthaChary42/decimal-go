package port_test

import (
	"math"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// sign is a helper that replicates Decimal.sign(n) from JavaScript.
// Returns 1 for positive, -1 for negative, 0 for zero, NaN for NaN.
// In Go, we return (int, bool) where bool=false means NaN.
func sign(d *decimal.Decimal) (int, bool) {
	if d.IsNaN() {
		return 0, false // NaN
	}
	if d.IsZero() {
		if d.IsNeg() {
			return 0, true // -0 → Go doesn't have -0 int, but JS returns -0
		}
		return 0, true // +0
	}
	if d.IsNeg() {
		return -1, true
	}
	return 1, true
}

// TestOriginal_Sign ports all assertions from test/modules/sign.js
func TestOriginal_Sign(t *testing.T) {
	// JS: Decimal.sign(n)
	// The JS function returns NaN for NaN, 1 for positive, -1 for negative, 0/-0 for zero.

	// t(NaN, NaN) — sign of NaN is NaN
	d, _ := decimal.New("NaN")
	if _, ok := sign(d); ok {
		t.Error("sign(NaN) should return NaN (ok=false)")
	}

	// t('NaN', NaN) — string NaN
	d, _ = decimal.New("NaN")
	if _, ok := sign(d); ok {
		t.Error("sign('NaN') should return NaN (ok=false)")
	}

	// t(Infinity, 1)
	d, _ = decimal.New("Infinity")
	if v, ok := sign(d); !ok || v != 1 {
		t.Errorf("sign(Infinity) = %d, ok=%v, want 1", v, ok)
	}

	// t(-Infinity, -1)
	d, _ = decimal.New("-Infinity")
	if v, ok := sign(d); !ok || v != -1 {
		t.Errorf("sign(-Infinity) = %d, ok=%v, want -1", v, ok)
	}

	// t('Infinity', 1)
	d, _ = decimal.New("Infinity")
	if v, ok := sign(d); !ok || v != 1 {
		t.Errorf("sign('Infinity') = %d, ok=%v, want 1", v, ok)
	}

	// t('-Infinity', -1)
	d, _ = decimal.New("-Infinity")
	if v, ok := sign(d); !ok || v != -1 {
		t.Errorf("sign('-Infinity') = %d, ok=%v, want -1", v, ok)
	}

	// T.assert(1 / Decimal.sign('0') === Infinity)  →  sign('0') is +0
	d, _ = decimal.New("0")
	if v, ok := sign(d); !ok || v != 0 {
		t.Errorf("sign('0') = %d, ok=%v, want 0", v, ok)
	}
	if d.IsNeg() {
		t.Error("sign('0') should be positive zero")
	}

	// T.assert(1 / Decimal.sign('-0') === -Infinity)  →  sign('-0') is -0
	d, _ = decimal.New("-0")
	if v, ok := sign(d); !ok || v != 0 {
		t.Errorf("sign('-0') = %d, ok=%v, want 0", v, ok)
	}
	// Verify it is negative zero by checking the sign bit
	if !d.IsNeg() {
		t.Error("sign('-0') should be negative zero")
	}

	// t('0', 0)
	d, _ = decimal.New("0")
	if v, ok := sign(d); !ok || v != 0 {
		t.Errorf("sign('0') = %d, want 0", v)
	}

	// t('-0', -0)  →  sign('-0') should be 0 (but negative zero)
	d, _ = decimal.New("-0")
	if v, ok := sign(d); !ok || v != 0 {
		t.Errorf("sign('-0') = %d, want 0", v)
	}

	// t('1', 1)
	d, _ = decimal.New("1")
	if v, ok := sign(d); !ok || v != 1 {
		t.Errorf("sign('1') = %d, want 1", v)
	}

	// t('-1', -1)
	d, _ = decimal.New("-1")
	if v, ok := sign(d); !ok || v != -1 {
		t.Errorf("sign('-1') = %d, want -1", v)
	}

	// t('9.99', 1)
	d, _ = decimal.New("9.99")
	if v, ok := sign(d); !ok || v != 1 {
		t.Errorf("sign('9.99') = %d, want 1", v)
	}

	// t('-9.99', -1)
	d, _ = decimal.New("-9.99")
	if v, ok := sign(d); !ok || v != -1 {
		t.Errorf("sign('-9.99') = %d, want -1", v)
	}

	// Suppress unused import warning
	_ = math.Copysign
}
