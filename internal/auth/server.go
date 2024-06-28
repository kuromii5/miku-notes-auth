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
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email string, password string) (models.TokenPair, error)
	Register(ctx context.Context, email string, password string) (int32, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	ValidateAccessToken(ctx context.Context, token string) (int32, error)
}

func RegisterServer(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.AuthResponse, error) {
	err := validateRegisterRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal register error")
	}

	// automatically log in after register
	tokens, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal login error")
	}

	return &sso.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.AuthResponse, error) {
	err := validateLoginRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// get the pair of tokens: access and refresh
	tokens, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
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
	accessToken, err := s.auth.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	return &sso.GetATResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *serverAPI) ValidateAccessToken(ctx context.Context, req *sso.ValidateATRequest) (*sso.ValidateATResponse, error) {
	userID, err := s.auth.ValidateAccessToken(ctx, req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &sso.ValidateATResponse{
		UserId: userID,
	}, nil
}
