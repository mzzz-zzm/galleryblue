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

// ============================================================
// Image queries
// ============================================================

// Image represents an image record from the database
type Image struct {
	ID               string
	OwnerID          string
	OwnerDisplayName string
	Filename         string
	ContentType      string
	Data             []byte
	Title            string
	Description      string
	CreatedAt        string
	UpdatedAt        string
}

// ImageInfo represents image metadata with thumbnail for gallery display
type ImageInfo struct {
	ID               string
	OwnerID          string
	OwnerDisplayName string
	Filename         string
	Title            string
	CreatedAt        string
	Thumbnail        []byte
}

// CreateImage inserts a new image with thumbnail and returns the generated ID
func CreateImage(ctx context.Context, ownerID, filename, contentType string, data, thumbnail []byte, title, description string) (string, error) {
	var imageID string
	err := DB.QueryRowContext(ctx,
		`INSERT INTO images (owner_id, filename, content_type, data, thumbnail, title, description) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		ownerID, filename, contentType, data, thumbnail, title, description,
	).Scan(&imageID)
	return imageID, err
}

// GetImageByID fetches a single image with owner info
func GetImageByID(ctx context.Context, id string) (*Image, error) {
	var img Image
	err := DB.QueryRowContext(ctx,
		`SELECT i.id, i.owner_id, COALESCE(u.display_name, u.email) as owner_name,
		        i.filename, i.content_type, i.data, COALESCE(i.title, ''), COALESCE(i.description, ''),
		        i.created_at::text
		 FROM images i
		 JOIN users u ON i.owner_id = u.id
		 WHERE i.id = $1`,
		id,
	).Scan(&img.ID, &img.OwnerID, &img.OwnerDisplayName, &img.Filename, &img.ContentType,
		&img.Data, &img.Title, &img.Description, &img.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &img, nil
}

// ListImages returns all images (public gallery)
func ListImages(ctx context.Context, limit, offset int) ([]ImageInfo, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	// Get total count
	var total int
	err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM images").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := DB.QueryContext(ctx,
		`SELECT i.id, i.owner_id, COALESCE(u.display_name, u.email) as owner_name,
		        i.filename, COALESCE(i.title, ''), i.created_at::text, i.thumbnail
		 FROM images i
		 JOIN users u ON i.owner_id = u.id
		 ORDER BY i.created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var images []ImageInfo
	for rows.Next() {
		var img ImageInfo
		if err := rows.Scan(&img.ID, &img.OwnerID, &img.OwnerDisplayName, &img.Filename, &img.Title, &img.CreatedAt, &img.Thumbnail); err != nil {
			return nil, 0, err
		}
		images = append(images, img)
	}
	return images, total, rows.Err()
}

// ListImagesByOwner returns images owned by a specific user
func ListImagesByOwner(ctx context.Context, ownerID string, limit, offset int) ([]ImageInfo, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var total int
	err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM images WHERE owner_id = $1", ownerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := DB.QueryContext(ctx,
		`SELECT i.id, i.owner_id, COALESCE(u.display_name, u.email) as owner_name,
		        i.filename, COALESCE(i.title, ''), i.created_at::text, i.thumbnail
		 FROM images i
		 JOIN users u ON i.owner_id = u.id
		 WHERE i.owner_id = $2
		 ORDER BY i.created_at DESC
		 LIMIT $1 OFFSET $3`,
		limit, ownerID, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var images []ImageInfo
	for rows.Next() {
		var img ImageInfo
		if err := rows.Scan(&img.ID, &img.OwnerID, &img.OwnerDisplayName, &img.Filename, &img.Title, &img.CreatedAt, &img.Thumbnail); err != nil {
			return nil, 0, err
		}
		images = append(images, img)
	}
	return images, total, rows.Err()
}

// UpdateImage updates image metadata (owner must be verified by caller)
func UpdateImage(ctx context.Context, imageID, title, description string) error {
	_, err := DB.ExecContext(ctx,
		"UPDATE images SET title = $1, description = $2, updated_at = NOW() WHERE id = $3",
		title, description, imageID,
	)
	return err
}

// DeleteImage removes an image (owner must be verified by caller)
func DeleteImage(ctx context.Context, imageID string) error {
	_, err := DB.ExecContext(ctx, "DELETE FROM images WHERE id = $1", imageID)
	return err
}

// GetImageOwner returns the owner_id for an image
func GetImageOwner(ctx context.Context, imageID string) (string, error) {
	var ownerID string
	err := DB.QueryRowContext(ctx, "SELECT owner_id FROM images WHERE id = $1", imageID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return ownerID, err
}

