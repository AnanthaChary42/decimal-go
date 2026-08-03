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

## D012: Transcendental Functions — `math/big.Float` Series Instead of Word-Array Taylor Expansion

**Decision:** Implement `Ln()`, `Exp()`, `Log()`, and `Pow()` using `math/big.Float` arithmetic (atanh-series for ln, Taylor reduction+squaring for exp) rather than translating the JS word-array `naturalLogarithm` and `naturalExponential` routines.

**Rationale:**

- decimal.js's custom `naturalLogarithm` and `naturalExponential` perform Taylor expansion directly on the base-1e7 digit arrays, interleaving rounding with coefficient arithmetic. Faithfully porting this requires reproducing hundreds of lines of intricate index-manipulation code that lacks clear invariants.
- `math/big.Float` provides arbitrary-precision binary floating-point with tunable mantissa bits, letting us match or exceed the precision of the JS originals with a fraction of the code. The guard-bit budget (`decimalDigits * 4 + 64`) ensures the final conversion through `decimalFromBigFloat` never loses a significant digit below the configured precision.
- Trade-off: introduces a dependency on `math/big`, which adds ~1 MB to the binary. This is acceptable because Go's standard library `math/big` is battle-tested, zero-allocation-per-call after setup, and the alternative (hand-porting 400+ lines of JS series code operating on word arrays) would be the project's highest-risk correctness surface.
- The same `math/big.Float` infrastructure is reused across `Ln`, `Exp`, `Log`, `Pow`, and all trigonometric/hyperbolic inverse functions, amortising the implementation cost.

## D013: Exact Rational Division — `divideFiniteExact` via `math/big.Int`

**Decision:** For ordinary decimal division (non-base-conversion), bypass the JS word-array long-division algorithm and instead compute the quotient as a scaled `big.Int` division followed by `finalise`.

**Rationale:**

- JS's `divide()` function performs long division on base-1e7 word arrays. Porting it faithfully for the base-1e7 path (logBase == 7) is error-prone: off-by-one in the remainder tracking caused rounding-boundary errors at word-array boundaries.
- `divideFiniteExact` extracts the integer significands as `big.Int`, scales the numerator to produce `sd + 1` significant digits, performs a single `QuoRem`, and feeds the quotient text through `ctx.New()` → `finalise()`. This is mathematically identical but structurally simpler.
- The original word-array `divide()` is retained for the base-conversion path (logBase == 1), where `toBaseString` needs digit-by-digit division in a non-decimal base. The conditional at `divide.go:116` dispatches between the two paths.
- Trade-off: the `big.Int` path allocates a temporary string for the quotient. This adds one allocation compared to the word-array path but eliminates the class of bugs caused by mixing rounding with base-1e7 carry propagation.

## D014: `ToFraction` and `ToNearest` — Exact Rational Arithmetic via `math/big`

**Decision:** Implement `ToFraction` using a continued-fraction algorithm on `big.Int` numerator/denominator pairs, and `ToNearest` using a `big.Int` quotient-remainder with explicit tie-breaking. The JS originals use word-array Decimal arithmetic for both.

**Rationale:**

- JS's `toFraction` builds the continued-fraction convergents using Decimal arithmetic (with `external = false`), which means the intermediate p₀/q₀ values carry base-1e7 word arrays and are subject to precision truncation. Using `big.Int` convergents is exact and avoids subtle truncation at large denominators.
- JS's `toNearest` divides, rounds the quotient to an integer, and multiplies back — all with Decimal operations that are precision-limited. The Go port extracts the exact rational value via `decimalRational()`, computes `QuoRem`, and rounds the integer quotient with a type switch over all 9 rounding modes. This is both more precise and easier to audit.
- The `decimalRational()` helper (compat.go:322) converts any finite `Decimal` to its exact `(numerator, denominator)` as `big.Int`, which is reused by `ToFraction`, `ToNearest`, and `baseDigitsRounded`. This factorisation has no JS equivalent.

## D015: Base Output Conversion — `big.Int` Significand Instead of Word-Array Division

