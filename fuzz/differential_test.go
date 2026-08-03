package fuzz

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	decimal "github.com/AnanthaChary42/decimal-go/src"
)

// TestDifferentialFuzz runs the Go decimal library and the original JS
// decimal.js side-by-side, comparing outputs on a shared public API.
//
// Usage:
//
//	go test ./fuzz/... -run TestDifferentialFuzz -timeout 120s -v
//
// Set FUZZ_DURATION to override the default 60s run time (e.g. "90s").
func TestDifferentialFuzz(t *testing.T) {
	duration := 60 * time.Second
	if d := os.Getenv("FUZZ_DURATION"); d != "" {
		parsed, err := time.ParseDuration(d)
		if err == nil {
			duration = parsed
		}
	}

	// Locate oracle.js relative to this test file.
	_, thisFile, _, _ := runtime.Caller(0)
	fuzzDir := filepath.Dir(thisFile)
	oraclePath := filepath.Join(fuzzDir, "oracle.js")

	// Verify oracle.js and node_modules exist.
	if _, err := os.Stat(oraclePath); os.IsNotExist(err) {
		t.Fatal("oracle.js not found — run: cd fuzz && npm install")
	}
	nodeModules := filepath.Join(fuzzDir, "node_modules", "decimal.js")
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		t.Fatal("decimal.js not installed — run: cd fuzz && npm install")
	}

	// Start the Node.js oracle as a child process.
	cmd := exec.Command("node", oraclePath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start Node.js oracle: %v", err)
	}
	defer func() {
		stdin.Close()
		cmd.Wait()
	}()

	scanner := bufio.NewScanner(stdout)
	// Increase scanner buffer for large JSON lines.
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	rng := rand.New(rand.NewSource(42))

	totalOps := 0
	divergences := 0
	startTime := time.Now()
	deadline := startTime.Add(duration)

	// Log file for the published fuzz report.
	logPath := filepath.Join(fuzzDir, "differential_fuzz_log.txt")
	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("failed to create log file: %v", err)
	}
	defer logFile.Close()

	logLine := func(format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		fmt.Fprintln(logFile, line)
		t.Log(line)
	}

	logLine("=== Differential Fuzz Report ===")
	logLine("Source:    github.com/MikeMcl/decimal.js (npm decimal.js)")
	logLine("Port:     github.com/AnanthaChary42/decimal-go")
	logLine("Oracle:   Node.js + original decimal.js via stdin/stdout pipe")
	logLine("Duration: %s", duration)
	logLine("Seed:     42")
	logLine("")
	logLine("--- Operations under test (shared public API) ---")
	logLine("  toString, valueOf, plus, minus, times, div, abs, neg,")
	logLine("  cmp, isNaN, isFinite, isZero, isNeg, isPos,")
	logLine("  toFixed(10), sqrt, floor, ceil, trunc, round, mod")
	logLine("")
	logLine("--- Fuzz progress ---")

	lastReport := startTime

	for time.Now().Before(deadline) {
		// Generate a random input pair.
		aStr := randomDecimalString(rng)
		bStr := randomDecimalString(rng)

		// Send to JS oracle.
		inputJSON, _ := json.Marshal(map[string]string{"a": aStr, "b": bStr})
		_, err := fmt.Fprintf(stdin, "%s\n", inputJSON)
		if err != nil {
			t.Fatalf("failed to write to oracle stdin: %v", err)
		}

		// Read JS result.
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				t.Fatalf("oracle stdout scan error: %v", err)
			}
			t.Fatal("oracle stdout closed unexpectedly")
		}

		var jsResult map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &jsResult); err != nil {
			t.Fatalf("failed to parse oracle output: %v: %s", err, scanner.Text())
		}

		if _, hasErr := jsResult["error"]; hasErr {
			// JS couldn't parse this input — skip (both sides reject it).
			totalOps++
			continue
		}

		// Compute Go results.
		goA, errA := decimal.New(aStr)
		goB, errB := decimal.New(bStr)

		if errA != nil || errB != nil {
			// Go rejected this input — check that JS also produced an error.
			if _, hasErr := jsResult["error"]; !hasErr {
				// Divergence: JS accepted but Go rejected.
				divergences++
				logLine("DIVERGENCE [parse]: Go rejected %q/%q but JS accepted", aStr, bStr)
			}
			totalOps++
			continue
		}

		// Compare each operation.
		checks := []struct {
			name  string
			goVal string
			jsKey string
		}{
			{"toString(a)", goA.String(), "a_str"},
			{"toString(b)", goB.String(), "b_str"},
			{"valueOf(a)", goA.ValueOf(), "a_valueOf"},
			{"valueOf(b)", goB.ValueOf(), "b_valueOf"},
			{"plus", goA.Plus(goB).String(), "plus"},
			{"minus", goA.Minus(goB).String(), "minus"},
			{"times", goA.Times(goB).String(), "times"},
			{"div", goA.Div(goB).String(), "div"},
			{"abs(a)", goA.Abs().String(), "abs_a"},
			{"neg(a)", goA.Neg().String(), "neg_a"},
		}

		for _, c := range checks {
			jsVal, ok := jsResult[c.jsKey]
			if !ok {
				continue
			}
			jsStr := fmt.Sprintf("%v", jsVal)
			if c.goVal != jsStr {
				divergences++
				logLine("DIVERGENCE [%s]: a=%q b=%q  go=%q  js=%q", c.name, aStr, bStr, c.goVal, jsStr)
			}
		}

		// Comparison (cmp).
		if jsVal, ok := jsResult["cmp"]; ok {
			cmpGo, cmpOk := goA.Cmp(goB)
			goStr := "NaN"
			if cmpOk {
				goStr = fmt.Sprintf("%d", cmpGo)
			}
			jsStr := fmt.Sprintf("%v", jsVal)
			if goStr != jsStr {
				divergences++
				logLine("DIVERGENCE [cmp]: a=%q b=%q  go=%q  js=%q", aStr, bStr, goStr, jsStr)
			}
		}

		// Predicates.
		predicateChecks := []struct {
			name  string
			goVal bool
			jsKey string
		}{
			{"isNaN(a)", goA.IsNaN(), "isNaN_a"},
			{"isFinite(a)", goA.IsFinite(), "isFinite_a"},
			{"isZero(a)", goA.IsZero(), "isZero_a"},
			{"isNeg(a)", goA.IsNeg(), "isNeg_a"},
			{"isPos(a)", goA.IsPos(), "isPos_a"},
		}
		for _, c := range predicateChecks {
			jsVal, ok := jsResult[c.jsKey]
			if !ok {
				continue
			}
			jsBool, _ := jsVal.(bool)
			if c.goVal != jsBool {
				divergences++
				logLine("DIVERGENCE [%s]: a=%q  go=%v  js=%v", c.name, aStr, c.goVal, jsBool)
			}
		}

		// toFixed(10) — only for finite a.
		if goA.IsFinite() {
			if jsVal, ok := jsResult["toFixed_a"]; ok {
				goFixed := goA.ToFixed(10)
				jsStr := fmt.Sprintf("%v", jsVal)
				if goFixed != jsStr {
					divergences++
					logLine("DIVERGENCE [toFixed(10)]: a=%q  go=%q  js=%q", aStr, goFixed, jsStr)
				}
			}
		}

		// sqrt — only for non-negative finite a.
		if goA.IsFinite() && !goA.IsNeg() {
			if jsVal, ok := jsResult["sqrt_a"]; ok {
				goSqrt := goA.Sqrt().String()
				jsStr := fmt.Sprintf("%v", jsVal)
				if goSqrt != jsStr {
					divergences++
					logLine("DIVERGENCE [sqrt]: a=%q  go=%q  js=%q", aStr, goSqrt, jsStr)
				}
			}
		}

		// floor, ceil, trunc, round — only for finite a.
		if goA.IsFinite() {
			roundChecks := []struct {
				name  string
				goVal string
				jsKey string
			}{
				{"floor", goA.Floor().String(), "floor_a"},
				{"ceil", goA.Ceil().String(), "ceil_a"},
				{"trunc", goA.Trunc().String(), "trunc_a"},
				{"round", goA.Round().String(), "round_a"},
			}
			for _, c := range roundChecks {
				if jsVal, ok := jsResult[c.jsKey]; ok {
					jsStr := fmt.Sprintf("%v", jsVal)
					if c.goVal != jsStr {
						divergences++
						logLine("DIVERGENCE [%s]: a=%q  go=%q  js=%q", c.name, aStr, c.goVal, jsStr)
					}
				}
			}
		}

		// mod — only if b is finite and non-zero.
		if goB.IsFinite() && !goB.IsZero() {
			if jsVal, ok := jsResult["mod"]; ok {
				goMod := goA.Mod(goB).String()
				jsStr := fmt.Sprintf("%v", jsVal)
				if goMod != jsStr {
					divergences++
					logLine("DIVERGENCE [mod]: a=%q b=%q  go=%q  js=%q", aStr, bStr, goMod, jsStr)
				}
			}
		}

		totalOps++

		// Progress report every 5 seconds.
		if time.Since(lastReport) >= 5*time.Second {
			elapsed := time.Since(startTime).Seconds()
			opsPerSec := float64(totalOps) / elapsed
			logLine("fuzz: elapsed: %.0fs, ops: %d (%.0f/sec), divergences: %d",
				elapsed, totalOps, opsPerSec, divergences)
			lastReport = time.Now()
		}
	}

	elapsed := time.Since(startTime)
	opsPerSec := float64(totalOps) / elapsed.Seconds()

	logLine("")
	logLine("--- Final Summary ---")
	logLine("Total input pairs tested: %d", totalOps)
	logLine("Total elapsed:            %s", elapsed.Round(time.Millisecond))
	logLine("Throughput:               %.0f ops/sec", opsPerSec)
	logLine("Divergences:              %d", divergences)
	logLine("")
	if divergences == 0 {
		logLine("RESULT: PASS — 0 divergences across %d differential comparisons in %s", totalOps, elapsed.Round(time.Second))
	} else {
		logLine("RESULT: FAIL — %d divergence(s) detected", divergences)
	}

	if divergences > 0 {
		t.Fatalf("Differential fuzz found %d divergence(s). See %s", divergences, logPath)
	}
}

