#!/bin/bash

# Run tests with coverage
echo "Running tests with coverage..."
go test -v -coverprofile=coverage.out ./...

# Display coverage summary
echo "Coverage Summary:"
go tool cover -func=coverage.out

# Generate HTML report
echo "Generating HTML report..."
go tool cover -html=coverage.out -o coverage.html

echo "Done. Open coverage.html to view implementation details."
