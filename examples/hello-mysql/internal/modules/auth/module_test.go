package auth

import (
	"testing"
	"time"

	"github.com/go-modkit/modkit/modkit/kernel"
)

func TestModule_Bootstrap(t *testing.T) {
	mod := NewModule(Options{})
	_, err := kernel.Bootstrap(mod)
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
}

func TestAuthModule_Definition(t *testing.T) {
	cfg := Config{Secret: "secret", Issuer: "issuer", TTL: time.Minute}
	def := NewModule(Options{Config: cfg}).(Module).Definition()

	if def.Name != "auth" {
		t.Fatalf("name = %q", def.Name)
	}
	if len(def.Controllers) != 1 {
		t.Fatalf("controllers = %d", len(def.Controllers))
	}
}
