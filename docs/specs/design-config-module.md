# Design Spec: modkit Core Config Module

**Status:** Draft
**Date:** 2026-02-07
**Author:** Sisyphus (AI Agent)
**Related PRD:** `docs/specs/prd-modkit-core.md` (Phase 2 `Config Module`)

## 1. Overview

This document specifies a reusable Config Module for core `modkit`.

Today, configuration loading exists only in example apps (`examples/hello-mysql/internal/config` and `examples/hello-mysql/internal/platform/config`). This spec defines a standard, reusable way to load typed configuration from environment variables into the DI container while preserving modkit constraints: explicit wiring, deterministic behavior, typed errors, and no reflection magic.

## 2. Goals

- Provide a first-class config pattern for modkit applications.
- Support typed config resolution via `module.Get[T]`.
- Make required keys, defaults, and parse failures explicit.
- Preserve module visibility/export semantics for config tokens.
- Keep v1 lightweight and standard-library-first.

## 3. Non-Goals

- Replacing full-feature config frameworks (Viper/Koanf scope).
- Implicit binding via struct tags/reflection.
- Dynamic hot-reload/watch mode in v1.
- Secret manager integrations (Vault/SSM/etc.) in v1.
- Owning full domain validation (business constraints stay in feature modules).

## 4. Design Principles

- **Explicit schema:** callers define exactly which keys are read and how they parse.
- **Deterministic behavior:** same source and schema always produce same results.
- **No hidden globals:** no package-level mutable config singleton.
- **Composable modules:** config is a normal module with normal exports.
- **Typed contextual errors:** include key, expected type, and wrapped parse/source error.
- **Secret-safe diagnostics:** values are never included for sensitive keys.

## 5. Proposed Package and API Shape

### 5.1. Package Location

- New core package: `modkit/config`

### 5.2. Core Source Abstraction

```go
type Source interface {
    Lookup(key string) (value string, ok bool)
}
```

v1 default source is environment (`os.LookupEnv`).

### 5.3. Core Option and Spec Types

```go
type Option func(*Builder)

type Builder struct {
    source Source
}

type ValueSpec[T any] struct {
    Key         string
    Required    bool
    AllowEmpty  bool
    Default     *T
    Sensitive   bool
    Description string
    Parse       func(raw string) (T, error)
}
```

Notes:
- No `schema any` in v1.
- Every parsed value is explicit via `ValueSpec[T]`.
- `Sensitive` controls diagnostics redaction only; it does not change parsing semantics.
- `AllowEmpty` controls whether an empty-after-trim value is treated as explicitly set (`true`) or unset (`false`, default).

### 5.4. Helper Constructors and Parsers

```go
func NewModule(opts ...Option) module.Module
func WithSource(src Source) Option
func WithTyped[T any](token module.Token, spec ValueSpec[T], export bool) Option

func ParseString(raw string) (string, error)
func ParseInt(raw string) (int, error)
func ParseFloat64(raw string) (float64, error)
func ParseBool(raw string) (bool, error)
func ParseDuration(raw string) (time.Duration, error)
func ParseCSV(raw string) ([]string, error)
```

Design intent:
- `NewModule` returns a regular module that registers providers for configured tokens.
- `WithTyped` is called once per token; this preserves explicitness and avoids reflection.
- Parsers are public helpers so apps can share tested behavior and compose custom parsers.

## 6. Token and Visibility Model

### 6.1. Token Convention

- Recommended prefix remains `config.`
- v1 recommendation: typed tokens only (for example `config.app`, `config.http`, `config.auth`)
- No built-in `config.raw` export in v1

### 6.2. Visibility

- Config providers follow standard `ModuleDef.Exports` semantics.
- Config module internals remain private unless explicitly exported.
- No visibility exceptions are added for config.

## 7. Loading and Parsing Semantics

### 7.1. Raw Value Resolution

For each key:

1. Read from `Source.Lookup`.
2. Trim surrounding spaces.
3. Treat missing as "unset".
4. Treat empty-after-trim as:
   - unset when `AllowEmpty == false` (default),
   - explicitly set when `AllowEmpty == true`.
5. If unset and default exists, use default.
6. If unset and required, return `MissingRequiredError`.
7. Otherwise parse via the spec parser.

### 7.2. Parser Semantics

- `ParseDuration` uses Go standard `time.ParseDuration` only (no aliases in v1).
- `ParseCSV` splits on `,`, trims each value, and drops empty items.
- Parse failures return `ParseError` with key and expected type context.

### 7.3. Failure Timing

modkit providers are lazy. Therefore:
- config load/parse fails when the config provider is first resolved.
- in typical apps this occurs during bootstrap when controllers/providers resolve config.
- if an app requires strict eager validation, it should add a startup provider that resolves required config tokens explicitly.

## 8. Error Model

```go
type MissingRequiredError struct {
    Key       string
    Token     module.Token
    Sensitive bool
}

type ParseError struct {
    Key       string
    Token     module.Token
    Type      string
    Sensitive bool
    Err       error
}

type InvalidSpecError struct {
    Token  module.Token
    Reason string
}

var (
    ErrMissingRequired = errors.New("config: missing required key")
    ErrParse           = errors.New("config: parse error")
    ErrInvalidSpec     = errors.New("config: invalid spec")
)
```

Requirements:
- Errors must support `errors.As` for typed inspection.
- Errors must support `errors.Is` against sentinel categories (`ErrMissingRequired`, `ErrParse`, `ErrInvalidSpec`).
- `MissingRequiredError`, `ParseError`, and `InvalidSpecError` should implement `Unwrap()` to return their category sentinel; `ParseError.Unwrap()` should also preserve the underlying parse/source cause.
- Error strings include key/token/type context, never secret values for sensitive keys.
- Invalid spec errors represent developer misconfiguration and should fail deterministically.

