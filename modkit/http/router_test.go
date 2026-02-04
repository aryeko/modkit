package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter_AllowsRoute(t *testing.T) {
	router := NewRouter()
	router.Handle(http.MethodGet, "/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	handler, ok := router.(http.Handler)
	if !ok {
		t.Fatalf("router does not implement http.Handler")
	}

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}
