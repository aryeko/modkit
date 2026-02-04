SHELL := /bin/sh

.PHONY: run test migrate compose-up compose-down sqlc

compose-up:
	docker compose up -d

compose-down:
	docker compose down -v

migrate:
	go run ./cmd/migrate

run: compose-up migrate
	go run ./cmd/api

sqlc:
	sqlc generate

test:
	go test ./...
