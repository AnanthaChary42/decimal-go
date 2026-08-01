package decimal

import (
	"math"
	"strings"
)

// ifloorDiv returns floor(a / b) for integer division.
// In JS: Math.floor(x.e / LOG_BASE).
func ifloorDiv(a, b int) int {
	if a >= 0 || a%b == 0 {
		return a / b
	}
	return a/b - 1
}

// tinyPow returns b^e for small bases and positive exponents.
// Equivalent to tinyPow in decimal.js.
func tinyPow(b, e int) float64 {
	n := float64(b)
	for e--; e > 0; e-- {
		n *= float64(b)
	}
	return n
}

// mathpow is math.Pow for int arguments.
func mathpow(base, exp int) float64 {
	return math.Pow(float64(base), float64(exp))
}

// digitsToString converts a digit array (base 1e7) to a decimal string.
// Equivalent to digitsToString in decimal.js.
func digitsToString(d []int32) string {
	return digitsToStringExact(d)
}

// digitsToStringExact is the core digit-to-string implementation.
//
// JS Algorithm:
//   - d[0] is written without padding (first word).
//   - d[1..lastWord-1] are padded to LOG_BASE (7) digits.
//   - d[lastWord]: padding zeros are added to str, but the word value itself
//     has trailing zeros stripped and is appended separately.
//   - For single-word arrays, just strip trailing zeros from d[0].
func digitsToStringExact(d []int32) string {
	if len(d) == 0 {
		return "0"
	}

	indexOfLastWord := len(d) - 1
	var sb strings.Builder
	w := d[0]

	if indexOfLastWord > 0 {
		// First word: no padding.
		sb.WriteString(i32toa(d[0]))

		// Middle words: padded to LOG_BASE digits.
		for i := 1; i < indexOfLastWord; i++ {
			ws := i32toa(d[i])
			k := LOG_BASE - len(ws)
			if k > 0 {
				sb.WriteString(getZeroString(k))
			}
			sb.WriteString(ws)
		}

		// Last word: add leading-zero padding to str, but NOT the word digits.
		// The word value (w) will be appended below after stripping trailing zeros.
		w = d[indexOfLastWord]
		ws := i32toa(w)
		k := LOG_BASE - len(ws)
		if k > 0 {
			sb.WriteString(getZeroString(k))
		}
	} else if w == 0 {
		return "0"
	}

	// Strip trailing zeros of last word.
	for w%10 == 0 && w != 0 {
		w /= 10
	}

	sb.WriteString(i32toa(w))
	return sb.String()
}

// i32toa converts an int32 to a string.
func i32toa(v int32) string {
	return strings.TrimLeft(itoa(int(v)), " ")
}

