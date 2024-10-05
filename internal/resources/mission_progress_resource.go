package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type MissionProgressResource struct {
	ID        uint   `json:"id"`
	Progress  uint   `json:"progress"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformMissionProgress(titleProgress domain.MissionProgress) *MissionProgressResource {
	return &MissionProgressResource{
		ID:        titleProgress.ID,
		Progress:  uint(titleProgress.Progress),
		Completed: titleProgress.Completed,
		CreatedAt: utils.FormatTimestamp(titleProgress.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(titleProgress.UpdatedAt),
	}
}

func TransformMissionProgresses(missionProgresses []domain.MissionProgress) []MissionProgressResource {
	var resources []MissionProgressResource

	for _, missionProgress := range missionProgresses {
		resources = append(resources, *TransformMissionProgress(missionProgress))
	}

	return resources
}
