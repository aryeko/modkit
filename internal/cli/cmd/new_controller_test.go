package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateNewController(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	moduleDir := filepath.Join(tmp, "internal", "modules", "user-service")
	if err := os.MkdirAll(moduleDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "module.go"), []byte("package userservice\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := createNewController("auth", "user-service"); err != nil {
		t.Fatalf("createNewController failed: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(moduleDir, "auth_controller.go"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, "package userservice") {
		t.Fatalf("expected sanitized package in controller file, got:\n%s", s)
	}
	if !strings.Contains(s, "type AuthController struct{}") {
		t.Fatalf("expected identifier-based controller type, got:\n%s", s)
	}
}

func TestCreateNewControllerFromCurrentModuleDir(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })

	moduleDir := filepath.Join(tmp, "internal", "modules", "users")
	if err := os.MkdirAll(moduleDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "module.go"), []byte("package users\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(moduleDir); err != nil {
		t.Fatal(err)
	}

	if err := createNewController("cache", ""); err != nil {
		t.Fatalf("createNewController from current dir failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(moduleDir, "cache_controller.go")); err != nil {
		t.Fatalf("expected controller file in current module dir, got %v", err)
	}
}

func TestCreateNewControllerInvalidName(t *testing.T) {
	if err := createNewController("bad/name", "users"); err == nil {
		t.Fatal("expected error for invalid controller name")
	}
}

func TestCreateNewControllerInvalidModuleName(t *testing.T) {
	if err := createNewController("auth", "../users"); err == nil {
		t.Fatal("expected error for invalid module name")
	}
}

func TestCreateNewControllerMissingModuleFile(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := createNewController("auth", "users"); err == nil {
		t.Fatal("expected error when module file is missing")
	}
}

func TestCreateNewControllerAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	moduleDir := filepath.Join(tmp, "internal", "modules", "users")
	if err := os.MkdirAll(moduleDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "module.go"), []byte("package users\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "auth_controller.go"), []byte("package users\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := createNewController("auth", "users"); err == nil {
		t.Fatal("expected error when controller file already exists")
	}
}
