package main

import "testing"

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
	t.Cleanup(func() { osExit = orig })

	called := false
	got := -1
	osExit = func(code int) {
		called = true
		got = code
	}

	main()

	if !called {
		t.Fatal("expected main to call osExit")
	}
	if got != 0 {
		t.Fatalf("expected exit code 0 from test binary args, got %d", got)
	}
}
