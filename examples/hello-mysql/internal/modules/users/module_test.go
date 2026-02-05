package users

import (
	"testing"

	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/auth"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/database"
)

func TestUsersModule_Definition_WiresAuth(t *testing.T) {
	mod := NewModule(Options{Database: &database.Module{}, Auth: auth.NewModule(auth.Options{})})
	def := mod.(Module).Definition()

	if def.Name != "users" {
		t.Fatalf("name = %q", def.Name)
	}
	if len(def.Imports) != 2 {
		t.Fatalf("imports = %d", len(def.Imports))
	}
}
