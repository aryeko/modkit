package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	errRun := fn()
	_ = w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String(), errRun
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

func TestCreateNewProviderInvalidName(t *testing.T) {
	if err := createNewProvider("../evil", "users"); err == nil {
		t.Fatal("expected error for invalid provider name")
	}
}
