package port_test

import (
	"strconv"
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

func TestOriginal_Pow(t *testing.T) {
	ctx := &decimal.Context{
		Precision: 40,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
	}

	tests := []struct {
		base, exp string
		sd        int
		rm        decimal.RoundingMode
		expected  string
	}{
		{"9", "0.5", 7, decimal.RoundHalfUp, "3"},
		{"9", "0.5", 26, decimal.RoundHalfUp, "3"},
		{"0.9999999999", "6", 39, decimal.RoundHalfUp, "0.999999999400000000149999999980000000001"},
		{"2.56", "6.5", 16, decimal.RoundDown, "450.3599627370496"},
		{"1.96", "1.5", 15, decimal.RoundDown, "2.744"},
		{"2.25", "9.5", 23, decimal.RoundDown, "2216.8378200531005859375"},
		{"11.05", "2.00000000000000007", 6, decimal.RoundHalfUp, "122.103"},
		{"10.5", "3.000000000000000002", 6, decimal.RoundHalfUp, "1157.63"},
		{"1.00000000000000000003", "4.00000005", 28, decimal.RoundHalfUp, "1.000000000000000000120000002"},
		{"6.0000005", "1.00000000000000006", 7, decimal.RoundHalfUp, "6.000001"},
		{"1.0000000000000000000005", "49.0000000000000000000002", 22, decimal.RoundHalfUp, "1.000000000000000000025"},
		{"15.333333333333333333", "28.33333333333333", 49, decimal.RoundHalfUp, "3917746643938779840069598486694964.98308568625045"},
		{"7.537714", "7.9", 21, decimal.RoundHalfUp, "8515169.08260507715975"},
		{"6.951", "9.225", 10, decimal.RoundHalfUp, "58598464.57"},
		{"6.01093", "9.8911", 36, decimal.RoundHalfUp, "50651225.3819968681522216250662534915"},
		{"8.7587", "4.23", 18, decimal.RoundHalfUp, "9694.37298592397372"},
		{"5.1749", "7.7267995", 19, decimal.RoundHalfUp, "328229.2815443039852"},
		{"0.16", "-0.9999999999999", 2, decimal.RoundHalfUp, "6.2"},
		{"0.4", "-20", 27, decimal.RoundHalfUp, "90949470.1772928237915039063"},
		{"0.5", "22", 15, decimal.RoundHalfUp, "0.000000238418579101563"},
		{"32", "0.4", 1, decimal.RoundHalfUp, "4"},
		{"4", "2.5", 11, decimal.RoundHalfUp, "32"},
		{"4", "5.5", 27, decimal.RoundHalfUp, "2048"},
		{"16", "23.5", 29, decimal.RoundHalfUp, "19807040628566084398385987584"},
		{"16", "26.5", 35, decimal.RoundHalfUp, "81129638414606681695789005144064"},
		{"25", "13.5", 39, decimal.RoundHalfUp, "7450580596923828125"},
		{"32", "28.2", 43, decimal.RoundHalfUp, "2787593149816327892691964784081045188247552"},
		{"32", "3.6", 35, decimal.RoundHalfUp, "262144"},
		{"25", "21.5", 31, decimal.RoundHalfUp, "1136868377216160297393798828125"},
		{"9", "8.5", 19, decimal.RoundHalfUp, "129140163"},
		{"4", "7.5", 13, decimal.RoundHalfUp, "32768"},
		{"4", "6.5", 10, decimal.RoundHalfUp, "8192"},
		{"6.034", "0.25964", 1, decimal.RoundHalfUp, "2"},
		{"9", "4.5", 16, decimal.RoundHalfUp, "19683"},
		{"9", "1.5", 5, decimal.RoundHalfUp, "27"},
		{"9.61", "3.5", 12, decimal.RoundHalfUp, "2751.2614111"},
		{"4", "6.5", 8, decimal.RoundHalfUp, "8192"},
		{"4", "7.5", 11, decimal.RoundHalfUp, "32768"},
		{"9", "4.5", 5, decimal.RoundHalfUp, "19683"},
		{"48.9262695992662373981", "1.0", 17, decimal.RoundDown, "48.926269599266237"},
		{"1.21", "0.5", 2, decimal.RoundDown, "1.1"},
		{"12.96", "0.5", 2, decimal.RoundFloor, "3.6"},
		{"3.24", "0.5", 2, decimal.RoundDown, "1.8"},
		{"70.56", "0.5", 2, decimal.RoundFloor, "8.4"},
		{"4.41", "6.5", 32, decimal.RoundFloor, "15447.2377739119461"},
		{"11.05", "2.00000000000000007", 6, decimal.RoundHalfUp, "122.103"},
		{"10.5", "3.000000000000000002", 6, decimal.RoundHalfUp, "1157.63"},
		{"1.00000000000000000003", "4.00000005", 28, decimal.RoundHalfUp, "1.000000000000000000120000002"},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"_"+tt.base+"^"+tt.exp, func(t *testing.T) {
			customCtx := *ctx
			customCtx.Precision = tt.sd
			customCtx.Rounding = tt.rm

			b, _ := customCtx.New(tt.base)
			e, _ := customCtx.New(tt.exp)

			got := b.Pow(e).ValueOf()
			if got != tt.expected {
				t.Errorf("Decimal(%q).Pow(%q) [sd=%d, rm=%d] = %q, want %q", tt.base, tt.exp, tt.sd, tt.rm, got, tt.expected)
			}
		})
	}
}
