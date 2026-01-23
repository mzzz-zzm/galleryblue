package handlers

import (
	"context"
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

	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}
	if user == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	return connect.NewResponse(&usersv1.GetUserResponse{
		Id:    user.ID,
		Name:  user.DisplayName,
		Email: user.Email,
	}), nil
}

// UpdateUser updates user information
func (s *UserServer) UpdateUser(
	ctx context.Context,
	req *connect.Request[usersv1.UpdateUserRequest],
) (*connect.Response[usersv1.UpdateUserResponse], error) {
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	currentPassword := req.Msg.CurrentPassword
	if currentPassword == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("current password is required"))
	}

	// Fetch user
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}
	if user == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("incorrect password"))
	}

	// Build updated values
	newDisplayName := user.DisplayName
	newEmail := user.Email
	newPasswordHash := user.PasswordHash

	// Check new display name uniqueness
	if req.Msg.NewDisplayName != nil && *req.Msg.NewDisplayName != "" {
		exists, err := db.DisplayNameExistsExcluding(ctx, *req.Msg.NewDisplayName, userID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
		}
		if exists {
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("display name already taken"))
		}
		newDisplayName = *req.Msg.NewDisplayName
	}

	// Check new email uniqueness
	if req.Msg.NewEmail != nil && *req.Msg.NewEmail != "" {
		exists, err := db.EmailExistsExcluding(ctx, *req.Msg.NewEmail, userID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
		}
		if exists {
			return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("email already taken"))
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
	if err := db.UpdateUser(ctx, userID, newDisplayName, newEmail, newPasswordHash); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user: %w", err))
	}

	return connect.NewResponse(&usersv1.UpdateUserResponse{
		UserId:      userID,
		DisplayName: newDisplayName,
		Email:       newEmail,
	}), nil
}
