package fuzz

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// FuzzDecimal tests string parsing and basic arithmetic invariants.
func FuzzDecimal(f *testing.F) {
	// Seed corpus
	f.Add("0.1", "0.2")
	f.Add("123.456", "789.012")
	f.Add("1e10", "-1e-10")
	f.Add("0", "1")

	f.Fuzz(func(t *testing.T, strA, strB string) {
		a, errA := decimal.New(strA)
		b, errB := decimal.New(strB)

		if errA != nil || errB != nil {
			return
		}

		// Invariant 1: String() parsing symmetry
		dStr, errStr := decimal.New(a.String())
		if errStr == nil && dStr.IsFinite() {
			if !dStr.Eq(a) {
				t.Fatalf("Parsing symmetry failed for %q: got %q, want %q", strA, dStr.String(), a.String())
			}
		}

		// Invariant 2: Add commutativity (a + b == b + a) for finite numbers
		if a.IsFinite() && b.IsFinite() {
			addAB := a.Plus(b)
			addBA := b.Plus(a)
			if !addAB.Eq(addBA) {
				t.Fatalf("Add commutativity failed: (%s + %s) = %s != (%s + %s) = %s",
					strA, strB, addAB.String(), strB, strA, addBA.String())
			}
		}

		// Invariant 3: Subtraction self-cancellation (a - a == 0) for finite numbers
		if a.IsFinite() {
			subAA := a.Minus(a)
			if !subAA.IsZero() {
				t.Fatalf("Self subtraction failed: (%s - %s) = %s, want 0", strA, strA, subAA.String())
			}
		}
	})
}
