help:
	@echo "Makefile for the testing tool"
	@echo "Available targets:"
	@echo " help - print help information"
	@echo " install - install required dependecies for the project"
	@echo " lint - run linter"
	@echo " test - run unit tests"
	@echo " build - build app binary"

install:
	echo "Install linter"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.42.0

lint:
	@echo "Run linter"
	# Lint golang
	golangci-lint run

test:
	@echo "Run unit tests"
	go test $$(go list ./... | grep -v gen) -timeout 30s --race -v -coverprofile coverage.txt -covermode atomic
	@echo "Code coverage"
	go tool cover -func coverage.txt

build:
	@echo "Build app binary"
	go build -ldflags "-s -w" -o out/server

run:
	LOG_LEVEL=debug go run .