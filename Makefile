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

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

# NOTE: Placeholder for buidling linux-binaries
# linux-binaries:
# 	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -