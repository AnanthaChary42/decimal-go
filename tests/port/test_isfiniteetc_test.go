package port_test

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_IsFiniteEtc ports all assertions from test/modules/isFiniteEtc.js
// JS: T.assert(actual) — checks truthiness of boolean expressions
// Config: precision: 20, rounding: 4, toExpNeg: -7, toExpPos: 21
func TestOriginal_IsFiniteEtc(t *testing.T) {
	// Helper to create Decimal
	d := func(s string) *decimal.Decimal {
		v, _ := decimal.New(s)
		return v
	}

	// === n = new Decimal(1) ===
	n := d("1")
	assertTrue(t, n.IsFinite(), "Decimal(1).IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal(1).IsNaN()")
	assertTrue(t, !n.IsNeg(), "!Decimal(1).IsNeg()")
	assertTrue(t, !n.IsZero(), "!Decimal(1).IsZero()")
	assertTrue(t, n.IsInt(), "Decimal(1).IsInt()")
	assertTrue(t, n.Eq(n), "Decimal(1).Eq(Decimal(1))")
	assertTrue(t, n.Eq(d("1")), "Decimal(1).Eq(1)")
	assertTrue(t, n.Eq(d("1.0")), "Decimal(1).Eq('1.0')")
	assertTrue(t, n.Eq(d("1.00")), "Decimal(1).Eq('1.00')")
	assertTrue(t, n.Eq(d("1.000")), "Decimal(1).Eq('1.000')")
	assertTrue(t, n.Eq(d("1.0000")), "Decimal(1).Eq('1.0000')")
	assertTrue(t, n.Eq(d("1.00000")), "Decimal(1).Eq('1.00000')")
	assertTrue(t, n.Eq(d("1.000000")), "Decimal(1).Eq('1.000000')")
	assertTrue(t, n.Eq(d("1")), "Decimal(1).Eq(new Decimal(1))")
	assertTrue(t, n.Eq(d("0x1")), "Decimal(1).Eq('0x1')")
	assertTrue(t, n.Eq(d("0o1")), "Decimal(1).Eq('0o1')")
	assertTrue(t, n.Eq(d("0b1")), "Decimal(1).Eq('0b1')")
	assertTrue(t, n.Gt(d("0.99999")), "Decimal(1).Gt(0.99999)")
	assertTrue(t, !n.Gte(d("1.1")), "!Decimal(1).Gte(1.1)")
	assertTrue(t, n.Lt(d("1.001")), "Decimal(1).Lt(1.001)")
	assertTrue(t, n.Lte(d("2")), "Decimal(1).Lte(2)")

	// === n = new Decimal('-0.1') ===
	n = d("-0.1")
	assertTrue(t, n.IsFinite(), "Decimal('-0.1').IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal('-0.1').IsNaN()")
	assertTrue(t, n.IsNeg(), "Decimal('-0.1').IsNeg()")
	assertTrue(t, !n.IsZero(), "!Decimal('-0.1').IsZero()")
	assertTrue(t, !n.IsInt(), "!Decimal('-0.1').IsInt()")
	assertTrue(t, !n.Eq(d("0.1")), "!Decimal('-0.1').Eq(0.1)")
	assertTrue(t, !n.Gt(d("-0.1")), "!Decimal('-0.1').Gt(-0.1)")
	assertTrue(t, n.Gte(d("-1")), "Decimal('-0.1').Gte(-1)")
	assertTrue(t, n.Lt(d("-0.01")), "Decimal('-0.1').Lt(-0.01)")
	assertTrue(t, !n.Lte(d("-1")), "!Decimal('-0.1').Lte(-1)")

	// === n = new Decimal(Infinity) ===
	n = d("Infinity")
	assertTrue(t, !n.IsFinite(), "!Decimal(Infinity).IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal(Infinity).IsNaN()")
	assertTrue(t, !n.IsNeg(), "!Decimal(Infinity).IsNeg()")
	assertTrue(t, !n.IsZero(), "!Decimal(Infinity).IsZero()")
	assertTrue(t, !n.IsInt(), "!Decimal(Infinity).IsInt()")
	assertTrue(t, n.Eq(d("Infinity")), "Decimal(Infinity).Eq('Infinity')")
	assertTrue(t, n.Gt(d("9e999")), "Decimal(Infinity).Gt('9e999')")
	assertTrue(t, n.Gte(d("Infinity")), "Decimal(Infinity).Gte(Infinity)")
	assertTrue(t, !n.Lt(d("Infinity")), "!Decimal(Infinity).Lt(Infinity)")
	assertTrue(t, n.Lte(d("Infinity")), "Decimal(Infinity).Lte(Infinity)")

	// === n = new Decimal('-Infinity') ===
	n = d("-Infinity")
	assertTrue(t, !n.IsFinite(), "!Decimal('-Infinity').IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal('-Infinity').IsNaN()")
	assertTrue(t, n.IsNeg(), "Decimal('-Infinity').IsNeg()")
	assertTrue(t, !n.IsZero(), "!Decimal('-Infinity').IsZero()")
	assertTrue(t, !n.IsInt(), "!Decimal('-Infinity').IsInt()")
	assertTrue(t, !n.Eq(d("Infinity")), "!Decimal('-Infinity').Eq(Infinity)")
	assertTrue(t, n.Eq(d("-Infinity")), "Decimal('-Infinity').Eq(-Infinity)")
	assertTrue(t, !n.Gt(d("-Infinity")), "!Decimal('-Infinity').Gt(-Infinity)")
	assertTrue(t, n.Gte(d("-Infinity")), "Decimal('-Infinity').Gte('-Infinity')")
	assertTrue(t, n.Lt(d("0")), "Decimal('-Infinity').Lt(0)")
	assertTrue(t, n.Lte(d("Infinity")), "Decimal('-Infinity').Lte(Infinity)")

	// === n = new Decimal('0.0000000') ===
	n = d("0.0000000")
	assertTrue(t, n.IsFinite(), "Decimal('0.0000000').IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal('0.0000000').IsNaN()")
	assertTrue(t, !n.IsNeg(), "!Decimal('0.0000000').IsNeg()")
	assertTrue(t, n.IsZero(), "Decimal('0.0000000').IsZero()")
	assertTrue(t, n.IsInt(), "Decimal('0.0000000').IsInt()")
	assertTrue(t, n.Eq(d("-0")), "Decimal('0.0000000').Eq(-0)")
	assertTrue(t, n.Gt(d("-0.000001")), "Decimal('0.0000000').Gt(-0.000001)")
	assertTrue(t, !n.Gte(d("0.1")), "!Decimal('0.0000000').Gte(0.1)")
	assertTrue(t, n.Lt(d("0.0001")), "Decimal('0.0000000').Lt(0.0001)")
	assertTrue(t, n.Lte(d("-0")), "Decimal('0.0000000').Lte(-0)")

	// === n = new Decimal(-0) → '-0' ===
	n = d("-0")
	assertTrue(t, n.IsFinite(), "Decimal(-0).IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal(-0).IsNaN()")
	assertTrue(t, n.IsNeg(), "Decimal(-0).IsNeg()")
	assertTrue(t, n.IsZero(), "Decimal(-0).IsZero()")
	assertTrue(t, n.IsInt(), "Decimal(-0).IsInt()")
	assertTrue(t, n.Eq(d("0.000")), "Decimal(-0).Eq('0.000')")
	assertTrue(t, n.Gt(d("-1")), "Decimal(-0).Gt(-1)")
	assertTrue(t, !n.Gte(d("0.1")), "!Decimal(-0).Gte(0.1)")
	assertTrue(t, !n.Lt(d("0")), "!Decimal(-0).Lt(0)")
	assertTrue(t, n.Lt(d("0.1")), "Decimal(-0).Lt(0.1)")
	assertTrue(t, n.Lte(d("0")), "Decimal(-0).Lte(0)")
	// t(n.valueOf() === '-0' && n.toString() === '0')
	if n.ValueOf() != "-0" {
		t.Errorf("Decimal(-0).ValueOf() = %q, want '-0'", n.ValueOf())
	}
	if n.String() != "0" {
		t.Errorf("Decimal(-0).String() = %q, want '0'", n.String())
	}

	// === n = new Decimal('NaN') ===
	n = d("NaN")
	assertTrue(t, !n.IsFinite(), "!Decimal('NaN').IsFinite()")
	assertTrue(t, n.IsNaN(), "Decimal('NaN').IsNaN()")
	assertTrue(t, !n.IsNeg(), "!Decimal('NaN').IsNeg()")
	assertTrue(t, !n.IsZero(), "!Decimal('NaN').IsZero()")
	assertTrue(t, !n.IsInt(), "!Decimal('NaN').IsInt()")
	assertTrue(t, !n.Eq(d("NaN")), "!Decimal('NaN').Eq(NaN)")
	assertTrue(t, !n.Eq(d("Infinity")), "!Decimal('NaN').Eq(Infinity)")
	assertTrue(t, !n.Gt(d("0")), "!Decimal('NaN').Gt(0)")
	assertTrue(t, !n.Gte(d("0")), "!Decimal('NaN').Gte(0)")
	assertTrue(t, !n.Lt(d("1")), "!Decimal('NaN').Lt(1)")
	assertTrue(t, !n.Lte(d("-0")), "!Decimal('NaN').Lte(-0)")
	assertTrue(t, !n.Lte(d("-1")), "!Decimal('NaN').Lte(-1)")

	// === n = new Decimal('-1.234e+2') ===
	n = d("-1.234e+2")
	assertTrue(t, n.IsFinite(), "Decimal('-1.234e+2').IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal('-1.234e+2').IsNaN()")
	assertTrue(t, n.IsNeg(), "Decimal('-1.234e+2').IsNeg()")
	assertTrue(t, !n.IsZero(), "!Decimal('-1.234e+2').IsZero()")
	assertTrue(t, !n.IsInt(), "!Decimal('-1.234e+2').IsInt()")
	assertTrue(t, n.Eq(d("-123.4")), "Decimal('-1.234e+2').Eq(-123.4)")
	assertTrue(t, n.Gt(d("-0xff")), "Decimal('-1.234e+2').Gt('-0xff')")
	assertTrue(t, n.Gte(d("-1.234e+3")), "Decimal('-1.234e+2').Gte('-1.234e+3')")
	assertTrue(t, n.Lt(d("-123.39999")), "Decimal('-1.234e+2').Lt(-123.39999)")
	assertTrue(t, n.Lte(d("-123.4e+0")), "Decimal('-1.234e+2').Lte('-123.4e+0')")

	// === n = new Decimal('5e-200') ===
	n = d("5e-200")
	assertTrue(t, n.IsFinite(), "Decimal('5e-200').IsFinite()")
	assertTrue(t, !n.IsNaN(), "!Decimal('5e-200').IsNaN()")
	assertTrue(t, !n.IsNeg(), "!Decimal('5e-200').IsNeg()")
	assertTrue(t, !n.IsZero(), "!Decimal('5e-200').IsZero()")
	assertTrue(t, !n.IsInt(), "!Decimal('5e-200').IsInt()")
	assertTrue(t, n.Eq(d("5e-200")), "Decimal('5e-200').Eq(5e-200)")
	assertTrue(t, n.Gt(d("5e-201")), "Decimal('5e-200').Gt(5e-201)")
	assertTrue(t, !n.Gte(d("1")), "!Decimal('5e-200').Gte(1)")
	assertTrue(t, n.Lt(d("6e-200")), "Decimal('5e-200').Lt(6e-200)")
	assertTrue(t, n.Lte(d("5.1e-200")), "Decimal('5e-200').Lte(5.1e-200)")
}

