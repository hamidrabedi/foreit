# Linting and Code Quality

forge uses multiple linting and static analysis tools to ensure code quality, security, and maintainability.

## Tools Used

### 1. golangci-lint

[golangci-lint](https://golangci-lint.run/) is a fast Go linters runner that aggregates multiple linters and runs them in parallel.

**Installation:**

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Configuration:**
The configuration is in `.golangci.yml` at the project root. We focus on **critical issues** (errors, security, bugs) rather than style preferences.

**Enabled Linters:**

- `govet` - Official Go linter
- `errcheck` - Checks for unchecked errors
- `staticcheck` - Advanced static analysis
- `unused` - Finds unused code
- `gosimple` - Simplifies code
- `gosec` - Security-oriented linter
- `gocritic` - Code quality checks (diagnostic, performance, not style)
- `gofmt` / `goimports` - Code formatting
- `misspell` - Spell checking
- And more...

**Usage:**

```bash
# Run linters
make lint

# Run with auto-fix
make lint-fix

# Or directly
golangci-lint run ./internal/...
```

### 2. go vet

`go vet` is Go's built-in static analysis tool.

**Usage:**

```bash
make vet
# Or
go vet ./internal/...
```

### 3. gosec

[gosec](https://github.com/securego/gosec) is a security-focused linter that scans for security vulnerabilities.

**Installation:**

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

**Usage:**

```bash
make sec
# Or
gosec ./internal/...
```

### 4. staticcheck

[staticcheck](https://staticcheck.io/) provides advanced static analysis beyond what `go vet` offers.

**Installation:**

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

**Usage:**

```bash
make staticcheck
# Or
staticcheck ./internal/...
```

## Makefile Commands

The Makefile provides convenient commands for running all linters:

```bash
# Install all linting tools
make install-tools

# Run golangci-lint
make lint

# Run golangci-lint with auto-fix
make lint-fix

# Run go vet
make vet

# Run gosec security scanner
make sec

# Run staticcheck
make staticcheck

# Run all checks (lint, vet, sec)
make check

# Run all checks including staticcheck
make check-all
```

## Configuration Philosophy

We focus on **critical issues** that affect:

- **Correctness** - Bugs, logic errors
- **Security** - Vulnerabilities, unsafe patterns
- **Reliability** - Error handling, resource leaks

We **ignore style preferences** like:

- Parameter type combinations
- Switch vs if statements
- Range value copying (for small structs)
- Named vs unnamed return values
- Builtin shadowing (when clear in context)

## Excluded Directories

The following directories are excluded from linting:

- `vendor/` - Dependencies
- `examples/` - Example code
- `cmd/forge/` - CLI tool code

## Excluded Files

- `*.gen.go` - Generated files
- `*_test.go` - Test files (for some linters)

## Common Issues and Solutions

### Unchecked Errors

If you see `errcheck` warnings about unchecked errors:

```go
// Bad
result, _ := someFunction()

// Good - handle the error
result, err := someFunction()
if err != nil {
    return err
}

// Or if error can't be handled meaningfully (e.g., HTTP responses)
// nolint:errcheck // HTTP response errors can't be handled meaningfully
_ = w.Write(data)
```

### Security Warnings (gosec)

For file permissions, if you need 0644 for config files:

```go
// nolint:gosec // G306: Config files need 0644 permissions for readability
if err := os.WriteFile(path, data, 0644); err != nil {
    return err
}
```

### Unused Code

If you have intentionally unused code (e.g., reserved for future use):

```go
// nolint:unused // Reserved for future middleware management
middleware []func(http.Handler) http.Handler
```

### Type Assertions

Always use the two-value form for type assertions:

```go
// Bad
value := x.(string)

// Good
value, ok := x.(string)
if !ok {
    // handle error
}
```

## CI/CD Integration

Linting should be integrated into your CI/CD pipeline. Example GitHub Actions workflow:

```yaml
name: Lint

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.25
      - name: Install golangci-lint
        run: |
          curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
      - name: Run golangci-lint
        run: |
          ./bin/golangci-lint run ./internal/...
      - name: Run go vet
        run: go vet ./internal/...
      - name: Run gosec
        run: |
          go install github.com/securego/gosec/v2/cmd/gosec@latest
          gosec ./internal/...
```

## Best Practices

1. **Run linters before committing** - Use `make check` before pushing code
2. **Fix auto-fixable issues** - Use `make lint-fix` to automatically fix formatting issues
3. **Don't ignore important warnings** - Address security and error handling issues
4. **Use nolint comments sparingly** - Only when you have a good reason
5. **Focus on critical issues** - Don't obsess over style preferences

## Resources

- [golangci-lint Documentation](https://golangci-lint.run/)
- [gosec Documentation](https://github.com/securego/gosec)
- [staticcheck Documentation](https://staticcheck.io/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
