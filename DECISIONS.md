# Design Decisions — decimal.js → Go Port

## D001: Digit Storage — Base 1e7 `[]int32`

**Decision:** Use the same base-1e7 digit array as decimal.js rather than wrapping `math/big.Int`.

**Rationale:**

- Near-identical internal representation → easier to verify behavioral equivalence.
- Direct translation of the division, multiply, and rounding algorithms.
- Trade-off: more code, but Port Mortem judges value "ported feel" and we control every edge case.

## D002: Configuration — Explicit `Context` struct

**Decision:** Use an explicit `Context` struct rather than global mutable config.

**Rationale:**

- JS uses mutable global config (`Decimal.precision = 50`). This is the biggest intentional divergence.
- More idiomatic Go, avoids data races in concurrent code.
- A default context is provided for convenience (`defaultCtx`).

## D003: Error Handling — `(Decimal, error)` returns

**Decision:** Constructors return `(*Decimal, error)`. NaN/Infinity are valid `Decimal` values (matching JS behavior), but parse failures return errors.

**Rationale:**

- JS throws `Error` on invalid input. Go returns error tuples.
- NaN/Infinity are representable as special `Decimal` values (s=0 for NaN, d=nil for Inf).

## D004: Package Structure — Root package

**Decision:** Use `package decimal` at the repo root rather than a `src/` subdirectory.

**Rationale:**

- More idiomatic Go (import as `github.com/AnanthaChary42/decimal-go`).
- The `src/` directory in the original plan was a placeholder.

## D005: NaN/Infinity Representation — Sign-based sentinels

**Decision:** Encode NaN as `s=0, d=nil` and Infinity as `s=±1, d=nil`.

**Rationale:**

- Simpler than adding a separate `kind` field.
- Matches the JS representation closely (JS uses `s=NaN` for NaN, `d=null` for Infinity).
- In Go we use `s=0` (neither positive nor negative) instead of float NaN.

## D006: `external` Variable — Package-level global

**Decision:** Keep the `external` variable as a package-level `bool`, matching JS.

**Rationale:**

- Many internal functions temporarily set `external = false` to suppress overflow/underflow checks.
- Converting to a context field would require threading it through dozens of functions.
- Trade-off: not goroutine-safe, but matches JS semantics exactly.
- Future: could be moved to context if concurrency is needed.

## D007: Division Algorithm NaN Handling

**Decision:** Check for zero-divided-by-zero (`xd[0] == 0 && yd[0] == 0`) specifically rather than checking equality of first digit arrays (`xd[0] == yd[0]`).

**Rationale:**

- The JS original checked `xd[0] == yd[0]` within a block where both arrays were already known to represent 0.
- Unconditional `xd[0] == yd[0]` caused valid equal-operand divisions (such as 5 / 5) to incorrectly evaluate to NaN.
- Specific zero-check accurately targets `0 / 0` and `Inf / Inf` while permitting exact equal number division.

## D008: Word Array Rounding (kVal Offset and Leading Digits)

**Decision:** Set `kVal = 1` when the rounding position `i == 0` in `finalise()`, preserving exact power-of-10 offsets for word array propagation.

**Rationale:**

- In JS `decimal.js`, when rounding at the exact word boundary (`i == 0`), incrementing the previous word requires `k = 1`.
- Calculating `10^(LOG_BASE - i)` when `i == 0` resulted in adding `10^7` instead of `1`, causing string representation corruptions (e.g. `Ceil(1.5)` producing `"10.000001"`).
- Restoring `kVal = 1` aligns rounding carry-overs with base-1e7 word storage boundaries.

## D009: Array Subtraction Loop Boundaries

**Decision:** Adjust the subtraction loop condition in `Sub()` from `i > k` to `i >= k` (down to index 0 inclusive).

**Rationale:**

- 0-indexed word arrays in Go stop prior to index 0 if the loop threshold is `i > 0`.
- Single-word digit operations (e.g., `3 - 2` or `1 - 1`) were skipping the subtraction step at index 0, returning unchanged operands.
- Including index 0 (`i >= k`) guarantees all digits are subtracted correctly down to the most significant word.

## D010: Native Go Fuzzing Harness

**Decision:** Implement `FuzzParseDecimal` using Go 1.18+ native `testing.F` in package `fuzz`.

**Rationale:**

- Provides automated coverage for decimal strings, exponential notation, hex/binary/octal parsing, and special float strings.
- Integrates seamlessly into standard Go tooling via `go test ./fuzz` and `make fuzz`.

## D011: IEEE 754 Special Value Float Handling

**Decision:** Use Go standard library `math.Inf(1)` and `math.Inf(-1)` in place of constant float division `1.0 / 0.0`.

**Rationale:**

- Go strict compiler prevents constant float division by zero at compile time.
- Standard library `math` primitives maintain standard IEEE 754 infinity semantics without compiler error.
