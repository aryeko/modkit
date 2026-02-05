package lifecycle

import (
	"context"
	"errors"
	"net/http"
)

// CleanupHook defines a shutdown cleanup function.
type CleanupHook func(ctx context.Context) error

// RunCleanup executes hooks in LIFO order and returns any combined errors.
func RunCleanup(ctx context.Context, hooks []CleanupHook) error {
	var joined error
	for i := len(hooks) - 1; i >= 0; i-- {
		if hooks[i] == nil {
			continue
		}
		if err := hooks[i](ctx); err != nil {
			joined = errors.Join(joined, err)
		}
	}
	return joined
}

// ShutdownServer shuts down the server, then runs cleanup hooks.
func ShutdownServer(ctx context.Context, server *http.Server, hooks []CleanupHook) error {
	shutdownErr := server.Shutdown(ctx)
	cleanupErr := RunCleanup(ctx, hooks)
	if shutdownErr != nil {
		return shutdownErr
	}
	return cleanupErr
}
