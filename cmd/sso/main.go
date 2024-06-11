package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kuromii5/sso-auth/internal/app"
	config "github.com/kuromii5/sso-auth/internal/cfg"
	"github.com/kuromii5/sso-auth/internal/lib/logger/handler"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()
	connStr := cfg.Postgres.ConnString()

	log := setupLogger(cfg.Env)
	log.Info("Starting application", slog.Any("config", cfg))

	app := app.New(log, cfg.GRPC.Port, connStr, cfg.JWT_SECRET, cfg.TokenTTL)

	go app.Server.MustRun()

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	s := <-shutdown
	log.Info("shutdown", slog.String("signal", s.String()))

	app.Server.Shutdown()
	log.Info("Server is stopped")
}

func setupLogger(env string) *slog.Logger {
	var h *handler.PrettyHandler
	var levelDebug slog.Level = -4
	var levelInfo slog.Level = 0

	switch env {
	case envLocal:
		h = handler.New(os.Stdout, &levelDebug)
	case envDev:
		h = handler.New(os.Stdout, &levelDebug)
	case envProd:
		h = handler.New(os.Stdout, &levelInfo)
	}

	return slog.New(h)
}
