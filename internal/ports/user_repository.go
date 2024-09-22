package ports

import "gcstatus/internal/domain"

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id uint) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
	FindUserByEmailOrNickname(EmailOrNickname string) (*domain.User, error)
	UpdateUserPassword(userID uint, hashedPassword string) error
	CreateWithProfile(user *domain.User) error
}
