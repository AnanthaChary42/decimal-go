# Design Decisions — decimal.js → Go Port

This document records every non-trivial architectural divergence between the original JavaScript `decimal.js` library and the Go `decimal-go` port, documenting the decision made, what the original implementation did, and the technical rationale for the divergence.

---

## D001: Digit Array Storage — Base 1e7 (`[]int32`) vs `math/big.Int`

- **What was decided**: Retain `decimal.js`'s custom base-1e7 word array (`[]int32`) representation (`BASE = 10,000,000`, `LOG_BASE = 7`) rather than wrapping Go's standard `math/big.Int`.
- **What JS original did**: Used a dynamic array of double-precision numbers (`d = [1234567, 8901234]`) operating in base-1e7.
- **Why the divergence**: Wrapping `math/big.Int` would require constant string/base conversions and obscure the exact multi-word division, multiplication, and rounding digit indexing. Keeping the base-1e7 word array enables 1:1 algorithmic parity with `decimal.js` internals while optimizing memory locality in Go slices.

---

## D002: Configuration Context — Immutable `Context` Struct vs Mutable Global State

- **What was decided**: Introduce an explicit `Context` struct (`*Context`) holding precision, rounding mode, and exponent thresholds (`ToExpNeg`, `ToExpPos`, `MinE`, `MaxE`), passed down to methods or attached to instances.
- **What JS original did**: Used global mutable properties on the constructor (`Decimal.precision = 20`, `Decimal.rounding = 4`, `Decimal.set(...)`).
- **Why the divergence**: Global mutable configuration in JavaScript causes data races and non-deterministic behavior in concurrent Go applications (`goroutines`). An explicit `Context` ensures thread-safe, deterministic decimal math.

---

## D003: Error Handling — `(Decimal, error)` Return Tuples vs JavaScript Exceptions

- **What was decided**: Constructor functions (`New`, `NewFromFloat64`, `ToFixed`, etc.) return `(*Decimal, error)`. Invalid inputs or precision limit violations return explicit `DecimalError` values (e.g., `ErrInvalidArg`).
- **What JS original did**: Threw runtime JavaScript `Error` objects (`throw Error("[DecimalError] Invalid argument")`).
- **Why the divergence**: Idiomatic Go avoids exception throwing (`panic`/`recover`) in favor of explicit `(value, error)` returns, giving callers compile-time control over error handling.

---

## D004: Package Layout — Modular `src/` Architecture

- **What was decided**: Organize core implementation files under `src/` (`package decimal`), test suites under `tests/port/`, kickoff hash under `tests/original/`, and fuzzing under `fuzz/`.
- **What JS original did**: Monolithic 4,953-line single file (`decimal.js`) at repository root.
- **Why the divergence**: Breaking the 4,953-line JavaScript file into domain-specific Go modules (`arithmetic.go`, `divide.go`, `parse.go`, `format.go`, `rounding.go`) improves code maintainability, test isolation, and Go tooling compatibility while keeping root clean.

---

## D005: NaN and Infinity Representation — Sign-Based Sentinel Encoding

- **What was decided**: Represent NaN as `s = 0, d = nil`, `+Infinity` as `s = 1, d = nil`, and `-Infinity` as `s = -1, d = nil`.
- **What JS original did**: Encoded NaN as `s = NaN` (floating-point NaN) and `d = null`, and Infinities as `s = ±1` with `d = null`.
- **Why the divergence**: Go's integer type `int` cannot store float `NaN`. Using integer `0` for NaN and a `nil` slice for Infinities avoids adding a separate `kind` enum byte, keeping the `Decimal` struct lightweight (24 bytes).

---

## D006: Standalone Knuth Long Division Engine (`divide.go`)

- **What was decided**: Extracted the Knuth Algorithm D long division logic into a standalone internal module `divide.go` supporting variable precision and modulo mode flags.
- **What JS original did**: Embedded division directly within a monolithic 260-line closure inside `decimal.js`.
- **Why the divergence**: Multi-word array division requires normalization scaling (`multiplyByInt`, `multiplyAndSubtract`, `addBack`). Isolating it into `divide.go` enables targeted profiling, benchmarking, and unit testing of division bottlenecks.

---

## D007: Native Go 1.18+ Fuzzing Engine (`fuzz/harness.go`)

- **What was decided**: Built a native Go `fuzz.F` test harness in `fuzz/harness.go` testing algebraic invariants (commutativity `a+b == b+a`, self-cancellation `a-a == 0`, and string parsing symmetry).
- **What JS original did**: Relied on external Python Hypothesis property scripts in `test/hypothesis/`.
- **Why the divergence**: Native Go fuzzing integrates directly with standard Go CLI tooling (`go test -fuzz`) and requires no Python environment dependencies.

---

## D008: Strongly-Typed IEEE 754-2008 Rounding Mode Enum

- **What was decided**: Defined a strongly-typed `RoundingMode int` type with constants (`RoundUp`, `RoundDown`, `RoundCeil`, `RoundFloor`, `RoundHalfUp`, `RoundHalfDown`, `RoundHalfEven`, `RoundHalfCeil`, `RoundHalfFloor`, `Euclid`).
- **What JS original did**: Assigned untyped numeric properties on `Decimal` (`Decimal.ROUND_UP = 0`, `Decimal.ROUND_DOWN = 1`, etc.).
- **Why the divergence**: Strong static typing prevents invalid integer values from being passed to rounding methods at compile time.

---

## D009: Precomputed High-Precision Constants (`LN10` and `PI`)

- **What was decided**: Embedded 1025-digit precision constants for `LN10` and `PI` as string literals in `constants.go`, lazily sliced and rounded based on active `Context` precision.
- **What JS original did**: Dynamically constructed or sliced internal constant strings within `getLn10` and `getPi`.
- **Why the divergence**: Static string literals eliminate dynamic calculation overhead and memory allocations during transcendental math initialization in Go.

---

## D010: Native `Float64()` Export Signature

- **What was decided**: `Float64()` returns `(float64, bool)` indicating both the converted IEEE 754 float64 value and whether the conversion was exact (no precision lost). Maps `Infinity` to `math.Inf(1)` and `math.Inf(-1)`, and `NaN` to `math.NaN()`.
- **What JS original did**: `toNumber()` returning a primitive JavaScript Number or `Infinity`/`NaN`.
- **Why the divergence**: Matches Go's standard library `strconv` convention where precision loss is explicitly reported via boolean flag.

---

## D011: Containerized One-Step Build & Automation Tooling

- **What was decided**: Provided a root `Makefile` and multi-stage `Dockerfile` supporting `make build`, `make test`, `make bench`, `make fuzz`, and `docker build -t decimal-go-port .`.
- **What JS original did**: Used `package.json` npm scripts (`npm test`).
- **Why the divergence**: Fulfills Port Mortem Rule 03 ("Standalone & runnable") across both Docker and local terminal environments without requiring Node/npm.
