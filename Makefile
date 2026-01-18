.SUFFIXES:

GO    ?= go
GIT   ?= git
REUSE ?= reuse

all: pre-commit

pre-commit: tidy fmt lint test/unit test/security clean # Runs all pre-commit checks.

commit: pre-commit # Commits the changes to the repository.
	$(GIT) commit -s

tidy: # Updates the go.mod file to ensure it matches the source code.
	$(GO) mod tidy

fmt: # Formats Go source files in this repository.
	$(GO) run 'github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest' fmt ./...

lint: # Runs golangci-lint on Go files.
	$(GO) run 'github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest' run ./...

test: test/unit test/security # Runs all tests.

test/unit: # Runs unit tests.
	$(GO) test -cover -race -vet all -mod readonly -short ./...

test/unit/coverage: # Runs unit tests with coverage information.
	$(GO) test -cover -coverprofile cover.out -race -vet all -mod readonly -short ./...
	$(GO) tool cover -html=cover.out

test/security: # Analyzes the codebase and looks for vulnerabilities affecting it.
	# https://github.com/golang/go/issues/73871
	GOEXPERIMENT= $(GO) run 'golang.org/x/vuln/cmd/govulncheck@latest' ./...

clean: # Cleans cache files from tests and deletes any build output.
	$(GO) clean -cache -fuzzcache -testcache

.PHONY: all pre-commit commit tidy fmt lint test test/unit \
	test/unit/coverage test/security clean
