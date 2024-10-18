package usecases

import (
	"errors"
	"gcstatus/internal/domain"
	self_errors "gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"gcstatus/internal/utils"
	"net/http"

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

func (s *UserService) GetUserByIDForAdmin(id uint) (*domain.User, error) {
	return s.repo.GetUserByIDForAdmin(id)
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) UpdateUserNickAndEmail(userID uint, request ports.UpdateNickAndEmailRequest) error {
	return s.repo.UpdateUserNickAndEmail(userID, request)
}

func (s *UserService) UpdateUserBasics(userID uint, request ports.UpdateUserBasicsRequest) error {
	return s.repo.UpdateUserBasics(userID, request)
}

func (s *UserService) AddExperience(userID uint, experienceAmount uint, awardTitleToUserFunc func(userID uint, titleID uint) error) error {
	return s.repo.AddExperience(userID, experienceAmount, awardTitleToUserFunc)
}

func (s *UserService) FindUserByEmailOrNickname(emailOrNickname string) (*domain.User, error) {
	user, err := s.repo.FindUserByEmailOrNickname(emailOrNickname)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindUserByEmailForAdmin(email string) (*domain.User, error) {
	user, err := s.repo.FindUserByEmailForAdmin(email)
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

func (s *UserService) AuthenticateUserForAdmin(emailOrNickname string, password string) (*domain.User, error) {
	user, err := s.FindUserByEmailForAdmin(emailOrNickname)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	if len(user.Roles) == 0 && len(user.Permissions) == 0 {
		return nil, self_errors.NewHttpError(http.StatusForbidden, "Failed to autheticates user: insuficcient permissions")
	}

	if user.Blocked {
		return nil, self_errors.NewHttpError(http.StatusForbidden, "You are blocked on GCStatus platform. If you think this is an error, please, contact support!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, self_errors.NewHttpError(http.StatusUnauthorized, "Failed to authenticates user: invalid credentials")
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
