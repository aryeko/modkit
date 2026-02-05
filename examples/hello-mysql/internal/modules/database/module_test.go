package database

import (
	"context"
	"testing"
)

func TestModuleDefinition_ProviderCleanupHook_CanceledContext(t *testing.T) {
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

func TestDatabaseModule_Definition_ProvidesDB(t *testing.T) {
	mod := NewModule(Options{DSN: "dsn"})
	def := mod.(*Module).Definition()
	if def.Name != "database" {
		t.Fatalf("expected name database, got %q", def.Name)
	}
	if len(def.Providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(def.Providers))
	}
	if def.Providers[0].Token != TokenDB {
		t.Fatalf("expected TokenDB, got %q", def.Providers[0].Token)
	}
	if def.Providers[0].Cleanup == nil {
		t.Fatal("expected cleanup hook")
	}
}
