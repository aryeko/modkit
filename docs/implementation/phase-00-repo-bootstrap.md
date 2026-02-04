# Phase 00 — Repo Bootstrap

## Assumptions (Initial State)
- GitHub repo exists: `aryeko/modkit`.
- Local repo cloned and `origin` points to `git@github.com-personal:aryeko/modkit.git`.
- `modkit_mvp_design_doc.md` exists at repo root and is committed.

## Requirements
- Initialize Go module at `github.com/aryeko/modkit`.
- Add baseline repository files:
  - `README.md` (MVP summary + quickstart stub)
  - `LICENSE`
  - `CONTRIBUTING.md`
  - `CODE_OF_CONDUCT.md`
  - `SECURITY.md`
- Add `docs/design/mvp.md` copied from `modkit_mvp_design_doc.md`.
- Add CI workflow: `go test ./...`.

## Design
- Follow `modkit_mvp_design_doc.md` Section 4 (Repo structure) as the target layout.
- This phase establishes scaffolding only; no public API implementations yet.

## Validation
Run:
- `go mod tidy`
- `go test ./...`

Expected:
- `go test` succeeds (no packages or only empty packages with passing tests).

## Commit
- One commit after validation, e.g. `chore: bootstrap repo structure`.