// itoa converts an int to its decimal string representation.
func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := false
	if v < 0 {
		neg = true
		v = -v
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// getZeroString returns a string of k zeros.
func getZeroString(k int) string {
	return strings.Repeat("0", k)
}

// getBase10Exponent calculates the base 10 exponent from base 1e7 exponent.
// Equivalent to getBase10Exponent in decimal.js.
func getBase10Exponent(digits []int32, e int) int {
	w := digits[0]
	// Add the number of digits of the first word.
	e *= LOG_BASE
	for w >= 10 {
		w /= 10
		e++
	}
	return e
}

// getPrecision returns the number of significant digits in a digit array.
// Equivalent to getPrecision in decimal.js.
func getPrecision(digits []int32) int {
	w := len(digits) - 1
	length := w*LOG_BASE + 1

	w2 := digits[w]

	// If non-zero...
	if w2 != 0 {
		// Subtract the number of trailing zeros of the last word.
		for w2%10 == 0 {
			w2 /= 10
			length--
		}

		// Add the number of digits of the first word.
		for w2 = digits[0]; w2 >= 10; w2 /= 10 {
			length++
		}
	}

	return length
}

// checkInt32 validates that i is an integer in [min, max].
func checkInt32(i, min, max int) error {
	if i < min || i > max {
		return newInvalidArgError(i)
	}
	return nil
}

// checkRoundingDigits checks rounding digits to determine if more precision is needed.
// Returns true if the result may need recalculation with higher precision.
func checkRoundingDigits(d []int32, i int, rm RoundingMode, repeating *int) bool {
	var di, k, r int
	var rd int32

	// Get the length of the first word of the array d.
	for k2 := d[0]; k2 >= 10; k2 /= 10 {
		i--
	}

	// Is the rounding digit in the first word of d?
	i--
	if i < 0 {
		i += LOG_BASE
		di = 0
	} else {
		di = iceil(i+1, LOG_BASE)
		i %= LOG_BASE
	}

	// i is the index (0-6) of the rounding digit.
	k = int(mathpow(10, LOG_BASE-i))
	if di < len(d) {
		rd = d[di] % int32(k)
	}

	if repeating == nil {
		if i < 3 {
			if i == 0 {
				rd = rd / 100
			} else if i == 1 {
				rd = rd / 10
			}
			r = boolToInt(int(rm) < 4 && rd == 99999 || int(rm) > 3 && rd == 49999 || rd == 50000 || rd == 0)
		} else {
			halfK := int32(k / 2)
			nextDigitCheck := int32(0)
			if di+1 < len(d) {
				nextDigitCheck = d[di+1] / int32(k) / 100
			}
			pow10 := int32(mathpow(10, i-2))
			r = boolToInt(
				(int(rm) < 4 && rd+1 == int32(k) || int(rm) > 3 && rd+1 == halfK) &&
					nextDigitCheck == pow10-1 ||
					(rd == halfK || rd == 0) && nextDigitCheck == 0)
		}
	} else {
		if i < 4 {
			if i == 0 {
				rd = rd / 1000
			} else if i == 1 {
				rd = rd / 100
			} else if i == 2 {
				rd = rd / 10
			}
			r = boolToInt((*repeating != 0 || int(rm) < 4) && rd == 9999 || *repeating == 0 && int(rm) > 3 && rd == 4999)
		} else {
			halfK := int32(k / 2)
			nextDigitCheck := int32(0)
			if di+1 < len(d) {
				nextDigitCheck = d[di+1] / int32(k) / 1000
			}
			pow10 := int32(mathpow(10, i-3))
			r = boolToInt(
				((*repeating != 0 || int(rm) < 4) && rd+1 == int32(k) ||
					(*repeating == 0 && int(rm) > 3) && rd+1 == halfK) &&
					nextDigitCheck == pow10-1)
		}
	}

	return r != 0
}

// iceil returns ceil(a / b).
func iceil(a, b int) int {
	if a%b == 0 {
		return a / b
	}
	return a/b + 1
}

// boolToInt converts bool to int (1 or 0).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// convertBase converts string of baseIn to an array of numbers of baseOut.
// E.g. convertBase("255", 10, 16) returns [15, 15].
func convertBase(str string, baseIn, baseOut int) []int32 {
	arr := []int32{0}
	strL := len(str)

	for i := 0; i < strL; i++ {
		// Multiply arr by baseIn.
		for j := len(arr) - 1; j >= 0; j-- {
			arr[j] *= int32(baseIn)
		}

		// Add the digit.
		ch := str[i]
		idx := strings.IndexByte(NUMERALS, ch)
		if idx < 0 {
			idx = 0
		}
		arr[0] += int32(idx)

		// Carry.
		for j := 0; j < len(arr); j++ {
			if arr[j] > int32(baseOut)-1 {
				if j+1 >= len(arr) {
					arr = append(arr, 0)
				}
				arr[j+1] += arr[j] / int32(baseOut)
				arr[j] %= int32(baseOut)
			}
		}
	}

	// Reverse.
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

// truncateArr truncates array to len, returning true if it was actually truncated.
func truncateArr(arr *[]int32, length int) bool {
	if len(*arr) > length {
		*arr = (*arr)[:length]
		return true
	}
	return false
}

// Global flag to track whether internal functions should bypass limits.
var external = true

// finalise rounds x to sd significant digits using rounding mode rm.
// Checks for overflow/underflow.
// Equivalent to finalise in decimal.js.
func finalise(x *Decimal, sd int, rm RoundingMode, isTruncated ...bool) *Decimal {
	truncated := false
	if len(isTruncated) > 0 {
		truncated = isTruncated[0]
	}

	ctx := x.getContext()

	// Don't round if sd is negative (used as sentinel for "no rounding").
	if sd < 0 {
		// Just check overflow/underflow.
		if external {
			if x.e > ctx.MaxE {
				x.d = nil
				x.e = 0
			} else if x.e < ctx.MinE {
				x.e = 0
				x.d = []int32{0}
			}
		}
		return x
	}

	xd := x.d

	// Infinity/NaN.
	if xd == nil {
		return x
	}

	// Get the length of the first word of the digits array xd.
	var digits int
	for k := xd[0]; k >= 10; k /= 10 {
		digits++
	}
	digits++

	i := sd - digits

	var rd int32
	var roundUp bool
	var xdi int
	var w int32
	var j int

	var isTrunc bool

	// Is the rounding digit in the first word of xd?
	if i < 0 {
		i += LOG_BASE
		j = sd
		xdi = 0
		w = xd[0]

		// Get the rounding digit at index j of w.
		rd = w / int32(mathpow(10, digits-j-1)) % 10
	} else {
		xdi = iceil(i+1, LOG_BASE)
		k := len(xd)
		if xdi >= k {
			if truncated {
				// Needed by naturalExponential, naturalLogarithm and squareRoot.
				for k <= xdi {
					xd = append(xd, 0)
					k++
				}
				x.d = xd
				w = 0
				rd = 0
				digits = 1
				i %= LOG_BASE
				j = i - LOG_BASE + 1
			} else {
				// No rounding needed.
				goto checkOverflow
			}
		} else {
			w = xd[xdi]

			// Get the number of digits of w.
			digits = 1
			for k2 := w; k2 >= 10; k2 /= 10 {
				digits++
			}

			// Get the index of rd within w.
			i %= LOG_BASE

			// Get the index of rd within w, adjusted for leading zeros.
			j = i - LOG_BASE + digits

			// Get the rounding digit at index j of w.
			if j < 0 {
				rd = 0
			} else {
				rd = w / int32(mathpow(10, digits-j-1)) % 10
			}
		}
	}

	// Are there any non-zero digits after the rounding digit?
	isTrunc = truncated || sd < 0 ||
		(xdi+1 < len(xd) && xd[xdi+1] != 0) // simplified check
	if !isTrunc && j >= 0 {
		// Check remaining digits in current word.
		rem := w % int32(mathpow(10, maxInt(0, digits-j-1)))
		if rem != 0 {
			isTrunc = true
		}
	} else if !isTrunc && j < 0 {
		if w != 0 {
			isTrunc = true
		}
	}

	roundUp = int(rm) < 4 &&
		(rd != 0 || isTrunc) && (rm == 0 || rm == RoundingMode(boolToInt(x.s < 0)*3+boolToInt(x.s >= 0)*2))

	if !roundUp {
		// Determine the rounding mode thresholds.
		rmNeg := RoundingMode(boolToInt(x.s < 0)*8 + boolToInt(x.s >= 0)*7)

		// Check whether the digit to the left of the rounding digit is odd.
		var leftDigit int32
		if i > 0 {
			if j > 0 {
				leftDigit = w / int32(mathpow(10, digits-j))
			}
		} else if xdi > 0 {
			leftDigit = xd[xdi-1]
		}
		isOddLeft := leftDigit%2 != 0

		roundUp = rd > 5 || (rd == 5 && (rm == 4 || isTrunc ||
			(rm == 6 && isOddLeft) || rm == rmNeg))
	}

	if sd < 1 || len(xd) == 0 || xd[0] == 0 {
		x.d = x.d[:0]
		if roundUp {
			// Convert sd to decimal places.
			sdAdj := sd - x.e - 1
			x.d = append(x.d[:0], int32(mathpow(10, (LOG_BASE-sdAdj%LOG_BASE)%LOG_BASE)))
			if sdAdj > 0 {
				x.e = -sdAdj
			} else {
				x.e = 0
			}
		} else {
			// Zero.
			x.d = append(x.d[:0], 0)
			x.e = 0
		}
		return x
	}

	// Remove excess digits.
	if i == 0 {
		xd = xd[:xdi]
		x.d = xd
		xdi--
	} else {
		xd = xd[:xdi+1]
		x.d = xd
		k := int32(mathpow(10, LOG_BASE-i))

		// E.g. 56700 becomes 56000 if 7 is the rounding digit.
		if j > 0 {
			xd[xdi] = (w / int32(mathpow(10, digits-j)) % int32(mathpow(10, j))) * k
		} else {
			xd[xdi] = 0
		}
	}

	if roundUp {
		for {
			// Is the digit to be rounded up in the first word of xd?
			if xdi == 0 {
				// i will be the length of xd[0] before k is added.
				iLen := 1
				for j2 := xd[0]; j2 >= 10; j2 /= 10 {
					iLen++
				}
				xd[0] += int32(mathpow(10, LOG_BASE-i))
				kLen := 1
				for j2 := xd[0]; j2 >= 10; j2 /= 10 {
					kLen++
				}

				if iLen != kLen {
					x.e++
					if xd[0] == BASE {
						xd[0] = 1
					}
				}
				break
			} else {
				xd[xdi] += 1
				if xd[xdi] != BASE {
					break
				}
				xd[xdi] = 0
				xdi--
				if xdi < 0 {
					// Need to prepend.
					x.d = append([]int32{1}, x.d...)
					x.e++
					break
				}
			}
		}
	}

	// Remove trailing zeros.
	for i2 := len(x.d) - 1; i2 >= 0 && x.d[i2] == 0; i2-- {
		x.d = x.d[:i2]
	}
	if len(x.d) == 0 {
		x.d = []int32{0}
		x.e = 0
	}

checkOverflow:
	if external {
		if x.e > ctx.MaxE {
			// Infinity.
			x.d = nil
			x.e = 0
		} else if x.e < ctx.MinE {
			// Zero.
			x.e = 0
			x.d = []int32{0}
		}
	}

	return x
}

// maxInt returns the larger of a and b.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// minInt returns the smaller of a and b.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// absInt returns the absolute value of an int.
func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// getLn10 returns LN10 to sd significant digits.
func getLn10(ctx *Context, sd int, pr ...int) (*Decimal, error) {
	if sd > LN10_PRECISION {
		external = true
		if len(pr) > 0 {
			ctx.Precision = pr[0]
		}
		return nil, ErrPrecisionLimit
	}
	x, _ := ctx.New(LN10)
	return finalise(x, sd, 1), nil
}

// getPi returns PI to sd significant digits.
func getPi(ctx *Context, sd int, rm RoundingMode) (*Decimal, error) {
	if sd > PI_PRECISION {
		return nil, ErrPrecisionLimit
	}
	x, _ := ctx.New(PI)
	return finalise(x, sd, rm, true), nil
}
