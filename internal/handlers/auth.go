package handlers

import (
	"context"
	"crypto/rand"
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

	// Check if email already exists
	exists, err := db.EmailExists(ctx, email)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}
	if exists {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("user with this email already exists"))
	}

	// Check if display_name already exists (if provided)
	if displayName != "" {
		exists, err = db.DisplayNameExists(ctx, displayName)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
		}
		if exists {
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("user with this display name already exists"))
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to hash password: %w", err))
	}

	// Insert user
	userID, err := db.CreateUser(ctx, email, string(hashedPassword), displayName)
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
	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}
	if user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid email or password"))
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid email or password"))
	}

	// Generate session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate token: %w", err))
	}
	sessionToken := hex.EncodeToString(tokenBytes)

	return connect.NewResponse(&usersv1.LoginResponse{
		SessionToken: sessionToken,
		UserId:       user.ID,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
	}), nil
}
