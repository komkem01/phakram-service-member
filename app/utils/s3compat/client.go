package s3compat

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	endpointURL     string
	region          string
	accessKeyID     string
	secretAccessKey string
	httpClient      *http.Client
}

func NewClient(endpointURL string, accessKeyID string, secretAccessKey string, region string, timeout time.Duration) *Client {
	trimmedEndpoint := strings.TrimRight(cleanEnvValue(endpointURL), "/")
	trimmedRegion := cleanEnvValue(region)
	if trimmedRegion == "" {
		trimmedRegion = "auto"
	}
	if strings.Contains(strings.ToLower(endpointHost(trimmedEndpoint)), "storageapi.dev") && strings.EqualFold(trimmedRegion, "auto") {
		trimmedRegion = "us-east-1"
	}
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	return &Client{
		endpointURL:     trimmedEndpoint,
		region:          trimmedRegion,
		accessKeyID:     cleanEnvValue(accessKeyID),
		secretAccessKey: cleanEnvValue(secretAccessKey),
		httpClient:      &http.Client{Timeout: timeout},
	}
}

func cleanEnvValue(value string) string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.TrimSuffix(trimmed, "\r")
	trimmed = strings.Trim(trimmed, "\"'")
	return strings.TrimSpace(trimmed)
}

func (c *Client) Enabled() bool {
	return c != nil && c.endpointURL != "" && c.accessKeyID != "" && c.secretAccessKey != ""
}

func (c *Client) EndpointURL() string {
	if c == nil {
		return ""
	}
	return c.endpointURL
}

func (c *Client) AccessKeyID() string {
	if c == nil {
		return ""
	}
	return c.accessKeyID
}

func (c *Client) SecretAccessKey() string {
	if c == nil {
		return ""
	}
	return c.secretAccessKey
}

func (c *Client) PublicObjectURL(bucket string, objectPath string) string {
	if c == nil || c.endpointURL == "" {
		return strings.TrimSpace(objectPath)
	}
	trimmedBucket := strings.Trim(strings.TrimSpace(bucket), "/")
	trimmedObject := strings.Trim(strings.TrimSpace(objectPath), "/")
	if trimmedBucket == "" || trimmedObject == "" {
		return strings.TrimSpace(objectPath)
	}
	target, err := c.buildTarget(trimmedBucket, trimmedObject)
	if err != nil {
		return strings.TrimSpace(objectPath)
	}
	return target.requestURL
}

func (c *Client) PutObject(ctx context.Context, bucket string, objectPath string, contentType string, payload []byte) error {
	if !c.Enabled() {
		return fmt.Errorf("object storage is not configured")
	}

	trimmedBucket := strings.Trim(strings.TrimSpace(bucket), "/")
	trimmedObject := strings.Trim(strings.TrimSpace(objectPath), "/")
	if trimmedBucket == "" || trimmedObject == "" {
		return fmt.Errorf("bucket and object path are required")
	}

	target, err := c.buildTarget(trimmedBucket, trimmedObject)
	if err != nil {
		return err
	}
	canonicalURI := target.canonicalURI
	requestURL := target.requestURL
	payloadHash := sha256Hex(payload)
	now := time.Now().UTC()
	dateStamp := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	host := target.host

	headers := map[string]string{
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
	}
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := buildCanonicalRequest(http.MethodPut, canonicalURI, "", headers, signedHeaders, payloadHash)
	signature := c.signature(dateStamp, amzDate, canonicalRequest)
	authorization := c.authorizationHeader(dateStamp, signedHeaders, signature)

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, requestURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Host", host)
	request.Header.Set("x-amz-content-sha256", payloadHash)
	request.Header.Set("x-amz-date", amzDate)
	request.Header.Set("Authorization", authorization)
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}
	request.Header.Set("Content-Type", contentType)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		trimmedBody := strings.TrimSpace(string(body))
		if strings.Contains(trimmedBody, "InvalidAccessKeyId") {
			return fmt.Errorf("object storage upload failed: %s (endpoint=%s, region=%s, access_key_id=%s)", trimmedBody, c.endpointURL, c.region, maskAccessKeyID(c.accessKeyID))
		}
		return fmt.Errorf("object storage upload failed: %s", trimmedBody)
	}

	return nil
}

