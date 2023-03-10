#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H') # only the hash

# --------------------------------------------------------------------------------------------------
# Build & Installation
# --------------------------------------------------------------------------------------------------
build:
	@echo "🐀 Buidling..."
	go build $(BUILD_FLAGS) -o bin/slash-refundd ./cmd/slash-refundd
	@echo "🛠 Build completed!"

install:
	@echo "🐀 Installing..."
	go install -mod=readonly ./cmd/slash-refundd
	@echo "🎉 Installation completed!"
	
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

# --------------------------------------------------------------------------------------------------
# All 
# --------------------------------------------------------------------------------------------------
all: test install