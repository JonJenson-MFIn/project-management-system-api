# Project Management System API Tests

This directory contains organized tests for the Project Management System API.

## Test Structure

```
tests/
├── README.md                 # This file
├── tdd/                     # Test-Driven Development tests
│   ├── tdd.go              # Auth directive tests (2 tests)
│   ├── resolver_tests.go   # Resolver functionality tests (3 tests)
│   └── middleware_tests.go # Middleware integration tests (2 tests)
└── benchmark/               # Performance benchmarks
    ├── benchmark.go         # Auth directive benchmarks (1 benchmark)
    ├── resolver_benchmarks.go # Resolver performance (2 benchmarks)
    └── api_benchmarks.go    # Full API endpoint benchmarks (1 benchmark)
```

**Total: 7 TDD tests + 4 benchmarks**

## Running Tests

### Run All TDD Tests
```bash
go test ./tests/tdd/...
```

### Run All Benchmarks
```bash
go test -bench=. ./tests/benchmark/...
```

### Run Specific Test Categories

#### TDD Tests
```bash
# All TDD tests
go test ./tests/tdd/...

# Specific TDD test file
go test ./tests/tdd/tdd.go
go test ./tests/tdd/resolver_tests.go
go test ./tests/tdd/middleware_tests.go
```

#### Benchmark Tests
```bash
# All benchmarks
go test -bench=. ./tests/benchmark/...

# Specific benchmark file
go test -bench=. ./tests/benchmark/benchmark.go
go test -bench=. ./tests/benchmark/resolver_benchmarks.go
go test -bench=. ./tests/benchmark/api_benchmarks.go
```

### Test Options

#### Verbose Output
```bash
go test -v ./tests/tdd/...
go test -v -bench=. ./tests/benchmark/...
```

#### Run Tests with Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./tests/tdd/...

# View coverage in browser
go tool cover -html=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out
```

#### Run Benchmarks with Memory Info
```bash
go test -bench=. -benchmem ./tests/benchmark/...
```

#### Run Benchmarks for Specific Duration
```bash
go test -bench=. -benchtime=10s ./tests/benchmark/...
```

### Quick Commands

```bash
# Run everything
go test ./tests/...

# Run TDD tests with coverage
go test -coverprofile=coverage.out ./tests/tdd/...

# Run benchmarks with memory allocation info
go test -bench=. -benchmem ./tests/benchmark/...
```
