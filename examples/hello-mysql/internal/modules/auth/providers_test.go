package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-modkit/modkit/modkit/module"
)

type noopResolver struct{}

func (noopResolver) Get(module.Token) (any, error) {
	return nil, nil
}

func TestAuthProviders_BuildsHandlerAndMiddleware(t *testing.T) {
	cfg := Config{Secret: "secret", Issuer: "issuer", TTL: time.Minute}
	defs := Providers(cfg)

	var handlerBuilt, mwBuilt bool
	for _, def := range defs {
		value, err := def.Build(noopResolver{})
		if err != nil {
			t.Fatalf("build: %v", err)
		}
		switch def.Token {
		case TokenHandler:
			_, handlerBuilt = value.(*Handler)
		case TokenMiddleware:
			_, mwBuilt = value.(func(http.Handler) http.Handler)
		}
	}
	if !handlerBuilt || !mwBuilt {
		t.Fatalf("handler=%v middleware=%v", handlerBuilt, mwBuilt)
	}
}
