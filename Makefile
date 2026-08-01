.PHONY: all build test bench fuzz docker-build clean

all: build test

build:
	@echo "==> Building decimal-go package..."
	go build ./src/...

test:
	@echo "==> Running Go test suites in tests/port/..."
	go test ./tests/port/... -v

bench:
	@echo "==> Running benchmarks..."
	go test ./src/... -bench=. -benchmem

fuzz:
	@echo "==> Running differential fuzzing harness (60s)..."
	go test ./fuzz/... -fuzz=FuzzDecimal -fuzztime=60s

docker-build:
	@echo "==> Building Docker image decimal-go-port..."
	docker build -t decimal-go-port .

clean:
	go clean
