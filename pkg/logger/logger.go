package logger

import (
	"os"

	"log/slog"

	offlog "github.com/kuromii5/miku-notes-auth/pkg/logger/off"
	prettylog "github.com/kuromii5/miku-notes-auth/pkg/logger/pretty"
)

var (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

// TODO: refactor json handler output
func New(env string) *slog.Logger {
	switch env {
	case local:
		return prettylog.NewTextLogger(os.Stdout, slog.LevelDebug)
	case dev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prod:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return offlog.New()
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
