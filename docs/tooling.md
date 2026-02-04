# Tooling

This repository uses standard Go tools plus common OSS linters.

## Format

```bash
make fmt
```

Runs:
- `gofmt -w .`
- `goimports -w .`

## Lint

```bash
make lint
```

Runs:
- `golangci-lint run`

See `.golangci.yml` for enabled linters and excluded paths.

## Vulnerability Scan

```bash
make vuln
```

Runs:
- `govulncheck ./...`

## Test

```bash
make test
```

Runs:
- `go test ./...`
