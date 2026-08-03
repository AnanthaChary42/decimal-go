#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');

function option(name) {
  const index = process.argv.indexOf(name);
  return index < 0 ? undefined : process.argv[index + 1];
}

const aggregatePath = option('--aggregate');
const powSqrtPath = option('--pow-sqrt');
const outputPath = option('--output');
if (!aggregatePath || !powSqrtPath || !outputPath) {
  throw new Error('--aggregate, --pow-sqrt, and --output are required');
}

function readResult(file) {
  return JSON.parse(fs.readFileSync(path.resolve(file), 'utf8'));
}

function countGoPortTests(directory) {
  let count = 0;
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    const entryPath = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      count += countGoPortTests(entryPath);
    } else if (entry.isFile() && entry.name.endsWith('_test.go')) {
      const source = fs.readFileSync(entryPath, 'utf8');
      count += (source.match(/^func\s+TestOriginal_/gm) || []).length;
    }
  }
  return count;
}

const aggregate = readResult(aggregatePath);
const powSqrt = readResult(powSqrtPath);
if (aggregate.status !== 'passed' || powSqrt.status !== 'passed') {
  throw new Error('Cannot write a passing verification record from a failed JavaScript source suite.');
}
if (aggregate.reference_library !== powSqrt.reference_library) {
  throw new Error('The aggregate and powSqrt suites used different decimal.js reference libraries.');
}
if (aggregate.reference_repository !== powSqrt.reference_repository) {
  throw new Error('The aggregate and powSqrt suites used different decimal.js source repositories.');
}

const repoRoot = path.resolve(__dirname, '..', '..');
const verification = {
  status: 'passed',
  original_js: {
    status: 'passed',
    module_files: aggregate.module_files + powSqrt.module_files,
    canonical_runner_modules: aggregate.module_files,
    standalone_modules: ['powSqrt.js'],
    assertions: aggregate.total + powSqrt.total,
    canonical_runner_assertions: aggregate.total,
    standalone_assertions: powSqrt.total,
    runner: 'test/test.js + test/modules/powSqrt.js',
    reference_library: aggregate.reference_library,
    reference_repository: aggregate.reference_repository,
  },
  port_go: {
    status: 'passed',
    command: 'go test -count=1 ./tests/port',
    top_level_test_functions: countGoPortTests(path.join(repoRoot, 'tests', 'port')),
    scope: 'All current Go port-package tests were executed.',
    parity_note: 'This records execution of the current Go port tests; it does not claim a one-to-one translation of every original JavaScript assertion.',
  },
};

fs.writeFileSync(path.resolve(outputPath), `${JSON.stringify(verification, null, 2)}\n`);
