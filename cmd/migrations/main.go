package main

import (
	"flag"
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

	var migrationsTable string
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Name of migrations table")
	flag.Parse()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		fmt.Sprint(cfg.Port),
		cfg.DBName,
		migrationsTable,
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
