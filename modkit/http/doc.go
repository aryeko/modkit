// Package http adapts controller instances to HTTP routing.
//
// Route registration is explicit: controllers implement RouteRegistrar and are
// invoked via RegisterRoutes. No reflection is used.
package http
