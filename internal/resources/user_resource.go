package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/s3"
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
	Wallet     *WalletResource  `json:"wallet"`
}

func TransformUser(user domain.User, s3Client s3.S3ClientInterface) UserResource {
	userResource := UserResource{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Nickname:   user.Nickname,
		Experience: user.Experience,
		Birthdate:  user.Birthdate.Format("2006-01-02T15:04:05"),
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02T15:04:05"),
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

	return userResource
}

func TransformUsers(users []domain.User, s3Client s3.S3ClientInterface) []UserResource {
	var resources []UserResource
	for _, user := range users {
		resources = append(resources, TransformUser(user, s3Client))
	}

	return resources
}
