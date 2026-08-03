#!/usr/bin/env node
// oracle.js — Differential fuzz oracle using the original decimal.js library.
//
// Protocol: Reads newline-delimited JSON from stdin. Each line is a JSON
// object with fields {a: string, b: string}. For each input, the oracle
// computes a set of operations using the original decimal.js and writes a
// single JSON line to stdout with the results.
//
// This runs as a long-lived child process — the Go test pipes batches in
// and reads results out without restarting Node for every pair.

'use strict';

const Decimal = require('decimal.js');
const readline = require('readline');

// Use default decimal.js configuration (precision 20, rounding ROUND_HALF_UP).
// This matches the Go DefaultContext().
Decimal.set({
  precision: 20,
  rounding: Decimal.ROUND_HALF_UP,
  toExpNeg: -7,
  toExpPos: 21,
  minE: -9e15,
  maxE: 9e15,
});

const rl = readline.createInterface({ input: process.stdin, terminal: false });

rl.on('line', (line) => {
  let input;
  try {
    input = JSON.parse(line);
  } catch (e) {
    // Malformed JSON — skip.
    process.stdout.write(JSON.stringify({ error: 'bad_json' }) + '\n');
    return;
  }

  const { a: aStr, b: bStr } = input;
  const result = {};

  try {
    // Parse both values through original decimal.js.
    const a = new Decimal(aStr);
    const b = new Decimal(bStr);

    // --- Shared Public API operations ---

    // 1. toString (parse roundtrip)
    result.a_str = a.toString();
    result.b_str = b.toString();

    // 2. valueOf
    result.a_valueOf = a.valueOf();
    result.b_valueOf = b.valueOf();

    // 3. Arithmetic
    result.plus = a.plus(b).toString();
    result.minus = a.minus(b).toString();
    result.times = a.times(b).toString();

    // Division: skip if b is zero/NaN to avoid messy Infinity comparisons
    // that depend on sign handling — we still test finite/finite division.
    result.div = a.div(b).toString();

    // 4. Unary
    result.abs_a = a.abs().toString();
    result.neg_a = a.neg().toString();

    // 5. Comparison
    try {
      const c = a.cmp(b);
      result.cmp = c.toString();
    } catch (e) {
      result.cmp = 'NaN';
    }

    // 6. Predicates
    result.isNaN_a = a.isNaN();
    result.isFinite_a = a.isFinite();
    result.isZero_a = a.isZero();
    result.isNeg_a = a.isNeg();
    result.isPos_a = a.isPos();

    // 7. Formatting (only for finite values where dp makes sense)
    if (a.isFinite()) {
      result.toFixed_a = a.toFixed(10);
    }

    // 8. Square root (for non-negative finite values)
    if (a.isFinite() && !a.isNeg()) {
      result.sqrt_a = a.sqrt().toString();
    }

    // 9. Floor / Ceil / Trunc / Round
    if (a.isFinite()) {
      result.floor_a = a.floor().toString();
      result.ceil_a = a.ceil().toString();
      result.trunc_a = a.trunc().toString();
      result.round_a = a.round().toString();
    }

    // 10. Mod (skip if b is zero)
    if (b.isFinite() && !b.isZero()) {
      result.mod = a.mod(b).toString();
    }

  } catch (e) {
    result.error = e.message || 'unknown_error';
  }

  process.stdout.write(JSON.stringify(result) + '\n');
});

rl.on('close', () => {
  process.exit(0);
});
