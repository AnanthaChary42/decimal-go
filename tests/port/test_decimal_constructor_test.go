package port_test

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// --- Shared test helpers ---

// newDec creates a Decimal from various types, matching JS `new Decimal(n)`.
func newDec(t *testing.T, v interface{}) *decimal.Decimal {
	t.Helper()
	switch val := v.(type) {
	case string:
		d, err := decimal.New(val)
		if err != nil {
			t.Fatalf("newDec(%q): %v", val, err)
		}
		return d
	case int:
		return decimal.NewFromInt64(int64(val))
	case int64:
		return decimal.NewFromInt64(val)
	case float64:
		if math.IsNaN(val) {
			d, _ := decimal.New("NaN")
			return d
		}
		if math.IsInf(val, 1) {
			d, _ := decimal.New("Infinity")
			return d
		}
		if math.IsInf(val, -1) {
			d, _ := decimal.New("-Infinity")
			return d
		}
		if val == 0 && math.Signbit(val) {
			d, _ := decimal.New("-0")
			return d
		}
		d, err := decimal.NewFromFloat64(val)
		if err != nil {
			t.Fatalf("newDec(%v): %v", val, err)
		}
		return d
	case *decimal.Decimal:
		return val
	default:
		t.Fatalf("newDec: unsupported type %T", v)
		return nil
	}
}

// assertEqualProps replicates JS T.assertEqualProps(digits, exponent, sign, n).
func assertEqualProps(t *testing.T, d []int32, e int, s int8, x *decimal.Decimal, label string) {
	t.Helper()
	xd := x.D()
	if len(xd) != len(d) {
		t.Errorf("[%s] d length: got %v, want %v", label, xd, d)
		return
	}
	for i := range d {
		if xd[i] != d[i] {
			t.Errorf("[%s] d[%d]: got %d, want %d", label, i, xd[i], d[i])
			return
		}
	}
	if x.E() != e {
		t.Errorf("[%s] e: got %d, want %d", label, x.E(), e)
	}
	if x.S() != s {
		t.Errorf("[%s] s: got %d, want %d", label, x.S(), s)
	}
}

// assertException verifies that decimal.New(s) returns an error.
func assertException(t *testing.T, s string, label string) {
	t.Helper()
	_, err := decimal.New(s)
	if err == nil {
		t.Errorf("[%s] expected error for %q, got nil", label, s)
	}
}

// assertPanics verifies that fn panics with a DecimalError.
func assertPanics(t *testing.T, fn func(), label string) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("[%s] expected panic, got none", label)
		}
	}()
	fn()
}

// assertEqualDecimal replicates JS T.assertEqualDecimal: x.eq(y) || (x.isNaN() && y.isNaN()).
func assertEqualDecimal(t *testing.T, x, y *decimal.Decimal, label string) {
	t.Helper()
	if x.Eq(y) || (x.IsNaN() && y.IsNaN()) {
		return
	}
	t.Errorf("[%s] not equal: x=%s, y=%s", label, x.ValueOf(), y.ValueOf())
}

// saveDefaultCtx saves the current default context state and restores it after the test.
func saveDefaultCtx(t *testing.T) {
	t.Helper()
	ctx := decimal.GetDefaultContext()
	saved := *ctx
	t.Cleanup(func() {
		*ctx = saved
	})
}

