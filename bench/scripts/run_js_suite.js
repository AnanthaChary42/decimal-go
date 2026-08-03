#!/usr/bin/env node
'use strict';

// Run the original decimal.js test files without modifying them. The source
// suite loads "../decimal" from test/setup.js, which resolves directly to the
// decimal.js implementation in the supplied original repository.
const fs = require('fs');
const path = require('path');

function option(name) {
  const index = process.argv.indexOf(name);
  return index < 0 ? undefined : process.argv[index + 1];
}

const suite = option('--suite');
const resultPath = option('--result');
const failuresOnly = process.argv.includes('--failures-only');
const sourceRepository = path.resolve(
  option('--js-repo') || path.join(__dirname, '..', '..', '..', 'decimal.js'),
);
const sourceTestRoot = path.join(sourceRepository, 'test');
const sourceDecimal = path.join(sourceRepository, 'decimal.js');

if (!['aggregate', 'pow-sqrt'].includes(suite)) {
  throw new Error('--suite must be aggregate or pow-sqrt');
}
if (!resultPath) throw new Error('--result is required');
if (!fs.existsSync(sourceTestRoot) || !fs.existsSync(sourceDecimal)) {
  throw new Error(
    `Original decimal.js repository was not found or incomplete: ${sourceRepository}`,
  );
}

const decimalPackage = require(path.join(sourceRepository, 'package.json'));

let transcript = '';
const originalWrite = process.stdout.write.bind(process.stdout);
process.stdout.write = function captureTestOutput(chunk, ...args) {
  transcript += Buffer.isBuffer(chunk) ? chunk.toString('utf8') : String(chunk);
  if (failuresOnly) return true;
  return originalWrite(chunk, ...args);
};

const script = suite === 'aggregate'
  ? path.join(sourceTestRoot, 'test.js')
  : path.join(sourceTestRoot, 'modules', 'powSqrt.js');

if (suite === 'pow-sqrt') {
  // powSqrt.js is not in test/test.js. Its loop condition refers to a global
  // `total`, while test/setup.js keeps its assertion counter private. Supply
  // the standalone driver's missing counter and advance it only after each
  // original assertion has run, so all 10,000 assertions execute unchanged.
  global.total = 0;
  require(path.join(sourceTestRoot, 'setup.js'));
  const assertEqual = global.T.assertEqual;
  global.T.assertEqual = function countPowSqrtAssertion(...args) {
    const outcome = assertEqual.apply(this, args);
    global.total += 1;
    return outcome;
  };
}

require(script);

process.stdout.write = originalWrite;

const matches = [...transcript.matchAll(/(?:In total,\s*)?(\d+) of (\d+) tests passed/g)];
if (matches.length === 0) {
  throw new Error(`Could not find a test summary while running ${script}`);
}

const lastMatch = matches[matches.length - 1];
const passed = Number(lastMatch[1]);
const total = Number(lastMatch[2]);
const result = {
  suite,
  source_script: path.relative(sourceRepository, script).replace(/\\/g, '/'),
  module_files: suite === 'aggregate' ? 60 : 1,
  reference_library: `decimal.js ${decimalPackage.version}`,
  reference_repository: sourceRepository,
  passed,
  total,
  status: passed === total ? 'passed' : 'failed',
};

fs.writeFileSync(path.resolve(resultPath), `${JSON.stringify(result, null, 2)}\n`);
if (failuresOnly) {
  const failures = transcript.match(
    /\n\s*Test number \d+ failed:[\s\S]*?(?=\n\s*Test number \d+ failed:|\s+Testing [^\n]*\.\.\.|\n\s*In total,|$)/g,
  ) || [];
  for (const failure of failures) originalWrite(`${failure.trimEnd()}\n`);
  const overallSummary = transcript.match(/\n\s*In total,\s*\d+ of \d+ tests passed[^\n]*/g);
  if (overallSummary && overallSummary.length > 0) {
    originalWrite(`${overallSummary[overallSummary.length - 1].trim()}\n`);
  } else {
    originalWrite(`${passed} of ${total} tests passed\n`);
  }
}
console.log(`BENCH_SUITE_RESULT=${JSON.stringify(result)}`);
if (result.status !== 'passed') process.exitCode = 1;
