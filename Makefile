# Variables for linter configuration
LINT_TIMEOUT = 10m

# Install golangci-lint if not already available
.PHONY: install-linter
install-linter:
	@echo "Checking for golangci-lint installation..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

# Run golangci-lint with timeout
.PHONY: lint
lint: install-linter
	@echo "Running golangci-lint with timeout $(LINT_TIMEOUT)..."
	golangci-lint run --timeout $(LINT_TIMEOUT)
