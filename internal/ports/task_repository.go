package ports

import "gcstatus/internal/domain"

type TaskRepository interface {
	GetTitleRequirementsByKey(actionKey string) ([]domain.TitleRequirement, error)
	GetOrCreateTitleProgress(userID, requirementID uint) (*domain.TitleProgress, error)
	UpdateTitleProgress(progress *domain.TitleProgress) error
	UserHasTitle(userID uint, titleID uint) (bool, error)
	AwardTitleToUser(userID uint, titleID uint) error
	GetMissionRequirementsByKey(actionKey string) ([]domain.MissionRequirement, error)
	GetOrCreateMissionProgress(userID, requirementID uint) (*domain.MissionProgress, error)
	UpdateMissionProgress(progress *domain.MissionProgress) error
}
