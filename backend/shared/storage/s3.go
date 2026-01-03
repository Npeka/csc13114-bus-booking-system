package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

// StorageService defines the interface for file storage operations
type StorageService interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, path string) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
	GetPublicURL(key string) string
}

// S3StorageService implements StorageService using S3-compatible storage
type S3StorageService struct {
	client    *s3.S3
	bucket    string
	cdnDomain string
	region    string
}

// S3Config holds the configuration for S3 storage
type S3Config struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
	Region    string
	UseSSL    bool
	CDNDomain string
}

// NewS3StorageService creates a new S3-compatible storage service
func NewS3StorageService(cfg S3Config) (*S3StorageService, error) {
	// Create AWS session
	awsConfig := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(!cfg.UseSSL),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3StorageService{
		client:    s3.New(sess),
		bucket:    cfg.Bucket,
		cdnDomain: cfg.CDNDomain,
		region:    cfg.Region,
	}, nil
}

// UploadFile uploads a file to S3 and returns the public URL
func (s *S3StorageService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, path string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)
	key := filepath.Join(path, filename)

	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Upload to S3
	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          strings.NewReader(string(fileBytes)),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(fileBytes))),
		ACL:           aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return s.GetPublicURL(key), nil
}

// DeleteFile deletes a file from S3
func (s *S3StorageService) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract key from URL
	key := s.extractKeyFromURL(fileURL)
	if key == "" {
		return fmt.Errorf("invalid file URL")
	}

	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetPublicURL returns the public URL for a given key
func (s *S3StorageService) GetPublicURL(key string) string {
	if s.cdnDomain != "" {
		return fmt.Sprintf("%s/%s", s.cdnDomain, key)
	}
	return fmt.Sprintf("https://%s.%s/%s", s.bucket, s.extractDomainFromEndpoint(), key)
}

// extractKeyFromURL extracts the S3 key from a full URL
func (s *S3StorageService) extractKeyFromURL(fileURL string) string {
	if s.cdnDomain != "" && strings.HasPrefix(fileURL, s.cdnDomain) {
		return strings.TrimPrefix(fileURL, s.cdnDomain+"/")
	}

	// Extract from S3 URL format
	parts := strings.Split(fileURL, "/")
	if len(parts) < 4 {
		return ""
	}

	return strings.Join(parts[3:], "/")
}

// extractDomainFromEndpoint extracts the domain from endpoint URL
func (s *S3StorageService) extractDomainFromEndpoint() string {
	endpoint := s.client.Endpoint
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return endpoint
}
