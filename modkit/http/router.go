package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router provides a minimal method-based handler registration API.
type Router interface {
	Handle(method string, pattern string, handler http.Handler)
}

// RouteRegistrar defines a controller that can register its HTTP routes.
type RouteRegistrar interface {
	RegisterRoutes(router Router)
}

type chiRouter struct {
	chi.Router
}

func (r *chiRouter) Handle(method string, pattern string, handler http.Handler) {
	r.Method(method, pattern, handler)
}

// NewRouter creates a chi router with baseline middleware for the HTTP adapter.
func NewRouter() Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	return &chiRouter{Router: router}
}

// RegisterRoutes invokes controller route registration functions.
func RegisterRoutes(router Router, controllers map[string]any) error {
	for name, controller := range controllers {
		registrar, ok := controller.(RouteRegistrar)
		if !ok {
			return &RouteRegistrationError{Name: name}
		}
		registrar.RegisterRoutes(router)
	}
	return nil
}
