package test_mocks

import (
	"gcstatus/internal/domain"
	"testing"

	"gorm.io/gorm"
)

func CreateDummyLevel(t *testing.T, dbConn *gorm.DB, overrides *domain.Level) (*domain.Level, error) {
	defaultLevel := domain.Level{
		Level:      1,
		Experience: 50,
		Coins:      100,
		Rewards:    []domain.Reward{},
	}

	if overrides != nil {
		if overrides.Level != 0 {
			defaultLevel.Level = overrides.Level
		}
		if overrides.Experience != 0 {
			defaultLevel.Experience = overrides.Experience
		}
		if overrides.Coins != 0 {
			defaultLevel.Coins = overrides.Coins
		}
		if len(overrides.Rewards) > 0 {
			defaultLevel.Rewards = overrides.Rewards
		}
	}

	if err := dbConn.Create(&defaultLevel).Error; err != nil {
		return nil, err
	}

	return &defaultLevel, nil
}
