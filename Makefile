SHELL=/bin/bash

.PHONY: all
all: deps lint build test

.PHONY: deps
deps:
	go mod download

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -o artifacts/svc .

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run
