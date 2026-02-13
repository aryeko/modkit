package httpserver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorResponseWriter struct {
	header http.Header
	status int
}

func (w *errorResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *errorResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *errorResponseWriter) Write(_ []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return 0, errors.New("write failed")
}

func TestHealthEncodeErrorReturnsServerError(t *testing.T) {
	w := &errorResponseWriter{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	controller := &HealthController{}
	controller.health(w, req)

	if w.status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.status)
	}
}
