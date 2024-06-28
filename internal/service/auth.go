package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kuromii5/sso-auth/internal/models"
	"github.com/kuromii5/sso-auth/internal/repo/postgres"
	"github.com/kuromii5/sso-auth/internal/service/tokens"
	"github.com/kuromii5/sso-auth/pkg/hasher"
	l "github.com/kuromii5/sso-auth/pkg/logger/err"
)

var (
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	tokenManager *tokens.TokenManager
}

// Postgres DB methods
type UserSaver interface {
	SaveUser(ctx context.Context, email string, hash []byte) (int32, error)
}
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenManager *tokens.TokenManager,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenManager: tokenManager,
	}
}

func (a *Auth) Register(ctx context.Context, email, password string) (int32, error) {
	const f = "auth.Register"

	log := a.log.With(slog.String("func", f))
	log.Info("registering new user")

	hash, err := hasher.HashPassword(password)
	if err != nil {
		log.Error("failed to generate password", l.Err(err))

		return 0, fmt.Errorf("%s:%w", f, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, hash)
	if err != nil {
		if errors.Is(err, postgres.ErrUserExists) {
			a.log.Warn("user already exists", l.Err(err))

			return 0, fmt.Errorf("%s:%w", f, ErrUserExists)
		}

		log.Error("failed to save user", l.Err(err))
		return 0, fmt.Errorf("%s:%v", f, err)
	}

	log.Info("successfully registered new user")

	return id, nil
}

func (a *Auth) Login(ctx context.Context, email, password, fingerprint string) (models.TokenPair, error) {
	const f = "auth.Login"

	log := a.log.With(slog.String("func", f))
	log.Info("trying to log in user")

	// get the user from db
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			a.log.Warn("user not found", l.Err(err))

			return models.TokenPair{}, fmt.Errorf("%s:%w", f, ErrInvalidCreds)
		}

		a.log.Error("failed to get user", l.Err(err))
		return models.TokenPair{}, fmt.Errorf("%s:%w", f, err)
	}

	// check password
	if err := hasher.CheckPassword(password, user.PasswordHash); err != nil {
		a.log.Warn("invalid credentials", l.Err(err))

		return models.TokenPair{}, fmt.Errorf("%s:%w", f, ErrInvalidCreds)
	}

	// generate new access token
	accessToken, err := a.tokenManager.NewAccessToken(ctx, user.ID)
	if err != nil {
		a.log.Error("failed to generate jwt access token", l.Err(err))

		return models.TokenPair{}, fmt.Errorf("%s:%w", f, err)
	}

	// generate new refresh token
	refreshToken, err := a.tokenManager.NewRefreshToken(ctx, user.ID, fingerprint)
	if err != nil {
		a.log.Error("failed to generate refresh token", l.Err(err))

		return models.TokenPair{}, fmt.Errorf("%s:%w", f, err)
	}

	log.Info("user logged in successfully")

	return models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) GetAccessToken(ctx context.Context, refreshToken, fingerprint string) (string, error) {
	const f = "service.GetAccessToken"

	log := a.log.With(slog.String("func", f))
	log.Info("attempting to generate new access token using refresh token")

	// Validate the refresh token
	userID, err := a.tokenManager.ValidateRefreshToken(ctx, refreshToken, fingerprint)
	if err != nil {
		log.Error("failed to validate refresh token", l.Err(err))

		return "", fmt.Errorf("%s:%w", f, err)
	}

	// Generate the access token
	accessToken, err := a.tokenManager.NewAccessToken(ctx, userID)
	if err != nil {
		log.Error("failed to create access token", l.Err(err))

		return "", fmt.Errorf("%s:%w", f, err)
	}

	log.Info("successfully generated new access token", slog.Int("user_id", int(userID)))

	return accessToken, nil
}

func (a *Auth) ValidateAccessToken(ctx context.Context, token string) (int32, error) {
	const f = "service.ValidateAccessToken"

	log := a.log.With(slog.String("func", f))
	log.Info("validating access token")

	// Validate the access token
	userID, err := a.tokenManager.ValidateAccessToken(ctx, token)
	if err != nil {
		log.Warn("failed to validate access token", l.Err(err))

		return 0, fmt.Errorf("%s:%w", f, err)
	}

	log.Info("access token validated successfully", slog.Int("user_id", int(userID)))

	return userID, nil
}

func (a *Auth) Logout(ctx context.Context, accessToken, fingerprint string) error {
	const f = "service.Logout"

	log := a.log.With(slog.String("func", f))
	log.Info("logging out user")

	// Validate the access token to get user id
	userID, err := a.tokenManager.ValidateAccessToken(ctx, accessToken)
	if err != nil {
		log.Warn("failed to validate access token", l.Err(err))

		return fmt.Errorf("%s:%w", f, err)
	}

	if err = a.tokenManager.Delete(ctx, userID, fingerprint); err != nil {
		log.Error("internal error", l.Err(err))

		return fmt.Errorf("%s:%w", f, err)
	}

	log.Info("successfully logged out user", slog.Int("user_id", int(userID)))

	return nil
}
