package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	config "github.com/kuromii5/sso-auth/internal/cfg"
)

func main() {
	// read config file
	cfg := config.LoadForMigrations()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		fmt.Sprint(cfg.Port),
		cfg.DBName,
		"migrations",
	)

	m, err := migrate.New(
		"file://migrations/",
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("migrations applied")
}
