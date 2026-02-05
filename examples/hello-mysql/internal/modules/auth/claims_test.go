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