**Decision:** Implement `ToBinary`, `ToOctal`, and `ToHex` by extracting the Decimal's exact rational value into `big.Int` numerator/denominator, computing the significand as a scaled `big.Int` in the target base, and applying rounding via `baseRoundUp()`.

**Rationale:**

- JS's `toStringBinary` relies on `divide()` in base-2/8/16 mode (logBase = 1), performing digit-by-digit long division in the target base. This is the most complex path through `divide()` and is the source of the `inexact` global flag.
- The Go port avoids this entirely: `baseDigitsRounded()` computes `floor(log_base(value))` using `floorLogBase()`, scales the numerator to produce `sd` significant digits, divides once, and rounds. The result is converted to the target-base string via `big.Int.Text(base)`.
- Trade-off: loses the ability to report `inexact` from the division loop. Instead, `baseRoundUp` detects inexactness from the `big.Int` remainder directly. This is functionally equivalent — both paths determine whether the result is exact or needs rounding.

## D016: `crypto/rand` for Random Instead of `Math.random()` / Web Crypto

**Decision:** `Context.Random()` uses `crypto/rand.Int()` unconditionally (standard library CSPRNG), rather than providing separate `Math.random()` and `crypto.getRandomValues()` paths.

**Rationale:**

- JS's `Decimal.random()` defaults to `Math.random()` (non-cryptographic PRNG) and optionally uses `crypto.getRandomValues()` when `config({ crypto: true })` is set.
- Go does not have a built-in equivalent of `Math.random()` that returns a float64 with enough bits. The `crypto/rand` package provides platform-native CSPRNG (`/dev/urandom` on Linux, `CryptGenRandom` on Windows) with no warm-up cost.
- Using CSPRNG unconditionally simplifies the implementation (no conditional dispatch on `ctx.Crypto`) and is strictly more secure. The `Crypto` config field is retained in the `Context` struct for API compatibility but does not change Random's behaviour.

## D017: Pointer-Based `ConfigOptions` to Distinguish Unset from Zero

**Decision:** Use `*int` and `*bool` fields in `ConfigOptions` rather than bare value types.

**Rationale:**

- In JS, `Decimal.config({ precision: 20 })` passes an object where only the `precision` key is present; other fields are `undefined`. The JS library checks `obj.precision !== void 0` to skip unset fields.
- In Go, a zero-valued `int` field (e.g., `Precision: 0`) is indistinguishable from "not set" without pointer indirection. Pointer fields let `applyConfig` check `if opts.Precision != nil` before applying, preserving the JS "partial update" semantics exactly.
- Helper constructors `IntPtr()` and `BoolPtr()` reduce boilerplate at call sites.

## D018: `RoundingMode` Named Type Instead of Raw Integers

**Decision:** Define `RoundingMode` as a named `int` type with `const` iota-like definitions (RoundUp=0 through Euclid=9), rather than using bare `int` constants.

**Rationale:**

- JS uses raw integer constants (`ROUND_UP = 0`, etc.) and relies on runtime range checks. Go's type system can provide compile-time safety: a function accepting `RoundingMode` will reject an arbitrary `int` at the call site without explicit conversion.
- The named type also makes function signatures self-documenting: `func (x *Decimal) ToFixed(dp int, rm ...RoundingMode) string` immediately conveys that `rm` is a rounding mode, not an arbitrary integer.
- The `Modulo` field in `Context` also uses `RoundingMode`, unifying the two concepts (rounding mode and modulo mode both use the same 0–9 range from decimal.js).

## D019: `goto checkOverflow` in `finalise` Instead of Refactored Early Return

**Decision:** Use a `goto checkOverflow` label in `finalise()` (helpers.go:299) when no rounding is needed but overflow/underflow checks still apply.

**Rationale:**

- JS's `finalise` uses nested `if/else` with fallthrough: when `xdi >= k` and `!isTruncated`, it falls through to the overflow check at the function's end without executing the rounding logic.
- Go lacks statement-level fallthrough in `if` blocks. Refactoring into multiple functions would require passing 8+ local variables. A `goto` to the shared overflow check preserves the original control flow with minimal indirection.
- This is the only `goto` in the codebase. It is forward-only (jumps down, never up), structurally equivalent to a labeled break, and avoids duplicating the 10-line overflow/underflow block.
