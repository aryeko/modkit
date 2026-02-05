# Module Re-exporting Design

**Goal:** Allow modules to re-export tokens from their imports with strict visibility validation and transitive propagation.

**Context:** Story 2.4 (Module Re-exporting - Implementation) in `docs/specs/epic-02-core-nest-compatibility.md`.

## Behavior & Semantics

- A module may list in `Exports` any token that is:
  - Provided by the module itself, or
  - Exported by a directly imported module.
- Transitive re-exporting is supported: if A imports B and B re-exports from C, A can access those tokens via B.
- Non-exported tokens remain private: re-exporting a token that an imported module does not export is rejected.
- Ambiguous re-exports are errors: if multiple imported modules export the same token, re-exporting that token is rejected.
- Global provider token uniqueness remains unchanged (duplicate provider tokens across modules are rejected in `BuildGraph`).
- No `All()` helper or module-level re-export shorthand is introduced.

## Architecture

The kernel already builds visibility by combining a module's own providers and the effective exports of its imports. We extend validation to treat "exports must be visible" as the rule, so imports can be re-exported only if they are already visible in the module. Effective exports are computed per module and then used to populate visibility for importers, enabling transitive access.

Ambiguity detection is handled when validating a module's exports: if a token is exported by more than one imported module and is re-exported, surface a deterministic error.

## Data Flow

1. Build module graph (`BuildGraph`) as today.
2. Build visibility:
   - Start with providers in the module.
   - Add tokens exported by imported modules (their effective exports).
3. Validate module exports against visibility:
   - If token not visible, return `ExportNotVisibleError`.
   - If token appears in multiple import export sets, return an error for ambiguous re-export.
4. Record module effective exports (explicit `Exports` list only) for downstream importers.

## Error Handling

- **Export not visible:** continue to use `ExportNotVisibleError` with the module name and token.
- **Ambiguous re-export:** return a clear error stating the token is exported by multiple imports and cannot be re-exported without disambiguation.

## Testing

Add unit tests covering:
- Valid re-export of imported token (visibility test).
- Invalid re-export of non-exported token (visibility test).
- Ambiguous re-export (visibility test).
- Transitive re-export: A imports B, B re-exports from C, A can access token (graph/visibility test).
- Negative transitive case: B imports C but does not export; A cannot access C token.

## Verification

- `go test ./modkit/kernel -run TestVisibility`
- `go test ./modkit/kernel -run TestGraph`

## Out of Scope

- `module.All()` export helper
- Documentation updates (handled in Story 2.5)
