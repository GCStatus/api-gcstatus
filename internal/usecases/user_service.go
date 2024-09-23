package usecases

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"gcstatus/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *domain.User) error {
	return s.repo.CreateUser(user)
}

func (s *UserService) CreateWithProfile(user *domain.User) error {
	return s.repo.CreateWithProfile(user)
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) FindUserByEmailOrNickname(emailOrNickname string) (*domain.User, error) {
	user, err := s.repo.FindUserByEmailOrNickname(emailOrNickname)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(emailOrNickname, password string) (*domain.User, error) {
	user, err := s.FindUserByEmailOrNickname(emailOrNickname)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) UpdateUserPassword(userID uint, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := s.repo.UpdateUserPassword(userID, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}
