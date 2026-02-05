package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/app"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/auth"
)

func TestBuildHandler_DocsRoute(t *testing.T) {
	h, err := BuildHandler(app.Options{
		HTTPAddr: ":8080",
		MySQLDSN: "root:password@tcp(localhost:3306)/app?parseTime=true&multiStatements=true",
		Auth: auth.Config{
			Secret:   "dev-secret-change-me",
			Issuer:   "hello-mysql",
			TTL:      time.Hour,
			Username: "demo",
			Password: "demo",
		},
	})
	if err != nil {
		t.Fatalf("build handler: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/docs/index.html", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
