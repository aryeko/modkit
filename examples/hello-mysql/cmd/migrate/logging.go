package main

import (
	"log/slog"

	modkitlogging "github.com/aryeko/modkit/modkit/logging"
)

func logMigrateComplete(logger modkitlogging.Logger) {
	if logger == nil {
		logger = modkitlogging.Nop()
	}
	logger.Info("migrations complete", slog.String("component", "migrate"))
}
