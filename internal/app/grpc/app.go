package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/kuromii5/sso-auth/internal/auth"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log    *slog.Logger
	server *grpc.Server
	port   int
}

func New(log *slog.Logger, port int, authGRPC auth.Auth) *GRPCApp {
	server := grpc.NewServer()

	auth.Register(server, authGRPC)

	return &GRPCApp{
		log:    log,
		server: server,
		port:   port,
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
