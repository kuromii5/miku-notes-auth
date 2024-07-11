package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// TokenTTL - Token time-to-live

type Config struct {
	Env      string         `yaml:"env" env:"ENV"`
	Postgres PostgresConfig `yaml:"postgres"`
	GRPC     GrpcConfig     `yaml:"grpc"`
	Tokens   TokensConfig   `yaml:"tokens"`
}

type PostgresConfig struct {
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DBName   string `yaml:"dbname" env:"POSTGRES_DBNAME" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
}

type TokensConfig struct {
	AccessTTL  time.Duration `yaml:"access_ttl" env:"TOKENS_ACCESS_TTL"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env:"TOKENS_REFRESH_TTL"`
	RedisAddr  string        `yaml:"redis_addr" env:"TOKENS_REDIS_ADDR"`
	Secret     string        `yaml:"secret" env:"TOKENS_SECRET"`
}

type GrpcConfig struct {
	Port            int    `yaml:"port" env:"GRPC_PORT"`
	ConnectionToken string `yaml:"connection_token" env:"GRPC_CONNECTION_TOKEN"`
}

func MustLoad() *Config {
	path := checkPath()

	return ReadConfig(path)
}

func ReadConfig(path string) *Config {
	var cfg Config

	var err error
	if path != "" {
		err = cleanenv.ReadConfig(path, &cfg)
	} else {
		if err := godotenv.Load(".env"); err != nil {
			panic("error loading .env")
		}

		err = cleanenv.ReadEnv(&cfg)
	}

	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		panic("couldn't read config")
	}

	return &cfg
}

func LoadForMigrations() *PostgresConfig {
	path := checkPath()

	// Define a struct that contains only the `postgres` field
	var config struct {
		Postgres PostgresConfig `yaml:"postgres"`
	}

	if path != "" {
		if err := cleanenv.ReadConfig(path, &config.Postgres); err != nil {
			panic("error reading config")
		}
	} else {
		if err := cleanenv.ReadEnv(&config); err != nil {
			panic("error reading .env")
		}
	}

	return &config.Postgres
}

func checkPath() string {
	path := fetchConfigPath()
	if path == "" {
		log.Println("config path is empty. Looking for .env file...")
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic("config file doesn't exist" + path)
		}
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
