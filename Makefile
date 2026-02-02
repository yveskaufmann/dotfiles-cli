.SHELL := /usr/bin/env bash
 
BINARY_NAME := dotfiles

.PHONY: build 
build:
	go build -o bin/$(BINARY_NAME)-v2 ./cmd/$(BINARY_NAME)

.PHONY: run 
 run: 
	go run ./cmd/$(BINARY_NAME)/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf bin/$(BINARY_NAME)-v2