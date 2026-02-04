# HTTP Adapter Notes

## Route Registration (No Reflection)

Controllers register routes explicitly by implementing:

```go
type RouteRegistrar interface {
	RegisterRoutes(router Router)
}
```

At bootstrap, controller instances are collected by the kernel and passed to the
HTTP adapter:

```go
err := RegisterRoutes(router, controllers)
```

`RegisterRoutes` type-asserts each controller to `RouteRegistrar` and invokes
`RegisterRoutes` to attach handlers to the router. This avoids reflection and
keeps routing explicit and testable.
