package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

type report struct {
	Runtime          string  `json:"runtime"`
	P99LatencyNS     int64   `json:"p99_latency_ns"`
	RSSKB            uint64  `json:"rss_kb"`
	ThroughputOpsSec float64 `json:"throughput_ops_sec"`
	Operations       int     `json:"operations"`
}

func main() {
	startupProbe := flag.Bool("startup-probe", false, "write a readiness line and exit")
	operations := flag.Int("ops", 100000, "number of mixed-operation iterations")
	flag.Parse()

	if *startupProbe {
		fmt.Println("ready")
		return
	}
	if *operations < 1 {
		fmt.Fprintln(os.Stderr, "--ops must be a positive integer")
		os.Exit(2)
	}

	ctx := &decimal.Context{
		Precision: 40,
		Rounding:  decimal.RoundHalfUp,
		ToExpNeg:  -9000000000000000,
		ToExpPos:  9000000000000000,
		MinE:      -9000000000000000,
		MaxE:      9000000000000000,
		Modulo:    decimal.RoundDown,
	}
	a, _ := ctx.New("123456.78901234567890123456789")
	b, _ := ctx.New("987654.32109876543210987654321")
	c, _ := ctx.New("0.000000123456789")
	d, _ := ctx.New("1.0000000000000001")
	divisor, _ := ctx.New("7.123456789")

	samples := make([]int64, *operations)
	peakRSS := currentRSSKB()
	started := monotonicNowNS()
	var sink string
	for i := 0; i < *operations; i++ {
		opStarted := monotonicNowNS()
		value := a.Add(b).Sub(c).Mul(d).Div(divisor)
		sink = value.ToFixed(12) + value.ToDP(10).ValueOf() + value.Round().ValueOf()
		samples[i] = monotonicNowNS() - opStarted

		if i&1023 == 0 {
			peakRSS = max(peakRSS, currentRSSKB())
		}
	}
	elapsedNS := monotonicNowNS() - started
	peakRSS = max(peakRSS, currentRSSKB())
	if sink == "" {
		panic("benchmark result was unexpectedly empty")
	}

	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	p99Index := (len(samples)*99 + 99) / 100
	if p99Index > 0 {
		p99Index--
	}

	result := report{
		Runtime:          "Go " + runtime.Version(),
		P99LatencyNS:     samples[p99Index],
		RSSKB:            peakRSS,
		ThroughputOpsSec: float64(*operations) * 1e9 / float64(elapsedNS),
		Operations:       *operations,
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
