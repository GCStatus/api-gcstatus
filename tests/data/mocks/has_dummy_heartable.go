package test_mocks

import (
	"gcstatus/internal/domain"
	"testing"

	"gorm.io/gorm"
)

func CreateDummyHeartable(t *testing.T, dbConn *gorm.DB, overrides *domain.Heartable) (*domain.Heartable, error) {
	defaultHeartable := domain.Heartable{
		HeartableID:   1,
		HeartableType: "games",
	}

	if overrides != nil {
		if overrides.HeartableID != 0 {
			defaultHeartable.HeartableID = overrides.HeartableID
		}
		if overrides.HeartableType != "" {
			defaultHeartable.HeartableType = overrides.HeartableType
		}
		if overrides.User.ID != 0 {
			defaultHeartable.User = overrides.User
		} else {
			user, err := CreateDummyUser(t, dbConn, &overrides.User)
			if err != nil {
				t.Fatalf("failed to create dummy user for comment: %+v", err)
			}

			defaultHeartable.User = *user
		}
	}

	if err := dbConn.Create(&defaultHeartable).Error; err != nil {
		return nil, err
	}

	return &defaultHeartable, nil
}
