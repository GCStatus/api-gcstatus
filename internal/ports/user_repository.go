package ports

import (
	"gcstatus/internal/domain"
)

type UpdateNickAndEmailRequest struct {
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserBasicsRequest struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id uint) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	FindUserByEmailOrNickname(EmailOrNickname string) (*domain.User, error)
	UpdateUserPassword(userID uint, hashedPassword string) error
	CreateWithProfile(user *domain.User) error
	UpdateUserNickAndEmail(userID uint, request UpdateNickAndEmailRequest) error
	UpdateUserBasics(userID uint, request UpdateUserBasicsRequest) error
	AddExperience(userID uint, experienceAmount uint, awardTitleToUserFunc func(userID uint, titleID uint) error) error
}
