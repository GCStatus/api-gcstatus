package ports

import "gcstatus/internal/domain"

type MissionRepository interface {
	GetAllForUser(userID uint) ([]*domain.Mission, error)
	CompleteMission(userID uint, missionID uint) error
	FindByID(missionID uint) (*domain.Mission, error)
}
