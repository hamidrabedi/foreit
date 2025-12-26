.PHONY: build test clean generate migrate runserver lint lint-fix vet sec staticcheck check check-all install-tools help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-tools: ## Install linting and analysis tools
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing gosec..."
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Installing staticcheck..."
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "Tools installed successfully!"

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@golangci-lint run ./internal/...

lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	@golangci-lint run --fix ./internal/...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./internal/...

sec: ## Run gosec security scanner
	@echo "Running gosec security scanner..."
	@gosec ./internal/...

staticcheck: ## Run staticcheck
	@echo "Running staticcheck..."
	@staticcheck ./internal/...

check: lint vet sec ## Run all checks (lint, vet, sec)

check-all: lint vet sec staticcheck ## Run all checks including staticcheck

build:
	go build -o bin/forge ./cmd/forge

test:
	go test ./...

clean:
	rm -rf bin/
	go clean

generate:
	go run ./cmd/forge generate

migrate:
	go run ./cmd/forge migrate

runserver:
	go run ./cmd/forge runserver

