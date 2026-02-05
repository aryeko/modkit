package lifecycle

import (
	"context"
	"net"
	"net/http"
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

func TestShutdown_WaitsForInFlightRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	cleanupCalled := make(chan struct{})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusOK)
		close(done)
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	server := &http.Server{Handler: handler}
	go func() {
		_ = server.Serve(ln)
	}()

	reqDone := make(chan struct{})
	go func() {
		_, _ = http.Get("http://" + ln.Addr().String())
		close(reqDone)
	}()

	<-started

	hooks := []CleanupHook{
		func(ctx context.Context) error {
			close(cleanupCalled)
			return nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- ShutdownServer(ctx, server, hooks)
	}()

	select {
	case <-cleanupCalled:
		t.Fatal("cleanup ran before in-flight request completed")
	default:
	}

	close(release)
	<-done
	<-reqDone

	select {
	case <-cleanupCalled:
	case <-time.After(time.Second):
		t.Fatal("cleanup did not run after in-flight request completed")
	}

	if err := <-shutdownDone; err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}
}
