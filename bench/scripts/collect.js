#!/usr/bin/env node
'use strict';

const childProcess = require('child_process');
const fs = require('fs');
const path = require('path');

function option(name) {
  const index = process.argv.indexOf(name);
  return index < 0 ? undefined : process.argv[index + 1];
}

const goBinOption = option('--go-bin');
const jsRepositoryOption = option('--js-repo');
const operations = Number(option('--ops') || '100000');
const output = option('--output') || path.join(__dirname, 'results.json');
const suiteVerificationPath = option('--suite-verification');
if (!goBinOption) throw new Error('--go-bin is required');
if (!jsRepositoryOption) throw new Error('--js-repo is required');
if (!Number.isInteger(operations) || operations < 1) throw new Error('--ops must be a positive integer');

const goBin = path.resolve(goBinOption);
const jsRepository = path.resolve(jsRepositoryOption);
const jsBenchmark = path.join(__dirname, 'js_bench.js');
const suiteVerification = suiteVerificationPath
  ? JSON.parse(fs.readFileSync(path.resolve(suiteVerificationPath), 'utf8'))
  : { status: 'not_run', reason: 'No suite-verification file was provided.' };

function startupMs(command, args) {
  return new Promise((resolve, reject) => {
    const started = process.hrtime.bigint();
    const child = childProcess.spawn(command, args, { stdio: ['ignore', 'pipe', 'pipe'] });
    let stderr = '';
    let resolved = false;
    child.stderr.on('data', (data) => { stderr += data; });
    child.stdout.once('data', () => {
      resolved = true;
      resolve(Number(process.hrtime.bigint() - started) / 1e6);
    });
    child.on('error', reject);
    child.on('close', (code) => {
      if (!resolved) reject(new Error(`${command} startup probe exited with ${code}: ${stderr.trim()}`));
    });
  });
}

function runReport(command, args) {
  return new Promise((resolve, reject) => {
    const child = childProcess.spawn(command, args, { stdio: ['ignore', 'pipe', 'pipe'] });
    let stdout = '';
    let stderr = '';
    child.stdout.on('data', (data) => { stdout += data; });
    child.stderr.on('data', (data) => { stderr += data; });
    child.on('error', reject);
    child.on('close', (code) => {
      if (code !== 0) {
        reject(new Error(`${command} benchmark exited with ${code}: ${stderr.trim()}`));
        return;
      }
      try {
        resolve(JSON.parse(stdout));
      } catch (error) {
        reject(new Error(`${command} did not emit benchmark JSON: ${error.message}\n${stdout}\n${stderr}`));
      }
    });
  });
}

function finiteRatio(numerator, denominator) {
  return denominator === 0 ? null : numerator / denominator;
}

function displayNumber(value, fractionDigits = 2) {
  return Number(value).toLocaleString('en-US', { maximumFractionDigits: fractionDigits });
}

async function main() {
  const jsStartup = await startupMs(process.execPath, [jsBenchmark, '--js-repo', jsRepository, '--startup-probe']);
  const goStartup = await startupMs(goBin, ['--startup-probe']);
  const originalJs = await runReport(process.execPath, [jsBenchmark, '--js-repo', jsRepository, '--ops', String(operations)]);
  const portGo = await runReport(goBin, ['--ops', String(operations)]);

  originalJs.startup_ms = jsStartup;
  portGo.startup_ms = goStartup;

  const latencySpeedup = finiteRatio(originalJs.p99_latency_ns, portGo.p99_latency_ns);
  const throughputIncrease = finiteRatio(portGo.throughput_ops_sec, originalJs.throughput_ops_sec);
  const startupSpeedup = finiteRatio(originalJs.startup_ms, portGo.startup_ms);
  const memoryReduction = originalJs.rss_kb === 0
    ? null
    : (1 - portGo.rss_kb / originalJs.rss_kb) * 100;

  const report = {
    benchmark_name: `decimal.js vs decimal-go (${operations.toLocaleString('en-US')} mixed operations)`,
    workload: 'add, sub, mul, div, fixed formatting, decimal-place rounding, whole-number rounding',
    measured_at: new Date().toISOString(),
    original_js: originalJs,
    port_go: portGo,
    suite_verification: suiteVerification,
    improvements: {
      latency_speedup: latencySpeedup === null ? null : `${latencySpeedup.toFixed(2)}x`,
      throughput_increase: throughputIncrease === null ? null : `${throughputIncrease.toFixed(2)}x`,
      memory_reduction_pct: memoryReduction === null ? null : `${memoryReduction.toFixed(1)}%`,
      startup_speedup: startupSpeedup === null ? null : `${startupSpeedup.toFixed(2)}x`,
    },
  };

  fs.writeFileSync(output, `${JSON.stringify(report, null, 2)}\n`);
  console.log(`\nMeasured ${operations.toLocaleString('en-US')} mixed operations per runtime`);
  console.table([
    {
      implementation: 'decimal.js',
      runtime: originalJs.runtime,
      p99_latency_ns: displayNumber(originalJs.p99_latency_ns, 0),
      rss_kb: displayNumber(originalJs.rss_kb, 0),
      startup_ms: displayNumber(originalJs.startup_ms),
      throughput_ops_sec: displayNumber(originalJs.throughput_ops_sec, 0),
    },
    {
      implementation: 'decimal-go',
      runtime: portGo.runtime,
      p99_latency_ns: displayNumber(portGo.p99_latency_ns, 0),
      rss_kb: displayNumber(portGo.rss_kb, 0),
      startup_ms: displayNumber(portGo.startup_ms),
      throughput_ops_sec: displayNumber(portGo.throughput_ops_sec, 0),
    },
  ]);
  if (suiteVerification.status === 'passed') {
    const js = suiteVerification.original_js;
    console.log(
      `Preflight verification passed: ${js.module_files} original JavaScript module files ` +
      `(${js.assertions.toLocaleString('en-US')} assertions) and all current Go port tests.`,
    );
  } else if (suiteVerification.status === 'skipped') {
    console.log('Preflight verification was skipped (performance-only run).');
  } else {
    console.log('Preflight verification was not recorded for this run.');
  }
  console.log(`Results written to ${output}`);
}

main().catch((error) => {
  console.error(`Benchmark failed: ${error.message}`);
  process.exitCode = 1;
});
