# Design Spec: P2 Graph Visualization and Devtools Direction

**Status:** Ready for implementation
**Date:** 2026-02-08
**Author:** Sisyphus (AI Agent)
**Related PRD:** `docs/specs/prd-modkit-core.md` (P2 Graph visualization, P2 Devtools direction decision)

## 1. Overview

This spec defines P2 as two coordinated tracks:

1. **P2.1 Graph Visualization**: add first-class graph export output (Mermaid and DOT) from the existing kernel graph.
2. **P2.2 Devtools Direction Decision**: make and codify a clear decision for built-in devtools scope (implement minimal scope or formally de-scope).

Primary outcome: improve architecture observability without violating modkit constraints (explicit behavior, no runtime magic, no hidden globals).

## 2. Motivation

The PRD explicitly calls out observability and graph introspection as core value, and Phase 3 still has both graph export and devtools unresolved.

Current state in repo:

- Kernel already has a complete graph model (`modkit/kernel/graph.go`) suitable for deterministic export.
- No graph export API exists yet (no Mermaid/DOT emitter in core packages).
- Compatibility docs still treat devtools as not planned while PRD says decision pending, so direction must be resolved.

## 3. Goals

1. Provide deterministic graph export from existing kernel graph state.
2. Keep graph export read-only and side-effect free.
3. Avoid new reflection/decorator/magic patterns.
4. Resolve devtools direction with a single, documented product decision.
5. Synchronize PRD and compatibility docs after decision.

## 4. Non-Goals

1. No runtime hot-reload tooling.
2. No UI/dashboard application in P2.
3. No automatic network-exposed debug endpoints by default.
4. No graph mutation/edit APIs.

## 5. P2.1 Graph Visualization

### 5.1 Proposed API

Add deterministic exporters in `modkit/kernel`:

```go
package kernel

type GraphFormat string

const (
    GraphFormatMermaid GraphFormat = "mermaid"
    GraphFormatDOT     GraphFormat = "dot"
)

func ExportGraph(g *Graph, format GraphFormat) (string, error)
func ExportAppGraph(app *App, format GraphFormat) (string, error)
```

Add typed errors:

- `ErrNilGraph` (reuse existing)
- `ErrNilApp` (for `ExportAppGraph(nil, ...)`)
- `UnsupportedGraphFormatError{Format}`

### 5.2 Export Semantics

1. Node set = all modules in `Graph.Modules`.
2. Directed edge `A -> B` means module A imports module B.
3. Output ordering is deterministic:
   - sort nodes by module name
   - sort each node's imports lexicographically
4. Export is pure serialization; it must not build providers/controllers.
5. Export includes root module annotation.
6. Exporters MUST NOT mutate `g.Modules`, `g.Nodes`, or `ModuleNode.Imports`; ordering is produced from copied slices only.
7. `ExportAppGraph` behavior:
   - `app == nil` returns `ErrNilApp`
   - `app.Graph == nil` returns `ErrNilGraph`

### 5.3 Encoding and Escaping Contract

1. Module names are serialized as labels and escaped per format.
2. DOT output always quotes node IDs and edge endpoints (`"name"`).
3. Mermaid output uses stable generated node IDs (`m0`, `m1`, ...) and escaped labels (`m0["module-name"]`) to avoid syntax breakage from arbitrary module names.
4. ID allocation is deterministic from sorted module names.
5. Root annotation is explicit in both formats:
   - Mermaid: add `classDef root` and `class <root-id> root`
   - DOT: add `"<root>" [shape=doublecircle];`

### 5.4 Output Shapes

Mermaid (example):

```text
graph TD
    m0["app"]
    m1["auth"]
    m2["users"]
    m3["db"]
    m0 --> m1
    m0 --> m2
    m2 --> m3
    classDef root stroke-width:3px;
    class m0 root;
```

DOT (example):

```text
digraph modkit {
    rankdir=LR;
    "app";
    "app" [shape=doublecircle];
    "users";
    "app" -> "users";
}
```

### 5.5 Test Requirements

1. Deterministic table-driven fixture tests for both formats.
2. Nil graph, nil app, and unsupported format error tests.
3. Edge-case coverage:
   - single-module graph
   - multi-import graph
   - re-export scenario should not alter edge semantics (imports-only edges)
