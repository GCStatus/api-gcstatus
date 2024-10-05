package di

import (
	"gcstatus/internal/domain"
	"log"

	"gorm.io/gorm"
)

func MigrateModels(dbConn *gorm.DB) {
	models := []interface{}{
		&domain.Reward{},
		&domain.Level{},
		&domain.Wallet{},
		&domain.User{},
		&domain.Profile{},
		&domain.PasswordReset{},
		&domain.Title{},
		&domain.TitleRequirement{},
		&domain.TitleProgress{},
		&domain.UserTitle{},
		&domain.TransactionType{},
		&domain.Transaction{},
		&domain.Notification{},
		&domain.Mission{},
		&domain.MissionRequirement{},
		&domain.MissionProgress{},
		&domain.UserMission{},
		&domain.UserMissionAssignment{},
		&domain.Genre{},
		&domain.Tag{},
		&domain.Platform{},
		&domain.Category{},
	}

	for _, model := range models {
		if err := dbConn.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate model %T: %v", model, err)
		}
	}
}
