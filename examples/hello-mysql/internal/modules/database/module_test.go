package database

import (
	"context"
	"testing"
)

func TestModuleDefinition_ProviderCleanupHook(t *testing.T) {
	def := Module{}.Definition()
	if len(def.Providers) == 0 {
		t.Fatal("expected at least one provider")
	}
	var cleanup func(ctx context.Context) error
	for _, provider := range def.Providers {
		if provider.Token == TokenDB {
			cleanup = provider.Cleanup
			break
		}
	}
	if cleanup == nil {
		t.Fatal("expected provider cleanup hook")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := cleanup(ctx); err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
