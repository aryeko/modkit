package logging

import (
	"log/slog"
	"os"

	modkitlogging "github.com/aryeko/modkit/modkit/logging"
)

func New() modkitlogging.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return modkitlogging.NewSlog(slog.New(handler))
}
