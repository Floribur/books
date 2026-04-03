DB_URL ?= postgres://postgres:postgres@localhost:5432/floslib?sslmode=disable

# Load .env if it exists
ifneq (,$(wildcard .env))
  include .env
  export
endif

.PHONY: dev migrate migrate-down sqlc build build-pi

dev:
	docker compose up -d
	migrate -path ./migrations -database "$(DB_URL)" up
	go run ./cmd/server

migrate:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

sqlc:
	sqlc generate

# Production build: React frontend first (populates frontend/dist/ for go:embed), then Go binary.
# CRITICAL ORDER: npm run build MUST precede go build — embed picks up whatever is in frontend/dist/
# at compile time. Running go build first embeds stale/empty dist.
build:
	cd frontend && npm run build
	go build -o flos-library ./cmd/server

# Cross-compile for Raspberry Pi 4/5 (ARMv8 / arm64).
# CGO_ENABLED=0 is REQUIRED — cgo cannot cross-compile on Windows without a cross-toolchain.
# Output: flos-library-linux-arm64 — ready to scp to Pi.
build-pi:
	cd frontend && npm run build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o flos-library-linux-arm64 ./cmd/server
