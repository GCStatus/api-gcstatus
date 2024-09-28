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
	}

	for _, model := range models {
		if err := dbConn.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate model %T: %v", model, err)
		}
	}
}
