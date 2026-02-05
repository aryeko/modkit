package main

import (
	"testing"
	"time"
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
