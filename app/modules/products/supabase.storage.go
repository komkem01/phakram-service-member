package products

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const maxProductImageFileSizeBytes = 5 * 1024 * 1024

type supabaseStorageClient struct {
	url           string
	serviceKey    string
	publicBucket  string
	privateBucket string
	httpClient    *http.Client
}

type uploadedProductImage struct {
	Path     string
	FileName string
	MIMEType string
	Size     int64
}

func newSupabaseStorageClient(conf SupabaseConfig) *supabaseStorageClient {
	url := strings.TrimRight(strings.TrimSpace(conf.URL), "/")
	if url == "" {
		url = strings.TrimRight(firstNonEmptyEnv("SUPABASE_URL"), "/")
	}

	serviceKey := strings.TrimSpace(conf.ServiceRoleKey)
	if serviceKey == "" {
		serviceKey = firstNonEmptyEnv("SUPABASE_SERVICE_ROLE_KEY", "SUPABASE_SERVICE_KEY", "SUPABASE_ANON_KEY")
	}

	publicBucket := strings.TrimSpace(conf.PublicBucket)
	if publicBucket == "" {
		publicBucket = firstNonEmptyEnv("SUPABASE_PUBLIC_BUCKET", "OBJECT_PUBLIC_BUCKET")
	}

	privateBucket := strings.TrimSpace(conf.PrivateBucket)
	if privateBucket == "" {
		privateBucket = firstNonEmptyEnv("SUPABASE_PRIVATE_BUCKET", "OBJECT_PRIVATE_BUCKET")
	}

	return &supabaseStorageClient{
		url:           url,
		serviceKey:    serviceKey,
		publicBucket:  publicBucket,
		privateBucket: privateBucket,
		httpClient:    &http.Client{Timeout: 20 * time.Second},
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

func (c *supabaseStorageClient) enabledPublic() bool {
	c.ensureRuntimeConfig()
	return c != nil && c.url != "" && c.serviceKey != "" && c.publicBucket != ""
}

func (c *supabaseStorageClient) missingPublicConfigFields() []string {
	c.ensureRuntimeConfig()
	if c == nil {
		return []string{"client"}
	}

	missing := make([]string, 0, 3)
	if c.url == "" {
		missing = append(missing, "url")
	}
	if c.serviceKey == "" {
		missing = append(missing, "service_role_key")
	}
	if c.publicBucket == "" {
		missing = append(missing, "public_bucket")
	}

	return missing
}

func (c *supabaseStorageClient) ensureRuntimeConfig() {
	if c == nil {
		return
	}

	if c.url != "" && c.serviceKey != "" && c.publicBucket != "" {
		return
	}

	loadDotEnvForSupabase()

	if c.url == "" {
		c.url = strings.TrimRight(firstNonEmptyEnv("SUPABASE_URL"), "/")
	}
	if c.serviceKey == "" {
		c.serviceKey = firstNonEmptyEnv("SUPABASE_SERVICE_ROLE_KEY", "SUPABASE_SERVICE_KEY", "SUPABASE_ANON_KEY")
	}
	if c.publicBucket == "" {
		c.publicBucket = firstNonEmptyEnv("SUPABASE_PUBLIC_BUCKET", "OBJECT_PUBLIC_BUCKET")
	}
	if c.privateBucket == "" {
		c.privateBucket = firstNonEmptyEnv("SUPABASE_PRIVATE_BUCKET", "OBJECT_PRIVATE_BUCKET")
	}
}

func loadDotEnvForSupabase() {
	paths := []string{
		".env",
		"phakram-service-member/.env",
		"../phakram-service-member/.env",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		for _, path := range discoverDotEnvCandidates(cwd) {
			if _, statErr := os.Stat(path); statErr == nil {
				_ = godotenv.Load(path)
				return
			}
		}
	}
}

func discoverDotEnvCandidates(start string) []string {
	trimmed := strings.TrimSpace(start)
	if trimmed == "" {
		return nil
	}

	candidates := make([]string, 0, 12)
	current := trimmed
	for {
		candidates = append(candidates,
			filepath.Join(current, ".env"),
			filepath.Join(current, "phakram-service-member", ".env"),
		)

		next := filepath.Dir(current)
		if next == current {
			break
		}
		current = next
	}

	return candidates
}

func (c *supabaseStorageClient) ResolveObjectURL(storedPath string) string {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") || strings.HasPrefix(trimmed, "data:") {
		return trimmed
	}
	if c == nil || c.url == "" || c.publicBucket == "" {
		return trimmed
	}
	bucket, objectPath, ok := splitBucketAndObjectPath(trimmed)
	if !ok {
		return trimmed
	}
	if bucket != c.publicBucket {
		return trimmed
	}
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.url, bucket, objectPath)
}

func (c *supabaseStorageClient) UploadProductImage(ctx context.Context, productID uuid.UUID, fileName string, encoded string) (*uploadedProductImage, error) {
	if !c.enabledPublic() {
		return nil, errors.New("supabase public storage is not configured")
	}

	data, mimeType, err := decodeBase64Image(encoded)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("image is empty")
	}
	if len(data) > maxProductImageFileSizeBytes {
		return nil, errors.New("image exceeds 5 MB")
	}
	if !isAllowedProductImageMIME(mimeType) {
		return nil, fmt.Errorf("unsupported image type: %s", mimeType)
	}

	ext := extensionByMIME(mimeType)
	safeName := strings.TrimSpace(fileName)
	if safeName == "" {
		safeName = fmt.Sprintf("product-%s%s", productID.String(), ext)
	}

	objectPath := fmt.Sprintf("products/%s/%s-%d%s", productID.String(), uuid.NewString(), time.Now().UnixMilli(), ext)
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.url, c.publicBucket, objectPath)

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

	return &uploadedProductImage{
		Path:     fmt.Sprintf("%s/%s", c.publicBucket, objectPath),
		FileName: safeName,
		MIMEType: mimeType,
		Size:     int64(len(data)),
	}, nil
}

func (c *supabaseStorageClient) DeleteProductImageObject(ctx context.Context, storedPath string) error {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return nil
	}
	if c == nil {
		return nil
	}
	c.ensureRuntimeConfig()
	if c.url == "" || c.serviceKey == "" {
		return nil
	}

	bucket, objectPath, ok := splitBucketAndObjectPath(trimmed)
	if !ok {
		return nil
	}

	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.url, bucket, objectPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	req.Header.Set("apikey", c.serviceKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("supabase delete failed: %s", strings.TrimSpace(string(body)))
	}

	return nil
}

func decodeBase64Image(input string) ([]byte, string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, "", errors.New("image is required")
	}

	mimeType := ""
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, "", errors.New("invalid data url format")
		}
		header := parts[0]
		if !strings.Contains(header, ";base64") {
			return nil, "", errors.New("image must be base64 data url")
		}
		mimeType = strings.TrimPrefix(strings.SplitN(header, ";", 2)[0], "data:")
		raw = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(raw)
		if err != nil {
			return nil, "", errors.New("invalid base64 image")
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

func isAllowedProductImageMIME(mimeType string) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

func extensionByMIME(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
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
