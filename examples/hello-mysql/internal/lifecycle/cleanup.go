package lifecycle

import (
	"context"
	"errors"
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
