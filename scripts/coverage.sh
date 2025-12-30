#!/bin/bash
# Generate test coverage report for forge framework

set -e

echo "Running tests with coverage..."
cd forge

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate HTML report
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Generate text summary
echo "Coverage summary:"
go tool cover -func=coverage.out | tail -1

echo ""
echo "Coverage report generated: forge/coverage.html"
echo "Open it in your browser to view detailed coverage."
