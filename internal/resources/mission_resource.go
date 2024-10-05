package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type MissionResource struct {
	ID                  uint                         `json:"id"`
	Mission             string                       `json:"mission"`
	Description         string                       `json:"description"`
	Coins               uint                         `json:"coins"`
	Experience          uint                         `json:"experience"`
	Status              string                       `json:"status"`
	Frequency           string                       `json:"frequency"`
	ResetTime           string                       `json:"reset_time"`
	CreatedAt           string                       `json:"created_at,omitempty"`
	MissionRequirements []MissionRequirementResource `json:"requirements"`
	UserMission         *UserMissionResource         `json:"user_mission"`
	Rewards             []RewardResource             `json:"rewards"`
}

func TransformMission(mission domain.Mission) MissionResource {
	missionResource := MissionResource{
		ID:          mission.ID,
		Mission:     mission.Mission,
		Description: mission.Description,
		Coins:       mission.Coins,
		Experience:  mission.Experience,
		Status:      mission.Status,
		Frequency:   mission.Frequency,
		ResetTime:   utils.FormatTimestamp(mission.ResetTime),
		CreatedAt:   utils.FormatTimestamp(mission.CreatedAt),
	}

	if len(mission.MissionRequirements) > 0 {
		missionResource.MissionRequirements = TransformMissionRequirements(mission.MissionRequirements)
	} else {
		missionResource.MissionRequirements = []MissionRequirementResource{}
	}

	if len(mission.UserMission) > 0 {
		missionResource.UserMission = TransformUserMission(mission.UserMission[0])
	} else {
		missionResource.UserMission = nil
	}

	if len(mission.Rewards) > 0 {
		missionResource.Rewards = TransformRewards(mission.Rewards)
	} else {
		missionResource.Rewards = []RewardResource{}
	}

	return missionResource
}

func TransformMissions(missions []*domain.Mission) []MissionResource {
	var resources []MissionResource

	for _, mission := range missions {
		resources = append(resources, TransformMission(*mission))
	}

	return resources
}
