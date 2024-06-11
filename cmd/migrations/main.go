package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dbHost, dbPort, dbUser, dbPassword, dbName, migrationsPath, migrationsTable string

	flag.StringVar(&dbHost, "db-host", "localhost", "Database host")
	flag.StringVar(&dbPort, "db-port", "5432", "Database port")
	flag.StringVar(&dbUser, "db-user", "postgres", "Database user")
	flag.StringVar(&dbPassword, "db-password", "", "Database password")
	flag.StringVar(&dbName, "db-name", "sso", "Database name")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Name of migrations table")
	flag.Parse()

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || migrationsPath == "" {
		log.Fatal("All database and migration path flags must be provided")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, migrationsTable)

	m, err := migrate.New(
		"file://"+migrationsPath,
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
