package audit

import (
	"testing"

	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/users"
)

func TestAuditModule_Definition_WiresUsersImport(t *testing.T) {
	usersMod := &users.Module{}
	mod := NewModule(Options{Users: usersMod})
	def := mod.(*Module).Definition()
	if def.Name != "audit" {
		t.Fatalf("expected name audit, got %q", def.Name)
	}
	if len(def.Imports) != 1 {
		t.Fatalf("expected 1 import, got %d", len(def.Imports))
	}
	if def.Imports[0].Definition().Name != "users" {
		t.Fatalf("expected users import, got %q", def.Imports[0].Definition().Name)
	}
}
