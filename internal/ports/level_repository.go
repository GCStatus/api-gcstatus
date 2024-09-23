package ports

import "gcstatus/internal/domain"

type LevelRepository interface {
	GetAll() ([]*domain.Level, error)
	FindById(id uint) (*domain.Level, error)
	FindByLevel(level uint) (*domain.Level, error)
}
