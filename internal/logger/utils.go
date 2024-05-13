package logger

import "log/slog"

func AsError(v error) slog.Attr {
	return slog.Any("err", v)
}
