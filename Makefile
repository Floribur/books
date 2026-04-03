DB_URL ?= postgres://postgres:postgres@localhost:5432/floslib?sslmode=disable

# Load .env if it exists
ifneq (,$(wildcard .env))
  include .env
  export
endif

.PHONY: dev migrate migrate-down sqlc build

dev:
	docker compose up -d
	go run ./cmd/server

migrate:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

sqlc:
	sqlc generate

build:
	go build -o flos-library ./cmd/server
