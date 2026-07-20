.DEFAULT_GOAL := help

# Application
application_name := topsigner

# Directories
target      := target
bin         := $(target)/bin
reports     := $(target)/reports

## build: Build the application.
.PHONY: build
build: $(bin)
	go build -o $(bin)/$(application_name) ./cmd/$(application_name)

## run: Run the application.
.PHONY: run
run:
	@go run ./cmd/$(application_name)

## deps: Install and verify dependencies.
.PHONY: deps
deps:
	@go mod download
	@go mod verify

## fmt: Format Go code.
.PHONY: fmt
fmt:
	@go fmt ./...

## test: Run tests.
.PHONY: test
test:
	@go test -v ./...	

## test-coverage: Run tests with coverage report.
.PHONY: test-coverage
test-coverage: | $(reports)
	@go test -v -covermode=atomic -coverprofile=$(reports)/cover.out ./...
	@go tool cover -html=$(reports)/cover.out -o=$(reports)/cover.html

## lint: Run linter.
.PHONY: lint
lint:
	@golangci-lint run

## lint/fix: Run linter and fix issues.
.PHONY: lint/fix
lint/fix:
	@golangci-lint run --fix

## help: Display available targets.
.PHONY: help
help: Makefile
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@sed -n 's/^## //p' $< | awk -F ':' '{printf "  %-20s%s\n",$$1,$$2}'
