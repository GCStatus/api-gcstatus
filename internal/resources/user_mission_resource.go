package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type UserMissionResource struct {
	Completed       bool   `json:"completed"`
	LastCompletedAt string `json:"last_completed_at,omitempty"`
}

func TransformUserMission(userMission domain.UserMission) *UserMissionResource {
	return &UserMissionResource{
		Completed:       userMission.Completed,
		LastCompletedAt: utils.FormatTimestamp(userMission.LastCompletedAt),
	}
}

func TransformUserMissions(userMissions []domain.UserMission) []*UserMissionResource {
	var resources []*UserMissionResource

	for _, userMission := range userMissions {
		resources = append(resources, TransformUserMission(userMission))
	}

	return resources
}
