# modkit MVP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement the full modkit MVP (core packages, HTTP adapter, example app, and docs/CI) in ordered phases with per-phase validation and commits.

**Architecture:** Follow `docs/implementation/master.md` and phase docs in order; shared architecture and semantics come from `modkit_mvp_design_doc.md`. Each phase is validated independently and committed before moving on.

**Tech Stack:** Go, `net/http`, `github.com/go-chi/chi/v5`, MySQL + `database/sql`, sqlc, Make.

### Task 1: Phase 00 — Repo Bootstrap (Scaffolding)

**Files:**
- Modify: `docs/implementation/phase-00-repo-bootstrap.md`
- Create: `go.mod`
- Create: `README.md`
- Create: `LICENSE`
- Create: `CONTRIBUTING.md`
- Create: `CODE_OF_CONDUCT.md`
- Create: `SECURITY.md`
- Create: `docs/design/mvp.md`
- Modify: `modkit_mvp_design_doc.md` (turn into pointer or remove)
- Create: `.github/workflows/ci.yml`

**Step 1: Write the failing test**
- Create a placeholder test to ensure `go test ./...` executes (if no packages exist yet).

**Step 2: Run test to verify it fails**
Run: `go test ./...`
Expected: FAIL due to missing module/packages.

**Step 3: Write minimal implementation**
- Initialize Go module `github.com/aryeko/modkit`.
- Add baseline repository docs and `docs/design/mvp.md` (canonical), update root design doc to be a pointer.
- Add CI workflow for `go test ./...`.

**Step 4: Run test to verify it passes**
Run: `go test ./...`
Expected: PASS (empty or minimal packages).

**Step 5: Commit**
```bash
git add go.mod README.md LICENSE CONTRIBUTING.md CODE_OF_CONDUCT.md SECURITY.md docs/design/mvp.md modkit_mvp_design_doc.md .github/workflows/ci.yml
git commit -m "chore: bootstrap repo structure"
```

### Task 2: Phase 01 — module Package

**Files:**
- Create: `modkit/module/module.go`
- Create: `modkit/module/token.go`
- Create: `modkit/module/provider.go`
- Create: `modkit/module/controller.go`
- Create: `modkit/module/errors.go`
- Create: `modkit/module/module_test.go`

**Step 1: Write the failing test**
- Add compile-only tests asserting exported types and errors are accessible.

**Step 2: Run test to verify it fails**
Run: `go test ./modkit/module/...`
Expected: FAIL (missing package).

**Step 3: Write minimal implementation**
- Implement module metadata types per `modkit_mvp_design_doc.md` Section 5.1.

**Step 4: Run test to verify it passes**
Run: `go test ./modkit/module/...`
Expected: PASS.

**Step 5: Commit**
```bash
git add modkit/module
git commit -m "feat: add module definitions"
```

### Task 3: Phase 02 — Kernel Graph + Container

**Files:**
- Create: `modkit/kernel/bootstrap.go`
- Create: `modkit/kernel/graph.go`
- Create: `modkit/kernel/visibility.go`
- Create: `modkit/kernel/container.go`
- Create: `modkit/kernel/errors.go`
- Create: `modkit/kernel/*_test.go`

**Step 1: Write the failing test**
- Implement kernel unit tests per `modkit_mvp_design_doc.md` Section 8.1.

**Step 2: Run test to verify it fails**
Run: `go test ./modkit/kernel/...`
Expected: FAIL (missing implementation).

**Step 3: Write minimal implementation**
- Implement graph build/validation, container, visibility enforcement, and `Bootstrap`.

**Step 4: Run test to verify it passes**
Run: `go test ./modkit/kernel/...`
Expected: PASS.

**Step 5: Commit**
```bash
git add modkit/kernel
git commit -m "feat: add kernel bootstrap and graph"
```

### Task 4: Phase 03 — HTTP Adapter

**Files:**
- Create: `modkit/http/server.go`
- Create: `modkit/http/router.go`
- Create: `modkit/http/middleware.go`
- Create: `modkit/http/errors.go`
- Create: `modkit/http/*_test.go`

**Step 1: Write the failing test**
- Minimal tests for router creation and Serve helper.

**Step 2: Run test to verify it fails**
Run: `go test ./modkit/http/...`
Expected: FAIL (missing implementation).

**Step 3: Write minimal implementation**
- Implement `NewRouter` and `Serve` using `chi` and `net/http`.

**Step 4: Run test to verify it passes**
Run: `go test ./modkit/http/...`
Expected: PASS.

**Step 5: Commit**
```bash
git add modkit/http
git commit -m "feat: add http adapter"
```

### Task 5: Phase 04 — Example App (hello-mysql)

**Files:**
- Create: `examples/hello-mysql/...` (per design doc layout)
- Create: `examples/hello-mysql/Makefile`
- Create: `examples/hello-mysql/sql/sqlc.yaml`
- Create: `examples/hello-mysql/sql/queries.sql`
- Create: `examples/hello-mysql/migrations/*`

**Step 1: Write the failing test**
- Add initial module/package tests and a smoke test scaffold.

**Step 2: Run test to verify it fails**
Run: `go test ./examples/hello-mysql/...`
Expected: FAIL (missing implementation).

**Step 3: Write minimal implementation**
- Implement modules, controllers, routes, sqlc integration, and Makefile targets.

**Step 4: Run test to verify it passes**
Run: `go test ./examples/hello-mysql/...`
Expected: PASS.

**Step 5: Commit**
```bash
git add examples/hello-mysql
git commit -m "feat: add hello-mysql example app"
```

### Task 6: Phase 05 — Docs + CI Completeness

**Files:**
- Modify: `README.md`
- Create: `docs/guides/getting-started.md`
- Create: `docs/guides/modules.md`
- Create: `docs/guides/testing.md`
- Modify: `.github/workflows/ci.yml`

**Step 1: Write the failing test**
- Not applicable (docs). Ensure validation commands are defined.

**Step 2: Run test to verify it fails**
Run: `go test ./...`
Expected: PASS (baseline). If failing, fix tests/CI config.

**Step 3: Write minimal implementation**
- Complete docs and CI per phase doc requirements.

**Step 4: Run test to verify it passes**
Run: `go test ./...`
Expected: PASS.

**Step 5: Commit**
```bash
git add README.md docs/guides .github/workflows/ci.yml
git commit -m "docs: complete guides and ci"
```
