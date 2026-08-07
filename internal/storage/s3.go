package storage

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client wraps the AWS S3 client for file storage operations.
type S3Client struct {
	client     *s3.Client
	bucketName string
}

// Config holds the S3 connection configuration.
type Config struct {
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

// allowedFolders defines the valid upload destinations.
var allowedFolders = map[string]bool{
	"images":   true,
	"articles": true,
	"projects": true,
}

// NewS3Client creates a new S3 client wrapper.
func NewS3Client(ctx context.Context, cfg Config) (*S3Client, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Endpoint,
			SigningRegion:     cfg.Region,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Client{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

// UploadFile uploads a file to S3 and returns the generated file path.
func (c *S3Client) UploadFile(ctx context.Context, folder, fileName, contentType string, reader io.Reader) (string, error) {
	if !allowedFolders[folder] {
		return "", fmt.Errorf("invalid folder: %s", folder)
	}

	path := generateFilePath(folder, fileName)

	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(path),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("uploading file: %w", err)
	}

	return path, nil
}

// UploadHTML uploads HTML content as a file to S3.
func (c *S3Client) UploadHTML(ctx context.Context, folder, fileName, htmlContent string) (string, error) {
	if !allowedFolders[folder] {
		return "", fmt.Errorf("invalid folder: %s", folder)
	}

	path := generateFilePath(folder, fileName)

	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(path),
		Body:        strings.NewReader(htmlContent),
		ContentType: aws.String("text/html; charset=utf-8"),
	})
	if err != nil {
		return "", fmt.Errorf("uploading html: %w", err)
	}

	return path, nil
}

// DeleteFile removes a single file from S3.
func (c *S3Client) DeleteFile(ctx context.Context, filePath string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("deleting file: %w", err)
	}
	return nil
}

// DeleteFiles removes multiple files from S3 in parallel.
func (c *S3Client) DeleteFiles(ctx context.Context, filePaths []string) error {
	for _, path := range filePaths {
		if err := c.DeleteFile(ctx, path); err != nil {
			return err
		}
	}
	return nil
}

// GetFile retrieves a file from S3 and returns its body and content type.
func (c *S3Client) GetFile(ctx context.Context, filePath string) (io.ReadCloser, string, error) {
	output, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, "", fmt.Errorf("getting file: %w", err)
	}

	contentType := "application/octet-stream"
	if output.ContentType != nil {
		contentType = *output.ContentType
	}

	return output.Body, contentType, nil
}

func generateFilePath(folder, fileName string) string {
	ts := time.Now().UnixMilli()
	ext := filepath.Ext(fileName)
	sanitized := sanitizeFileName(strings.TrimSuffix(fileName, ext))
	nanoID := randomNanoID(8)
	return fmt.Sprintf("%s/%d_%s_%s%s", folder, ts, nanoID, sanitized, ext)
}

func sanitizeFileName(name string) string {
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
	}
	s := result.String()
	s = strings.Trim(s, "_")
	if s == "" {
		s = "file"
	}
	return s
}

func randomNanoID(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// Fallback: use time-based seed
		for i := range b {
			b[i] = alphabet[time.Now().UnixNano()%int64(len(alphabet))]
		}
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}
