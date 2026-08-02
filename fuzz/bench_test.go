package fuzz

import (
	"testing"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// BenchmarkNew measures constructor performance.
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decimal.New("123456.789012345")
	}
}

// BenchmarkPlus measures addition performance.
func BenchmarkPlus(b *testing.B) {
	a, _ := decimal.New("123456.789012345")
	c, _ := decimal.New("987654.321098765")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Plus(c)
	}
}

// BenchmarkTimes measures multiplication performance.
func BenchmarkTimes(b *testing.B) {
	a, _ := decimal.New("123456.789012345")
	c, _ := decimal.New("987654.321098765")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Times(c)
	}
}

// BenchmarkDiv measures division performance.
func BenchmarkDiv(b *testing.B) {
	a, _ := decimal.New("123456.789012345")
	c, _ := decimal.New("987654.321098765")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Div(c)
	}
}

// BenchmarkString measures toString performance.
func BenchmarkString(b *testing.B) {
	a, _ := decimal.New("123456.789012345")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.String()
	}
}
