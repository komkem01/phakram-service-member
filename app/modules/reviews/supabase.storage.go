package reviews

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const maxReviewImageFileSizeBytes = 5 * 1024 * 1024

type supabaseStorageClient struct {
	url           string
	serviceKey    string
	publicBucket  string
	reviewBucket  string
	privateBucket string
	httpClient    *http.Client
}

type uploadedReviewImage struct {
	Path     string
	FileName string
	MIMEType string
	Size     int64
}

func newSupabaseStorageClient(conf SupabaseConfig) *supabaseStorageClient {
	reviewBucket := strings.TrimSpace(conf.ReviewBucket)
	if reviewBucket == "" {
		reviewBucket = strings.TrimSpace(conf.PublicBucket)
	}

	return &supabaseStorageClient{
		url:           strings.TrimRight(strings.TrimSpace(conf.URL), "/"),
		serviceKey:    strings.TrimSpace(conf.ServiceRoleKey),
		publicBucket:  strings.TrimSpace(conf.PublicBucket),
		reviewBucket:  reviewBucket,
		privateBucket: strings.TrimSpace(conf.PrivateBucket),
		httpClient:    &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *supabaseStorageClient) enabledForPublic() bool {
	return c != nil && c.url != "" && c.serviceKey != "" && c.reviewBucket != ""
}

func (c *supabaseStorageClient) ResolveObjectURL(storedPath string) string {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") || strings.HasPrefix(trimmed, "data:") {
		return trimmed
	}
	if c == nil || c.url == "" {
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

	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.url, bucket, objectPath)
}

func (c *supabaseStorageClient) UploadReviewImage(ctx context.Context, productID uuid.UUID, reviewID uuid.UUID, fileName string, encoded string) (*uploadedReviewImage, error) {
	if !c.enabledForPublic() {
		return nil, errors.New("supabase public storage is not configured")
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
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.url, c.reviewBucket, objectPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	req.Header.Set("apikey", c.serviceKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-upsert", "true")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("supabase upload failed: %s", strings.TrimSpace(string(body)))
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
	case "image/jpeg", "image/png", "image/webp":
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
	default:
		ext := filepath.Ext(mimeType)
		if ext != "" {
			return ext
		}
		return ".bin"
	}
}
