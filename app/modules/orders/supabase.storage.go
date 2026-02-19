package orders

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

const maxSlipFileSizeBytes = 5 * 1024 * 1024

type supabaseStorageClient struct {
	url           string
	serviceKey    string
	publicBucket  string
	privateBucket string
	httpClient    *http.Client
}

type uploadedSlipObject struct {
	Path     string
	FileName string
	MIMEType string
	Size     int64
}

func newSupabaseStorageClient(conf SupabaseConfig) *supabaseStorageClient {
	return &supabaseStorageClient{
		url:           strings.TrimRight(strings.TrimSpace(conf.URL), "/"),
		serviceKey:    strings.TrimSpace(conf.ServiceRoleKey),
		publicBucket:  strings.TrimSpace(conf.PublicBucket),
		privateBucket: strings.TrimSpace(conf.PrivateBucket),
		httpClient:    &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *supabaseStorageClient) enabledForPrivate() bool {
	return c != nil && c.url != "" && c.serviceKey != "" && c.privateBucket != ""
}

func (c *supabaseStorageClient) UploadPaymentSlip(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID, fileName string, encoded string) (*uploadedSlipObject, error) {
	if !c.enabledForPrivate() {
		return nil, errors.New("supabase storage is not configured")
	}

	data, mimeType, err := decodeBase64Image(encoded)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("slip image is empty")
	}
	if len(data) > maxSlipFileSizeBytes {
		return nil, errors.New("slip image exceeds 5 MB")
	}

	if !isAllowedImageMIME(mimeType) {
		return nil, fmt.Errorf("unsupported slip image type: %s", mimeType)
	}

	ext := extensionByMIME(mimeType)
	safeName := strings.TrimSpace(fileName)
	if safeName == "" {
		safeName = fmt.Sprintf("payment-slip-%s%s", orderID.String(), ext)
	}

	objectPath := fmt.Sprintf("payments/%s/%s-%d%s", orderID.String(), paymentID.String(), time.Now().UnixMilli(), ext)
	endpoint := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.url, c.privateBucket, objectPath)

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

	return &uploadedSlipObject{
		Path:     fmt.Sprintf("%s/%s", c.privateBucket, objectPath),
		FileName: safeName,
		MIMEType: mimeType,
		Size:     int64(len(data)),
	}, nil
}

func decodeBase64Image(input string) ([]byte, string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, "", errors.New("slip image is required")
	}

	mimeType := ""
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, "", errors.New("invalid data url format")
		}
		header := parts[0]
		if !strings.Contains(header, ";base64") {
			return nil, "", errors.New("slip image must be base64 data url")
		}
		mimeType = strings.TrimPrefix(strings.SplitN(header, ";", 2)[0], "data:")
		raw = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(raw)
		if err != nil {
			return nil, "", errors.New("invalid base64 slip image")
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

func isAllowedImageMIME(mimeType string) bool {
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
		ext := filepath.Ext(mimeType)
		if ext != "" {
			return ext
		}
		return ".bin"
	}
}
