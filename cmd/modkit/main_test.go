package main

import (
	"os"
	"testing"
)

func TestRunHelp(t *testing.T) {
	if code := run([]string{"modkit", "--help"}); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRunInvalidCommand(t *testing.T) {
	if code := run([]string{"modkit", "nope"}); code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestMainInvokesExit(t *testing.T) {
	orig := osExit
	origArgs := os.Args
	t.Cleanup(func() {
		osExit = orig
		os.Args = origArgs
	})

	called := false
	got := -1
	osExit = func(code int) {
		called = true
		got = code
	}

	os.Args = []string{"modkit", "--help"}
	main()

	if !called {
		t.Fatal("expected main to call osExit")
	}
	if got != 0 {
		t.Fatalf("expected exit code 0, got %d", got)
	}
}
