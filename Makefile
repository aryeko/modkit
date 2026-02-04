SHELL := /bin/sh

.PHONY: fmt lint vuln test

fmt:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run

vuln:
	govulncheck ./...

test:
	go test ./...
