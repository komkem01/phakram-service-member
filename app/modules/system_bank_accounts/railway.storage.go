package systembankaccounts

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

const maxSystemBankQRCodeFileSizeBytes = 5 * 1024 * 1024

type railwayStorageClient struct {
	s3            *s3compat.Client
	publicBucket  string
	privateBucket string
}

type uploadedSystemBankQR struct {
	Path     string
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

	return &railwayStorageClient{
		s3:            s3compat.NewClient(endpointURL, accessKeyID, secretAccessKey, region, 20*time.Second),
		publicBucket:  publicBucket,
		privateBucket: privateBucket,
	}
}

func firstNonEmptyEnv(names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return value
		}
	}
	return ""
}

func (c *railwayStorageClient) enabledPublic() bool {
	return c != nil && c.s3 != nil && c.s3.Enabled() && c.publicBucket != ""
}

func (c *railwayStorageClient) ResolveObjectURL(storedPath string) string {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") || strings.HasPrefix(trimmed, "data:") {
		return trimmed
	}
	if c == nil || c.s3 == nil || c.publicBucket == "" {
		return trimmed
	}

	bucket, objectPath, ok := splitBucketAndObjectPath(trimmed)
	if !ok {
		return trimmed
	}

	presignedURL, err := c.s3.PresignGetObject(bucket, objectPath, 24*time.Hour)
	if err == nil && strings.TrimSpace(presignedURL) != "" {
		return presignedURL
	}

	publicURL := c.s3.PublicObjectURL(bucket, objectPath)
	if publicURL != "" && publicURL != strings.TrimSpace(objectPath) {
		return publicURL
	}

	return trimmed
}

func (c *railwayStorageClient) UploadSystemBankQRCode(ctx context.Context, systemBankAccountID uuid.UUID, encoded string) (*uploadedSystemBankQR, error) {
	if !c.enabledPublic() {
		return nil, errors.New("railway public storage is not configured")
	}

	data, mimeType, err := decodeSystemBankBase64Image(encoded)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("qr image is empty")
	}
	if len(data) > maxSystemBankQRCodeFileSizeBytes {
		return nil, errors.New("qr image exceeds 5 MB")
	}
	if !isAllowedSystemBankImageMIME(mimeType) {
		return nil, fmt.Errorf("unsupported qr image type: %s", mimeType)
	}

	ext := extensionBySystemBankMIME(mimeType)
	objectPath := fmt.Sprintf("system-bank-accounts/%s/qr-%d%s", systemBankAccountID.String(), time.Now().UnixMilli(), ext)
	if err := c.s3.PutObject(ctx, c.publicBucket, objectPath, mimeType, data); err != nil {
		return nil, err
	}

	return &uploadedSystemBankQR{
		Path:     fmt.Sprintf("%s/%s", c.publicBucket, objectPath),
		MIMEType: mimeType,
		Size:     int64(len(data)),
	}, nil
}

func decodeSystemBankBase64Image(input string) ([]byte, string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, "", errors.New("qr image is required")
	}

	mimeType := ""
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, "", errors.New("invalid data url format")
		}
		header := parts[0]
		if !strings.Contains(header, ";base64") {
			return nil, "", errors.New("qr image must be base64 data url")
		}
		mimeType = strings.TrimPrefix(strings.SplitN(header, ";", 2)[0], "data:")
		raw = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(raw)
		if err != nil {
			return nil, "", errors.New("invalid base64 qr image")
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

func isAllowedSystemBankImageMIME(mimeType string) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/png", "image/webp", "image/heic", "image/heif", "image/heic-sequence", "image/heif-sequence":
		return true
	default:
		return false
	}
}

func extensionBySystemBankMIME(mimeType string) string {
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

func splitBucketAndObjectPath(storedPath string) (string, string, bool) {
	trimmed := strings.Trim(strings.TrimSpace(storedPath), "/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	bucket := strings.TrimSpace(parts[0])
	objectPath := strings.TrimSpace(parts[1])
	if bucket == "" || objectPath == "" {
		return "", "", false
	}
	return bucket, objectPath, true
}