func TestOriginal_DecimalConstructor(t *testing.T) {
	saveDefaultCtx(t)

	decimal.Config(decimal.ConfigOptions{
		Precision: decimal.IntPtr(40),
		Rounding:  decimal.IntPtr(4),
		ToExpNeg:  decimal.IntPtr(-9000000000000000),
		ToExpPos:  decimal.IntPtr(9000000000000000),
		MaxE:      decimal.IntPtr(9000000000000000),
		MinE:      decimal.IntPtr(-9000000000000000),
		Crypto:    decimal.BoolPtr(false),
		Modulo:    decimal.IntPtr(1),
	})

	// --- assertEqualProps tests from number inputs ---
	t.Run("number_inputs", func(t *testing.T) {
		type propsCase struct {
			d []int32; e int; s int8; v interface{}
		}
		cases := []propsCase{
			{[]int32{0}, 0, 1, float64(0)},
			{[]int32{0}, 0, -1, math.Copysign(0, -1)},
			{[]int32{1}, 0, -1, float64(-1)},
			{[]int32{10}, 1, -1, float64(-10)},
			{[]int32{1}, 0, 1, float64(1)},
			{[]int32{10}, 1, 1, float64(10)},
			{[]int32{100}, 2, 1, float64(100)},
			{[]int32{1000}, 3, 1, float64(1000)},
			{[]int32{10000}, 4, 1, float64(10000)},
			{[]int32{100000}, 5, 1, float64(100000)},
			{[]int32{1000000}, 6, 1, float64(1000000)},
			{[]int32{1}, 7, 1, float64(10000000)},
			{[]int32{10}, 8, 1, float64(100000000)},
			{[]int32{100}, 9, 1, float64(1000000000)},
			{[]int32{1000}, 10, 1, float64(10000000000)},
			{[]int32{10000}, 11, 1, float64(100000000000)},
			{[]int32{100000}, 12, 1, float64(1000000000000)},
			{[]int32{1000000}, 13, 1, float64(10000000000000)},
			{[]int32{1}, 14, -1, float64(-100000000000000)},
			{[]int32{10}, 15, -1, float64(-1000000000000000)},
			// Small fractions
			{[]int32{1000000}, -1, 1, 1e-1},
			{[]int32{100000}, -2, -1, -1e-2},
			{[]int32{10000}, -3, 1, 1e-3},
			{[]int32{1000}, -4, -1, -1e-4},
			{[]int32{100}, -5, 1, 1e-5},
			{[]int32{10}, -6, -1, -1e-6},
			{[]int32{1}, -7, 1, 1e-7},
			{[]int32{1000000}, -8, 1, 1e-8},
			{[]int32{100000}, -9, -1, -1e-9},
			{[]int32{10000}, -10, 1, 1e-10},
			{[]int32{1000}, -11, -1, -1e-11},
			{[]int32{100}, -12, 1, 1e-12},
			{[]int32{10}, -13, -1, -1e-13},
			{[]int32{1}, -14, 1, 1e-14},
			{[]int32{1000000}, -15, 1, 1e-15},
			{[]int32{100000}, -16, -1, -1e-16},
			{[]int32{10000}, -17, 1, 1e-17},
			{[]int32{1000}, -18, -1, -1e-18},
			{[]int32{100}, -19, 1, 1e-19},
			{[]int32{10}, -20, -1, -1e-20},
			{[]int32{1}, -21, 1, 1e-21},
		}
		for i, c := range cases {
			x := newDec(t, c.v)
			assertEqualProps(t, c.d, c.e, c.s, x, "num_"+strconv.Itoa(i))
		}
	})

	// --- assertEqualProps tests from string inputs ---
	t.Run("string_inputs", func(t *testing.T) {
		type propsCase struct {
			d []int32; e int; s int8; v string
		}
		cases := []propsCase{
			{[]int32{9}, 0, 1, "9"},
			{[]int32{99}, 1, -1, "-99"},
			{[]int32{999}, 2, 1, "999"},
			{[]int32{9999}, 3, -1, "-9999"},
			{[]int32{99999}, 4, 1, "99999"},
			{[]int32{999999}, 5, -1, "-999999"},
			{[]int32{9999999}, 6, 1, "9999999"},
			{[]int32{9, 9999999}, 7, -1, "-99999999"},
			{[]int32{99, 9999999}, 8, 1, "999999999"},
			{[]int32{999, 9999999}, 9, -1, "-9999999999"},
			{[]int32{9999, 9999999}, 10, 1, "99999999999"},
			{[]int32{99999, 9999999}, 11, -1, "-999999999999"},
			{[]int32{999999, 9999999}, 12, 1, "9999999999999"},
			{[]int32{9999999, 9999999}, 13, -1, "-99999999999999"},
			{[]int32{9, 9999999, 9999999}, 14, 1, "999999999999999"},
			{[]int32{99, 9999999, 9999999}, 15, -1, "-9999999999999999"},
			{[]int32{999, 9999999, 9999999}, 16, 1, "99999999999999999"},
			{[]int32{9999, 9999999, 9999999}, 17, -1, "-999999999999999999"},
			{[]int32{99999, 9999999, 9999999}, 18, 1, "9999999999999999999"},
			{[]int32{999999, 9999999, 9999999}, 19, -1, "-99999999999999999999"},
			{[]int32{9999999, 9999999, 9999999}, 20, 1, "999999999999999999999"},
		}
		for i, c := range cases {
			x, _ := decimal.New(c.v)
			assertEqualProps(t, c.d, c.e, c.s, x, "str_"+strconv.Itoa(i))
		}
	})

	// --- Base conversion: valueOf equality tests ---
	t.Run("base_conversion_static", func(t *testing.T) {
		type valCase struct {
			expected, input string
		}
		cases := []valCase{
			// Binary
			{"0", "0b0"}, {"0", "0B0"}, {"-5", "-0b101"}, {"5", "+0b101"},
			{"1.5", "0b1.1"}, {"-1.5", "-0b1.1"},
			{"18181", "0b100011100000101.00"},
			{"-12.5", "-0b1100.10"},
			{"343872.5", "0b1010011111101000000.10"},
			{"-328.28125", "-0b101001000.010010"},
			{"-341919.144535064697265625", "-0b1010011011110011111.0010010100000000010"},
			{"97.10482025146484375", "0b1100001.000110101101010110000"},
			{"-120914.40625", "-0b11101100001010010.01101"},
			{"8080777260861123367657", "0b1101101100000111101001111111010001111010111011001010100101001001011101001"},
			// Octal
			{"8", "0o10"}, {"-8.5", "-0O010.4"}, {"8.5", "+0O010.4"},
			{"-262144.000000059604644775390625", "-0o1000000.00000001"},
			{"572315667420.390625", "0o10250053005734.31"},
			// Hex
			{"1", "0x00001"}, {"255", "0xff"},
			{"-15.5", "-0Xf.8"}, {"15.5", "+0Xf.8"},
			{"-16777216.00000000023283064365386962890625", "-0x1000000.00000001"},
			{"325927753012307620476767402981591827744994693483231017778102969592507", "0xc16de7aa5bf90c3755ef4dea45e982b351b6e00cd25a82dcfe0646abb"},
		}
		for _, c := range cases {
			d, err := decimal.New(c.input)
			if err != nil {
				t.Errorf("New(%q): %v", c.input, err)
				continue
			}
			if got := d.ValueOf(); got != c.expected {
				t.Errorf("New(%q).valueOf() = %q, want %q", c.input, got, c.expected)
			}
		}
	})

	// --- Random base conversion tests (deterministic PRNG) ---
	t.Run("base_conversion_random", func(t *testing.T) {
		rng := rand.New(rand.NewSource(42))
		for i := 0; i < 127; i++ {
			k := rng.Int63n(0x20000000000000) / int64(math.Pow(10, float64(rng.Intn(16))))
			expected := strconv.FormatInt(k, 10)

			// Binary
			d, err := decimal.New("0b" + strconv.FormatInt(k, 2))
			if err != nil {
				t.Errorf("binary parse error: %v", err)
			} else if d.ValueOf() != expected {
				t.Errorf("0b%s: got %s, want %s", strconv.FormatInt(k, 2), d.ValueOf(), expected)
			}
			// Octal
			k2 := rng.Int63n(0x20000000000000) / int64(math.Pow(10, float64(rng.Intn(16))))
			expected2 := strconv.FormatInt(k2, 10)
			d, err = decimal.New("0o" + strconv.FormatInt(k2, 8))
			if err != nil {
				t.Errorf("octal parse error: %v", err)
			} else if d.ValueOf() != expected2 {
				t.Errorf("0o%s: got %s, want %s", strconv.FormatInt(k2, 8), d.ValueOf(), expected2)
			}
			// Hex
			k3 := rng.Int63n(0x20000000000000) / int64(math.Pow(10, float64(rng.Intn(16))))
			expected3 := strconv.FormatInt(k3, 10)
			d, err = decimal.New("0x" + strconv.FormatInt(k3, 16))
			if err != nil {
				t.Errorf("hex parse error: %v", err)
			} else if d.ValueOf() != expected3 {
				t.Errorf("0x%s: got %s, want %s", strconv.FormatInt(k3, 16), d.ValueOf(), expected3)
			}
		}
	})

	// --- NaN, Infinity, special value parsing ---
	t.Run("special_values", func(t *testing.T) {
		type valCase struct{ expected string; v interface{} }
		cases := []valCase{
			{"NaN", math.NaN()},
			{"NaN", -math.NaN()},
			{"NaN", "NaN"}, {"NaN", "-NaN"}, {"NaN", "+NaN"},
			{"Infinity", math.Inf(1)},
			{"-Infinity", math.Inf(-1)},
			{"Infinity", "Infinity"}, {"-Infinity", "-Infinity"}, {"Infinity", "+Infinity"},
		}
		for _, c := range cases {
			d := newDec(t, c.v)
			if got := d.ValueOf(); got != c.expected {
				t.Errorf("valueOf(%v) = %q, want %q", c.v, got, c.expected)
			}
		}
	})

	// --- Whitespace/invalid string rejection (assertException) ---
	t.Run("invalid_strings", func(t *testing.T) {
		invalids := []string{
			" NaN", "NaN ", " NaN ", " -NaN", " +NaN", "-NaN ", "+NaN ", ".NaN", "NaN.",
			" Infinity", "Infinity ", " Infinity ", " -Infinity", " +Infinity", ".Infinity", "Infinity.",
			" 0", "0 ", " 0 ", "0-", " -0", "-0 ", "+0 ", " +0", " .0", "0. ",
			"+-0", "-+0", "--0", "++0", ".-0", ".+0", "0 .", ". 0", "..0",
			"+.-0", "-.+0", "+. 0", ".0.",
			" 1", "1 ", " 1 ", "1-", " -1", "-1 ", " +1", "+1 ", ".1.",
			"+-1", "-+1", "--1", "++1", ".-1", ".+1", "1 .", ". 1", "..1",
			"+.-1", "-.+1", "+. 1", "-. 1", "1..", "+1..", "-1..",
			"-.1.", "+.1.", ".-10.", ".+10.", ". 10.",
			"", " ", "nan", "23e", "e4", "ff", "0xg", "0Xfi",
			"++45", "--45", "9.99--", "9.99++", "0 0",
		}
		for _, s := range invalids {
			assertException(t, s, s)
		}
	})

	// --- Zero and one parsing ---
	t.Run("zero_one_parsing", func(t *testing.T) {
		type valCase struct{ expected, input string }
		cases := []valCase{
			{"0", "0"}, {"-0", "-0"}, {"0", "0."}, {"-0", "-0."},
			{"0", "0.0"}, {"-0", "-0.0"}, {"0", "0.00000000"},
			{"-0", "-0.0000000000000000000000"},
			{"1", "1"}, {"-1", "-1"}, {"0.1", ".1"}, {"-0.1", "-.1"},
			{"0.1", "+.1"}, {"1", "1."}, {"1", "1.0"}, {"-1", "-1."},
			{"1", "+1."}, {"-1", "-1.0000"}, {"1", "1.0000"},
			{"1", "1.00000000"}, {"-1", "-1.000000000000000000000000"},
			{"1", "+1.000000000000000000000000"},
			{"123.456789", "123.456789"}, {"-123.456789", "-123.456789"},
			{"123.456789", "+123.456789"},
		}
		for _, c := range cases {
			d, err := decimal.New(c.input)
			if err != nil {
				t.Errorf("New(%q): %v", c.input, err)
				continue
			}
			if got := d.ValueOf(); got != c.expected {
				t.Errorf("New(%q).valueOf() = %q, want %q", c.input, got, c.expected)
			}
		}

		// Number inputs for 0 and 123.456789
		d := newDec(t, float64(0))
		if d.ValueOf() != "0" {
			t.Errorf("valueOf(0) = %q", d.ValueOf())
		}
		d = newDec(t, math.Copysign(0, -1))
		if d.ValueOf() != "-0" {
			t.Errorf("valueOf(-0) = %q", d.ValueOf())
		}
		d = newDec(t, 123.456789)
		if d.ValueOf() != "123.456789" {
			t.Errorf("valueOf(123.456789) = %q", d.ValueOf())
		}
		d = newDec(t, -123.456789)
		if d.ValueOf() != "-123.456789" {
			t.Errorf("valueOf(-123.456789) = %q", d.ValueOf())
		}
	})
}
