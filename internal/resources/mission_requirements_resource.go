package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type MissionRequirementResource struct {
	ID              uint                     `json:"id"`
	Task            string                   `json:"task"`
	Description     string                   `json:"description"`
	Goal            int                      `json:"goal"`
	CreatedAt       string                   `json:"created_at"`
	UpdatedAt       string                   `json:"updated_at"`
	MissionProgress *MissionProgressResource `json:"progress"`
}

func TransformMissionRequirement(missionRequirement domain.MissionRequirement) MissionRequirementResource {
	missionRequirementResource := MissionRequirementResource{
		ID:          missionRequirement.ID,
		Task:        missionRequirement.Task,
		Description: missionRequirement.Description,
		Goal:        missionRequirement.Goal,
		CreatedAt:   utils.FormatTimestamp(missionRequirement.CreatedAt),
		UpdatedAt:   utils.FormatTimestamp(missionRequirement.UpdatedAt),
	}

	if missionRequirement.MissionProgress.ID != 0 {
		missionRequirementResource.MissionProgress = TransformMissionProgress(missionRequirement.MissionProgress)
	} else {
		missionRequirementResource.MissionProgress = nil
	}

	return missionRequirementResource
}

func TransformMissionRequirements(missionRequirements []domain.MissionRequirement) []MissionRequirementResource {
	if len(missionRequirements) == 0 {
		return []MissionRequirementResource{}
	}

	var resources []MissionRequirementResource

	for _, missionRequirement := range missionRequirements {
		resources = append(resources, TransformMissionRequirement(missionRequirement))
	}

	return resources
}
