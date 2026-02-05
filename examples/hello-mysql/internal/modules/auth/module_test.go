package auth

import (
	"testing"

	"github.com/go-modkit/modkit/modkit/kernel"
)

func TestModule_Bootstrap(t *testing.T) {
	mod := NewModule(Options{})
	_, err := kernel.Bootstrap(mod)
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
}
