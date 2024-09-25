package tests

import (
	"context"
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockS3Client struct{}

func (m *MockS3Client) GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	return fmt.Sprintf("https://mock-presigned-url.com/%s", objectKey), nil
}

func (m *MockS3Client) RemoveFile(ctx context.Context, fileName string) error {
	return nil
}

func (s *MockS3Client) UploadFile(ctx context.Context, folder, fileName string, fileContent []byte) (string, error) {
	return fmt.Sprintf("https://mock-presigned-url.com/%s/%s", folder, fileName), nil
}

func TestTransformProfile(t *testing.T) {
	staticTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    domain.Profile
		expected resources.ProfileResource
	}{
		"normal profile": {
			input: domain.Profile{
				ID:        1,
				Share:     true,
				Photo:     "photo-key",
				Phone:     "5511928342813",
				Facebook:  "https://facebook.com/any",
				Instagram: "https://instagram.com/any",
				Twitter:   "https://twitter.com/any",
				Youtube:   "https://youtube.com/any",
				Twitch:    "https://twitch.com/any",
				Github:    "https://github.com/any",
				CreatedAt: staticTime,
				UpdatedAt: staticTime,
				UserID:    1,
			},
			expected: resources.ProfileResource{
				ID:        1,
				Share:     true,
				Photo:     "https://mock-presigned-url.com/photo-key",
				Phone:     "5511928342813",
				Facebook:  "https://facebook.com/any",
				Instagram: "https://instagram.com/any",
				Twitter:   "https://twitter.com/any",
				Youtube:   "https://youtube.com/any",
				Twitch:    "https://twitch.com/any",
				Github:    "https://github.com/any",
				CreatedAt: staticTime.Format("2006-01-02T15:04:05"),
				UpdatedAt: staticTime.Format("2006-01-02T15:04:05"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockS3Client := &MockS3Client{}
			profileResource := resources.TransformProfile(test.input, mockS3Client)

			assert.Equal(t, test.expected.ID, profileResource.ID)
			assert.Equal(t, test.expected.Share, profileResource.Share)
			assert.Equal(t, test.expected.Photo, profileResource.Photo)
			assert.Equal(t, test.expected.Phone, profileResource.Phone)
			assert.Equal(t, test.expected.Facebook, profileResource.Facebook)
			assert.Equal(t, test.expected.Instagram, profileResource.Instagram)
			assert.Equal(t, test.expected.Twitter, profileResource.Twitter)
			assert.Equal(t, test.expected.Youtube, profileResource.Youtube)
			assert.Equal(t, test.expected.Twitch, profileResource.Twitch)
			assert.Equal(t, test.expected.Github, profileResource.Github)
			assert.Equal(t, test.expected.CreatedAt, profileResource.CreatedAt)
			assert.Equal(t, test.expected.UpdatedAt, profileResource.UpdatedAt)
		})
	}
}
