package resources_admin

import (
	"context"
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"log"
	"time"
)

type UserResource struct {
	ID          uint                       `json:"id"`
	Name        string                     `json:"name"`
	Email       string                     `json:"email"`
	Nickname    string                     `json:"nickname"`
	Birthdate   string                     `json:"birthdate"`
	CreatedAt   string                     `json:"created_at"`
	UpdatedAt   string                     `json:"updated_at"`
	Profile     *resources.ProfileResource `json:"profile"`
	Roles       []RoleResource             `json:"roles"`
	Permissions []PermissionResource       `json:"permissions"`
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
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Birthdate:   utils.FormatTimestamp(user.Birthdate),
		CreatedAt:   utils.FormatTimestamp(user.CreatedAt),
		UpdatedAt:   utils.FormatTimestamp(user.UpdatedAt),
		Permissions: []PermissionResource{},
		Roles:       []RoleResource{},
	}

	if user.Profile.ID != 0 {
		userResource.Profile = resources.TransformProfile(user.Profile, s3Client)
	}

	if len(user.Roles) > 0 {
		for _, roleable := range user.Roles {
			userResource.Roles = append(userResource.Roles, TransformRole(roleable.Role))
		}
	}

	if len(user.Permissions) > 0 {
		for _, permissionable := range user.Permissions {
			userResource.Permissions = append(userResource.Permissions, TransformPermission(permissionable.Permission))
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
