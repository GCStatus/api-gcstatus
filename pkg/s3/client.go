package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	envconfig "gcstatus/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3ClientInterface interface {
	UploadFile(ctx context.Context, folder, fileName string, fileContent []byte) (string, error)
	GetPresignedURL(ctx context.Context, fileName string, expiration time.Duration) (string, error)
	RemoveFile(ctx context.Context, fileName string) error
}

type S3Client struct {
	client     *s3.Client
	bucketName string
}

var GlobalS3Client S3ClientInterface

func NewS3Client() *S3Client {
	env := envconfig.LoadConfig()

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(env.AwsBucketRegion))
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config: %v", err))
	}

	s3Client := &S3Client{
		client:     s3.NewFromConfig(awsCfg),
		bucketName: env.AwsBucket,
	}

	return s3Client
}

func (s *S3Client) UploadFile(ctx context.Context, folder, fileName string, fileContent []byte) (string, error) {
	uploader := manager.NewUploader(s.client)

	fullPath := fmt.Sprintf("%s/%s", folder, fileName)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fullPath),
		Body:   bytes.NewReader(fileContent),
		ACL:    types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	return fullPath, nil
}

func (s *S3Client) GetPresignedURL(ctx context.Context, fileName string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.URL, nil
}

func (s *S3Client) RemoveFile(ctx context.Context, fileName string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		if isObjectNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func isObjectNotFoundError(err error) bool {
	var notFoundErr *types.NoSuchKey
	return err != nil && (errors.As(err, &notFoundErr))
}
