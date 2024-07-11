package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kuromii5/miku-notes-auth/internal/app"
	"github.com/kuromii5/miku-notes-auth/internal/config"
	"github.com/kuromii5/miku-notes-auth/pkg/logger"
)

func main() {
	// read config file
	cfg := config.MustLoad()

	// make connection string for postgres
	postgresConnStr := cfg.Postgres.ConnString()

	// setup logger for logs
	log := logger.New(cfg.Env)

	log.Info("Starting application", slog.Any("config", cfg))

	// initialize application
	app := app.New(
		log,
		cfg.GRPC.Port,
		cfg.GRPC.ConnectionToken,
		postgresConnStr,
		cfg.Tokens.Secret,
		cfg.Tokens.RedisAddr,
		cfg.Tokens.AccessTTL,
		cfg.Tokens.RefreshTTL,
	)

	// run the server as goroutine
	go app.Server.MustRun()

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT) // terminate, interrupt
	s := <-shutdown
	log.Info("shutdown", slog.String("signal", s.String()))

	app.Server.Shutdown()
	log.Info("Server is stopped")
}
