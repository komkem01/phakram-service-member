package reviews

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"phakram/app/utils/s3compat"
	"strings"
	"time"

	"github.com/google/uuid"
)

const maxReviewImageFileSizeBytes = 5 * 1024 * 1024
const publicReviewImageSignedURLExpiresInHours = 24

type railwayStorageClient struct {
	s3            *s3compat.Client
	publicBucket  string
	reviewBucket  string
	privateBucket string
}

type uploadedReviewImage struct {
	Path     string
	FileName string
	MIMEType string
	Size     int64
}

func newRailwayStorageClient(conf RailwayConfig) *railwayStorageClient {
	endpointURL := strings.TrimRight(strings.TrimSpace(conf.URL), "/")
	if endpointURL == "" {
		endpointURL = strings.TrimRight(firstNonEmptyEnv("OBJECT_ENDPOINT_URL"), "/")
	}

	secretAccessKey := strings.TrimSpace(conf.ServiceRoleKey)
	if secretAccessKey == "" {
		secretAccessKey = firstNonEmptyEnv("OBJECT_SECRET_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY", "S3_SECRET_ACCESS_KEY", "RAILWAY_STORAGE_SECRET_ACCESS_KEY", "RAILWAY_SECRET_ACCESS_KEY", "SECRET_ACCESS_KEY")
	}

	accessKeyID := firstNonEmptyEnv("OBJECT_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID", "S3_ACCESS_KEY_ID", "RAILWAY_STORAGE_ACCESS_KEY_ID", "RAILWAY_ACCESS_KEY_ID", "ACCESS_KEY_ID")
	region := firstNonEmptyEnv("OBJECT_REGION", "AWS_REGION", "AWS_DEFAULT_REGION")

	publicBucket := strings.TrimSpace(conf.PublicBucket)
	if publicBucket == "" {
		publicBucket = firstNonEmptyEnv("OBJECT_PUBLIC_BUCKET")
	}

	privateBucket := strings.TrimSpace(conf.PrivateBucket)
	if privateBucket == "" {
		privateBucket = firstNonEmptyEnv("OBJECT_PRIVATE_BUCKET")
	}

	reviewBucket := strings.TrimSpace(conf.ReviewBucket)
	if reviewBucket == "" {
		reviewBucket = firstNonEmptyEnv("OBJECT_REVIEW_BUCKET")
	}
	if reviewBucket == "" {
		reviewBucket = publicBucket
	}

	return &railwayStorageClient{
		s3:            s3compat.NewClient(endpointURL, accessKeyID, secretAccessKey, region, 20*time.Second),
		publicBucket:  publicBucket,
		reviewBucket:  reviewBucket,
		privateBucket: privateBucket,
	}
}

func firstNonEmptyEnv(names ...string) string {
	for _, name := range names {
		trimmedName := strings.TrimSpace(name)
		if trimmedName == "" {
			continue
		}
		if value := strings.TrimSpace(os.Getenv(trimmedName)); value != "" {
			return value
		}
	}
	return ""
}

func (c *railwayStorageClient) enabledForPublic() bool {
	return c != nil && c.s3 != nil && c.s3.Enabled() && c.reviewBucket != ""
}

func (c *railwayStorageClient) ResolveObjectURL(storedPath string) string {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") || strings.HasPrefix(trimmed, "data:") {
		return trimmed
	}
	if c == nil || c.s3 == nil {
		return trimmed
	}

	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return trimmed
	}
	bucket := strings.TrimSpace(parts[0])
	objectPath := strings.TrimLeft(parts[1], "/")
	if objectPath == "" {
		return trimmed
	}

	// Try presigned URL first
	presignedURL, err := c.s3.PresignGetObject(bucket, objectPath, time.Duration(publicReviewImageSignedURLExpiresInHours)*time.Hour)
	if err == nil && strings.TrimSpace(presignedURL) != "" {
		return presignedURL
	}

	// Fallback to public object URL
	publicURL := c.s3.PublicObjectURL(bucket, objectPath)
	if publicURL != "" && publicURL != strings.TrimSpace(objectPath) {
		return publicURL
	}

	// Last fallback: return the original path
	return trimmed
}

func (c *railwayStorageClient) UploadReviewImage(ctx context.Context, productID uuid.UUID, reviewID uuid.UUID, fileName string, encoded string) (*uploadedReviewImage, error) {
	if !c.enabledForPublic() {
		return nil, errors.New("railway public storage is not configured")
	}

	data, mimeType, err := decodeReviewBase64Image(encoded)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("review image is empty")
	}
	if len(data) > maxReviewImageFileSizeBytes {
		return nil, errors.New("review image exceeds 5 MB")
	}
	if !isAllowedReviewImageMIME(mimeType) {
		return nil, fmt.Errorf("unsupported review image type: %s", mimeType)
	}

	ext := extensionByReviewMIME(mimeType)
	safeName := strings.TrimSpace(fileName)
	if safeName == "" {
		safeName = fmt.Sprintf("review-%s%s", reviewID.String(), ext)
	}

	objectPath := fmt.Sprintf("reviews/%s/%s/%s-%d%s", productID.String(), reviewID.String(), uuid.NewString(), time.Now().UnixMilli(), ext)
	if err := c.s3.PutObject(ctx, c.reviewBucket, objectPath, mimeType, data); err != nil {
		return nil, err
	}

	return &uploadedReviewImage{
		Path:     fmt.Sprintf("%s/%s", c.reviewBucket, objectPath),
		FileName: safeName,
		MIMEType: mimeType,
		Size:     int64(len(data)),
	}, nil
}

func decodeReviewBase64Image(input string) ([]byte, string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, "", errors.New("review image is required")
	}

	mimeType := ""
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, "", errors.New("invalid data url format")
		}
		header := parts[0]
		if !strings.Contains(header, ";base64") {
			return nil, "", errors.New("review image must be base64 data url")
		}
		mimeType = strings.TrimPrefix(strings.SplitN(header, ";", 2)[0], "data:")
		raw = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(raw)
		if err != nil {
			return nil, "", errors.New("invalid base64 review image")
		}
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(decoded)
	}
	if mediaType, _, err := mime.ParseMediaType(mimeType); err == nil {
		mimeType = mediaType
	}

	return decoded, strings.ToLower(strings.TrimSpace(mimeType)), nil
}

func isAllowedReviewImageMIME(mimeType string) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/png", "image/webp", "image/heic", "image/heif", "image/heic-sequence", "image/heif-sequence":
		return true
	default:
		return false
	}
}

func extensionByReviewMIME(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/heic", "image/heic-sequence":
		return ".heic"
	case "image/heif", "image/heif-sequence":
		return ".heif"
	default:
		ext := filepath.Ext(mimeType)
		if ext != "" {
			return ext
		}
		return ".bin"
	}
}
