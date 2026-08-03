# Docker environment for decimal-go-Port with Go 1.22 and Node.js (for JS oracle fuzzing/benchmarks)
FROM golang:1.22-alpine

# Install Node.js and npm required for differential fuzzing oracle (decimal.js)
RUN apk add --no-cache nodejs npm bash

WORKDIR /app

# Cache Go modules
COPY go.mod ./
RUN go mod download

# Cache Node modules for fuzzing oracle
COPY fuzz/package*.json ./fuzz/
RUN cd fuzz && npm install --production

# Copy codebase
COPY src/ ./src/
COPY tests/ ./tests/
COPY fuzz/ ./fuzz/
COPY bench/ ./bench/
COPY Makefile LICENSE README.md DECISIONS.md .port-mortem.toml ./

# Build source to verify compilation
RUN go build ./src/...

# Default command: Run Go test suite
CMD ["go", "test", "./tests/port/...", "-v"]

