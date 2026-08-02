package port_test

import (
	"math/rand"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestOriginal_PowSqrt ports JS powSqrt.js using a pre-seeded deterministic PRNG loop.
// The JS test uses Math.random() seeded nondeterministically; here we fix the seed to 42
// for deterministic, reproducible runs while preserving the same logic.
func TestOriginal_PowSqrt(t *testing.T) {
	saveDefaultCtx(t)

	decimal.Config(decimal.ConfigOptions{
		ToExpNeg: decimal.IntPtr(-7),
		ToExpPos: decimal.IntPtr(21),
		MinE:     decimal.IntPtr(-9000000000000000),
		MaxE:     decimal.IntPtr(9000000000000000),
	})

	ctx := decimal.GetDefaultContext()
	rng := rand.New(rand.NewSource(42))
	total := 0

	for total < 10000 {
		// Get a random value in [0,1) with a random number of significant digits in [1, 40].
		sd := rng.Intn(40) + 1
		e := ctx.Random(sd).ToExponential(-1)

		// Change exponent to a non-zero value of random length in (-9e15, 9e15).
		eIdx := -1
		for i := 0; i < len(e); i++ {
			if e[i] == 'e' {
				eIdx = i
				break
			}
		}

		prefix := e[:eIdx+1]
		sign := ""
		if rng.Float64() < 0.5 {
			sign = "-"
		}
		n := rng.Int63n(9000000000000000)
		nStr := func() string {
			s := ""
			if n == 0 {
				return "0"
			}
			tmp := n
			for tmp > 0 {
				s = string(rune('0'+tmp%10)) + s
				tmp /= 10
			}
			return s
		}()

		cutLen := rng.Intn(len(nStr))
		expStr := nStr[cutLen:]
		if expStr == "" || expStr == "0" {
			expStr = "1"
		}

		rStr := prefix + sign + expStr
		r, err := decimal.New(rStr)
		if err != nil {
			total++
			continue
		}

		// Random rounding mode.
		ctx.Rounding = decimal.RoundingMode(rng.Intn(9))

		// Random precision in [1, 40].
		ctx.Precision = rng.Intn(40) + 1

		halfD, _ := decimal.New("0.5")
		p := r.Pow(halfD) // r.pow(0.5)
		s := r.Sqrt()     // r.sqrt()

		if p.ValueOf() != s.ValueOf() {
			t.Errorf("pow(0.5) != sqrt(): r=%s, pow=%s, sqrt=%s (precision=%d, rounding=%d)",
				rStr, p.ValueOf(), s.ValueOf(), ctx.Precision, ctx.Rounding)
		}

		total++
	}
}
