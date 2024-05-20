package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

func NewLogger() Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
