package storages

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const signedURLExpiresInSeconds = 60 * 60

type supabaseStorageClient struct {
	url           string
	serviceKey    string
	publicBucket  string
	privateBucket string
	httpClient    *http.Client
}

type signedURLResponse struct {
	SignedURL string `json:"signedURL"`
	SignedUrl string `json:"signedUrl"`
	Path      string `json:"path"`
}

func newSupabaseStorageClient(conf SupabaseConfig) *supabaseStorageClient {
	return &supabaseStorageClient{
		url:           strings.TrimRight(strings.TrimSpace(conf.URL), "/"),
		serviceKey:    strings.TrimSpace(conf.ServiceRoleKey),
		publicBucket:  strings.TrimSpace(conf.PublicBucket),
		privateBucket: strings.TrimSpace(conf.PrivateBucket),
		httpClient:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *supabaseStorageClient) enabled() bool {
	return c != nil && c.url != "" && c.serviceKey != ""
}

func (c *supabaseStorageClient) ResolveObjectURL(ctx context.Context, storedPath string) (string, error) {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return trimmed, nil
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") || strings.HasPrefix(trimmed, "data:") {
		return trimmed, nil
	}
	if !c.enabled() {
		return trimmed, nil
	}

	bucket, objectPath, ok := splitBucketAndObjectPath(trimmed)
	if !ok {
		return trimmed, nil
	}

	if c.publicBucket != "" && bucket == c.publicBucket {
		return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.url, bucket, objectPath), nil
	}

	if c.privateBucket != "" && bucket == c.privateBucket {
		return c.createSignedURL(ctx, bucket, objectPath)
	}

	return trimmed, nil
}

func (c *supabaseStorageClient) createSignedURL(ctx context.Context, bucket string, objectPath string) (string, error) {
	endpoint := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", c.url, bucket, objectPath)
	payload, _ := json.Marshal(map[string]int{"expiresIn": signedURLExpiresInSeconds})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	req.Header.Set("apikey", c.serviceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("supabase sign failed: %s", strings.TrimSpace(string(body)))
	}

	result := new(signedURLResponse)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return "", err
	}

	signed := strings.TrimSpace(result.SignedURL)
	if signed == "" {
		signed = strings.TrimSpace(result.SignedUrl)
	}
	if signed == "" {
		return "", fmt.Errorf("supabase sign response missing signedURL")
	}
	if strings.HasPrefix(signed, "http://") || strings.HasPrefix(signed, "https://") {
		return signed, nil
	}
	return c.url + signed, nil
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