## 9. Security Considerations

- Never include raw secret values in logs or error strings.
- Respect `Sensitive` for key-level redaction.
- If debug/inspection output is introduced later, default it to redacted values.
- Docs should continue recommending env/secret-store injection instead of committed files.

## 10. Integration Pattern

```go
type AppConfig struct {
    HTTPAddr  string
    JWTSecret string
    JWTTTL    time.Duration
}

const TokenAppConfig module.Token = "config.app"

func newConfigModule() module.Module {
    return config.NewModule(
        config.WithTyped(TokenAppConfig,
            config.ValueSpec[AppConfig]{
                Key:      "APP_CONFIG", // illustrative: app may also compose smaller tokens
                Required: true,
                Parse:    parseAppConfig,
            },
            true,
        ),
    )
}

func (m *AppModule) Definition() module.ModuleDef {
    cfgModule := newConfigModule()

    return module.ModuleDef{
        Name:    "app",
        Imports: []module.Module{cfgModule},
        Providers: []module.ProviderDef{{
            Token: "app.service",
            Build: func(r module.Resolver) (any, error) {
                cfg, err := module.Get[AppConfig](r, TokenAppConfig)
                if err != nil {
                    return nil, err
                }
                return NewService(cfg), nil
            },
        }},
    }
}
```

Note: the example above shows API shape only. Real-world usage will usually define explicit specs per env key and compose a typed struct in provider build logic.

## 11. Testing Strategy

### 11.1. Unit Tests

- Source resolution: present/missing/whitespace/default/required.
- Parser helpers: valid and invalid cases for each supported type.
- Error typing: `errors.As` for missing/parse/spec errors.
- Redaction behavior for sensitive keys.

### 11.2. Integration Tests

- Bootstrapping succeeds with valid env.
- First config resolution fails with descriptive error for missing required key.
- Visibility checks for non-exported config tokens.

### 11.3. Compatibility Tests

- Migration tests preserve existing `hello-mysql` env key behavior.
- Duration parsing remains compatible with current `JWT_TTL` usage (`time.ParseDuration`).

## 12. Adoption and Migration Plan

1. Add `modkit/config` with source abstraction, parser helpers, and typed spec plumbing.
2. Add docs guide for recommended token naming and module wiring.
3. Migrate `examples/hello-mysql` incrementally to core config helpers.
4. Keep small app-local wrappers only where domain-specific transformation is needed.

## 13. Acceptance Criteria

This PRD item is complete when all are true:

1. A core `modkit/config` package exists and is documented.
2. Apps can resolve typed config via `module.Get[T]`.
3. Missing/invalid required config returns typed descriptive errors.
4. Sensitive keys are redacted from diagnostic/error value surfaces.
5. At least one example app demonstrates the core pattern.

## 14. Resolved v1 Decisions

1. **Raw token exposure:** typed-only by default; no built-in `config.raw` export in v1.
2. **Duration parsing:** strict Go duration format only (`time.ParseDuration`).
3. **`.env` support:** out of core scope; callers may wrap `Source` externally.
4. **Validation ownership:** config module loads/parses only; feature modules validate domain constraints.

## 15. Future Enhancements (Not in v1)

- Multi-source layering (env + file + flags).
- Optional eager validation helper for strict startup guarantees.
- Secret manager source adapters.
- Schema export for documentation generation.

## 16. Implementation Blueprint (v1)

This section defines a minimal delivery plan mapped to concrete files and test coverage.

### 16.1. New Package Files

Create `modkit/config/` with:

- `source.go`
  - `type Source interface { Lookup(key string) (string, bool) }`
  - default env source implementation (`os.LookupEnv` wrapper)
- `parse.go`
  - `ParseString`, `ParseInt`, `ParseFloat64`, `ParseBool`, `ParseDuration`, `ParseCSV`
- `errors.go`
  - `MissingRequiredError`, `ParseError`, `InvalidSpecError` with wrapped errors where applicable
- `module.go`
  - option plumbing (`Option`, builder)
  - `NewModule(opts ...Option) module.Module`
  - `WithSource`, `WithTyped`
  - provider creation from typed specs, export handling

### 16.2. Unit Tests

Add:

- `modkit/config/parse_test.go`
  - table-driven tests for all parser helpers
  - whitespace behavior and invalid parse assertions
- `modkit/config/errors_test.go`
  - error string context and `errors.As`/`errors.Is` behavior
- `modkit/config/module_test.go`
  - required/default behavior
  - missing required error typing
  - parse failure typing and context
  - sensitive redaction safety checks
  - token export behavior through kernel visibility

### 16.3. Integration Tests (Kernel + Config)

Add `modkit/config/integration_test.go`:

- boot app with config module and successful typed resolution
- verify non-exported config token is not visible to importer (`TokenNotVisibleError`)
- verify first config resolution fails for required missing key with typed error

### 16.4. Example Migration (hello-mysql)

Incremental migration target:

1. Keep `examples/hello-mysql/internal/platform/config` public API stable (`Load()` remains).
2. Internally re-implement parsing with `modkit/config` helpers first.
3. Optionally introduce module-level config tokens for consumers in `internal/modules/app` after behavior parity tests pass.

This keeps runtime behavior stable while demonstrating adoption.

### 16.5. Documentation Updates

- Add `docs/guides/configuration.md` with env-first typed config examples.
- Cross-link from `README.md` guide list.
