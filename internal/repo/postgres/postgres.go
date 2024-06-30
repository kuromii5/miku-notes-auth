package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kuromii5/sso-auth/internal/models"
	"github.com/lib/pq"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("application not found")
)

type DB struct {
	db *sql.DB
}

func New(connString string) (*DB, error) {
	const f = "postgres.NewDB"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", f, err)
	}

	return &DB{db: db}, nil
}

func (d *DB) SaveUser(ctx context.Context, email string, passwordHash []byte) (int32, error) {
	const f = "postgres.SaveUser"

	query := "INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id"

	var userID int32
	err := d.db.QueryRowContext(ctx, query, email, passwordHash).Scan(&userID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique violation
				return 0, fmt.Errorf("%s:%w", f, ErrUserExists)
			}
		}

		return 0, fmt.Errorf("%s:%w", f, err)
	}

	return userID, nil
}

func (d *DB) User(ctx context.Context, email string) (models.User, error) {
	const f = "postgres.User"

	query := "SELECT id, email, pass_hash, created_at, updated_at FROM users WHERE email = $1"

	var user models.User
	err := d.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s:%w", f, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s:%w", f, err)
	}

	return user, nil
}
