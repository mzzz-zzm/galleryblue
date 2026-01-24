package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"

	"connectrpc.com/connect"
	"golang.org/x/image/draw"

	"github.com/mzzz-zzm/galleryblue/internal/db"
	usersv1 "github.com/mzzz-zzm/galleryblue/gen/go/users/v1"
)

const maxImageSize = 5 * 1024 * 1024 // 5MB
const thumbnailMaxWidth = 300
const thumbnailMaxHeight = 200

// ImageServer implements the ImageService
type ImageServer struct{}

// generateThumbnail creates a smaller version of the image for gallery display
func generateThumbnail(data []byte, maxWidth, maxHeight int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Calculate new dimensions maintaining aspect ratio
	newWidth := maxWidth
	newHeight := maxHeight

	widthRatio := float64(maxWidth) / float64(origWidth)
	heightRatio := float64(maxHeight) / float64(origHeight)

	if widthRatio < heightRatio {
		newHeight = int(float64(origHeight) * widthRatio)
	} else {
		newWidth = int(float64(origWidth) * heightRatio)
	}

	// Create thumbnail
	thumbnail := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(thumbnail, thumbnail.Bounds(), img, img.Bounds(), draw.Over, nil)

	// Encode as JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 70}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return buf.Bytes(), nil
}

// UploadImage uploads a new image (authenticated user becomes owner)
func (s *ImageServer) UploadImage(
	ctx context.Context,
	req *connect.Request[usersv1.UploadImageRequest],
) (*connect.Response[usersv1.UploadImageResponse], error) {
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	// Validate request
	if req.Msg.Filename == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("filename is required"))
	}
	if req.Msg.ContentType != "image/jpeg" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("only JPEG images are supported"))
	}
	if len(req.Msg.Data) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("image data is required"))
	}
	if len(req.Msg.Data) > maxImageSize {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("image too large (max 5MB)"))
	}

	// Generate thumbnail
	thumbnail, err := generateThumbnail(req.Msg.Data, thumbnailMaxWidth, thumbnailMaxHeight)
	if err != nil {
		// Log error but continue without thumbnail
		fmt.Printf("Warning: failed to generate thumbnail: %v\n", err)
		thumbnail = nil
	}

	imageID, err := db.CreateImage(ctx, userID, req.Msg.Filename, req.Msg.ContentType,
		req.Msg.Data, thumbnail, req.Msg.Title, req.Msg.Description)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create image: %w", err))
	}

	return connect.NewResponse(&usersv1.UploadImageResponse{
		ImageId: imageID,
	}), nil
}

// GetImage retrieves a single image by ID (public)
func (s *ImageServer) GetImage(
	ctx context.Context,
	req *connect.Request[usersv1.GetImageRequest],
) (*connect.Response[usersv1.GetImageResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("image id is required"))
	}

	img, err := db.GetImageByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}
	if img == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("image not found"))
	}

	return connect.NewResponse(&usersv1.GetImageResponse{
		Id:               img.ID,
		OwnerId:          img.OwnerID,
		OwnerDisplayName: img.OwnerDisplayName,
		Filename:         img.Filename,
		ContentType:      img.ContentType,
		Data:             img.Data,
		Title:            img.Title,
		Description:      img.Description,
		CreatedAt:        img.CreatedAt,
	}), nil
}

// ListImages returns all images (public gallery)
func (s *ImageServer) ListImages(
	ctx context.Context,
	req *connect.Request[usersv1.ListImagesRequest],
) (*connect.Response[usersv1.ListImagesResponse], error) {
	images, total, err := db.ListImages(ctx, int(req.Msg.Limit), int(req.Msg.Offset))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}

	var pbImages []*usersv1.ImageInfo
	for _, img := range images {
		pbImages = append(pbImages, &usersv1.ImageInfo{
			Id:               img.ID,
			OwnerId:          img.OwnerID,
			OwnerDisplayName: img.OwnerDisplayName,
			Filename:         img.Filename,
			Title:            img.Title,
			CreatedAt:        img.CreatedAt,
			Thumbnail:        img.Thumbnail,
		})
	}

	return connect.NewResponse(&usersv1.ListImagesResponse{
		Images: pbImages,
		Total:  int32(total),
	}), nil
}

// ListMyImages returns images owned by the current user
func (s *ImageServer) ListMyImages(
	ctx context.Context,
	req *connect.Request[usersv1.ListMyImagesRequest],
) (*connect.Response[usersv1.ListMyImagesResponse], error) {
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	images, total, err := db.ListImagesByOwner(ctx, userID, int(req.Msg.Limit), int(req.Msg.Offset))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}

	var pbImages []*usersv1.ImageInfo
	for _, img := range images {
		pbImages = append(pbImages, &usersv1.ImageInfo{
			Id:               img.ID,
			OwnerId:          img.OwnerID,
			OwnerDisplayName: img.OwnerDisplayName,
			Filename:         img.Filename,
			Title:            img.Title,
			CreatedAt:        img.CreatedAt,
			Thumbnail:        img.Thumbnail,
		})
	}

	return connect.NewResponse(&usersv1.ListMyImagesResponse{
		Images: pbImages,
		Total:  int32(total),
	}), nil
}

// UpdateImage updates image metadata (owner only)
func (s *ImageServer) UpdateImage(
	ctx context.Context,
	req *connect.Request[usersv1.UpdateImageRequest],
) (*connect.Response[usersv1.UpdateImageResponse], error) {
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("image id is required"))
	}

	// Verify ownership
	ownerID, err := db.GetImageOwner(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}
	if ownerID == "" {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("image not found"))
	}
	if ownerID != userID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("you can only edit your own images"))
	}

	// Get current values for fields not being updated
	img, err := db.GetImageByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}

	newTitle := img.Title
	newDescription := img.Description
	if req.Msg.Title != nil {
		newTitle = *req.Msg.Title
	}
	if req.Msg.Description != nil {
		newDescription = *req.Msg.Description
	}

	if err := db.UpdateImage(ctx, req.Msg.Id, newTitle, newDescription); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update image: %w", err))
	}

	return connect.NewResponse(&usersv1.UpdateImageResponse{
		Id:          req.Msg.Id,
		Title:       newTitle,
		Description: newDescription,
	}), nil
}

// DeleteImage removes an image (owner only)
func (s *ImageServer) DeleteImage(
	ctx context.Context,
	req *connect.Request[usersv1.DeleteImageRequest],
) (*connect.Response[usersv1.DeleteImageResponse], error) {
	userID := req.Header().Get("X-User-ID")
	if userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("image id is required"))
	}

	// Verify ownership
	ownerID, err := db.GetImageOwner(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("database error: %w", err))
	}
	if ownerID == "" {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("image not found"))
	}
	if ownerID != userID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("you can only delete your own images"))
	}

	if err := db.DeleteImage(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete image: %w", err))
	}

	return connect.NewResponse(&usersv1.DeleteImageResponse{
		Success: true,
	}), nil
}