4. Non-mutation invariant test: exporter does not alter graph/module/import order in memory.
5. Backward-compat check: no behavior change to existing bootstrap/visibility tests.

### 5.6 Documentation Deliverables

1. Add `docs/guides/graph-visualization.md` with examples for both formats.
2. Cross-link from `README.md` guide index.
3. Add one example snippet using `kernel.ExportAppGraph` after bootstrap.

## 6. P2.2 Devtools Direction Decision

### 6.1 Decision Inputs

Decision must evaluate:

1. Alignment with no-magic and explicit-architecture goals.
2. Security surface (default exposure risk).
3. Maintenance burden of long-term API support.
4. Overlap with standard Go tooling (`pprof`, `delve`, logging, metrics).
5. User value beyond graph export delivered in P2.1.

### 6.2 Decision Options

#### Option A: Minimal Built-in Devtools (Explicit Opt-In)

Scope:

- Provide optional helper that mounts debug endpoints only when user explicitly enables it.
- Initial endpoint set (minimal):
  - `GET /debug/modkit/graph?format=mermaid|dot`
  - `GET /debug/modkit/modules` (module names + imports)

Constraints:

1. Must be disabled by default.
2. Must require explicit route mounting by user code.
3. Must not expose provider instance values/secrets.

#### Option B: Formal De-Scope for Core

Scope:

- No built-in devtools endpoints in modkit core.
- Recommend standard Go tooling and optional external adapters.
- Keep roadmap focused on graph export + docs recipes.

### 6.3 P2 Recommendation

**Recommendation: Option B (Formal De-Scope for core in v1).**

Rationale:

1. Maintains minimal stable API surface in core.
2. Avoids security and support burden for framework-owned debug endpoints.
3. Aligns with current compatibility guide direction and Go ecosystem norms.
4. Still delivers observability via P2.1 graph export + docs recipes.

### 6.4 Decision Record Deliverable

Produce/update one ADR-style record in docs (or dedicated section) containing:

1. Decision date
2. Chosen option
3. Rationale and rejected alternatives
4. Consequences for roadmap/docs

## 7. Execution Plan

### Story P2.1 - Graph Export API

1. Implement `ExportGraph` and `ExportAppGraph` in kernel.
2. Add format/error types and tests.
3. Add docs guide + README cross-link.

### Story P2.2 - Devtools Decision Finalization

1. Lock Option B decision for v1 core scope (no built-in devtools endpoints).
2. Update PRD and compatibility guide wording to reflect this finalized decision.
3. Add explicit "de-scoped from core" language and recommended alternatives (`pprof`, `delve`, optional adapters).

## 8. Acceptance Criteria

P2 is complete when all are true:

1. Graph export supports Mermaid and DOT with deterministic output.
2. Export API has typed error behavior for nil app/nil graph/unsupported format cases.
3. Graph export tests cover deterministic fixtures, edge cases, and non-mutation invariants.
4. New graph visualization guide exists and is linked from README.
5. Devtools direction is finalized as Option B and documented in an ADR-style record.
6. PRD + compatibility docs are synchronized to the same Option B decision language.

## 9. Verification Plan

1. `lsp_diagnostics` clean on all edited files (or manual sanity checks when unavailable).
2. Run tests at minimum for affected package(s), then full required project verification before completion:
   - `make fmt && make lint && make vuln && make test && make test-coverage`
3. Confirm no unexpected API regressions in existing kernel tests.
4. Verify docs links and wording consistency across:
   - `docs/specs/prd-modkit-core.md`
   - `docs/guides/nestjs-compatibility.md`
   - `README.md`

## 10. Risks and Mitigations

1. **Risk:** Non-deterministic graph output causes flaky tests.
   - **Mitigation:** enforce lexical ordering in exporter and snapshot assertions.
2. **Risk:** Devtools decision remains ambiguous after implementation.
   - **Mitigation:** require ADR-style decision record and synchronized wording updates.
3. **Risk:** Core API bloat from debug features.
   - **Mitigation:** keep v1 recommendation as de-scope for devtools endpoints.

## 11. Out-of-Scope Follow-ups

1. Interactive graph UI or web dashboard.
2. Runtime graph diffing across deployments.
3. Full framework-managed diagnostics suite.
