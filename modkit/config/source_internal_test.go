package config

import "testing"

func TestEnvSourceLookup(t *testing.T) {
	t.Setenv("CONFIG_SOURCE_TEST", "value")

	src := envSource{}
	v, ok := src.Lookup("CONFIG_SOURCE_TEST")
	if !ok {
		t.Fatalf("expected key to exist")
	}
	if v != "value" {
		t.Fatalf("unexpected value: %q", v)
	}

	if _, ok := src.Lookup("CONFIG_SOURCE_TEST_MISSING"); ok {
		t.Fatalf("expected missing key")
	}
}
