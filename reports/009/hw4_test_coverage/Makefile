SOURCE_FILES ?= ./...
TEST_PATTERN ?= .
TEST_OPTS ?=

.PHONY: setup
setup: ## Install dev tools
	@GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	@go get -u golang.org/x/tools/cmd/goimports

.PHONY: test
test: ## Run all the tests
	@go test $(TEST_OPTS) -v -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	@go tool cover -html=coverage.out

.PHONY: tidy
tidy: ## Prune any no-longer-needed dependencies
	@go mod tidy

.PHONY: fmt
fmt: tidy ## Run gofmt and goimports on all go files
	@find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

.PHONY: lint
lint: fmt ## Run all the linters
	golangci-lint run --enable-all

.PHONY: clean
clean: ## Remove all generated files
	@rm -f coverage.out
	@go clean

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
default: help
