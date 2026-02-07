package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("MYSQL_DSN", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_ISSUER", "")
	t.Setenv("JWT_TTL", "")
	t.Setenv("AUTH_USERNAME", "")
	t.Setenv("AUTH_PASSWORD", "")

	cfg := Load()

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q", cfg.HTTPAddr)
	}
	if cfg.JWTSecret != "dev-secret-change-me" {
		t.Fatalf("JWTSecret = %q", cfg.JWTSecret)
	}
}

func TestEnvOrDefault_TrimsSpace(t *testing.T) {
	t.Setenv("JWT_ISSUER", "   ")
	if got := envOrDefault("JWT_ISSUER", "hello-mysql"); got != "hello-mysql" {
		t.Fatalf("envOrDefault = %q", got)
	}
}

func TestEnvOrDefault_UsesDefaultWhenUnset(t *testing.T) {
	const key = "JWT_ISSUER_UNSET"
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv failed: %v", err)
	}

	if got := envOrDefault(key, "hello-mysql"); got != "hello-mysql" {
		t.Fatalf("envOrDefault = %q", got)
	}
}

func TestEnvOrDefault_ReturnsTrimmedValueWhenSet(t *testing.T) {
	t.Setenv("JWT_ISSUER", "  issuer-v1  ")

	if got := envOrDefault("JWT_ISSUER", "hello-mysql"); got != "issuer-v1" {
		t.Fatalf("envOrDefault = %q", got)
	}
}
