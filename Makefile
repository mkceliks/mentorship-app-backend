LINT_TIMEOUT = 10m

.PHONY: install-linter
install-linter:
	@echo "Checking for golangci-lint installation..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

.PHONY: lint
lint: install-linter
	@echo "Running golangci-lint with timeout $(LINT_TIMEOUT)..."
	golangci-lint run --timeout $(LINT_TIMEOUT)
