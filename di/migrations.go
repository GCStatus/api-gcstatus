package di

import (
	"gcstatus/internal/domain"
	"log"

	"gorm.io/gorm"
)

func MigrateModels(dbConn *gorm.DB) {
	models := []interface{}{
		&domain.Level{},
		&domain.User{},
		&domain.Profile{},
		&domain.PasswordReset{},
	}

	for _, model := range models {
		if err := dbConn.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate model %T: %v", model, err)
		}
	}
}
