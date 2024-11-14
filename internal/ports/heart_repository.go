package ports

import "gcstatus/internal/domain"

type HeartTogglePayload struct {
	HeartableID   uint   `json:"heartable_id" binding:"required"`
	HeartableType string `json:"heartable_type" binding:"required"`
}

type HeartRepositry interface {
	FindForUser(heartableID uint, heartableType string, userID uint) (*domain.Heartable, error)
	Create(*domain.Heartable) error
	Delete(id uint) error
}
