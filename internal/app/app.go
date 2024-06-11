package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/kuromii5/sso-auth/internal/app/grpc"
	postgres "github.com/kuromii5/sso-auth/internal/db"
	grpcauth "github.com/kuromii5/sso-auth/internal/services/grpcauth"
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

	authService := grpcauth.New(log, db, db, secret, tokenTTL)
	app := grpcapp.New(log, port, authService)

	return &App{Server: app}
}
