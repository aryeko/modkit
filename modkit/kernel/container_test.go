package kernel_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aryeko/modkit/modkit/kernel"
	"github.com/aryeko/modkit/modkit/module"
)

func TestContainerMissingTokenError(t *testing.T) {
	modA := mod("A", nil, nil, nil, nil)

	app, err := kernel.Bootstrap(modA)
	if err != nil {
		t.Fatalf("Bootstrap failed: %v", err)
	}

	_, err = app.Container.Get(module.Token("missing"))
	if err == nil {
		t.Fatalf("expected missing token error")
	}

	var notFound *kernel.ProviderNotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("unexpected error type: %T", err)
	}
	if notFound.Token != module.Token("missing") {
		t.Fatalf("unexpected token: %q", notFound.Token)
	}
}

func TestContainerSingletonConcurrent(t *testing.T) {
	shared := module.Token("shared")
	var calls int32
	started := make(chan struct{})
	release := make(chan struct{})

	modA := mod("A", nil,
		[]module.ProviderDef{{
			Token: shared,
			Build: func(r module.Resolver) (any, error) {
				if atomic.AddInt32(&calls, 1) == 1 {
					close(started)
				}
				<-release
				return "value", nil
			},
		}},
		nil,
		nil,
	)

	app, err := kernel.Bootstrap(modA)
	if err != nil {
		t.Fatalf("Bootstrap failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = app.Container.Get(shared)
		}()
	}

	<-started
	close(release)
	wg.Wait()

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected single build call, got %d", got)
	}
}
