GO := /usr/local/go/bin/go

build:
	mkdir -p ./bin
	$(GO) build -o ./bin/ ./cmd/munit

.PHONY: build
