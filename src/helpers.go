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

// mathpow is math.Pow for int arguments.
func mathpow(base, exp int) float64 {
	return math.Pow(float64(base), float64(exp))
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

	var directedRM RoundingMode
	var halfDirectedRM RoundingMode
	var isTrunc bool

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
	var kVal int32

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
	// JS: isTruncated = isTruncated || sd < 0 ||
	//   xd[xdi + 1] !== void 0 || (j < 0 ? w : w % mathpow(10, digits - j - 1));
	// Note: xd[xdi + 1] !== void 0 means the element EXISTS (index in bounds), not non-zero.
	isTrunc = truncated || sd < 0 || xdi+1 < len(xd)
	if !isTrunc {
		if j < 0 {
			isTrunc = w != 0
		} else {
			rem := w % int32(mathpow(10, maxInt(0, digits-j-1)))
			isTrunc = rem != 0
		}
	}

	// JS: roundUp = rm < 4
	//   ? (rd || isTruncated) && (rm == 0 || rm == (x.s < 0 ? 3 : 2))
	//   : rd > 5 || rd == 5 && (rm == 4 || isTruncated || rm == 6 &&
	//     ((i > 0 ? j > 0 ? w / mathpow(10, digits - j) : 0 : xd[xdi - 1]) % 10) & 1 ||
	//       rm == (x.s < 0 ? 8 : 7));

	if x.s < 0 {
		directedRM = 3 // RoundFloor
	} else {
		directedRM = 2 // RoundCeil
	}

	if int(rm) < 4 {
		// Directed rounding modes: RoundUp(0), RoundDown(1), RoundCeil(2), RoundFloor(3)
		roundUp = (rd != 0 || isTrunc) && (rm == 0 || rm == directedRM)
	} else {
		// Half-rounding modes: RoundHalfUp(4), RoundHalfDown(5), RoundHalfEven(6),
		// RoundHalfCeil(7), RoundHalfFloor(8), Euclid(9)
		if x.s < 0 {
			halfDirectedRM = 8 // RoundHalfFloor
		} else {
			halfDirectedRM = 7 // RoundHalfCeil
		}

		// Get the digit to the left of the rounding digit.
		var leftDigit int32
		if i > 0 {
			if j > 0 {
				leftDigit = w / int32(mathpow(10, digits-j))
			} else {
				leftDigit = 0
			}
		} else {
			if xdi > 0 {
				leftDigit = xd[xdi-1]
			} else {
				leftDigit = 0
			}
		}

		roundUp = rd > 5 || (rd == 5 && (rm == 4 || isTrunc || (rm == 6 &&
			(leftDigit%10)&1 != 0) ||
			rm == halfDirectedRM))
	}

	if sd < 1 || len(xd) == 0 || xd[0] == 0 {
		x.d = x.d[:0]
		if roundUp {
			// Convert sd to decimal places.
			sdAdj := sd - x.e - 1
			x.d = append(x.d[:0], int32(mathpow(10, (LOG_BASE-sdAdj%LOG_BASE)%LOG_BASE)))
			// JS: x.e = -sd || 0  (i.e. if -sd is 0, use 0; otherwise use -sd)
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
		kVal = 1
		xdi--
	} else {
		xd = xd[:xdi+1]
		x.d = xd
		kVal = int32(mathpow(10, LOG_BASE-i))

		// E.g. 56700 becomes 56000 if 7 is the rounding digit.
		if j > 0 {
			xd[xdi] = (w / int32(mathpow(10, digits-j)) % int32(mathpow(10, j))) * kVal
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
				xd[0] += kVal
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
				xd[xdi] += kVal
				if xd[xdi] != BASE {
					break
				}
				xd[xdi] = 0
				xdi--
				kVal = 1
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
