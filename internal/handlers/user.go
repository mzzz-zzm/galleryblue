package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"golang.org/x/crypto/bcrypt"

	"github.com/mzzz-zzm/galleryblue/internal/db"
	usersv1 "github.com/mzzz-zzm/galleryblue/gen/go/users/v1"
)

// UserServer implements the UserService
type UserServer struct{}

// GetUser retrieves a user by ID
func (s *UserServer) GetUser(
	ctx context.Context,
	req *connect.Request[usersv1.GetUserRequest],
) (*connect.Response[usersv1.GetUserResponse], error) {
	userID := req.Msg.Id

	if userID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user id is required"))
	}

	var id, displayName, email string
	err := db.DB.QueryRowContext(ctx,
		"SELECT id, display_name, email FROM users WHERE id = $1",
		userID,
	).Scan(&id, &displayName, &email)
	if err == sql.ErrNoRows {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	} else if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}

	return connect.NewResponse(&usersv1.GetUserResponse{
		Id:    id,
		Name:  displayName,
		Email: email,
	}), nil
}

// UpdateUser updates user information
func (s *UserServer) UpdateUser(
	ctx context.Context,
	req *connect.Request[usersv1.UpdateUserRequest],
) (*connect.Response[usersv1.UpdateUserResponse], error) {
	// Get user ID from context (should be set by auth middleware)
	// For now, we'll extract it from a header
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	currentPassword := req.Msg.CurrentPassword
	if currentPassword == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("current password is required"))
	}

	// Verify current password
	var passwordHash, currentDisplayName, currentEmail string
	err := db.DB.QueryRowContext(ctx,
		"SELECT password_hash, display_name, email FROM users WHERE id = $1",
		userID,
	).Scan(&passwordHash, &currentDisplayName, &currentEmail)
	if err == sql.ErrNoRows {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	} else if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(currentPassword)); err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("incorrect password"))
	}

	// Build update query dynamically based on provided fields
	newDisplayName := currentDisplayName
	newEmail := currentEmail
	newPasswordHash := passwordHash

	// Check new display name uniqueness
	if req.Msg.NewDisplayName != nil && *req.Msg.NewDisplayName != "" {
		var existingID string
		err = db.DB.QueryRowContext(ctx,
			"SELECT id FROM users WHERE display_name = $1 AND id != $2",
			*req.Msg.NewDisplayName, userID,
		).Scan(&existingID)
		if err == nil {
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("display name already taken"))
		} else if err != sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
		}
		newDisplayName = *req.Msg.NewDisplayName
	}

	// Check new email uniqueness
	if req.Msg.NewEmail != nil && *req.Msg.NewEmail != "" {
		var existingID string
		err = db.DB.QueryRowContext(ctx,
			"SELECT id FROM users WHERE email = $1 AND id != $2",
			*req.Msg.NewEmail, userID,
		).Scan(&existingID)
		if err == nil {
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("email already taken"))
		} else if err != sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
		}
		newEmail = *req.Msg.NewEmail
	}

	// Hash new password if provided
	if req.Msg.NewPassword != nil && *req.Msg.NewPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Msg.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to hash password: %w", err))
		}
		newPasswordHash = string(hashed)
	}

	// Update user
	_, err = db.DB.ExecContext(ctx,
		"UPDATE users SET display_name = $1, email = $2, password_hash = $3, updated_at = NOW() WHERE id = $4",
		newDisplayName, newEmail, newPasswordHash, userID,
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user: %w", err))
	}

	return connect.NewResponse(&usersv1.UpdateUserResponse{
		UserId:      userID,
		DisplayName: newDisplayName,
		Email:       newEmail,
	}), nil
}
