#!/usr/bin/env node
'use strict';

const path = require('path');

function option(name) {
  const index = process.argv.indexOf(name);
  return index < 0 ? undefined : process.argv[index + 1];
}

const jsRepository = path.resolve(
  option('--js-repo') || path.join(__dirname, '..', '..', '..', 'decimal.js'),
);
const Decimal = require(path.join(jsRepository, 'decimal.js'));

if (process.argv.includes('--startup-probe')) {
  process.stdout.write('ready\n');
  process.exit(0);
}

const operations = Number(option('--ops') || '100000');
if (!Number.isInteger(operations) || operations < 1) {
  throw new Error('--ops must be a positive integer');
}

Decimal.set({
  precision: 40,
  rounding: Decimal.ROUND_HALF_UP,
  toExpNeg: -9000000000000000,
  toExpPos: 9000000000000000,
  minE: -9000000000000000,
  maxE: 9000000000000000,
});

const a = new Decimal('123456.78901234567890123456789');
const b = new Decimal('987654.32109876543210987654321');
const c = new Decimal('0.000000123456789');
const d = new Decimal('1.0000000000000001');
const divisor = new Decimal('7.123456789');
const samples = new Array(operations);
let peakRss = process.memoryUsage().rss;
let sink = '';

const started = process.hrtime.bigint();
for (let i = 0; i < operations; i++) {
  const opStarted = process.hrtime.bigint();
  const value = a.plus(b).minus(c).times(d).div(divisor);
  sink = value.toFixed(12) + value.toDecimalPlaces(10).valueOf() + value.round().valueOf();
  samples[i] = Number(process.hrtime.bigint() - opStarted);

  if ((i & 1023) === 0) {
    peakRss = Math.max(peakRss, process.memoryUsage().rss);
  }
}
const elapsedNs = Number(process.hrtime.bigint() - started);
peakRss = Math.max(peakRss, process.memoryUsage().rss);
samples.sort((left, right) => left - right);

if (sink.length === 0) throw new Error('benchmark result was unexpectedly empty');

const p99Index = Math.max(0, Math.ceil(samples.length * 0.99) - 1);
process.stdout.write(JSON.stringify({
  runtime: `Node.js ${process.version}`,
  p99_latency_ns: samples[p99Index],
  rss_kb: Math.ceil(peakRss / 1024),
  throughput_ops_sec: operations * 1e9 / elapsedNs,
  operations,
}) + '\n');
