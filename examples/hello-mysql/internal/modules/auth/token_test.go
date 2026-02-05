package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestIssueToken_InvalidConfig(t *testing.T) {
	_, err := IssueToken(Config{Secret: "", TTL: time.Minute}, User{ID: "demo"})
	if err == nil {
		t.Fatal("expected error")
	}

	_, err = IssueToken(Config{Secret: "secret", TTL: 0}, User{ID: "demo"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestIssueToken_EmptyUserClaims(t *testing.T) {
	cfg := Config{
		Secret: "secret",
		Issuer: "issuer",
		TTL:    time.Minute,
	}
	token, err := IssueToken(cfg, User{})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	parsed, err := parseToken(token, cfg, time.Now())
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected map claims")
	}
	if _, ok := claims["sub"]; ok {
		t.Fatalf("unexpected sub claim: %v", claims["sub"])
	}
	if _, ok := claims["email"]; ok {
		t.Fatalf("unexpected email claim: %v", claims["email"])
	}
}
