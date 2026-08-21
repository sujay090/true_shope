.PHONY: build run

build:
	@go build -o bin/api/main ./cmd/api/

run: build
	@./bin/api/main

migrate-up:
	@go run ./api/migrate up

migrate-down:
	@go run ./api/migrate down