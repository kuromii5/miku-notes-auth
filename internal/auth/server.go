package auth

import (
	"context"
	"errors"

	sso "github.com/kuromii5/sso-auth/generated"
	"github.com/kuromii5/sso-auth/internal/models"
	"github.com/kuromii5/sso-auth/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth            Auth
	connectionToken string
}

//go:generate mockgen -source=server.go -destination=mock/server.go
type Auth interface {
	Register(ctx context.Context, email, password string) (int32, error)
	Login(ctx context.Context, email, password, fingerprint string) (models.TokenPair, error)
	GetAccessToken(ctx context.Context, refreshToken, fingerprint string) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (int32, error)
	Logout(ctx context.Context, accessToken, fingerprint string) error
}

func RegisterServer(auth Auth, connectionToken string) *grpc.Server {
	server := &serverAPI{auth: auth, connectionToken: connectionToken}

	// Register the interceptor with the gRPC server
	interceptor := grpc.UnaryInterceptor(server.validateBearerTokenInterceptor)
	gRPC := grpc.NewServer(interceptor)

	sso.RegisterAuthServer(gRPC, server)

	return gRPC
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.AuthResponse, error) {
	// validate given data
	err := validateRegisterRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// register user in the system
	_, err = s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal register error")
	}

	// automatically log in after register
	tokens, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetFingerprint())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal login error")
	}

	return &sso.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.AuthResponse, error) {
	// validate given data
	err := validateLoginRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// get the pair of tokens: access and refresh
	tokens, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetFingerprint())
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal login error")
	}

	return &sso.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverAPI) GetAccessToken(ctx context.Context, req *sso.GetATRequest) (*sso.GetATResponse, error) {
	// get the access token
	accessToken, err := s.auth.GetAccessToken(ctx, req.GetRefreshToken(), req.GetFingerprint())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	return &sso.GetATResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *serverAPI) ValidateAccessToken(ctx context.Context, req *sso.ValidateATRequest) (*sso.ValidateATResponse, error) {
	// validate access token
	userID, err := s.auth.ValidateAccessToken(ctx, req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &sso.ValidateATResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *sso.LogoutRequest) (*sso.LogoutResponse, error) {
	// log out
	if err := s.auth.Logout(ctx, req.GetAccessToken(), req.GetFingerprint()); err != nil {
		return nil, status.Error(codes.Internal, "failed to log out")
	}

	return &sso.LogoutResponse{}, nil
}
