#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H') # only the hash

# --------------------------------------------------------------------------------------------------
# Test
# --------------------------------------------------------------------------------------------------
test:
	@echo "🐀 Running tests..."
	go test ./x/...
	@echo "✅ Tests completed!"

# --------------------------------------------------------------------------------------------------
# Linting
# --------------------------------------------------------------------------------------------------
lint: 
	@echo "🐀 Running linter..."
	golangci-lint run --timeout=10m
	@echo "🪄 Linting completed!"	