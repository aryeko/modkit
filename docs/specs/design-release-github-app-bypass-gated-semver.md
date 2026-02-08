# Design Spec: GitHub App Bypass + CI-Gated Semantic Tagging + GoReleaser Changelog

## Status

- Type: design spec
- Scope: `go-modkit/modkit`
- Goal: standard, CI-gated release flow with automatic changelog and artifact publishing

## Constraints

- Keep existing cache setup in workflows exactly as-is.
- Go version remains pinned to `1.25.7` in CI.
- Public repository: remove `CODECOV_TOKEN` secret usage and workflow references.
- Release must be gated on successful CI for `main`.
- Use GitHub App token for release automation under ruleset bypass.

## Desired Release Architecture

1. `ci` workflow runs on pull requests and pushes to `main`.
2. `release-semantic` workflow runs only after `ci` completes successfully on `main` push.
3. `release-artifacts` workflow runs on pushed tags matching `v*`.
4. Changelog and release notes are generated automatically by GoReleaser.

This split keeps responsibilities clear:

- Semantic release: version calculation + tag orchestration.
- GoReleaser: GitHub Release content + artifacts + changelog generation.

## Organization-Level Secrets and Variables

Use org-level credentials (shared across repos):

- Variable: `RELEASE_APP_ID=2823689`
- Secret: `RELEASE_APP_PRIVATE_KEY=<PEM>`

Workflow usage requirements:

- Read app credentials from `${{ vars.RELEASE_APP_ID }}` and `${{ secrets.RELEASE_APP_PRIVATE_KEY }}`.
- Scope generated app token to target repo when creating token:
  - `owner: go-modkit`
  - `repositories: modkit`

## Ruleset / Bypass Requirements

- GitHub App `go-modkit-release` must be a ruleset bypass actor for `main`.
- Keep required checks in place for normal contributor flow.
- Bypass is only for release automation using app token.

## Workflow Specs

### A) CI Workflow (`.github/workflows/ci.yml`)

Keep existing jobs/caching unchanged. Apply only these updates:

1. Keep coverage generation output at `.coverage/coverage.out`.
2. Upload coverage with `codecov/codecov-action@v5` without token input.
3. Set strict failure behavior: `fail_ci_if_error: true`.
4. Use OIDC for Codecov auth in public repo:
   - Add `use_oidc: true` on Codecov step.
   - Ensure permissions include `id-token: write` and `contents: read`.
5. Remove any reference to `${{ secrets.CODECOV_TOKEN }}`.

### B) Semantic Gate Workflow (`.github/workflows/release-semantic.yml`)

Trigger:

- `on.workflow_run.workflows: ["ci"]` (must match workflow `name` exactly)
- `types: [completed]`
- `branches: [main]`

Job gate condition:

- `github.event.workflow_run.conclusion == 'success'`
- `github.event.workflow_run.event == 'push'`
- `github.event.workflow_run.head_branch == 'main'`

Core steps:

1. Create GitHub App installation token via `actions/create-github-app-token@v2` using org var/secret and owner/repositories scoping.
2. Checkout repository at exact tested SHA: `${{ github.event.workflow_run.head_sha }}` with app token.
3. Fetch `origin/main` and fail if it moved since CI finished.
   - Guard condition: `origin/main` commit must equal `workflow_run.head_sha`.
4. Run `go-semantic-release/action@v1` with app token.

Required controls:

- `permissions: contents: write`
- `concurrency` single lane for semantic releases on main (`cancel-in-progress: false`)

Notes:

- Merge queue is recommended for branch discipline, but does not replace the SHA drift guard in `workflow_run` release architecture.

### C) Artifacts Workflow (`.github/workflows/release-artifacts.yml`)

Trigger:

- `on.push.tags: ["v*"]`

Core steps:

1. Checkout with `fetch-depth: 0`.
2. Keep current Go setup and cache blocks unchanged.
3. Run `goreleaser/goreleaser-action@v6` with `args: release --clean`.

Required controls:

- `permissions: contents: write`
- concurrency by tag reference:
  - group: `release-artifacts-${{ github.ref }}`
  - `cancel-in-progress: false`

## `.goreleaser.yml` Changelog Ownership

GoReleaser is the source of truth for release notes/changelog.

Required change:

- Enable changelog generation and set:

```yaml
changelog:
  use: github
```

Additional rules:

- Do not disable changelog in `.goreleaser.yml`.
- Keep existing build/archive/checksum settings unless needed for compatibility.

## Documentation Update

Update `docs/guides/release-process.md` to describe:

1. CI hard gate via `workflow_run` from `ci`.
2. Semantic tagging with GitHub App token and SHA drift guard.
3. Tag-driven GoReleaser publish flow.
4. Org-level app secret/variable convention.
5. GoReleaser-generated changelog ownership.

## Acceptance Criteria

1. PRs and pushes to `main` run `ci` with existing caching unchanged.
2. Codecov upload runs without `CODECOV_TOKEN` and fails CI on upload errors.
3. Semantic release runs only after successful `ci` on `main` push.
4. Semantic release exits without releasing if `origin/main` advanced beyond tested SHA.
5. Tag `vX.Y.Z` triggers artifact workflow.
6. GoReleaser publishes assets and generates release notes/changelog using `changelog.use: github`.
