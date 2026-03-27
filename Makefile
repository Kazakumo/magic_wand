.PHONY: run tui test lint up down build

BINARY := magic_wand
CMD     := ./cmd/magic_wand

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

tui: run

test:
	go test ./...

lint:
	golangci-lint run ./...

up:
	docker compose up -d

down:
	docker compose down
