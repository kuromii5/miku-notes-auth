package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// TokenTTL - Token time-to-live

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	Postgres PostgresConfig `yaml:"postgres"`
	GRPC     GrpcConfig     `yaml:"grpc"`
	Tokens   TokensConfig   `yaml:"tokens" env-required:"true"`
}

type PostgresConfig struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type TokensConfig struct {
	AccessTTL  time.Duration `yaml:"access_ttl"`
	RefreshTTL time.Duration `yaml:"refresh_ttl"`
	RedisAddr  string        `yaml:"redis_addr"`
	Secret     string        `yaml:"secret"`
}

type GrpcConfig struct {
	Port            int    `yaml:"port"`
	ConnectionToken string `yaml:"connection_token"`
}

func MustLoad() *Config {
	path := checkPath()

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("error reading config")
	}

	return &cfg
}

func LoadForMigrations() *PostgresConfig {
	path := checkPath()

	// Define a struct that contains only the `postgres` field
	var config struct {
		Postgres PostgresConfig `yaml:"postgres"`
	}

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		fmt.Println(err)
		panic("error reading config for postgres")
	}

	return &config.Postgres
}

func checkPath() string {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty.")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist" + path)
	}

	return path
}

func fetchConfigPath() string {
	var result string

	flag.StringVar(&result, "config", "", "path to cfg file")
	flag.Parse()

	return result
}

func (pc PostgresConfig) ConnString() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		pc.User, pc.Password, pc.Host, pc.Port, pc.DBName, pc.SSLMode)
}
