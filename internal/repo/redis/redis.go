package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("user not found")

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

func (t *TokenStorage) Set(ctx context.Context, userID int64, token string, expires time.Duration) error {
	const f = "redis.Set"

	value := fmt.Sprintf("%d", userID)

	if err := t.client.Set(ctx, token, value, expires).Err(); err != nil {
		log.Printf("Couldn't set refresh token for uid %d, token:%s\n%v\n", userID, token, err)

		return fmt.Errorf("%s:%w", f, err)
	}

	log.Printf("Successfully set and verified refresh token for userID/token: %d/%s\n", userID, token)

	return nil
}

func (t *TokenStorage) UserID(ctx context.Context, token string) (string, error) {
	const f = "redis.UserID"

	userIdStr, err := t.client.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("%s:%w", f, ErrNotFound)
		}

		return "", fmt.Errorf("%s:%w", f, err)
	}

	return userIdStr, nil
}

func (t *TokenStorage) Delete(ctx context.Context, userID int64, token string) error {
	const f = "redis.Delete"

	key := fmt.Sprintf("%d", userID)
	if err := t.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/token: %d/%s\n", userID, token)

		return fmt.Errorf("%s:%w", f, err)
	}

	return nil
}
