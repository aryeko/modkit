package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestUserFromClaims_RejectsEmpty(t *testing.T) {
	if _, ok := userFromClaims(jwt.MapClaims{}); ok {
		t.Fatal("expected false")
	}
}

func TestUserFromClaims_SubOrEmail(t *testing.T) {
	user, ok := userFromClaims(jwt.MapClaims{"sub": "demo"})
	if !ok {
		t.Fatal("expected true for sub claim")
	}
	if user.ID != "demo" || user.Email != "" {
		t.Fatalf("unexpected user: %+v", user)
	}

	user, ok = userFromClaims(jwt.MapClaims{"email": "demo@example.com"})
	if !ok {
		t.Fatal("expected true for email claim")
	}
	if user.Email != "demo@example.com" || user.ID != "" {
		t.Fatalf("unexpected user: %+v", user)
	}
}
