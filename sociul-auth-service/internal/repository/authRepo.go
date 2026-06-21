package repository

import (
	"context"

	"sociul-auth-service/internal/models"
	"sociul-auth-service/internal/sentinel"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

// Check if a user exists
func (r *AuthRepository) UserExists(ctx context.Context, username string) (bool, error) {
	sqlQuery := "SELECT 1 FROM users WHERE username=$1"
	var result string

	err := r.pool.QueryRow(ctx, sqlQuery, username).Scan(&result)
	if err == pgx.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Creates a user in the database
func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) error {
	sqlQuery := "INSERT INTO users (id, username, password, email) VALUES ($1, $2, $3, $4)"

	_, err := r.pool.Exec(ctx, sqlQuery, user.Id, user.Username, user.Password, user.Email)
	return err
}

// Mark a user as verified
func (r *AuthRepository) MarkVerified(ctx context.Context, email string) error {
	sqlQuery := "UPDATE users SET is_verified = True, updated_at = NOW() WHERE email=$1"

	_, err := r.pool.Exec(ctx, sqlQuery, email)
	return err
}

// Update a user's password
func (r *AuthRepository) UpdatePassword(ctx context.Context, email, password string) error {
	sqlQuery := "UPDATE users SET password = $1, updated_at = NOW() WHERE email=$2"

	_, err := r.pool.Exec(ctx, sqlQuery, password, email)
	return err
}

// Deletes a user from the database
func (r *AuthRepository) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	sqlQuery := "DELETE FROM users WHERE id=$1"

	_, err := r.pool.Exec(ctx, sqlQuery, uuid)
	return err
}

// Returns uuid and password for a user
func (r *AuthRepository) GetUserCreds(ctx context.Context, username string) (uuid.UUID, string, error) {
	sqlQuery := "SELECT id, password FROM users WHERE username=$1"
	var id uuid.UUID
	var password string

	err := r.pool.QueryRow(ctx, sqlQuery, username).Scan(&id, &password)
	if err == pgx.ErrNoRows {
		return uuid.Nil, "", sentinel.ErrUserNotFound
	}
	if err != nil {
		return uuid.Nil, "", err
	}
	return id, password, nil
}
