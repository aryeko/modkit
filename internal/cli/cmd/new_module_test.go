package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateNewModule(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := createNewModule("user-service"); err != nil {
		t.Fatalf("createNewModule failed: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(tmp, "internal", "modules", "user-service", "module.go"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, "package userservice") {
		t.Fatalf("expected sanitized package, got:\n%s", s)
	}
	if !strings.Contains(s, "type UserServiceModule struct{}") {
		t.Fatalf("expected exported identifier, got:\n%s", s)
	}
}

func TestCreateNewModuleInvalidName(t *testing.T) {
	if err := createNewModule("../evil"); err == nil {
		t.Fatal("expected error for invalid name")
	}
}
