package data_test

import (
	"gcstatus/internal/domain"
	"testing"

	"gorm.io/gorm"
)

func CreateDefaultLevels(dbConn *gorm.DB) error {
	levels := []domain.Level{
		{ID: 1, Level: 1, Experience: 0, Coins: 0},
		{ID: 2, Level: 2, Experience: 50, Coins: 100},
		{ID: 3, Level: 3, Experience: 100, Coins: 200},
	}
	for _, level := range levels {
		if err := dbConn.FirstOrCreate(&level, level).Error; err != nil {
			return err
		}
	}
	return nil
}

func Seed(t *testing.T, dbConn *gorm.DB) {
	if err := CreateDefaultLevels(dbConn); err != nil {
		t.Fatalf("failed to seed database: %+v", err)
	}
}
