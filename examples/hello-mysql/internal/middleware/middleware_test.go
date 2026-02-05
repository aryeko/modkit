package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_AddsHeaders(t *testing.T) {
	cors := NewCORS(CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: nil,
	})

	handler := cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("expected allow origin header to be set")
	}
	if rec.Header().Get("Access-Control-Allow-Methods") != "GET, POST" {
		t.Fatalf("expected allow methods header to be set")
	}
	if _, ok := rec.Header()["Access-Control-Allow-Headers"]; !ok {
		t.Fatalf("expected allow headers header to be set")
	}
}
