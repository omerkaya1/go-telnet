BUILD= $(CURDIR)/bin
$(shell mkdir -p $(BUILD))
export GO111MODULE=on
export GOPATH=$(go env GOPATH)

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get -u golang.org/x/tools
	go get -u golang.org/x/lint/golint

.PHONY: mod
mod: ## Runs mod
	go mod verify
	go mod vendor
	go mod tidy

.PHONY: test
test: setup ## Runs all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s ./...

.PHONY: cover
cover: test ## Runs all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: fmt
fmt: setup ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: lint
lint: setup ## Runs all the linters
	golint ./...

.PHONY: build
build: ## Builds the project
	go build -o $(BUILD)/go-telnet $(CURDIR)/cmd/go-telnet

.PHONY: clean
clean: ## Remove temporary files
	go clean $(CURDIR)/cmd/go-telnet
	rm -rf $(BUILD)
	rm -rf coverage.txt

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
