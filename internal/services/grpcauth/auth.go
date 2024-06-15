package grpcauth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	postgres "github.com/kuromii5/sso-auth/internal/db"
	"github.com/kuromii5/sso-auth/internal/lib/jwt"
	l "github.com/kuromii5/sso-auth/internal/lib/logger/err"
	"github.com/kuromii5/sso-auth/internal/models"
	"golang.org/x/crypto/bcrypt"
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
	secret       string
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, hash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	secret string,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		secret:       secret,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const f = "auth.RegisterNewUser"

	log := a.log.With(slog.String("func", f))
	log.Info("registering new user")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password", l.Err(err))
		return 0, fmt.Errorf("%s:%v", f, err)
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

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const f = "auth.Login"

	log := a.log.With(slog.String("func", f))
	log.Info("trying to log in user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			a.log.Warn("user not found", l.Err(err))
			return "", fmt.Errorf("%s:%w", f, ErrInvalidCreds)
		}

		a.log.Error("failed to get user", l.Err(err))
		return "", fmt.Errorf("%s:%w", f, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Warn("invalid credentials", l.Err(err))
		return "", fmt.Errorf("%s:%w", f, ErrInvalidCreds)
	}

	token, err := jwt.NewJWT(user, a.tokenTTL, a.secret)
	if err != nil {
		a.log.Error("failed to generate jwt", l.Err(err))
		return "", fmt.Errorf("%s:%w", f, err)
	}

	log.Info("user logged in successfully")

	return token, nil
}

func (a *Auth) IsAdmin(ctx context.Context, uid int64) (bool, error) {
	const f = "auth.IsAdmin"

	log := a.log.With(slog.Int64("user_id", uid))
	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, uid)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			a.log.Warn("user not found", l.Err(err))
			return false, fmt.Errorf("%s:%w", f, ErrInvalidCreds)
		}
		return false, fmt.Errorf("%s:%w", f, err)
	}

	log.Info("successful if user is admin check", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
