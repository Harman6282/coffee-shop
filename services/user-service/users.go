package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Insert(ctx context.Context, user RegisterRequest) error
	GetByEmail(ctx context.Context, email string) (string, error)
	UpdateLastLogin(ctx context.Context, email string) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) Insert(ctx context.Context, user RegisterRequest) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var existingID string
	err := u.db.QueryRowContext(
		ctx,
		`SELECT id FROM users WHERE email = $1`,
		user.Email,
	).Scan(&existingID)
	if err == nil {
		return errors.New("email already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	id := uuid.NewString()

	_, err = u.db.ExecContext(
		ctx,
		`INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`,
		id,
		user.Email,
		user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (string, error) {
	var hashedPassword string

	err := u.db.QueryRowContext(
		ctx,
		`SELECT password FROM users WHERE email = $1`,
		email,
	).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	return hashedPassword, nil
}

func (u *userRepo) UpdateLastLogin(ctx context.Context, email string) error {
	_, err := u.db.ExecContext(
		ctx,
		`UPDATE users SET last_login = NOW() WHERE email = $1`,
		email,
	)
	if err != nil {
		return err
	}

	return nil
}
