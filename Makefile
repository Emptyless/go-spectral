MAKEFLAGS := --no-print-directory --silent

default: help

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z\._-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.PHONY: help

install: ### dependencies using npm
	npm ci

build: ### dist using webpack
	npm run webpack

fmt: ## Format go code & tidy the go.mod file
	go mod tidy
	@gofumpt -l -w .
	@golangci-lint run --fix
.PHONY: fmt

t: test
test: ## Run unit tests, alias: t
	go test ./...
.PHONY: test

ci: fmt test ## simulate pipeline checks
.PHONY: ci

tools: ## Install extra tools for development
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: tools