// TestOriginal_IsFiniteEtc_Equals ports the equals comparison section
func TestOriginal_IsFiniteEtc_Equals(t *testing.T) {
	d := func(s string) *decimal.Decimal {
		v, _ := decimal.New(s)
		return v
	}

	n := d("1")
	assertTrue(t, n.Eq(n), "n.Eq(n)")
	assertTrue(t, n.Eq(d("1")), "n.Eq(1)")
	assertTrue(t, n.Eq(d("1e+0")), "n.Eq('1e+0')")
	assertTrue(t, !n.Eq(d("-1")), "!n.Eq(-1)")
	assertTrue(t, !n.Eq(d("0.1")), "!n.Eq(0.1)")

	assertTrue(t, !d("NaN").Eq(d("0")), "!NaN.Eq(0)")
	assertTrue(t, !d("Infinity").Eq(d("0")), "!Infinity.Eq(0)")
	assertTrue(t, !d("0.1").Eq(d("0")), "!0.1.Eq(0)")
	assertTrue(t, !d("1000000001").Eq(d("1000000000")), "!1e9+1.Eq(1e9)")
	assertTrue(t, !d("999999999").Eq(d("1000000000")), "!1e9-1.Eq(1e9)")
	assertTrue(t, d("1000000001").Eq(d("1000000001")), "1e9+1.Eq(1e9+1)")
	assertTrue(t, d("1").Eq(d("1")), "1.Eq(1)")
	assertTrue(t, !d("1").Eq(d("-1")), "!1.Eq(-1)")
	assertTrue(t, !d("NaN").Eq(d("NaN")), "!NaN.Eq(NaN)")
}