// randomDecimalString generates a random decimal string for fuzzing.
func randomDecimalString(rng *rand.Rand) string {
	kind := rng.Intn(100)
	switch {
	case kind < 3:
		// Special values.
		specials := []string{"0", "-0", "NaN", "Infinity", "-Infinity"}
		return specials[rng.Intn(len(specials))]
	case kind < 8:
		// Small integers.
		n := rng.Intn(201) - 100 // [-100, 100]
		return fmt.Sprintf("%d", n)
	case kind < 15:
		// Exponential notation.
		coeff := rng.Float64()*20 - 10 // [-10, 10]
		exp := rng.Intn(620) - 310     // [-310, 309]
		return fmt.Sprintf("%ge%d", coeff, exp)
	case kind < 30:
		// Integer with many digits.
		digits := rng.Intn(30) + 1
		return randomDigitString(rng, digits)
	default:
		// Decimal with fractional part.
		intDigits := rng.Intn(15) + 1
		fracDigits := rng.Intn(20) + 1
		sign := ""
		if rng.Intn(2) == 0 {
			sign = "-"
		}
		return sign + randomDigitString(rng, intDigits) + "." + randomDigitString(rng, fracDigits)
	}
}

// randomDigitString generates a string of n random digits (no leading zeros
// except for a single "0").
func randomDigitString(rng *rand.Rand, n int) string {
	if n <= 0 {
		return "0"
	}
	var sb strings.Builder
	sb.Grow(n)
	// First digit: 1-9 (no leading zero for multi-digit).
	if n == 1 {
		sb.WriteByte(byte('0' + rng.Intn(10)))
	} else {
		sb.WriteByte(byte('1' + rng.Intn(9)))
		for i := 1; i < n; i++ {
			sb.WriteByte(byte('0' + rng.Intn(10)))
		}
	}
	return sb.String()
}
