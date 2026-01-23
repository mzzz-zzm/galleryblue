package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"golang.org/x/crypto/bcrypt"

	"github.com/mzzz-zzm/galleryblue/internal/db"
	usersv1 "github.com/mzzz-zzm/galleryblue/gen/go/users/v1"
)

// AuthServer implements the AuthService
type AuthServer struct{}

// Register creates a new user account
func (s *AuthServer) Register(
	ctx context.Context,
	req *connect.Request[usersv1.RegisterRequest],
) (*connect.Response[usersv1.RegisterResponse], error) {
	email := req.Msg.Email
	password := req.Msg.Password
	displayName := req.Msg.DisplayName

	// Validate input
	if email == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email and password are required"))
	}

	// Check if user already exists
	var existingID string
	err := db.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("user with this email already exists"))
	} else if err != sql.ErrNoRows {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}

	// Check if display_name already exists (if provided)
	if displayName != "" {
		err = db.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE display_name = $1", displayName).Scan(&existingID)
		if err == nil {
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("user with this display name already exists"))
		} else if err != sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to hash password: %w", err))
	}

	// Insert user
	var userID string
	err = db.DB.QueryRowContext(ctx,
		"INSERT INTO users (email, password_hash, display_name) VALUES ($1, $2, $3) RETURNING id",
		email, string(hashedPassword), displayName,
	).Scan(&userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create user: %w", err))
	}

	return connect.NewResponse(&usersv1.RegisterResponse{
		UserId:      userID,
		DisplayName: displayName,
		Email:       email,
	}), nil
}

// Login authenticates a user and returns a session token
func (s *AuthServer) Login(
	ctx context.Context,
	req *connect.Request[usersv1.LoginRequest],
) (*connect.Response[usersv1.LoginResponse], error) {
	email := req.Msg.Email
	password := req.Msg.Password

	// Validate input
	if email == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email and password are required"))
	}

	// Fetch user
	var userID, passwordHash, displayName string
	err := db.DB.QueryRowContext(ctx,
		"SELECT id, password_hash, display_name FROM users WHERE email = $1",
		email,
	).Scan(&userID, &passwordHash, &displayName)
	if err == sql.ErrNoRows {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid email or password"))
	} else if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid email or password"))
	}

	// Generate session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate token: %w", err))
	}
	sessionToken := hex.EncodeToString(tokenBytes)

	// TODO: Store session token in database or cache (Redis) for validation
	// For now, we just return the token (stateless approach - could use JWT instead)

	return connect.NewResponse(&usersv1.LoginResponse{
		SessionToken: sessionToken,
		UserId:       userID,
		DisplayName:  displayName,
		Email:        email,
	}), nil
}
