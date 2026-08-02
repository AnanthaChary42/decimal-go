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
├── DECISIONS.md                     # architectural decision records (D001–D011)
├── Dockerfile                       # single command → runnable artifact
├── Makefile                         # make build / make test / make bench / make fuzz
├── go.mod
├── .port-mortem.toml                # track letter, source repo URL, kickoff hash reference
├── src/                             # package decimal implementation
│   ├── arithmetic.go
│   ├── comparison.go
│   ├── constants.go
│   ├── decimal.go
│   ├── divide.go
│   ├── errors.go
│   ├── format.go
│   ├── helpers.go
│   ├── parse.go
│   └── rounding.go
├── tests/
│   ├── original/
│   │   └── test-suite.sha256        # SHA256 hash of original JS test suite, pinned at kickoff
│   └── port/
│       ├── test_original_test.go    # original test assertions ported 1:1
│       └── decimal_port_test.go     # additional Go-idiomatic tests
├── fuzz/
│   ├── harness.go                   # differential fuzz harness
│   └── log.txt                      # actual 60s+ run output
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
# Direct Go command: go test ./src/... -bench=. -benchmem
```

### Run Differential Fuzzer (60s+)

```bash
make fuzz
# Direct Go command: go test ./fuzz/... -fuzz=FuzzDecimal -fuzztime=60s
```

### Docker Containerized Build & Test

```bash
docker build -t decimal-go-port .
```
