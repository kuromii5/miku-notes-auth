package tokens

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("")
)

type TokenManager struct {
	log *slog.Logger

	accessTTL  time.Duration
	refreshTTL time.Duration
	secret     string

	refreshTokenSetter  RefreshTokenSetter
	refreshTokenDeleter RefreshTokenDeleter
	userGetter          UserGetter
}

// Redis methods that are being used here
type RefreshTokenSetter interface {
	Set(ctx context.Context, userID int64, token string, expires time.Duration) error
}
type RefreshTokenDeleter interface {
	Delete(ctx context.Context, userID int64, token string) error
}
type UserGetter interface {
	UserID(ctx context.Context, token string) (string, error)
}

func New(
	log *slog.Logger,
	secret string,
	accessTTL, refreshTTL time.Duration,
	refreshTokenSetter RefreshTokenSetter,
	refreshTokenDeleter RefreshTokenDeleter,
	userGetter UserGetter,
) *TokenManager {
	return &TokenManager{
		log:                 log,
		accessTTL:           accessTTL,
		refreshTTL:          refreshTTL,
		secret:              secret,
		refreshTokenSetter:  refreshTokenSetter,
		refreshTokenDeleter: refreshTokenDeleter,
		userGetter:          userGetter,
	}
}

func (t *TokenManager) CreateAccessToken(_ context.Context, userID int64) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(t.accessTTL).Unix(),
	})

	return jwtToken.SignedString([]byte(t.secret))
}

func (t *TokenManager) NewRefreshToken(ctx context.Context, userID int64) (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	refreshToken := base64.URLEncoding.EncodeToString(b)

	// save token
	if err = t.refreshTokenSetter.Set(ctx, userID, refreshToken, t.refreshTTL); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (t *TokenManager) ValidateRefreshToken(ctx context.Context, token string) (int64, error) {
	userIDStr, err := t.userGetter.UserID(ctx, token)
	if err != nil {
		return 0, err
	}

	// convert string to int64
	id, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *TokenManager) ValidateAccessToken(ctx context.Context, token string) (int64, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secret), nil
	}

	accessToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, keyFunc)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := accessToken.Claims.(*jwt.StandardClaims)
	if !ok || !accessToken.Valid {
		return 0, fmt.Errorf("invalid token claims")
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return userID, nil
}
