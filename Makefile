.PHONY: build run

build:
	@go build -o bin/api/main ./cmd/api/

run: build
	@./bin/api/main