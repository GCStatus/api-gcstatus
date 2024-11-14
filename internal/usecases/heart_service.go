package usecases

import (
	"errors"
	"gcstatus/internal/domain"
	httpErr "gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"gcstatus/internal/resources"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HeartService struct {
	repo ports.HeartRepositry
}

func NewHeartService(repo ports.HeartRepositry) *HeartService {
	return &HeartService{repo: repo}
}

func (h *HeartService) ToggleHeartable(heartableID uint, heartableType string, userID uint) (resources.Response, *httpErr.HttpError) {
	heart, err := h.repo.FindForUser(heartableID, heartableType, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return resources.Response{}, httpErr.NewHttpError(http.StatusInternalServerError, "Failed to find item heart for user.")
	}

	if heart != nil {
		if err := h.repo.Delete(heart.ID); err != nil {
			return resources.Response{}, httpErr.NewHttpError(http.StatusInternalServerError, "Failed to remove the heart from item.")
		}

		return resources.Response{
			Data: gin.H{"message": "Heart removed successfully"},
		}, nil
	}

	newHeart := domain.Heartable{
		UserID:        userID,
		HeartableID:   heartableID,
		HeartableType: heartableType,
	}

	if err := h.repo.Create(&newHeart); err != nil {
		return resources.Response{}, httpErr.NewHttpError(http.StatusInternalServerError, "Failed to heart the given item.")
	}

	return resources.Response{
		Data: gin.H{"message": "Heart added successfully"},
	}, nil
}
