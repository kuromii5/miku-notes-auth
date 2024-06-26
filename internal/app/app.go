package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/kuromii5/sso-auth/internal/app/grpc"
	postgres "github.com/kuromii5/sso-auth/internal/db"
	"github.com/kuromii5/sso-auth/internal/service"
)

type App struct {
	Server *grpcapp.GRPCApp
}

func New(
	log *slog.Logger,
	port int,
	dbPath string,
	secret string,
	tokenTTL time.Duration,
) *App {
	db, err := postgres.New(dbPath)
	if err != nil {
		panic(err)
	}

	authService := service.New(log, db, db, secret, tokenTTL)
	app := grpcapp.New(log, port, authService)

	return &App{Server: app}
}
