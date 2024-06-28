package models

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
