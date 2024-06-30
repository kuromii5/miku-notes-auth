package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	config "github.com/kuromii5/sso-auth/internal/cfg"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// read full config
	cfg := config.MustLoad()

	// read postgres config
	postgresCfg := cfg.Postgres

	// dbUrl for postgres
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		postgresCfg.User,
		postgresCfg.Password,
		postgresCfg.Host,
		fmt.Sprint(postgresCfg.Port),
		postgresCfg.DBName,
		"migrations",
	)

	if err := clearRedisCache(cfg.Tokens.RedisAddr); err != nil {
		panic("Error clearing Redis cache")
	}

	log.Println("Successfully cleared Redis")

	if err := clearPostgresDB(dbURL); err != nil {
		panic("Error clearing Postgres DB")
	}

	log.Println("Successfully cleared Postgres")
}

func clearRedisCache(redisAddr string) error {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush redis: %w", err)
	}

	return nil
}

func clearPostgresDB(dbURL string) error {
	m, err := migrate.New(
		"file://migrations/",
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	return nil
}
