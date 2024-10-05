package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type MissionService struct {
	repo ports.MissionRepository
}

func NewMissionService(repo ports.MissionRepository) *MissionService {
	return &MissionService{repo: repo}
}

func (h *MissionService) FindByID(missionID uint) (*domain.Mission, error) {
	return h.repo.FindByID(missionID)
}

func (h *MissionService) GetAllForUser(userID uint) ([]*domain.Mission, error) {
	return h.repo.GetAllForUser(userID)
}

func (h *MissionService) CompleteMission(userID uint, missionID uint) error {
	return h.repo.CompleteMission(userID, missionID)
}
