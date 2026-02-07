package ast

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddProvider(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "module.go")
	original := `package users

import "github.com/go-modkit/modkit/modkit/module"

type Module struct{}

func (m *Module) Definition() module.ModuleDef {
	return module.ModuleDef{
		Name: "users",
		Providers: []module.ProviderDef{},
	}
}
`
	if err := os.WriteFile(file, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := AddProvider(file, "users.auth", "buildAuth"); err != nil {
		t.Fatalf("AddProvider failed: %v", err)
	}

	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `Token: "users.auth"`) {
		t.Fatalf("expected token in output:\n%s", s)
	}
	if !strings.Contains(s, `Build: buildAuth`) {
		t.Fatalf("expected build func in output:\n%s", s)
	}
}

func TestAddProviderNoProvidersField(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "module.go")
	content := `package users

import "github.com/go-modkit/modkit/modkit/module"

type Module struct{}

func (m *Module) Definition() module.ModuleDef {
	return module.ModuleDef{Name: "users"}
}
`
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := AddProvider(file, "users.auth", "buildAuth"); err == nil {
		t.Fatal("expected error when Providers field is missing")
	}
}

func TestAddProviderParseError(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "module.go")
	if err := os.WriteFile(file, []byte("package users\nfunc ("), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := AddProvider(file, "users.auth", "buildAuth"); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestAddProviderNoDefinitionMethod(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "module.go")
	content := `package users

type Module struct{}
`
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := AddProvider(file, "users.auth", "buildAuth"); err == nil {
		t.Fatal("expected error when Definition method is missing")
	}
}