func maskAccessKeyID(value string) string {
	trimmed := cleanEnvValue(value)
	if trimmed == "" {
		return "<empty>"
	}
	if len(trimmed) <= 8 {
		return "****"
	}
	return trimmed[:4] + "..." + trimmed[len(trimmed)-4:]
}

func (c *Client) DeleteObject(ctx context.Context, bucket string, objectPath string) error {
	if !c.Enabled() {
		return nil
	}

	trimmedBucket := strings.Trim(strings.TrimSpace(bucket), "/")
	trimmedObject := strings.Trim(strings.TrimSpace(objectPath), "/")
	if trimmedBucket == "" || trimmedObject == "" {
		return nil
	}

	target, err := c.buildTarget(trimmedBucket, trimmedObject)
	if err != nil {
		return err
	}
	canonicalURI := target.canonicalURI
	requestURL := target.requestURL
	payloadHash := sha256Hex(nil)
	now := time.Now().UTC()
	dateStamp := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	host := target.host

	headers := map[string]string{
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
	}
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := buildCanonicalRequest(http.MethodDelete, canonicalURI, "", headers, signedHeaders, payloadHash)
	signature := c.signature(dateStamp, amzDate, canonicalRequest)
	authorization := c.authorizationHeader(dateStamp, signedHeaders, signature)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Host", host)
	request.Header.Set("x-amz-content-sha256", payloadHash)
	request.Header.Set("x-amz-date", amzDate)
	request.Header.Set("Authorization", authorization)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return fmt.Errorf("object storage delete failed: %s", strings.TrimSpace(string(body)))
	}

	return nil
}

func (c *Client) PresignGetObject(bucket string, objectPath string, expires time.Duration) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("object storage is not configured")
	}

	trimmedBucket := strings.Trim(strings.TrimSpace(bucket), "/")
	trimmedObject := strings.Trim(strings.TrimSpace(objectPath), "/")
	if trimmedBucket == "" || trimmedObject == "" {
		return "", fmt.Errorf("bucket and object path are required")
	}

	if expires <= 0 {
		expires = 15 * time.Minute
	}
	expiresSeconds := int(expires.Seconds())
	if expiresSeconds > 604800 {
		expiresSeconds = 604800
	}

	now := time.Now().UTC()
	dateStamp := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	target, err := c.buildTarget(trimmedBucket, trimmedObject)
	if err != nil {
		return "", err
	}
	host := target.host
	canonicalURI := target.canonicalURI
	credentialScope := c.credentialScope(dateStamp)

	queryParams := map[string]string{
		"X-Amz-Algorithm":     "AWS4-HMAC-SHA256",
		"X-Amz-Credential":    c.accessKeyID + "/" + credentialScope,
		"X-Amz-Date":          amzDate,
		"X-Amz-Expires":       strconv.Itoa(expiresSeconds),
		"X-Amz-SignedHeaders": "host",
	}

	canonicalQuery := canonicalQueryString(queryParams)
	headers := map[string]string{"host": host}
	canonicalRequest := buildCanonicalRequest(http.MethodGet, canonicalURI, canonicalQuery, headers, "host", "UNSIGNED-PAYLOAD")
	signature := c.signature(dateStamp, amzDate, canonicalRequest)

	return target.requestURL + "?" + canonicalQuery + "&X-Amz-Signature=" + signature, nil
}

type requestTarget struct {
	host         string
	canonicalURI string
	requestURL   string
}

func (c *Client) buildTarget(bucket string, objectPath string) (*requestTarget, error) {
	parsed, err := url.Parse(c.endpointURL)
	if err != nil {
		return nil, err
	}

	basePath := strings.TrimSuffix(parsed.Path, "/")

	hostname := parsed.Hostname()
	port := parsed.Port()

	hostnameForRequest := hostname
	canonicalURI := buildCanonicalURI(bucket, objectPath)
	if shouldUseVirtualHostedStyle(hostname) {
		hostnameForRequest = bucket + "." + hostname
		canonicalURI = buildCanonicalObjectURI(objectPath)
	}

	hostForHeader := hostnameForRequest
	if port != "" {
		hostForHeader = hostForHeader + ":" + port
	}

	uriWithBase := canonicalURI
	if basePath != "" {
		uriWithBase = strings.TrimSuffix(basePath, "/") + canonicalURI
	}

	u := &url.URL{
		Scheme: parsed.Scheme,
		Host:   hostForHeader,
		Path:   uriWithBase,
	}

	return &requestTarget{
		host:         hostForHeader,
		canonicalURI: uriWithBase,
		requestURL:   u.String(),
	}, nil
}

