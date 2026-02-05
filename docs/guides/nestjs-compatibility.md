# NestJS Compatibility Guide

This guide maps NestJS concepts to modkit. It highlights what is implemented, what is intentionally skipped, and the Go-idiomatic alternatives.

## Feature Matrix

| Category | NestJS Feature | modkit Status | Notes |
|----------|----------------|---------------|-------|
| Modules | Module definition | ✅ Implemented | `ModuleDef` struct vs `@Module()` decorator |
| Modules | Imports | ✅ Implemented | Same concept |
| Modules | Exports | ✅ Implemented | Same concept |
| Modules | Providers | ✅ Implemented | Same concept |
| Modules | Controllers | ✅ Implemented | Same concept |
| Modules | Global modules | ⏭️ Skipped | Prefer explicit imports |
| Modules | Dynamic modules | ⏭️ Different | Use constructor functions with options |
| Modules | Module re-exporting | 🔄 Epic 02 | Export tokens from imported modules |
| Providers | Singleton scope | ✅ Implemented | Default and only scope |
| Providers | Request scope | ⏭️ Skipped | Use `context.Context` instead |
| Providers | Transient scope | ⏭️ Skipped | Use factory functions if needed |
| Providers | useClass | ✅ Implemented | `Build` function returning a concrete type |
| Providers | useValue | ✅ Implemented | `Build` returns a static value |
| Providers | useFactory | ✅ Implemented | `Build` is the factory |
| Providers | useExisting | ⏭️ Skipped | Use token aliases in `Build` |
| Providers | Async providers | ⏭️ Different | Go is sync; use goroutines if needed |
| Lifecycle | onModuleInit | ⏭️ Skipped | Put init logic in `Build()` |
| Lifecycle | onApplicationBootstrap | ⏭️ Skipped | Controllers built = app bootstrapped |
| Lifecycle | onModuleDestroy | 🔄 Epic 02 | Via `io.Closer` interface |
| Lifecycle | beforeApplicationShutdown | ⏭️ Skipped | Covered by `io.Closer` |
| Lifecycle | onApplicationShutdown | 🔄 Epic 02 | `App.Close()` method |
| Lifecycle | enableShutdownHooks | ⏭️ Different | Use `signal.NotifyContext` |
| HTTP | Controllers | ✅ Implemented | `RouteRegistrar` interface |
| HTTP | Route decorators | ⏭️ Different | Explicit `RegisterRoutes()` method |
| HTTP | Middleware | ✅ Implemented | Standard `func(http.Handler) http.Handler` |
| HTTP | Guards | ⏭️ Different | Implement as middleware |
| HTTP | Interceptors | ⏭️ Different | Implement as middleware |
| HTTP | Pipes | ⏭️ Different | Validate in handler or middleware |
| HTTP | Exception filters | ⏭️ Different | Error handling middleware |
| Other | CLI scaffolding | ❌ Not planned | Go boilerplate is minimal |
| Other | Devtools | ❌ Not planned | Use standard Go tooling |
| Other | Microservices | ❌ Not planned | Out of scope |
| Other | WebSockets | ❌ Not planned | Use `gorilla/websocket` directly |
| Other | GraphQL | ❌ Not planned | Use `gqlgen` directly |

## Justifications

### Global Modules

**NestJS:** `@Global()` makes a module’s exports available everywhere without explicit imports.

**modkit:** Skipped.

**Why:** Global modules hide dependencies, which conflicts with Go’s explicit import and visibility conventions.

**Alternative:** Construct the module once and import it explicitly wherever needed.

### Dynamic Modules

**NestJS:** `forRoot()`/`forRootAsync()` return module definitions at runtime.

**modkit:** Different.

**Why:** Go favors explicit constructors over dynamic metadata.

**Alternative:** Use module constructor functions that accept options and return a configured module instance.

### Request Scope

**NestJS:** Providers can be scoped per request.

**modkit:** Skipped.

**Why:** Per-request DI adds hidden lifecycle complexity in Go.

**Alternative:** Pass `context.Context` explicitly and construct request-specific values in handlers or middleware.

