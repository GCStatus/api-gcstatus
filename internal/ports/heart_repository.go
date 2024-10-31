package ports

import "gcstatus/internal/domain"

type HeartRepositry interface {
	FindForUser(heartableID uint, heartableType string, userID uint) (*domain.Heartable, error)
	Create(*domain.Heartable) error
	Delete(id uint) error
}
