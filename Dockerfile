# Multi-stage Docker build for decimal-go-Port
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy Go module definition
COPY go.mod ./

# Copy all source and test directories
COPY src/ ./src/
COPY tests/ ./tests/
COPY fuzz/ ./fuzz/
COPY bench/ ./bench/
COPY LICENSE README.md DECISIONS.md .port-mortem.toml ./

# Build package and run tests
RUN go build ./src/...
RUN go test ./tests/port/... -v

CMD ["go", "test", "./tests/port/...", "-v"]
