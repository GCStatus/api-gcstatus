package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type UserRepositoryMySQL struct {
	db *gorm.DB
}

func NewUserRepositoryMySQL(db *gorm.DB) ports.UserRepository {
	return &UserRepositoryMySQL{db: db}
}

func (repo *UserRepositoryMySQL) CreateUser(user *domain.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepositoryMySQL) CreateWithProfile(user *domain.User) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		user.Profile.UserID = user.ID
		if err := tx.Create(&user.Profile).Error; err != nil {
			return err
		}

		user.Wallet.UserID = user.ID
		if err := tx.Create(&user.Wallet).Error; err != nil {
			return err
		}

		return nil
	})
}

func (repo *UserRepositoryMySQL) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := repo.db.Preload("Wallet").Preload("Level").Preload("Profile").First(&user, id).Error
	return &user, err
}

func (repo *UserRepositoryMySQL) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	err := repo.db.Find(&users).Error
	return users, err
}

func (repo *UserRepositoryMySQL) FindUserByEmailOrNickname(emailOrNickname string) (*domain.User, error) {
	var user domain.User
	err := repo.db.Where("nickname = ? OR email = ?", emailOrNickname, emailOrNickname).First(&user).Error
	return &user, err
}

func (repo *UserRepositoryMySQL) UpdateUserPassword(userID uint, hashedPassword string) error {
	return repo.db.Model(&domain.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
