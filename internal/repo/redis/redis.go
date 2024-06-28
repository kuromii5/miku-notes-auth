package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrTokenNotFound = errors.New("refresh token for user not found")

type TokenStorage struct {
	client *redis.Client
}

func New(addr string) *TokenStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return &TokenStorage{client: rdb}
}

func (t *TokenStorage) Set(ctx context.Context, userID int32, token string, expires time.Duration) error {
	const f = "redis.Set"

	userIdStr := fmt.Sprintf("%d", userID)
	if err := t.client.Set(ctx, token, userIdStr, expires).Err(); err != nil {
		return fmt.Errorf("%s:%w", f, err)
	}

	return nil
}

func (t *TokenStorage) UserID(ctx context.Context, token string) (string, error) {
	const f = "redis.UserID"

	userIdStr, err := t.client.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("%s:%w", f, ErrTokenNotFound)
		}

		return "", fmt.Errorf("%s:%w", f, err)
	}

	return userIdStr, nil
}

func (t *TokenStorage) Delete(ctx context.Context, userID int32, token string) error {
	const f = "redis.Delete"

	key := fmt.Sprintf("%d", userID)
	if err := t.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("%s:%w", f, err)
	}

	return nil
}
