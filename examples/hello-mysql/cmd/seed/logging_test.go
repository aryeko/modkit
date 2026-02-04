package main

import (
	"log/slog"
	"testing"

	modkitlogging "github.com/aryeko/modkit/modkit/logging"
)

type captureLogger struct {
	messages []string
}

func (c *captureLogger) Debug(string, ...slog.Attr) {}
func (c *captureLogger) Info(msg string, _ ...slog.Attr) {
	c.messages = append(c.messages, msg)
}
func (c *captureLogger) Error(string, ...slog.Attr) {}
func (c *captureLogger) With(...slog.Attr) modkitlogging.Logger { return c }

func TestLogSeedComplete(t *testing.T) {
	logger := &captureLogger{}
	logSeedComplete(logger)

	if len(logger.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(logger.messages))
	}
	if logger.messages[0] != "seed complete" {
		t.Fatalf("unexpected message: %s", logger.messages[0])
	}
}
