package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCreateNewApp(t *testing.T) {
	tmp := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0o750); err != nil {
		t.Fatal(err)
	}
	shim := filepath.Join(binDir, "go")
	content := "#!/bin/sh\nexit 0\n"
	if runtime.GOOS == "windows" {
		shim = filepath.Join(binDir, "go.bat")
		content = "@echo off\r\nexit /b 0\r\n"
	}
	if err := os.WriteFile(shim, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })
	if err := os.Setenv("PATH", binDir+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatal(err)
	}

	if err := createNewApp("demo"); err != nil {
		t.Fatalf("createNewApp failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmp, "demo", "go.mod")); err != nil {
		t.Fatalf("expected go.mod, got %v", err)
	}
}

func TestCreateNewAppInvalidName(t *testing.T) {
	if err := createNewApp("../bad"); err == nil {
		t.Fatal("expected error for invalid app name")
	}
}
