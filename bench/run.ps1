[CmdletBinding()]
param(
    [ValidateRange(1, 1000000000)]
    [int]$Ops = 100000
)

$ErrorActionPreference = 'Stop'

$benchDir = $PSScriptRoot
$repoRoot = Split-Path -Parent $benchDir
# Set this to the absolute path of the original decimal.js repository clone.
# The Go repository path is derived from this script's location.
$jsRepository = 'C:\Users\Lenovo\projects\port_mortem\decimal.js'
$jsDecimal = Join-Path $jsRepository 'decimal.js'
$jsTests = Join-Path $jsRepository 'test'
$scriptsDir = Join-Path $benchDir 'scripts'
$tempDir = Join-Path $benchDir '.tmp'
$goBinary = Join-Path $benchDir 'bin\decimal-go-bench.exe'
$collector = Join-Path $scriptsDir 'collect.js'
$results = Join-Path $benchDir 'results.json'
$jsSuiteRunner = Join-Path $scriptsDir 'run_js_suite.js'
$suiteVerifier = Join-Path $scriptsDir 'create_suite_verification.js'
$aggregateResult = Join-Path $tempDir 'aggregate_js_suite.json'
$powSqrtResult = Join-Path $tempDir 'pow_sqrt_js_suite.json'
$suiteVerification = Join-Path $tempDir 'suite_verification.json'

if (-not (Get-Command node -ErrorAction SilentlyContinue)) {
    throw 'Node.js is required. Install Node.js, then run this script again.'
}
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw 'Go is required. Install Go, then run this script again.'
}

if ($Clean) {
    Write-Host '==> Clearing the prior benchmark executable, result file, and Go build cache...'
    if (Test-Path -LiteralPath $goBinary) {
        Remove-Item -LiteralPath $goBinary -Force
    }
    if (Test-Path -LiteralPath $results) {
        Remove-Item -LiteralPath $results -Force
    }
    & go clean -cache
    if ($LASTEXITCODE -ne 0) {
        throw "go clean failed with exit code $LASTEXITCODE"
    }
}

foreach ($transientFile in @($aggregateResult, $powSqrtResult, $suiteVerification)) {
    if (Test-Path -LiteralPath $transientFile) {
        Remove-Item -LiteralPath $transientFile -Force
    }
}

if (-not (Test-Path -LiteralPath $jsDecimal) -or -not (Test-Path -LiteralPath $jsTests)) {
    throw "Original decimal.js repository is required at $jsRepository"
}

New-Item -ItemType Directory -Force -Path (Split-Path -Parent $goBinary) | Out-Null
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

Push-Location $repoRoot
try {
    Write-Host '==> Verifying all 61 original JavaScript test-module files (60 aggregate + standalone powSqrt)...'
    & node $jsSuiteRunner --js-repo $jsRepository --suite aggregate --failures-only --result $aggregateResult
    if ($LASTEXITCODE -ne 0) {
        throw "Original JavaScript aggregate suite failed with exit code $LASTEXITCODE"
    }

    & node $jsSuiteRunner --js-repo $jsRepository --suite pow-sqrt --failures-only --result $powSqrtResult
    if ($LASTEXITCODE -ne 0) {
        throw "Original JavaScript powSqrt suite failed with exit code $LASTEXITCODE"
    }

    Write-Host '==> Verifying all current decimal-go port tests...'
    & go test -count=1 ./tests/port
    if ($LASTEXITCODE -ne 0) {
        throw "Go port test suite failed with exit code $LASTEXITCODE"
    }

    & node $suiteVerifier --aggregate $aggregateResult --pow-sqrt $powSqrtResult --output $suiteVerification
    if ($LASTEXITCODE -ne 0) {
        throw "Could not write suite-verification record (exit code $LASTEXITCODE)"
    }

    Write-Host "==> Measuring decimal.js and decimal-go ($Ops mixed operations each)..."
    & go build -o $goBinary '.\bench\cmd\go_bench'
    if ($LASTEXITCODE -ne 0) {
        throw "go build failed with exit code $LASTEXITCODE"
    }

    & node $collector --js-repo $jsRepository --go-bin $goBinary --ops $Ops --output $results --suite-verification $suiteVerification
    if ($LASTEXITCODE -ne 0) {
        throw "benchmark collector failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}
