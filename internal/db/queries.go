package db

import (
	"context"
	"database/sql"
)

// User represents a user record from the database
type User struct {
	ID           string
	Email        string
	PasswordHash string
	DisplayName  string
}

// EmailExists checks if a user with the given email already exists
func EmailExists(ctx context.Context, email string) (bool, error) {
	var id string
	err := DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DisplayNameExists checks if a user with the given display name already exists
func DisplayNameExists(ctx context.Context, name string) (bool, error) {
	var id string
	err := DB.QueryRowContext(ctx, "SELECT id FROM users WHERE display_name = $1", name).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DisplayNameExistsExcluding checks if another user (not the given ID) has this display name
func DisplayNameExistsExcluding(ctx context.Context, name, excludeUserID string) (bool, error) {
	var id string
	err := DB.QueryRowContext(ctx, "SELECT id FROM users WHERE display_name = $1 AND id != $2", name, excludeUserID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// EmailExistsExcluding checks if another user (not the given ID) has this email
func EmailExistsExcluding(ctx context.Context, email, excludeUserID string) (bool, error) {
	var id string
	err := DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1 AND id != $2", email, excludeUserID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := DB.QueryRowContext(ctx,
		"SELECT id, email, password_hash, display_name FROM users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByID fetches a user by ID
func GetUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := DB.QueryRowContext(ctx,
		"SELECT id, email, password_hash, display_name FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser inserts a new user and returns the generated ID
func CreateUser(ctx context.Context, email, passwordHash, displayName string) (string, error) {
	var userID string
	err := DB.QueryRowContext(ctx,
		"INSERT INTO users (email, password_hash, display_name) VALUES ($1, $2, $3) RETURNING id",
		email, passwordHash, displayName,
	).Scan(&userID)
	return userID, err
}

// UpdateUser updates user fields
func UpdateUser(ctx context.Context, userID, displayName, email, passwordHash string) error {
	_, err := DB.ExecContext(ctx,
		"UPDATE users SET display_name = $1, email = $2, password_hash = $3, updated_at = NOW() WHERE id = $4",
		displayName, email, passwordHash, userID,
	)
	return err
}
