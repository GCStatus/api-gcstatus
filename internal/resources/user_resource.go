package resources

import (
	"gcstatus/internal/domain"
)

type UserResource struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Nickname  string           `json:"nickname"`
	Birthdate string           `json:"birthdate"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	Profile   *ProfileResource `json:"profile,omitempty"`
}

func TransformUser(user domain.User) UserResource {
	userResource := UserResource{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Birthdate: user.Birthdate.Format("2006-01-02T15:04:05"),
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	if user.Profile.ID != 0 {
		userResource.Profile = TransformProfile(user.Profile)
	}

	return userResource
}

func TransformUsers(users []domain.User) []UserResource {
	var resources []UserResource
	for _, user := range users {
		resources = append(resources, TransformUser(user))
	}

	return resources
}