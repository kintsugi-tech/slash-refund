#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H') # only the hash

test:
	@echo "🐀 Running tests..."
	go test ./x/...
	@echo "✅ Tests completed!"