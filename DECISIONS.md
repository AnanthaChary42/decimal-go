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