// TestOriginal_IsFiniteEtc_Comparisons ports the comparison methods section
func TestOriginal_IsFiniteEtc_Comparisons(t *testing.T) {
	d := func(s string) *decimal.Decimal {
		v, _ := decimal.New(s)
		return v
	}

	assertTrue(t, !d("NaN").Gt(d("NaN")), "!NaN.Gt(NaN)")
	assertTrue(t, !d("NaN").Lt(d("NaN")), "!NaN.Lt(NaN)")
	assertTrue(t, d("0xa").Lte(d("0xff")), "0xa.Lte(0xff)")
	assertTrue(t, d("0xb").Gte(d("0x9")), "0xb.Gte(0x9)")

	assertTrue(t, !d("10").Gt(d("10")), "!10.Gt(10)")
	assertTrue(t, !d("10").Lt(d("10")), "!10.Lt(10)")
	assertTrue(t, !d("NaN").Lt(d("NaN")), "!NaN.Lt(NaN)")
	assertTrue(t, !d("Infinity").Lt(d("-Infinity")), "!Inf.Lt(-Inf)")
	assertTrue(t, !d("Infinity").Lt(d("Infinity")), "!Inf.Lt(Inf)")
	assertTrue(t, d("Infinity").Lte(d("Infinity")), "Inf.Lte(Inf)")
	assertTrue(t, !d("NaN").Gte(d("NaN")), "!NaN.Gte(NaN)")
	assertTrue(t, d("Infinity").Gte(d("Infinity")), "Inf.Gte(Inf)")
	assertTrue(t, d("Infinity").Gte(d("-Infinity")), "Inf.Gte(-Inf)")
	assertTrue(t, !d("NaN").Gte(d("-Infinity")), "!NaN.Gte(-Inf)")
	assertTrue(t, d("-Infinity").Gte(d("-Infinity")), "-Inf.Gte(-Inf)")

	assertTrue(t, !d("2").Gt(d("10")), "!2.Gt(10)")
	assertTrue(t, !d("10").Lt(d("2")), "!10.Lt(2)")
	assertTrue(t, d("255").Lte(d("0xff")), "255.Lte(0xff)")
	assertTrue(t, d("0xa").Gte(d("0x9")), "0xa.Gte(0x9)")
	assertTrue(t, !d("0").Lte(d("NaN")), "!0.Lte(NaN)")
	assertTrue(t, !d("0").Gte(d("NaN")), "!0.Gte(NaN)")
	assertTrue(t, !d("NaN").Lte(d("NaN")), "!NaN.Lte(NaN)")
	assertTrue(t, !d("NaN").Gte(d("NaN")), "!NaN.Gte(NaN)")
	assertTrue(t, !d("0").Lte(d("-Infinity")), "!0.Lte(-Inf)")
	assertTrue(t, d("0").Gte(d("-Infinity")), "0.Gte(-Inf)")
	assertTrue(t, d("0").Lte(d("Infinity")), "0.Lte(Inf)")
	assertTrue(t, !d("0").Gte(d("Infinity")), "!0.Gte(Inf)")
	assertTrue(t, d("10").Lte(d("20")), "10.Lte(20)")
	assertTrue(t, !d("10").Gte(d("20")), "!10.Gte(20)")

	// Precision comparisons with small exponents
	assertTrue(t, !d("1.23001e-2").Lt(d("1.23e-2")), "!1.23001e-2.Lt(1.23e-2)")
	assertTrue(t, d("1.23e-2").Lt(d("1.23001e-2")), "1.23e-2.Lt(1.23001e-2)")
	assertTrue(t, !d("1e-2").Lt(d("9.999999e-3")), "!1e-2.Lt(9.999999e-3)")
	assertTrue(t, d("9.999999e-3").Lt(d("1e-2")), "9.999999e-3.Lt(1e-2)")

	assertTrue(t, !d("1.23001e+2").Lt(d("1.23e+2")), "!1.23001e+2.Lt(1.23e+2)")
	assertTrue(t, d("1.23e+2").Lt(d("1.23001e+2")), "1.23e+2.Lt(1.23001e+2)")
	assertTrue(t, d("9.999999e+2").Lt(d("1e+3")), "9.999999e+2.Lt(1e+3)")
	assertTrue(t, !d("1e+3").Lt(d("9.9999999e+2")), "!1e+3.Lt(9.9999999e+2)")

	// lte with precision
	assertTrue(t, !d("1.23001e-2").Lte(d("1.23e-2")), "!1.23001e-2.Lte(1.23e-2)")
	assertTrue(t, d("1.23e-2").Lte(d("1.23001e-2")), "1.23e-2.Lte(1.23001e-2)")
	assertTrue(t, !d("1e-2").Lte(d("9.999999e-3")), "!1e-2.Lte(9.999999e-3)")
	assertTrue(t, d("9.999999e-3").Lte(d("1e-2")), "9.999999e-3.Lte(1e-2)")

	assertTrue(t, !d("1.23001e+2").Lte(d("1.23e+2")), "!1.23001e+2.Lte(1.23e+2)")
	assertTrue(t, d("1.23e+2").Lte(d("1.23001e+2")), "1.23e+2.Lte(1.23001e+2)")
	assertTrue(t, d("9.999999e+2").Lte(d("1e+3")), "9.999999e+2.Lte(1e+3)")
	assertTrue(t, !d("1e+3").Lte(d("9.9999999e+2")), "!1e+3.Lte(9.9999999e+2)")

	// gt with precision
	assertTrue(t, d("1.23001e-2").Gt(d("1.23e-2")), "1.23001e-2.Gt(1.23e-2)")
	assertTrue(t, !d("1.23e-2").Gt(d("1.23001e-2")), "!1.23e-2.Gt(1.23001e-2)")
	assertTrue(t, d("1e-2").Gt(d("9.999999e-3")), "1e-2.Gt(9.999999e-3)")
	assertTrue(t, !d("9.999999e-3").Gt(d("1e-2")), "!9.999999e-3.Gt(1e-2)")

	assertTrue(t, d("1.23001e+2").Gt(d("1.23e+2")), "1.23001e+2.Gt(1.23e+2)")
	assertTrue(t, !d("1.23e+2").Gt(d("1.23001e+2")), "!1.23e+2.Gt(1.23001e+2)")
	assertTrue(t, !d("9.999999e+2").Gt(d("1e+3")), "!9.999999e+2.Gt(1e+3)")
	assertTrue(t, d("1e+3").Gt(d("9.9999999e+2")), "1e+3.Gt(9.9999999e+2)")

	// gte with precision
	assertTrue(t, d("1.23001e-2").Gte(d("1.23e-2")), "1.23001e-2.Gte(1.23e-2)")
	assertTrue(t, !d("1.23e-2").Gte(d("1.23001e-2")), "!1.23e-2.Gte(1.23001e-2)")
	assertTrue(t, d("1e-2").Gte(d("9.999999e-3")), "1e-2.Gte(9.999999e-3)")
	assertTrue(t, !d("9.999999e-3").Gte(d("1e-2")), "!9.999999e-3.Gte(1e-2)")

	assertTrue(t, d("1.23001e+2").Gte(d("1.23e+2")), "1.23001e+2.Gte(1.23e+2)")
	assertTrue(t, !d("1.23e+2").Gte(d("1.23001e+2")), "!1.23e+2.Gte(1.23001e+2)")
	assertTrue(t, !d("9.999999e+2").Gte(d("1e+3")), "!9.999999e+2.Gte(1e+3)")
	assertTrue(t, d("1e+3").Gte(d("9.9999999e+2")), "1e+3.Gte(9.9999999e+2)")
}

// TestOriginal_IsFiniteEtc_IsInteger ports the isInteger section
func TestOriginal_IsFiniteEtc_IsInteger(t *testing.T) {
	d := func(s string) *decimal.Decimal {
		v, _ := decimal.New(s)
		return v
	}

	assertTrue(t, !d("1.0000000000000000000001").IsInt(), "!1.000...001.IsInt()")
	assertTrue(t, !d("0.999999999999999999999").IsInt(), "!0.999...9.IsInt()")
	assertTrue(t, d("4e4").IsInt(), "4e4.IsInt()")
	assertTrue(t, d("-4e4").IsInt(), "-4e4.IsInt()")
}

// assertTrue is a helper for asserting boolean conditions
func assertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("assertion failed: %s", msg)
	}
}
