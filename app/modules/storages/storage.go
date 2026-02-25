package storages

import (
	"context"
	"fmt"
	"os"
	"phakram/app/utils/s3compat"
	"strings"
	"time"
)

const signedURLExpiresInSeconds = 60 * 60

type supabaseStorageClient struct {
	s3            *s3compat.Client
	publicBucket  string
	privateBucket string
}

func newSupabaseStorageClient(conf SupabaseConfig) *supabaseStorageClient {
	endpointURL := strings.TrimRight(strings.TrimSpace(conf.URL), "/")
	if endpointURL == "" {
		endpointURL = strings.TrimRight(firstNonEmptyEnv("OBJECT_ENDPOINT_URL", "SUPABASE_URL"), "/")
	}

	secretAccessKey := strings.TrimSpace(conf.ServiceRoleKey)
	if secretAccessKey == "" {
		secretAccessKey = firstNonEmptyEnv("OBJECT_SECRET_ACCESS_KEY", "SUPABASE_SERVICE_ROLE_KEY", "SUPABASE_SERVICE_KEY", "SUPABASE_ANON_KEY")
	}

	accessKeyID := firstNonEmptyEnv("OBJECT_ACCESS_KEY_ID")
	region := firstNonEmptyEnv("OBJECT_REGION", "AWS_REGION", "AWS_DEFAULT_REGION")

	publicBucket := strings.TrimSpace(conf.PublicBucket)
	if publicBucket == "" {
		publicBucket = firstNonEmptyEnv("OBJECT_PUBLIC_BUCKET", "SUPABASE_PUBLIC_BUCKET")
	}

	privateBucket := strings.TrimSpace(conf.PrivateBucket)
	if privateBucket == "" {
		privateBucket = firstNonEmptyEnv("OBJECT_PRIVATE_BUCKET", "SUPABASE_PRIVATE_BUCKET")
	}

	return &supabaseStorageClient{
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

func (c *supabaseStorageClient) enabled() bool {
	return c != nil && c.s3 != nil && c.s3.Enabled()
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
		return c.s3.PublicObjectURL(bucket, objectPath), nil
	}

	if c.privateBucket != "" && bucket == c.privateBucket {
		return c.createSignedURL(ctx, bucket, objectPath)
	}

	return trimmed, nil
}

func (c *supabaseStorageClient) createSignedURL(ctx context.Context, bucket string, objectPath string) (string, error) {
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
