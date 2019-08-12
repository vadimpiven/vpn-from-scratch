.PHONY: setup
setup: ## Install dev tools
	@GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	@go get -u golang.org/x/tools/cmd/goimports

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

.PHONY: mysqlup
mysqlup: ## Start mysql for tests
	@docker container run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=hw6_db_explorer -d --name mysql-hw6_db_explorer mysql

.PHONY: mysqldown
mysqldown: ## Stop mysql after tests
	@docker container stop mysql-hw6_db_explorer && docker container rm mysql-hw6_db_explorer

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
default: help
