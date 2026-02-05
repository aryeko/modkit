package main

import (
	"testing"
	"time"

	"github.com/go-modkit/modkit/examples/hello-mysql/internal/platform/config"
)

func TestParseJWTTTL_DefaultOnInvalid(t *testing.T) {
	got := parseJWTTTL("bad-value")
	if got != time.Hour {
		t.Fatalf("ttl = %v", got)
	}
}

func TestParseJWTTTL_Valid(t *testing.T) {
	got := parseJWTTTL("30m")
	if got != 30*time.Minute {
		t.Fatalf("ttl = %v", got)
	}
}

func TestParseJWTTTL_RejectsNonPositive(t *testing.T) {
	for _, value := range []string{"0s", "-1s"} {
		got := parseJWTTTL(value)
		if got != time.Hour {
			t.Fatalf("ttl for %q = %v", value, got)
		}
	}
}

func TestBuildAuthConfig(t *testing.T) {
	cfg := config.Config{
		JWTSecret:    "secret",
		JWTIssuer:    "issuer",
		AuthUsername: "demo",
		AuthPassword: "s3cret",
	}
	ttl := 2 * time.Minute

	got := buildAuthConfig(cfg, ttl)

	if got.Secret != cfg.JWTSecret {
		t.Fatalf("secret = %q", got.Secret)
	}
	if got.Issuer != cfg.JWTIssuer {
		t.Fatalf("issuer = %q", got.Issuer)
	}
	if got.TTL != ttl {
		t.Fatalf("ttl = %v", got.TTL)
	}
	if got.Username != cfg.AuthUsername {
		t.Fatalf("username = %q", got.Username)
	}
	if got.Password != cfg.AuthPassword {
		t.Fatalf("password = %q", got.Password)
	}
}

func TestBuildAppOptions(t *testing.T) {
	cfg := config.Config{
		HTTPAddr:     ":9999",
		MySQLDSN:     "dsn",
		JWTSecret:    "secret",
		JWTIssuer:    "issuer",
		AuthUsername: "demo",
		AuthPassword: "s3cret",
	}
	ttl := 3 * time.Minute

	opts := buildAppOptions(cfg, ttl)

	if opts.HTTPAddr != cfg.HTTPAddr {
		t.Fatalf("http addr = %q", opts.HTTPAddr)
	}
	if opts.MySQLDSN != cfg.MySQLDSN {
		t.Fatalf("mysql dsn = %q", opts.MySQLDSN)
	}
	if opts.Auth.Secret != cfg.JWTSecret {
		t.Fatalf("auth secret = %q", opts.Auth.Secret)
	}
	if opts.Auth.Issuer != cfg.JWTIssuer {
		t.Fatalf("auth issuer = %q", opts.Auth.Issuer)
	}
	if opts.Auth.TTL != ttl {
		t.Fatalf("auth ttl = %v", opts.Auth.TTL)
	}
	if opts.Auth.Username != cfg.AuthUsername {
		t.Fatalf("auth username = %q", opts.Auth.Username)
	}
	if opts.Auth.Password != cfg.AuthPassword {
		t.Fatalf("auth password = %q", opts.Auth.Password)
	}
}