### Transient Scope

**NestJS:** Providers can be created on every injection.

**modkit:** Skipped.

**Why:** It encourages implicit, hidden object graphs.

**Alternative:** Use factory functions in `Build` or plain constructors where you need new instances.

### useExisting

**NestJS:** Alias one provider token to another with `useExisting`.

**modkit:** Skipped.

**Why:** Aliasing is simple and explicit in Go.

**Alternative:** Resolve the original token in `Build` and return it under the new token.

### Async Providers

**NestJS:** Providers can be async and awaited.

**modkit:** Different.

**Why:** Go initialization is synchronous; async is explicit and opt-in.

**Alternative:** Start goroutines in `Build` or expose `Start()` methods explicitly.

### onModuleInit

**NestJS:** Lifecycle hook invoked after a module is initialized.

**modkit:** Skipped.

**Why:** Initialization belongs in the constructor/build path in Go.

**Alternative:** Put setup logic in `Build()` or explicit `Start()` methods.

### onApplicationBootstrap

**NestJS:** Hook after the app finishes bootstrapping.

**modkit:** Skipped.

**Why:** modkit bootstraps deterministically when controllers are built.

**Alternative:** Use explicit post-bootstrap calls in `main`.

### beforeApplicationShutdown

**NestJS:** Hook before shutdown.

**modkit:** Skipped.

**Why:** Cleanup is modeled via `io.Closer` in Go.

**Alternative:** Implement `Close()` on providers and call `App.Close()`.

### enableShutdownHooks

**NestJS:** Enables signal handling for graceful shutdown.

**modkit:** Different.

**Why:** Go already provides standard signal handling primitives.

**Alternative:** Use `signal.NotifyContext` and call `App.Close()` and `http.Server.Shutdown()` explicitly.

### Route Decorators

**NestJS:** Decorators define routes on methods.

**modkit:** Different.

**Why:** Go does not use decorators or reflection for routing.

**Alternative:** Implement `RegisterRoutes(router)` and bind handlers explicitly.

### Guards

**NestJS:** Guard hooks control route access.

**modkit:** Different.

**Why:** Go middleware is the standard control point.

**Alternative:** Implement guards as `func(http.Handler) http.Handler` middleware.

### Interceptors

**NestJS:** Wrap request/response with cross-cutting logic.

**modkit:** Different.

**Why:** Go uses middleware for cross-cutting concerns.

**Alternative:** Implement as middleware or handler wrappers.

### Pipes

**NestJS:** Transform/validate input via pipes.

**modkit:** Different.

**Why:** Go favors explicit validation near the handler.

**Alternative:** Validate in handlers or middleware using standard libraries.

### Exception Filters

**NestJS:** Centralized exception handling layer.

**modkit:** Different.

**Why:** Errors are values in Go; handling is explicit.

**Alternative:** Use error-handling middleware or helpers that return `Problem Details` responses.

### CLI Scaffolding

**NestJS:** CLI generates boilerplate.

**modkit:** Not planned.

**Why:** Go projects are minimal and tooling is already strong.

**Alternative:** Use `go generate` or project templates if you want scaffolding.

### Devtools

**NestJS:** Devtools for inspection and debugging.

**modkit:** Not planned.

**Why:** Go relies on standard tooling and observability.

**Alternative:** Use `pprof`, logging, and standard debug tools.

### Microservices

**NestJS:** Built-in microservices framework.

**modkit:** Not planned.

**Why:** Out of modkit’s scope as a minimal backend framework.

**Alternative:** Use dedicated Go libraries for gRPC, NATS, or Kafka.

### WebSockets

**NestJS:** WebSocket gateway abstractions.

**modkit:** Not planned.

**Why:** Go already has solid standalone libraries.

**Alternative:** Use `gorilla/websocket` directly.

### GraphQL

**NestJS:** GraphQL module and decorators.

**modkit:** Not planned.

**Why:** Go has a strong, explicit GraphQL ecosystem.

**Alternative:** Use `gqlgen` directly.
