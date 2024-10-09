package resources

import (
	"context"
	"gcstatus/internal/domain"
	"gcstatus/pkg/s3"
	"gcstatus/pkg/utils"
	"log"
	"time"
)

type UserResource struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Level      uint             `json:"level"`
	Experience uint             `json:"experience"`
	Nickname   string           `json:"nickname"`
	Birthdate  string           `json:"birthdate"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Profile    *ProfileResource `json:"profile,omitempty"`
	Title      *TitleResource   `json:"title,omitempty"`
	Wallet     *WalletResource  `json:"wallet"`
}

type MinimalUserResource struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Photo     *string `json:"photo"`
	Email     string  `json:"email"`
	Nickname  string  `json:"nickname"`
	CreatedAt string  `json:"created_at"`
}

func TransformUser(user domain.User, s3Client s3.S3ClientInterface) UserResource {
	userResource := UserResource{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Nickname:   user.Nickname,
		Experience: user.Experience,
		Birthdate:  utils.FormatTimestamp(user.Birthdate),
		CreatedAt:  utils.FormatTimestamp(user.CreatedAt),
		UpdatedAt:  utils.FormatTimestamp(user.UpdatedAt),
	}

	if user.Profile.ID != 0 {
		userResource.Profile = TransformProfile(user.Profile, s3Client)
	}

	if user.Level.ID != 0 {
		userResource.Level = user.Level.Level
	}

	if user.Wallet.ID != 0 {
		userResource.Wallet = TransformWallet(&user.Wallet)
	}

	if len(user.Titles) > 0 {
		for _, userTitle := range user.Titles {
			if userTitle.Enabled {
				userResource.Title = &TitleResource{
					ID:          userTitle.Title.ID,
					Title:       userTitle.Title.Title,
					Description: userTitle.Title.Description,
				}

				break
			}
		}
	}

	return userResource
}

func TransformUsers(users []domain.User, s3Client s3.S3ClientInterface) []UserResource {
	var resources []UserResource
	for _, user := range users {
		resources = append(resources, TransformUser(user, s3Client))
	}

	return resources
}

func TransformMinimalUser(user domain.User, s3Client s3.S3ClientInterface) MinimalUserResource {
	userResource := MinimalUserResource{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Nickname:  user.Nickname,
		CreatedAt: utils.FormatTimestamp(user.CreatedAt),
	}

	if user.Profile.Photo != "" {
		url, err := s3Client.GetPresignedURL(context.TODO(), user.Profile.Photo, time.Hour*3)
		if err != nil {
			log.Printf("Error generating presigned URL: %v", err)
		} else {
			userResource.Photo = &url
		}
	}

	return userResource
}
