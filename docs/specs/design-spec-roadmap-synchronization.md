# Design Spec: P1 Spec/Roadmap Synchronization

**Status:** Ready for implementation
**Date:** 2026-02-08
**Author:** Sisyphus (AI Agent)
**Related PRD:** `docs/specs/prd-modkit-core.md` (P1 Spec/roadmap synchronization)

## 1. Overview

This spec defines a focused P1 docs pass to reconcile roadmap/spec status with shipped code and current product direction.

Goal: a contributor reading PRD + epics + key guides should get the same project-state answer as repository evidence (code, tests, workflows, README).

## 2. Motivation

Status and checklist drift is now a productivity and prioritization problem:

1. Roadmap priority is already explicitly set to P1 in PRD, but supporting specs are not fully reconciled.
2. Some docs still describe shipped items as pending/not planned.
3. Epic checklists remain largely unchecked despite implementation evidence.

When status docs lag reality, planning confidence and onboarding quality drop.

## 3. Problem Statement

Current documentation has concrete synchronization gaps.

### 3.1 Confirmed Mismatch Targets

1. `docs/specs/epic-01-examples-enhancement.md`
   - Story/task checklists are mostly unchecked while example README and repository behavior show many delivered capabilities.
2. `docs/specs/epic-02-core-nest-compatibility.md`
   - Acceptance criteria for graceful shutdown/re-export/docs remain unchecked although implemented and documented.
3. `docs/guides/nestjs-compatibility.md`
   - Matrix still says CLI scaffolding "Not planned" and lifecycle/re-export as "This Epic", conflicting with shipped state.
4. `docs/specs/design-release-versioning-sdlc-cli.md`
   - Status/checklist still draft/pending while release pipeline and CLI artifacts are in place.
5. `README.md`
   - Guide index does not link configuration guide despite config module/spec being implemented.

## 4. Goals

1. Normalize spec status language across roadmap docs.
2. Reconcile epic acceptance checklists with evidence-based completion state.
3. Remove contradictory roadmap statements between guides/specs/README.
4. Keep all changes docs-only and additive (no feature work).
5. Establish repeatable synchronization policy so drift does not recur.

## 5. Non-Goals

1. No new runtime/framework feature implementation.
2. No broad content rewrite of technical guides beyond status alignment.
3. No retrospective perfection pass on every historical note.
4. No weakening of explicit architecture constraints for doc simplicity.

## 6. Canonical Status Model

All roadmap/spec docs in scope must use one of these status labels:

- `Draft`
- `Active`
- `Implemented (v1)`
- `Ready for implementation`
- `Superseded`

Rules:

1. `Implemented (v1)` means shipped with repository evidence and expected CI/docs wiring.
2. `Ready for implementation` means design locked with acceptance criteria but code not delivered.
3. `Draft` means substantive open design uncertainty remains.
4. Every roadmap/spec file updated in this pass gets `Last Reviewed: <execution-date YYYY-MM-DD>`.

## 7. Source-of-Truth Precedence

When documentation conflicts, reconciliation decisions follow this order:

1. Shipped code and tests in repository.
2. CI/release workflow behavior (`.github/workflows`, `.goreleaser.yml`).
3. PRD roadmap and explicit priority list.
4. Feature guides and epic/spec narrative text.

Checklist items are marked complete only if the same session can point to evidence paths.

## 8. Scope and Deliverables

### 8.1 In-Scope Files

1. `docs/specs/epic-01-examples-enhancement.md`
2. `docs/specs/epic-02-core-nest-compatibility.md`
3. `docs/guides/nestjs-compatibility.md`
4. `docs/specs/design-release-versioning-sdlc-cli.md`
5. `README.md`
6. `docs/specs/prd-modkit-core.md` (required for synchronization summary and final consistency cross-links)

### 8.2 Required Deliverables

