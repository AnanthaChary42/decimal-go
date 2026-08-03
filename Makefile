.PHONY: all build test bench fuzz docker-build clean

BENCH_OPS ?= 100000
JS_REPO ?= ../decimal.js
GOOS := $(shell go env GOOS)

ifeq ($(GOOS),windows)
GO_BENCH_BIN := bench/bin/decimal-go-bench.exe
else
GO_BENCH_BIN := bench/bin/decimal-go-bench
endif

all: build test

build:
	@echo "==> Building decimal-go package..."
	go build ./src/...

test:
	@echo "==> Running Go test suites in tests/port/..."
	go test ./tests/port/... -v

bench:
	@test -f "$(JS_REPO)/decimal.js" && test -d "$(JS_REPO)/test" || (echo "Original decimal.js repository is required at $(JS_REPO)"; exit 1)
	@mkdir -p bench/bin bench/.tmp
	@rm -f bench/.tmp/aggregate_js_suite.json bench/.tmp/pow_sqrt_js_suite.json bench/.tmp/suite_verification.json
	@echo "==> Verifying all 61 original JavaScript test-module files (60 aggregate + standalone powSqrt)..."
	node bench/scripts/run_js_suite.js --js-repo "$(JS_REPO)" --suite aggregate --failures-only --result bench/.tmp/aggregate_js_suite.json
	node bench/scripts/run_js_suite.js --js-repo "$(JS_REPO)" --suite pow-sqrt --failures-only --result bench/.tmp/pow_sqrt_js_suite.json
	@echo "==> Verifying all current decimal-go port tests..."
	go test -count=1 ./tests/port
	node bench/scripts/create_suite_verification.js --aggregate bench/.tmp/aggregate_js_suite.json --pow-sqrt bench/.tmp/pow_sqrt_js_suite.json --output bench/.tmp/suite_verification.json
	@echo "==> Measuring decimal.js and decimal-go ($(BENCH_OPS) mixed operations each)..."
	go build -o $(GO_BENCH_BIN) ./bench/cmd/go_bench
	node bench/scripts/collect.js --js-repo "$(JS_REPO)" --go-bin $(GO_BENCH_BIN) --ops $(BENCH_OPS) --output bench/results.json --suite-verification bench/.tmp/suite_verification.json

fuzz:
	@echo "==> Running fuzzing harness (60s)..."
	go test ./fuzz/... -fuzz=FuzzDecimal -fuzztime=60s

diff-fuzz:
	@echo "==> Running differential fuzzing harness (60s)..."
	go test ./fuzz/... -run TestDifferentialFuzz -timeout 120s -v

docker-build:
	@echo "==> Building Docker image decimal-go-port..."
	docker build -t decimal-go-port .

docker-test:
	@echo "==> Running Go test suite inside Docker..."
	docker run --rm decimal-go-port go test ./tests/port/... -v

docker-diff-fuzz:
	@echo "==> Running differential fuzzing inside Docker..."
	docker run --rm decimal-go-port go test ./fuzz/... -run TestDifferentialFuzz -v -timeout 120s


clean:
	go clean
