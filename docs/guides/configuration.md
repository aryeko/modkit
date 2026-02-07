# Configuration

This guide describes the recommended configuration pattern for modkit apps using the core `modkit/config` package design.

If your project has not adopted `modkit/config` yet, use this guide as the migration target from ad-hoc `os.Getenv` calls.

## Design Goals

- Keep configuration explicit and testable.
- Parse once, then inject typed values via DI.
- Fail with contextual typed errors on missing/invalid required values.
- Preserve normal module visibility rules for config tokens.

## Core Pattern

Use a dedicated config module that exports typed tokens.

```go
package app

import (
    "time"

    "github.com/go-modkit/modkit/modkit/config"
    "github.com/go-modkit/modkit/modkit/module"
)

const (
    TokenHTTPAddr module.Token = "config.http_addr"
    TokenJWTTTL   module.Token = "config.jwt_ttl"
)

func NewConfigModule() module.Module {
    defaultAddr := ":8080"
    defaultTTL := 1 * time.Hour

    return config.NewModule(
        config.WithTyped(TokenHTTPAddr, config.ValueSpec[string]{
            Key:      "HTTP_ADDR",
            Default:  &defaultAddr,
            Parse:    config.ParseString,
            Required: false,
        }, true),
        config.WithTyped(TokenJWTTTL, config.ValueSpec[time.Duration]{
            Key:      "JWT_TTL",
            Default:  &defaultTTL,
            Parse:    config.ParseDuration,
            Required: false,
        }, true),
    )
}
```

When an app needs multiple independent config modules, set distinct names with `config.WithModuleName("...")` to avoid duplicate module names in the graph.

## Consuming Typed Config

Resolve config in providers/controllers using `module.Get[T]`:

```go
module.ProviderDef{
    Token: "auth.service",
    Build: func(r module.Resolver) (any, error) {
        ttl, err := module.Get[time.Duration](r, TokenJWTTTL)
        if err != nil {
            return nil, err
        }
        return NewAuthService(ttl), nil
    },
}
```

## Required Values and Sensitive Keys

Use `Required: true` for values that must exist, and `Sensitive: true` for secret-bearing keys:

```go
secretSpec := config.ValueSpec[string]{
    Key:       "JWT_SECRET",
    Required:  true,
    Sensitive: true,
    Parse:     config.ParseString,
}
```

Behavior:
- Missing required values return `MissingRequiredError`.
- Parse failures return `ParseError` with key/type context.
- Sensitive values are never included in diagnostic value surfaces.

Use `errors.Is` with sentinels for category checks:

- `config.ErrMissingRequired`
- `config.ErrParse`
- `config.ErrInvalidSpec`

Use `errors.As` when you need structured fields (`Key`, `Token`, `Type`, `Sensitive`).

## Parser Helpers

The core helpers cover common env value types:

- `ParseString`
- `ParseInt`
- `ParseFloat64`
- `ParseBool`
- `ParseDuration` (Go `time.ParseDuration` format)
- `ParseCSV` (comma-separated, trimmed, empty entries removed)

## Empty Value Semantics

By default, empty-after-trim values are treated as unset.

- `AllowEmpty: false` (default): empty values follow required/default behavior.
- `AllowEmpty: true`: empty values are treated as explicitly set and passed to `Parse`.

## Visibility and Exports

Config tokens behave exactly like any other provider token.

- Export only the tokens importers need.
- Keep internal/raw configuration private.
- If a token is not exported, importers receive `TokenNotVisibleError`.

## Testing Recommendations

- Use `t.Setenv` for key-by-key deterministic setup.
- Test default, missing, empty, whitespace, and invalid parse cases.
- Assert typed errors with `errors.As`.
- Verify secret-bearing keys do not leak raw values in errors.

## Migration from App-Local Config

For existing apps that use `os.Getenv` directly:

1. Keep existing `Load()` surface to avoid breaking callers.
2. Move parsing logic into `modkit/config` specs and helpers.
3. Export typed tokens and update consumers gradually.
4. Remove duplicated env parsing utilities once parity is verified.

## Related Docs

- [Modules](modules.md)
- [Providers](providers.md)
- [Testing](testing.md)
- [Error Handling](error-handling.md)
- [Design Spec: Core Config Module](../specs/design-config-module.md)
