package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/kuromii5/miku-notes-auth/internal/app/grpc"
	"github.com/kuromii5/miku-notes-auth/internal/repo/postgres"
	"github.com/kuromii5/miku-notes-auth/internal/repo/redis"
	"github.com/kuromii5/miku-notes-auth/internal/service"
	"github.com/kuromii5/miku-notes-auth/internal/service/tokens"
)

type App struct {
	Server *grpcapp.GRPCApp
}

func New(
	log *slog.Logger,
	port int,
	connToken string,
	dbPath string,
	secret string,
	redisAddr string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *App {
	db, err := postgres.New(dbPath)
	if err != nil {
		panic(err)
	}

	// define refresh token storage and manager
	// it's just part of authService
	tokenStorage := redis.New(redisAddr)
	tokenManager := tokens.New(log, secret, accessTTL, refreshTTL, tokenStorage, tokenStorage, tokenStorage)

	authService := service.New(log, db, db, tokenManager)
	app := grpcapp.New(log, port, connToken, authService)

	return &App{Server: app}
}
