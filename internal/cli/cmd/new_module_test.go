package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
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

	b, err := os.ReadFile(filepath.Join(tmp, "internal", "modules", "user-service", "module.go")) //nolint:gosec
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

func TestCreateNewModuleAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	modulePath := filepath.Join(tmp, "internal", "modules", "users", "module.go")
	if err := os.MkdirAll(filepath.Dir(modulePath), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(modulePath, []byte("package users\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := createNewModule("users"); err == nil {
		t.Fatal("expected error when module file already exists")
	}
}

func TestCreateNewModuleMkdirFail(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission semantics differ on windows")
	}

	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmp, "internal"), 0o500); err != nil {
		t.Fatal(err)
	}

	if err := createNewModule("users"); err == nil {
		t.Fatal("expected error when destination directory cannot be created")
	}
}

func TestCreateNewModuleRunE(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := newModuleCmd.RunE(&cobra.Command{}, []string{"billing"}); err != nil {
		t.Fatalf("RunE failed: %v", err)
	}
}
