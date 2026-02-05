# Lifecycle & Cleanup Patterns Design

**Goal:** Add an optional cleanup hook to core providers and demonstrate graceful shutdown, cleanup order, and context cancellation in `hello-mysql`.

**Architecture:**
- Extend core provider metadata with an optional `Cleanup func(ctx context.Context) error` to keep lifecycle explicit without new interfaces.
- `hello-mysql` registers cleanup hooks (DB close) in provider definitions and runs them on shutdown in LIFO order.
- Use `http.Server.Shutdown` with a timeout context and run cleanup hooks after server shutdown.

**Components:**
- Core: `modkit/module/provider.go` adds optional cleanup field and accessor path for registered hooks.
- Example DB module: register cleanup hook to close the DB.
- App main: graceful shutdown handler invokes LIFO cleanup hooks.
- Users service: add a long-running operation that honors `ctx.Done()` for cancellation.
- Tests: lifecycle tests for LIFO cleanup and shutdown behavior; cancellation tests for users service.

**Error Handling:**
- Cleanup is best-effort: collect errors but run all hooks.
- Cancellation returns `context.Canceled` or `context.DeadlineExceeded` as appropriate.

**Testing:**
- Unit tests for cleanup hook registration ordering.
- Integration-ish tests in `hello-mysql` for shutdown + cleanup LIFO.
- Users service test for context cancellation.

**Notes:**
- Avoid flaky sleeps; use channels and bounded timeouts.
- Document cleanup order and hook registration in README.
