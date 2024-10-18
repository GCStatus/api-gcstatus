package testutils

import (
	"context"
	"fmt"
	"time"
)

// MockS3Client is a mock implementation of the S3Client interface.
type MockS3Client struct{}

func (m *MockS3Client) GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	return fmt.Sprintf("https://mock-presigned-url.com/%s", objectKey), nil
}

func (m *MockS3Client) RemoveFile(ctx context.Context, fileName string) error {
	return nil
}

func (m *MockS3Client) UploadFile(ctx context.Context, folder, fileName string, fileContent []byte) (string, error) {
	return fmt.Sprintf("https://mock-presigned-url.com/%s/%s", folder, fileName), nil
}
