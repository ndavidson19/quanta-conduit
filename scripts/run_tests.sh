#!/bin/bash

set -e

# Make sure everything builds
go build ./...

# Run all tests
go test ./... -v -coverprofile=coverage.out

# Display coverage report
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

echo "Tests completed. Coverage report generated in coverage.html"

# Open the coverage report in the default browser (optional)
if [[ "$OSTYPE" == "darwin"* ]]; then
    open coverage.html
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    xdg-open coverage.html
fi