func shouldUseVirtualHostedStyle(hostname string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(hostname))
	if trimmed == "" {
		return false
	}
	if strings.Contains(trimmed, "storageapi.dev") {
		return true
	}
	if net.ParseIP(trimmed) != nil {
		return false
	}
	if trimmed == "localhost" || strings.HasSuffix(trimmed, ".localhost") {
		return false
	}
	return false
}

func (c *Client) signature(dateStamp string, amzDate string, canonicalRequest string) string {
	credentialScope := c.credentialScope(dateStamp)
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")

	kDate := hmacSHA256([]byte("AWS4"+c.secretAccessKey), dateStamp)
	kRegion := hmacSHA256(kDate, c.region)
	kService := hmacSHA256(kRegion, "s3")
	kSigning := hmacSHA256(kService, "aws4_request")
	return hex.EncodeToString(hmacSHA256(kSigning, stringToSign))
}

func (c *Client) authorizationHeader(dateStamp string, signedHeaders string, signature string) string {
	return fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		c.accessKeyID,
		c.credentialScope(dateStamp),
		signedHeaders,
		signature,
	)
}

func (c *Client) credentialScope(dateStamp string) string {
	return dateStamp + "/" + c.region + "/s3/aws4_request"
}

func buildCanonicalRequest(method string, canonicalURI string, canonicalQuery string, headers map[string]string, signedHeaders string, payloadHash string) string {
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, strings.ToLower(strings.TrimSpace(key)))
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(":")
		builder.WriteString(strings.TrimSpace(headers[key]))
		builder.WriteString("\n")
	}

	return strings.Join([]string{
		method,
		canonicalURI,
		canonicalQuery,
		builder.String(),
		signedHeaders,
		payloadHash,
	}, "\n")
}

func canonicalQueryString(values map[string]string) string {
	type pair struct {
		key   string
		value string
	}

	pairs := make([]pair, 0, len(values))
	for key, value := range values {
		pairs = append(pairs, pair{key: awsPercentEncode(key), value: awsPercentEncode(value)})
	}

	sort.Slice(pairs, func(i int, j int) bool {
		if pairs[i].key == pairs[j].key {
			return pairs[i].value < pairs[j].value
		}
		return pairs[i].key < pairs[j].key
	})

	parts := make([]string, 0, len(pairs))
	for _, item := range pairs {
		parts = append(parts, item.key+"="+item.value)
	}
	return strings.Join(parts, "&")
}

func endpointHost(endpoint string) string {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return strings.TrimSpace(endpoint)
	}
	return parsed.Host
}

func buildCanonicalObjectURI(objectPath string) string {
	segments := make([]string, 0, 8)
	for _, segment := range strings.Split(strings.Trim(strings.TrimSpace(objectPath), "/"), "/") {
		trimmed := strings.TrimSpace(segment)
		if trimmed == "" {
			continue
		}
		segments = append(segments, awsPathEncode(trimmed))
	}
	if len(segments) == 0 {
		return "/"
	}
	return "/" + strings.Join(segments, "/")
}

func buildCanonicalURI(bucket string, objectPath string) string {
	segments := []string{awsPathEncode(bucket)}
	for _, segment := range strings.Split(strings.Trim(strings.TrimSpace(objectPath), "/"), "/") {
		trimmed := strings.TrimSpace(segment)
		if trimmed == "" {
			continue
		}
		segments = append(segments, awsPathEncode(trimmed))
	}
	return "/" + strings.Join(segments, "/")
}

func awsPathEncode(value string) string {
	encoded := url.PathEscape(value)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

func awsPercentEncode(value string) string {
	encoded := url.QueryEscape(value)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

func sha256Hex(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, payload string) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(payload))
	return mac.Sum(nil)
}
