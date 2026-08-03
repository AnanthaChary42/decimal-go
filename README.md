# decimal-go (Port Mortem Track H: JavaScript → Go)

JavaScript natively lacks an arbitrary-precision decimal type, forcing financial and high-precision applications to rely on floating-point arithmetic or userland libraries like `decimal.js`. In Node.js/V8 environments, these dynamic object-heavy implementations suffer from continuous GC pressure, JIT warmup variability, and thread-safety hazards due to global mutable configuration. Porting `decimal.js` to Go eliminates these bottlenecks by leveraging Go's strong static typing, zero-overhead value semantics, compile-time memory control, and native concurrency primitives (`goroutines`), providing a high-throughput, deterministic, thread-safe arbitrary-precision decimal math engine for system-level and financial applications.

## Track H Self-Justification Statement

This repository is submitted under **Track H (Open Pair: JavaScript → Go)**:

- **Source Repository**: [MikeMcl/decimal.js](https://github.com/MikeMcl/decimal.js) (v10.6.0, 4,953 lines of JavaScript, MIT License).
- **Target Language**: Go 1.22 (Standard Library).
- **Justification**: While `decimal.js` is the gold standard for JavaScript decimal arithmetic, its dynamic design relies on mutable global state (`Decimal.precision = 20`) and prototype-based object allocation. This Go port preserves the exact base-1e7 word array (`[]int32`) internal algorithm and IEEE 754-2008 rounding mode semantics while introducing immutable `Context` structs, explicit `(Decimal, error)` error handling, and complete independence from JavaScript runtime environments.

## Repository Structure

```
decimal-go-Port/
├── .git/
├── .gitignore
├── LICENSE                          # MIT, matching original decimal.js license
├── README.md                        # migration rationale + build instructions
├── DECISIONS.md                     # architectural decision records (D001–D019)
├── UNSAFE.md                        # zero-unsafe escape-hatch inventory
├── Dockerfile                       # single command → runnable artifact
├── Makefile                         # make build / make test / make bench / make fuzz
├── go.mod
├── .port-mortem.toml                # track letter, source repo URL, kickoff hash reference
├── src/                             # package decimal implementation
│   ├── advanced_transcendental.go   # atan, atan2, sinh, cosh, tanh, asinh, acosh, atanh
│   ├── arithmetic.go                # add, sub, mul, pow, sqrt, cbrt
│   ├── base_format.go               # toBinary, toOctal, toHex
│   ├── comparison.go                # cmp, eq, gt, lt, clamp
│   ├── compat.go                    # JS-compatible dynamic argument wrappers
│   ├── config.go                    # context configuration, min/max/sign/sum/random
│   ├── constants.go                 # BASE, LOG_BASE, LN10, PI constants
│   ├── decimal.go                   # core Decimal type, Context, constructors
│   ├── divide.go                    # word-array long division (base conversion)
│   ├── errors.go                    # DecimalError type
│   ├── exact.go                     # exact rational division via big.Int
│   ├── format.go                    # string, toFixed, toExponential, toPrecision
│   ├── helpers.go                   # digit manipulation, finalise, rounding logic
│   ├── parse.go                     # decimal/hex/binary/octal string parser
│   ├── rounding.go                  # floor, ceil, trunc, round, toDP, toSD
│   ├── transcendental.go            # ln, exp, log, pow via math/big.Float
│   └── trigonometric.go             # sin, cos, tan, asin, acos
├── tests/
│   ├── original/
│   │   └── test-suite.sha256        # SHA256 hash of original JS test suite, pinned at kickoff
│   └── port/
│       ├── test_original_test.go    # original test assertions ported 1:1
│       └── decimal_port_test.go     # additional Go-idiomatic tests
├── fuzz/
│   ├── differential_test.go         # differential fuzz: Go port vs JS oracle (stdin/stdout pipe)
│   ├── harness_test.go              # property-based fuzz: parsing symmetry, commutativity
│   ├── oracle.js                    # Node.js oracle that evaluates decimal.js operations
│   ├── package.json                 # npm dependency on decimal.js for the oracle
│   └── differential_fuzz_log.txt    # actual 60s+ differential run output (693k ops, 0 divergences)
└── bench/
    ├── methodology.md               # measurement methodology & confounders
    └── results.json                 # p99, RSS, startup, throughput — original vs port
```

## Quick Start

All commands can be executed directly from the repo root without changing directories:

### Build Package

```bash
make build
# Direct Go command: go build ./src/...
```

### Run Test Suite

```bash
make test
# Direct Go command: go test ./tests/port/... -v
```

### Run Performance Benchmarks

```bash
make bench
```

### Run Differential Fuzzer (60s+)

The differential fuzzer runs the Go port and the original `decimal.js` library side-by-side, feeding identical random inputs to both and comparing outputs across 21 operations (toString, valueOf, plus, minus, times, div, abs, neg, cmp, predicates, toFixed, sqrt, floor, ceil, trunc, round, mod). It uses a Node.js child process as the JS oracle.

```bash
# Requires: cd fuzz && npm install (one-time setup for the JS oracle)
go test ./fuzz/... -run TestDifferentialFuzz -timeout 120s -v
# Override duration: FUZZ_DURATION=90s go test ./fuzz/... -run TestDifferentialFuzz -timeout 120s -v
```

The latest run completed **693,744 input pairs** in 60s with **0 divergences** (see `fuzz/differential_fuzz_log.txt`).

### Run Property-Based Fuzzer (Go native)

```bash
make fuzz
# Direct Go command: go test ./fuzz/... -fuzz=FuzzDecimal -fuzztime=60s
```

### Docker Containerized Build & Test

```bash
docker build -t decimal-go-port .
```

