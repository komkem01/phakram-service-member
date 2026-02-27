package storages

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"phakram/app/utils/s3compat"
	"strings"
	"time"
)

const signedURLExpiresInSeconds = 60 * 60

type railwayStorageClient struct {
	s3            *s3compat.Client
	publicBucket  string
	privateBucket string
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
		s3:            s3compat.NewClient(endpointURL, accessKeyID, secretAccessKey, region, 15*time.Second),
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

func (c *railwayStorageClient) enabled() bool {
	return c != nil && c.s3 != nil && c.s3.Enabled()
}

func (c *railwayStorageClient) ResolveObjectURL(ctx context.Context, storedPath string) (string, error) {
	trimmed := strings.TrimSpace(storedPath)
	if trimmed == "" {
		return trimmed, nil
	}
	if strings.HasPrefix(trimmed, "data:") {
		return trimmed, nil
	}
	if !c.enabled() {
		return trimmed, nil
	}

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		resolved, ok := c.resolveAbsoluteStorageURL(ctx, trimmed)
		if ok {
			return resolved, nil
		}
		return trimmed, nil
	}

	bucket, objectPath, ok := splitBucketAndObjectPath(trimmed)
	if !ok {
		return trimmed, nil
	}

	return c.resolveBucketObjectURL(ctx, bucket, objectPath, trimmed), nil
}

func (c *railwayStorageClient) resolveBucketObjectURL(ctx context.Context, bucket string, objectPath string, fallback string) string {
	signedURL, signErr := c.createSignedURL(ctx, bucket, objectPath)
	if signErr == nil && strings.TrimSpace(signedURL) != "" {
		return signedURL
	}

	if c.publicBucket != "" && bucket == c.publicBucket {
		publicURL := c.s3.PublicObjectURL(bucket, objectPath)
		if strings.TrimSpace(publicURL) != "" && strings.TrimSpace(publicURL) != strings.TrimSpace(objectPath) {
			return publicURL
		}
	}

	publicURL := c.s3.PublicObjectURL(bucket, objectPath)
	if strings.TrimSpace(publicURL) != "" && strings.TrimSpace(publicURL) != strings.TrimSpace(objectPath) {
		return publicURL
	}

	return fallback
}

func (c *railwayStorageClient) resolveAbsoluteStorageURL(ctx context.Context, rawURL string) (string, bool) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed == nil {
		return rawURL, false
	}

	if parsed.Query().Has("X-Amz-Signature") || parsed.Query().Has("x-amz-signature") {
		return rawURL, true
	}

	if !c.isStorageHost(parsed.Host) {
		return rawURL, false
	}

	objectPath := strings.Trim(strings.TrimSpace(parsed.Path), "/")
	if objectPath == "" {
		return rawURL, true
	}

	type candidate struct {
		bucket string
		path   string
	}
	candidates := make([]candidate, 0, 3)
	seen := map[string]struct{}{}
	appendCandidate := func(bucket string, path string) {
		normalizedBucket := strings.TrimSpace(bucket)
		normalizedPath := strings.Trim(strings.TrimSpace(path), "/")
		if normalizedBucket == "" || normalizedPath == "" {
			return
		}
		key := normalizedBucket + "|" + normalizedPath
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		candidates = append(candidates, candidate{bucket: normalizedBucket, path: normalizedPath})
	}

	if bucket, path, ok := splitBucketAndObjectPath(objectPath); ok {
		appendCandidate(bucket, path)
	}
	if c.privateBucket != "" {
		appendCandidate(c.privateBucket, objectPath)
	}
	if c.publicBucket != "" {
		appendCandidate(c.publicBucket, objectPath)
	}

	for _, item := range candidates {
		resolved := c.resolveBucketObjectURL(ctx, item.bucket, item.path, "")
		if strings.TrimSpace(resolved) != "" {
			return resolved, true
		}
	}

	return rawURL, true
}

func (c *railwayStorageClient) isStorageHost(host string) bool {
	trimmedHost := strings.ToLower(strings.TrimSpace(host))
	if trimmedHost == "" {
		return false
	}
	if colonIndex := strings.Index(trimmedHost, ":"); colonIndex > -1 {
		trimmedHost = strings.TrimSpace(trimmedHost[:colonIndex])
	}

	if strings.Contains(trimmedHost, "storageapi.dev") {
		return true
	}

	endpointHost := ""
	if c != nil && c.s3 != nil {
		if endpointURL := strings.TrimSpace(c.s3.EndpointURL()); endpointURL != "" {
			if parsedEndpoint, parseErr := url.Parse(endpointURL); parseErr == nil && parsedEndpoint != nil {
				endpointHost = strings.ToLower(strings.TrimSpace(parsedEndpoint.Hostname()))
			}
		}
	}

	if endpointHost != "" {
		return trimmedHost == endpointHost
	}

	return false
}

func (c *railwayStorageClient) createSignedURL(ctx context.Context, bucket string, objectPath string) (string, error) {
	if c == nil || c.s3 == nil {
		return "", fmt.Errorf("object storage is not configured")
	}
	_ = ctx
	return c.s3.PresignGetObject(bucket, objectPath, time.Duration(signedURLExpiresInSeconds)*time.Second)
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
