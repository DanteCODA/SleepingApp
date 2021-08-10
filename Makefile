.PHONY: all
all: build
FORCE: ;

SHELL  := env LIBRARY_ENV=$(LIBRARY_ENV) $(SHELL)
LIBRARY_ENV ?= dev

BIN_DIR = $(PWD)/bin

.PHONY: build

clean:
	rm -rf bin/*

dependencies:
	go mod download

build: dependencies build-api

build-api: 
	GOARCH=amd64 GOOS=linux go build -tags $(LIBRARY_ENV) -o ./bin/lambda/main api/lambda/main.go

build-cmd:
	go build -tags $(LIBRARY_ENV) -o ./bin/cmd/main cmd/main.go

ci: dependencies test	

test:
	go test -tags testing ./...

fmt: ## gofmt and g