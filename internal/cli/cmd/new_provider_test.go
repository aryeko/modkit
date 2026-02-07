package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func captureStdout(t *testing.T, fn func() error) (out string, errRun error) {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = w
	var panicV any

	defer func() {
		os.Stdout = orig
		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		out = buf.String()
		_ = r.Close()
		if panicV != nil {
			panic(panicV)
		}
	}()

	func() {
		defer func() {
			panicV = recover()
		}()
		errRun = fn()
	}()

	return out, errRun
}

func TestCreateNewProvider(t *testing.T) {
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

	out, err := captureStdout(t, func() error {
		return createNewProvider("auth", "user-service")
	})
	if err != nil {
		t.Fatalf("createNewProvider failed: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(moduleDir, "auth.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "package userservice") {
		t.Fatalf("expected sanitized package in provider file, got:\n%s", string(b))
	}
	if !strings.Contains(out, `Token: "userservice.auth"`) {
		t.Fatalf("expected module.component token in output, got:\n%s", out)
	}
}

func TestCreateNewProviderFromCurrentModuleDir(t *testing.T) {
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

	if _, err := captureStdout(t, func() error {
		return createNewProvider("cache", "")
	}); err != nil {
		t.Fatalf("createNewProvider from current dir failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(moduleDir, "cache.go")); err != nil {
		t.Fatalf("expected provider file in current module dir, got %v", err)
	}
}

func TestCreateNewProviderInvalidName(t *testing.T) {
	if err := createNewProvider("../evil", "users"); err == nil {
		t.Fatal("expected error for invalid provider name")
	}
}

func TestCreateNewProviderInvalidModuleName(t *testing.T) {
	if err := createNewProvider("auth", "../users"); err == nil {
		t.Fatal("expected error for invalid module name")
	}
}

func TestCreateNewProviderMissingModuleFile(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := createNewProvider("auth", "users"); err == nil {
		t.Fatal("expected error when module file is missing")
	}
}

func TestCreateNewProviderAlreadyExists(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(moduleDir, "auth.go"), []byte("package users\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := createNewProvider("auth", "users"); err == nil {
		t.Fatal("expected error when provider file already exists")
	}
}

func TestCreateNewProviderGetwdFailure(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(tmp); err != nil {
		t.Fatal(err)
	}

	if err := createNewProvider("auth", ""); err == nil {
		t.Fatal("expected error when cwd cannot be resolved")
	}
}

func TestCreateNewProviderRunE(t *testing.T) {
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

	cmd := &cobra.Command{}
	cmd.Flags().String("module", "", "")
	if err := cmd.Flags().Set("module", "users"); err != nil {
		t.Fatal(err)
	}
	if err := newProviderCmd.RunE(cmd, []string{"billing"}); err != nil {
		t.Fatalf("RunE failed: %v", err)
	}
}

func TestCreateNewProviderCreateFileFailure(t *testing.T) {
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
	if err := os.Chmod(moduleDir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(moduleDir, 0o750) })

	if err := createNewProvider("auth", "users"); err == nil {
		t.Fatal("expected error when provider file cannot be created")
	}
}
