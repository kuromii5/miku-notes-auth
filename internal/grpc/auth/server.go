package auth

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	ssov1 "github.com/kuromii5/proto-auth/gen/go/sso"
	grpcauth "github.com/kuromii5/sso-auth/internal/services/grpcauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	validate *validator.Validate
	auth     Auth
}

type Auth interface {
	Login(ctx context.Context, email string, password string) (string, error)
	RegisterNewUser(ctx context.Context, email string, password string) (int64, error)
}

func Register(gRPC *grpc.Server, auth Auth) {
	validate := validator.New()
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth, validate: validate})
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	err := s.validateRegisterRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, grpcauth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal register error")
	}

	// automatically log in after register
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, grpcauth.ErrInvalidCreds) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal login error")
	}

	return &ssov1.RegisterResponse{Token: token}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	err := s.validateLoginRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())

	// TODO: add

	if err != nil {
		if errors.Is(err, grpcauth.ErrInvalidCreds) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal login error")
	}

	// Set the JWT token in the metadata
	md := metadata.Pairs("authorization", "Bearer "+token)
	grpc.SetHeader(ctx, md)

	return &ssov1.LoginResponse{}, nil
}
