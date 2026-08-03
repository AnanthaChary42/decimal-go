# Benchmark Methodology: decimal.js vs decimal-go

## Overview

`make bench` first verifies correctness, then measures the checked-out original Node.js source at `../decimal.js/decimal.js` and the native `decimal-go` port with the same mixed-decimal workload. It writes the machine-specific results to `bench/results.json` and prints a comparison table.

On Windows, set `$jsRepository` at the top of `bench/run.ps1` to the absolute location of the original `decimal.js` clone. The Go location is automatically the repository containing that script. On systems using GNU Make, override `JS_REPO`, for example `make bench JS_REPO=/path/to/decimal.js`.

## Required test-suite preflight

By default, every benchmark run performs these checks before any timing begins:

- **Original JavaScript suite:** Runs the unmodified source `test/test.js`, which loads 60 module files, then separately runs the standalone `test/modules/powSqrt.js`. Together this executes all **61 JavaScript test-module files**.
- **Go port suite:** Runs `go test -count=1 ./tests/port`, which executes every current Go port test (currently 70 top-level `TestOriginal_` functions, including split-module and regression tests).

The generated `bench/results.json` includes a `suite_verification` object with the JavaScript assertion totals, source repository path/version, and the Go test command. A failed preflight stops the benchmark; timing results are never produced for a failed suite.

The launcher prints only failed JavaScript assertions and each suite's final total, so failure details are not lost in routine per-module progress output.

`powSqrt.js` is not included by the original aggregate runner and has an undeclared, non-incrementing `total` loop counter. The benchmark runner supplies that missing standalone-driver counter and advances it after each original assertion, causing its intended 10,000 unchanged assertions to execute.

This is a correctness gate, not a timing workload: the performance numbers remain the fixed mixed-operation loop described below. It also does not claim that every JavaScript assertion already has a one-to-one Go translation; it records that all original JavaScript tests and all currently present Go port tests were executed.

## Metrics

- **Workload**: A repeatable mixed-operation iteration: addition, subtraction, multiplication, division, fixed-point formatting, decimal-place rounding, and whole-number rounding.
- **Runtime**: The installed Node.js and Go versions collected at run time.
- **P99 latency**: The 99th percentile duration of one complete mixed-operation iteration, measured with a monotonic clock.
- **RSS**: The maximum sampled resident working set of the benchmark process. Windows uses `GetProcessMemoryInfo`; the non-Windows Go fallback reports Go runtime memory because the standard library has no portable RSS API.
- **Startup**: Wall-clock time from spawning a fresh process until its startup probe writes a readiness line. The Go executable is built before startup is measured, so compilation is excluded.
- **Throughput**: Completed mixed-operation iterations divided by the total workload time.

Run the default 100,000-iteration comparison:

```powershell
make bench
```

On Windows without GNU Make, run the included PowerShell launcher instead:

```powershell
powershell -ExecutionPolicy Bypass -File .\bench\run.ps1
```

For a clean build-cache run (the prior executable and result file are removed, and `go clean -cache` is run first):

```powershell
powershell -ExecutionPolicy Bypass -File .\bench\run.ps1 -Clean
```

`-Clean` is usually unnecessary between benchmark repetitions. It can make the first fresh executable startup unrepresentative because Windows security scanning may occur when a new `.exe` is created.

Run one million iterations:

```powershell
make bench BENCH_OPS=1000000
```

Or:

```powershell
powershell -ExecutionPolicy Bypass -File .\bench\run.ps1 -Ops 1000000
```

## Confounders

1. **V8 JIT warmup**: Node.js optimises JavaScript while it executes, so first-run latency and steady-state throughput are different observations.
2. **Garbage collection**: Both implementations allocate during the workload. P99 timing intentionally includes pauses a caller can observe.
3. **Host environment**: Compare measurements only on the same machine and with no significant competing workload.

## Results

No fixed benchmark figures are stored here. Every `make bench` or `run.ps1` run replaces `bench/results.json` with measured `runtime`, `p99_latency_ns`, `rss_kb`, `startup_ms`, and `throughput_ops_sec` values for both implementations, plus derived ratios and the preflight verification record.
