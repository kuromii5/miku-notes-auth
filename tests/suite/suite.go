package suite

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	sso "github.com/kuromii5/sso-auth/generated"
	"github.com/kuromii5/sso-auth/internal/config"
	"github.com/kuromii5/sso-auth/internal/service/tokens"
	mock_tokens "github.com/kuromii5/sso-auth/internal/service/tokens/mock"
	"github.com/kuromii5/sso-auth/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Suite struct {
	*testing.T
	Cfg          *config.Config
	AuthClient   sso.AuthClient
	TokenManager *tokens.TokenManager
	Mocks        *Mocks
}
type Mocks struct {
	RefreshTokenSetter  *mock_tokens.MockRefreshTokenSetter
	RefreshTokenDeleter *mock_tokens.MockRefreshTokenDeleter
	UserGetter          *mock_tokens.MockUserGetter
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	ctrl := gomock.NewController(t)

	cfg := config.ReadConfig("../config/local_test.yaml")
	log := logger.New(cfg.Env)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	cc, err := grpc.NewClient(
		net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPC.Port)),
		opts...,
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	mockRefreshTokenSetter := mock_tokens.NewMockRefreshTokenSetter(ctrl)
	mockRefreshTokenDeleter := mock_tokens.NewMockRefreshTokenDeleter(ctrl)
	mockUserGetter := mock_tokens.NewMockUserGetter(ctrl)

	tokenManager := tokens.New(
		log,
		cfg.Tokens.Secret,
		cfg.Tokens.AccessTTL,
		cfg.Tokens.RefreshTTL,
		mockRefreshTokenSetter,
		mockRefreshTokenDeleter,
		mockUserGetter,
	)

	// Add the microservice authorization token to the context
	md := metadata.Pairs("authorization", "Bearer "+cfg.GRPC.ConnectionToken)
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	return ctxWithToken, &Suite{
		T:            t,
		Cfg:          cfg,
		AuthClient:   sso.NewAuthClient(cc),
		TokenManager: tokenManager,
		Mocks: &Mocks{
			RefreshTokenSetter:  mockRefreshTokenSetter,
			RefreshTokenDeleter: mockRefreshTokenDeleter,
			UserGetter:          mockUserGetter,
		},
	}
}
