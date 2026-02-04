# Repository Guidelines

This file provides short, focused guidance for contributors and AI agents. Keep instructions concise and non-conflicting. For path-specific guidance, prefer scoped instruction files rather than growing this document.

## Project Structure
- Core library packages: `modkit/` (`module`, `kernel`, `http`, `logging`).
- Example apps: `examples/` (see `examples/hello-mysql/README.md`).
- Documentation: `docs/guides/` (user guides), `docs/reference/` (API reference).
- Internal plans: `.github/internal/plans/` (roadmap and implementation tracking).

## Tooling & Commands
- Format: `make fmt` (runs `gofmt`, `goimports`).
- Lint: `make lint` (runs `golangci-lint`).
- Vulnerability scan: `make vuln` (runs `govulncheck`).
- Tests: `make test` and `go test ./examples/hello-mysql/...`.

## Coding Conventions
- Use `gofmt` formatting and standard Go naming.
- Packages are lowercase, short, and stable.
- Keep exported API minimal; prefer explicit errors over panics.

## Testing Guidance
- Use Go’s `testing` package and keep tests close to code.
- Name tests `TestXxx` and use table-driven tests where it clarifies cases.
- Integration tests should be deterministic; keep external dependencies isolated.

## Commit & PR Hygiene
- Use conventional prefixes: `feat:`, `fix:`, `docs:`, `chore:`.
- One logical change per commit.
- PRs should include summary + validation commands run.

## Pull Request Requirements
When creating PRs, follow `.github/pull_request_template.md` exactly:

1. **Type section**: Check ALL types that apply (a PR can be `feat` + `docs` + `chore`).
2. **Validation section**: Run ALL commands (`make fmt && make lint && make test`) and paste results.
3. **Checklist section**: Verify EVERY item before submitting. All boxes must be checked.
4. **Breaking Changes**: If controller keys, function signatures, or public API changes, document it.

Before submitting:
```bash
make fmt      # Format code
make lint     # Run linter (must pass)
make test     # Run tests (must pass)
```

If any command fails, fix the issue before creating the PR.

## Agent Instruction Layout
- Agent instructions can be stored in `AGENTS.md` files; the closest `AGENTS.md` in the directory tree takes precedence.
- Keep instructions scoped and avoid conflicts across files.
