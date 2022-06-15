#!/usr/bin/make -f

build:
ifeq ($(OS),Windows_NT)
	go build  -o build/teleport-data-analytics.exe .
else
	go build  -o build/teleport-data-analytics .
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

install:
	go build -o teleport-data-analytics && mv teleport-data-analytics $(GOPATH)/bin

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs goimports -w -local github.com/teleport-network/teleport-data-analytics


setup: build-linux
	@docker build -t teleport-data-analytics .
	@rm -rf ./build
