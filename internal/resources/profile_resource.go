package resources

import (
	"context"
	"gcstatus/internal/domain"
	"gcstatus/pkg/s3"
	"log"
	"time"
)

type ProfileResource struct {
	ID        uint   `json:"id"`
	Share     bool   `json:"share"`
	Photo     string `json:"photo,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Twitch    string `json:"twitch,omitempty"`
	Github    string `json:"github,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformProfile(profile domain.Profile, s3Client s3.S3ClientInterface) *ProfileResource {
	url, err := s3Client.GetPresignedURL(context.TODO(), profile.Photo, time.Hour*2)
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
	}

	return &ProfileResource{
		ID:        profile.ID,
		Share:     profile.Share,
		Photo:     url,
		Phone:     profile.Phone,
		Facebook:  profile.Facebook,
		Instagram: profile.Instagram,
		Twitter:   profile.Twitter,
		Youtube:   profile.Youtube,
		Twitch:    profile.Twitch,
		Github:    profile.Github,
		CreatedAt: profile.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt: profile.UpdatedAt.Format("2006-01-02T15:04:05"),
	}
}
