package auth

import (
	"testing"
	"time"
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
