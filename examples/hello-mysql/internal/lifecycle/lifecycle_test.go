package lifecycle

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestShutdown_InvokesCleanupHooksInLIFO(t *testing.T) {
	calls := make([]string, 0, 2)
	hooks := []CleanupHook{
		func(ctx context.Context) error {
			calls = append(calls, "first")
			return nil
		},
		func(ctx context.Context) error {
			calls = append(calls, "second")
			return nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := RunCleanup(ctx, hooks); err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if got, want := strings.Join(calls, ","), "second,first"; got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
