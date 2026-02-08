# ADR-0001: Devtools Scope for modkit Core v1

- **Status:** Accepted
- **Date:** 2026-02-08
- **Deciders:** modkit maintainers
- **Related:** `docs/specs/design-p2-graph-visualization-and-devtools.md`

## Context

The roadmap requires a P2 decision on devtools scope in core. We needed a concrete answer that preserves modkit goals:

- explicit architecture and no runtime magic
- minimal stable core API surface
- safe defaults with no accidental debug exposure

Two options were considered:

1. Minimal built-in devtools endpoints (explicit opt-in)
2. Formal de-scope from core, with guidance toward standard Go tooling

## Decision

Choose **Option B: formal de-scope from modkit core for v1**.

modkit core will not ship framework-owned debug endpoints in v1.

## Rationale

1. Keeps core focused on architecture and DI primitives, not runtime hosting concerns.
2. Avoids a default security/support burden around debug surfaces.
3. Maintains a smaller long-term compatibility surface.
4. Aligns with Go ecosystem norms (`pprof`, `delve`, logging/metrics tooling).
5. Observability value is still delivered by graph export (`ExportGraph` / `ExportAppGraph`) and docs recipes.

## Rejected Alternative

### Option A: Minimal Built-in Devtools Endpoints

Rejected for v1 because even opt-in framework endpoints add ownership burden for security hardening, compatibility guarantees, and maintenance across runtime contexts.

## Consequences

- Core roadmap language is updated from "decision pending" to "de-scoped from core (v1)".
- NestJS compatibility guidance is updated to reflect this finalized decision.
- Users are directed to standard Go diagnostics (`pprof`, `delve`) and optional external adapters for custom debug endpoints.
- Future reconsideration can happen in a separate proposal if there is strong demand and clear scope boundaries.
