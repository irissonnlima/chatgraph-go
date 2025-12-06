#!/bin/bash

# Coverage script for chatgraph-go
# Generates coverage report in coverage/ directory

set -e

COVERAGE_DIR="coverage"

# Create coverage directory
mkdir -p "$COVERAGE_DIR"

echo "Running tests with coverage..."
go test -coverprofile="$COVERAGE_DIR/coverage.out" ./... || true

echo "Generating HTML report..."
go tool cover -html="$COVERAGE_DIR/coverage.out" -o "$COVERAGE_DIR/coverage.html"

echo ""
echo "✅ Coverage report generated:"
echo "   - $COVERAGE_DIR/coverage.out"
echo "   - $COVERAGE_DIR/coverage.html"
echo ""

# Show summary
echo "Coverage summary:"
go tool cover -func="$COVERAGE_DIR/coverage.out" | tail -1
