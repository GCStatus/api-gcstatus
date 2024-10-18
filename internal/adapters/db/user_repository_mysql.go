package db

import (
	"encoding/json"
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"log"

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
	err := repo.db.Preload("Wallet").
		Preload("Level").
		Preload("Profile").
		Preload("Titles.Title").
		First(&user, id).Error
	return &user, err
}

func (repo *UserRepositoryMySQL) GetUserByIDForAdmin(id uint) (*domain.User, error) {
	var user domain.User
	err := repo.db.
		Preload("Profile").
		Preload("Permissions.Permission").
		Preload("Roles.Role.Permissions.Permission").
		First(&user, id).Error
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

func (repo *UserRepositoryMySQL) FindUserByEmailForAdmin(email string) (*domain.User, error) {
	var user domain.User
	err := repo.db.Model(&domain.User{}).
		Preload("Roles").
		Preload("Permissions").
		Where("email = ?", email).
		First(&user).
		Error
	return &user, err
}

func (repo *UserRepositoryMySQL) UpdateUserPassword(userID uint, hashedPassword string) error {
	return repo.db.Model(&domain.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (repo *UserRepositoryMySQL) UpdateUserNickAndEmail(userID uint, request ports.UpdateNickAndEmailRequest) error {
	updateFields := map[string]interface{}{
		"email":    request.Email,
		"nickname": request.Nickname,
	}

	if err := repo.db.Model(&domain.User{}).Where("id = ?", userID).Updates(updateFields).Error; err != nil {
		return fmt.Errorf("failed to update nick or email: %w", err)
	}

	return nil
}

func (repo *UserRepositoryMySQL) UpdateUserBasics(userID uint, request ports.UpdateUserBasicsRequest) error {
	updateFields := map[string]interface{}{
		"name":      request.Name,
		"birthdate": request.Birthdate,
	}

	if err := repo.db.Model(&domain.User{}).Where("id = ?", userID).Updates(updateFields).Error; err != nil {
		return fmt.Errorf("failed to update user basic informations: %+s", err.Error())
	}

	return nil
}

func (h *UserRepositoryMySQL) AddExperience(
	userID uint,
	experienceGained uint,
	awardTitleToUserFunc func(userID uint, titleID uint) error,
) error {
	var user domain.User
	if err := h.db.Preload("Level").Preload("Wallet").First(&user, userID).Error; err != nil {
		return err
	}

	var levels []domain.Level
	if err := h.db.Preload("Rewards").Order("level ASC").Find(&levels).Error; err != nil {
		return err
	}

	user.Experience += experienceGained

	for {
		var currentLevel *domain.Level
		for _, level := range levels {
			if level.ID == user.LevelID {
				currentLevel = &level
				break
			}
		}

		var nextLevel *domain.Level
		if currentLevel != nil {
			for _, level := range levels {
				if level.Level == currentLevel.Level+1 {
					nextLevel = &level
					break
				}
			}
		}

		if nextLevel == nil {
			break
		}

		if user.Experience >= nextLevel.Experience {
			user.Level = *nextLevel
			user.LevelID = nextLevel.ID
			user.Experience -= nextLevel.Experience
			user.Wallet.Amount += int(nextLevel.Coins)

			h.createTransactionForLevelUp(*nextLevel, user.ID)
			h.createNotificationForLevelUp(*nextLevel, user.ID)

			for _, reward := range nextLevel.Rewards {
				if reward.RewardableType == "titles" {
					if err := awardTitleToUserFunc(user.ID, reward.RewardableID); err != nil {
						return fmt.Errorf("error awarding title: %w", err)
					}

					h.createNotificationForRewardTitle(reward, user.ID)
				}
			}

			if err := h.db.Model(&user).Updates(map[string]any{
				"LevelID":    user.LevelID,
				"Experience": user.Experience,
			}).Error; err != nil {
				return fmt.Errorf("error updating user level in the database: %w", err)
			}

			if err := h.db.Model(&user.Wallet).Updates(map[string]any{
				"Amount": user.Wallet.Amount,
			}).Error; err != nil {
				return fmt.Errorf("failed to update coins for user wallet: %w", err)
			}
		} else {
			break
		}
	}

	if err := h.db.Save(&user).Error; err != nil {
		return fmt.Errorf("error saving user at the end: %w", err)
	}

	return nil
}

func (h *UserRepositoryMySQL) createTransactionForLevelUp(nextLevel domain.Level, userID uint) {
	transaction := &domain.Transaction{
		Amount:            nextLevel.Coins,
		Description:       fmt.Sprintf("Received coins from level up to Level %d.", nextLevel.Level),
		UserID:            userID,
		TransactionTypeID: domain.AdditionTransactionTypeID,
	}

	if err := h.db.Create(&transaction).Error; err != nil {
		log.Printf("Failed to save the transaction of level up: %+v", err)
	}
}

func (h *UserRepositoryMySQL) createNotificationForLevelUp(nextLevel domain.Level, userID uint) {
	notificationContent := &domain.NotificationData{
		Title:     fmt.Sprintf("Congratulations! You've reached level %d!", nextLevel.Level),
		ActionUrl: "/profile/?section=levels",
		Icon:      "FaMedal",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Printf("Failed to marshal notification content: %+v", err)
	}

	notification := &domain.Notification{
		Type:   "NewLevelNotification",
		Data:   string(dataJson),
		UserID: userID,
	}

	if err := h.db.Create(&notification).Error; err != nil {
		log.Printf("Failed to save the title award notification: %+v", err)
	}
}

func (h *UserRepositoryMySQL) createNotificationForRewardTitle(reward domain.Reward, userID uint) {
	notificationContent := &domain.NotificationData{
		Title:     fmt.Sprintf("Congratulations! You've unlocked the title: %d", reward.RewardableID),
		ActionUrl: "/profile/?section=titles",
		Icon:      "FaMedal",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Printf("Failed to marshal notification content: %+v", err)
	}

	notification := &domain.Notification{
		Type:   "NewTitleNotification",
		Data:   string(dataJson),
		UserID: userID,
	}

	if err := h.db.Create(&notification).Error; err != nil {
		log.Printf("Failed to save the title award notification: %+v", err)
	}
}
