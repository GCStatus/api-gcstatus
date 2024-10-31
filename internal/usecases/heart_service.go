package usecases

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type HeartService struct {
	repo ports.HeartRepositry
}

func NewHeartService(repo ports.HeartRepositry) *HeartService {
	return &HeartService{repo: repo}
}

func (h *HeartService) ToggleHeartable(heartableID uint, heartableType string, userID uint) error {
	heart, err := h.repo.FindForUser(heartableID, heartableType, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if heart != nil {
		return h.repo.Delete(heart.ID)
	}

	newHeart := domain.Heartable{
		UserID:        userID,
		HeartableID:   heartableID,
		HeartableType: heartableType,
	}

	return h.repo.Create(&newHeart)
}
