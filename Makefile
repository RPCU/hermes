.PHONY: build test lint clean help

BINARY_NAME=hermes

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o $(BINARY_NAME) .

test: ## Run tests
	go test -v ./...

lint: ## Run linter
	golangci-lint run

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)
	rm -rf dist/

run-dry: ## Run in dry-run mode
	go run . --dry-run --verbose
