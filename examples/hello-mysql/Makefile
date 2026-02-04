SHELL := /bin/sh

.PHONY: run test migrate compose-up compose-up-db compose-up-app compose-down sqlc

TEST_PACKAGES ?= ./...

compose-up:
	docker compose up -d mysql app

compose-up-db:
	docker compose up -d mysql

compose-up-app:
	docker compose up -d app

compose-down:
	docker compose down -v

migrate:
	go run ./cmd/migrate

run: compose-up-db migrate compose-up-app

sqlc:
	sqlc generate

test:
	go test $(TEST_PACKAGES)
