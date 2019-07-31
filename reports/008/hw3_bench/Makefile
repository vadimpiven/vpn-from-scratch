.PHONY: setup
setup: ## Install dev tools
	@GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/psampaz/go-mod-outdated

.PHONY: tidy
tidy: ## Prune any no-longer-needed dependencies
	@go mod tidy

.PHONY: fmt
fmt: tidy ## Run gofmt and goimports
	@find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

.PHONY: lint
lint: fmt ## Check for outdated dependencies and run all the linters
	go list -u -m -json all | go-mod-outdated -update -direct; golangci-lint run --enable-all

.PHONY: test
test: ## Run all the tests
	go test -v

.PHONY: bench
bench: ## Run all the benchmarks
	go test -bench . -benchmem

.PHONY: profile
profile: clean ## Writes CPU and memory usage profiles for BenchmarkFast
	@GOGC=off go test -bench BenchmarkFast -run=^$ . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1

.PHONY: profcpu
profcpu: ## Visualise CPU profile
	pprof -http=:8080 cpu.out

.PHONY: profmem
profmem: ## Visualise memory profile
	pprof -http=:8080 mem.out

.PHONY: clean
clean: ## Remove CPU and memory usage profiles
	@rm -f cpu.out mem.out 2> /dev/null

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
default: help