1. Updated status/checklist state in each in-scope roadmap/spec file; wording/link consistency updates in in-scope guides/README.
2. Explicit wording alignment for CLI, lifecycle, re-export, and devtools direction.
3. Configuration guide link added to `README.md` guide list.
4. One concise "synchronization summary" section in PRD (or linked from PRD) listing what was reconciled and when.

## 9. Reconciliation Method

For each in-scope file:

1. Extract status/checklist claims.
2. Map each claim to repository evidence paths.
3. Update claim state (`[x]`, `[ ]`, status label, wording) based on evidence.
4. Keep unresolved items explicit; do not force-close ambiguous criteria.
5. Add/refresh `Last Reviewed` line for roadmap/spec files.

## 10. Story Breakdown

### Story P1.1 - Epic 02 State Reconciliation

Update `docs/specs/epic-02-core-nest-compatibility.md` checklists for:

- `App.Close`/`App.CloseContext`
- closer ordering/error aggregation/idempotency evidence
- module re-export criteria
- compatibility guide presence and cross-link state

### Story P1.2 - Epic 01 Checklist Reconciliation

Update `docs/specs/epic-01-examples-enhancement.md` to reflect shipped example/testing/docs items, leaving genuinely open work unchecked.

### Story P1.3 - Compatibility Matrix Truth Pass

Update `docs/guides/nestjs-compatibility.md` statuses to match current project state:

- CLI scaffolding no longer "Not planned"
- lifecycle/re-export wording moved from "This Epic" to implemented state
- devtools wording aligned with PRD decision-pending language

### Story P1.4 - SDLC/Release Spec State Update

Update `docs/specs/design-release-versioning-sdlc-cli.md` from stale draft checklist to implemented-state wording where evidence exists.

### Story P1.5 - README/Guide Link Consistency

Add/verify links so README and guides match implemented modules (including configuration guide).

## 11. Acceptance Criteria

P1 is complete when all are true:

1. All in-scope roadmap/spec files have synchronized status language and `Last Reviewed` metadata.
2. No known contradiction remains for CLI/lifecycle/re-export/devtools status across PRD, epics, and compatibility guide.
3. In-scope epic checklist updates are evidence-backed with a claim table in PR description: `claim -> evidence path(s) -> resolution`.
4. `README.md` guide index includes configuration guide and no broken roadmap-related links.
5. Remaining open items are explicitly tracked as open (not silently dropped).
   - Required evidence: an "Open Items Ledger" section in PR description with `item -> reason open -> next owner/step`.

## 12. Verification Plan

1. Run targeted grep checks for stale phrases with expected outcomes:
   - "Not planned" (CLI/devtools where inconsistent)
   - "This Epic" where feature is already implemented
   - outdated unchecked checklist items for shipped behavior
   - Expected: no contradictory hits remain in in-scope files after edits.
2. Manual read-through of all in-scope files after edits.
3. `lsp_diagnostics` clean for each edited markdown file when diagnostics are available; otherwise, complete a markdown sanity pass (headings, list formatting, link syntax).
4. `git diff` review confirms docs-only scope.
5. Before PR create/update, run repository-required verification commands and report outcomes:
   - `make fmt && make lint && make vuln && make test && make test-coverage`

## 13. Risks and Mitigations

1. **Risk:** Over-checking items without enough evidence.
   - **Mitigation:** Require path-level evidence before checking any box.
2. **Risk:** Under-checking and preserving stale alarms.
   - **Mitigation:** Apply source-of-truth precedence and explicit decision log.
3. **Risk:** Reintroducing contradictions in future updates.
   - **Mitigation:** Add lightweight synchronization policy note to PRD/docs-spec process.

## 14. Out-of-Scope Follow-ups

1. P2 Graph visualization spec/implementation.
2. P2 Devtools direction final decision and corresponding spec.
3. Broader docs quality rewrite unrelated to roadmap-state correctness.

## 15. Implementation Notes

This is intentionally a docs-only operational spec. It should be executed as a short, auditable pass with one PR containing:

1. Updated statuses/checklists/wording
2. Evidence-backed rationale in PR description
3. No behavioral code changes
