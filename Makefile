BUILDDIR=./build
REPLTARGET=pbrepl

.PHONY: build test build-repl deps

build-repl:
	go build -o $(BUILDDIR)/$(REPLTARGET) ./cmd/bundler-repl

build: build-repl

deps:
	go mod download

test:
	go test -v ./...
