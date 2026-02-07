package config_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-modkit/modkit/modkit/config"
	"github.com/go-modkit/modkit/modkit/module"
)

func TestErrors(t *testing.T) {
	t.Run("missing required", func(t *testing.T) {
		err := &config.MissingRequiredError{Key: "JWT_SECRET", Token: module.Token("config.jwt_secret"), Sensitive: true}
		if !strings.Contains(err.Error(), "JWT_SECRET") {
			t.Fatalf("missing key in error: %q", err.Error())
		}
		if strings.Contains(strings.ToLower(err.Error()), "secret-value") {
			t.Fatalf("error leaked value: %q", err.Error())
		}
	})

	t.Run("parse unwrap", func(t *testing.T) {
		inner := errors.New("bad int")
		err := &config.ParseError{
			Key:   "RATE_LIMIT_BURST",
			Token: module.Token("config.rate_limit_burst"),
			Type:  "int",
			Err:   inner,
		}

		if !errors.Is(err, inner) {
			t.Fatalf("expected ParseError to unwrap")
		}

		var pe *config.ParseError
		if !errors.As(err, &pe) {
			t.Fatalf("expected ParseError via errors.As")
		}
	})

	t.Run("invalid spec", func(t *testing.T) {
		err := &config.InvalidSpecError{Token: module.Token("config.jwt_ttl"), Reason: "parse function must not be nil"}
		if !strings.Contains(err.Error(), "invalid config spec") {
			t.Fatalf("unexpected error string: %q", err.Error())
		}
	})
}
