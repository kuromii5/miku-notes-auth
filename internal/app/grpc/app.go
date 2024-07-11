package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/kuromii5/miku-notes-auth/internal/auth"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log             *slog.Logger
	server          *grpc.Server
	port            int
	connectionToken string
}

func New(log *slog.Logger, port int, connectionToken string, authGRPC auth.Auth) *GRPCApp {
	server := auth.RegisterServer(authGRPC, connectionToken)

	return &GRPCApp{
		log:             log,
		server:          server,
		port:            port,
		connectionToken: connectionToken,
	}
}

func (a *GRPCApp) run() error {
	const f = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", f, err)
	}

	a.log.Info("Starting gRPC server",
		slog.Int("port", a.port),
		slog.String("func", f),
		slog.String("addr", l.Addr().String()),
		slog.String("connection token", a.connectionToken),
	)

	if err := a.server.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", f, err)
	}

	return nil
}

func (a *GRPCApp) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *GRPCApp) Shutdown() {
	const f = "grpcapp.Stop"

	a.log.Info("Stopping gRPC server",
		slog.String("f", f),
	)

	a.server.GracefulStop()
}
