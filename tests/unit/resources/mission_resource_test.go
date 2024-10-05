package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"testing"
	"time"
)

func TestTransformMission(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputMission domain.Mission
		expected     resources.MissionResource
	}{
		"normal title": {
			inputMission: domain.Mission{
				ID:          1,
				Mission:     "Mission 1",
				Description: "Mission 1",
				Status:      "available",
				ForAll:      true,
				Coins:       10,
				Experience:  100,
				Frequency:   "one-time",
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Rewards: []domain.Reward{
					{
						ID:             1,
						SourceableID:   1,
						SourceableType: "missions",
						RewardableID:   1,
						RewardableType: "titles",
						CreatedAt:      fixedTime,
						UpdatedAt:      fixedTime,
					},
				},
				UserMission: []domain.UserMission{
					{
						ID:              1,
						Completed:       true,
						LastCompletedAt: fixedTime,
						CreatedAt:       fixedTime,
						UpdatedAt:       fixedTime,
					},
				},
				MissionRequirements: []domain.MissionRequirement{
					{
						ID:          1,
						Task:        "Do something",
						Key:         "do_something",
						Goal:        10,
						Description: "Do something.",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
			},
			expected: resources.MissionResource{
				ID:          1,
				Mission:     "Mission 1",
				Description: "Mission 1",
				Status:      "available",
				Coins:       10,
				Experience:  100,
				Frequency:   "one-time",
				ResetTime:   utils.FormatTimestamp(fixedTime),
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				Rewards: resources.TransformRewards([]domain.Reward{
					{
						ID:             1,
						SourceableID:   1,
						SourceableType: "missions",
						RewardableID:   1,
						RewardableType: "titles",
						CreatedAt:      fixedTime,
						UpdatedAt:      fixedTime,
					},
				}),
				UserMission: resources.TransformUserMission(domain.UserMission{
					ID:              1,
					Completed:       true,
					LastCompletedAt: fixedTime,
					CreatedAt:       fixedTime,
					UpdatedAt:       fixedTime,
				}),
				MissionRequirements: resources.TransformMissionRequirements([]domain.MissionRequirement{
					{
						ID:          1,
						Task:        "Do something",
						Key:         "do_something",
						Goal:        10,
						Description: "Do something.",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				}),
			},
		},
		"missing requirements": {
			inputMission: domain.Mission{
				ID:          2,
				Mission:     "Mission 2",
				Description: "Mission 2",
				Status:      "available",
				ForAll:      true,
				Coins:       10,
				Experience:  100,
				Frequency:   "one-time",
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Rewards: []domain.Reward{
					{
						ID:             1,
						SourceableID:   1,
						SourceableType: "missions",
						RewardableID:   1,
						RewardableType: "titles",
						CreatedAt:      fixedTime,
						UpdatedAt:      fixedTime,
					},
				},
				UserMission: []domain.UserMission{
					{
						ID:              1,
						Completed:       true,
						LastCompletedAt: fixedTime,
						CreatedAt:       fixedTime,
						UpdatedAt:       fixedTime,
					},
				},
			},
			expected: resources.MissionResource{
				ID:          2,
				Mission:     "Mission 2",
				Description: "Mission 2",
				Status:      "available",
				Coins:       10,
				Experience:  100,
				Frequency:   "one-time",
				ResetTime:   utils.FormatTimestamp(fixedTime),
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				Rewards: resources.TransformRewards([]domain.Reward{
					{
						ID:             1,
						SourceableID:   1,
						SourceableType: "missions",
						RewardableID:   1,
						RewardableType: "titles",
						CreatedAt:      fixedTime,
						UpdatedAt:      fixedTime,
					},
				}),
				UserMission: resources.TransformUserMission(domain.UserMission{
					ID:              1,
					Completed:       true,
					LastCompletedAt: fixedTime,
					CreatedAt:       fixedTime,
					UpdatedAt:       fixedTime,
				}),
				MissionRequirements: []resources.MissionRequirementResource{},
			},
		},
		"missing requirements and rewards": {
			inputMission: domain.Mission{
				ID:          2,
				Mission:     "Mission 2",
				Description: "Mission 2",
				Status:      "available",
				ForAll:      true,
				Coins:       10,
				Experience:  100,
				Frequency:   "one-time",
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				UserMission: []domain.UserMission{
					{
						ID:              1,
						Completed:       true,
						LastCompletedAt: fixedTime,
						CreatedAt:       fixedTime,
						UpdatedAt:       fixedTime,
					},
				},
			},
			expected: resources.MissionResource{
				ID:          2,
				Mission:     "Mission 2",
				Description: "Mission 2",
				Status:      "available",
				Coins:       10,
				Experience:  100,
				Frequency:   "one-time",
				ResetTime:   utils.FormatTimestamp(fixedTime),
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				Rewards:     []resources.RewardResource{},
				UserMission: resources.TransformUserMission(domain.UserMission{
					ID:              1,
					Completed:       true,
					LastCompletedAt: fixedTime,
					CreatedAt:       fixedTime,
					UpdatedAt:       fixedTime,
				}),
				MissionRequirements: []resources.MissionRequirementResource{},
			},
		},
		"missing requirements, rewards and user mission": {
			inputMission: domain.Mission{
				ID:          2,
				Mission:     "Mission 2",
				Description: "Mission 2",
				Status:      "available",
				ForAll:      true,
				Coins:       10,
				Experience:  100,
				Frequency:   "one-time",
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				UserMission: []domain.UserMission{
					{
						ID:              1,
						Completed:       true,
						LastCompletedAt: fixedTime,
						CreatedAt:       fixedTime,
						UpdatedAt:       fixedTime,
					},
				},
			},
			expected: resources.MissionResource{
				ID:                  2,
				Mission:             "Mission 2",
				Description:         "Mission 2",
				Status:              "available",
				Coins:               10,
				Experience:          100,
				Frequency:           "one-time",
				ResetTime:           utils.FormatTimestamp(fixedTime),
				CreatedAt:           utils.FormatTimestamp(fixedTime),
				Rewards:             []resources.RewardResource{},
				UserMission:         nil,
				MissionRequirements: []resources.MissionRequirementResource{},
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			missionResource := resources.TransformMission(test.inputMission)

			if missionResource.ID != test.expected.ID {
				t.Errorf("Expected ID %d, got %d", test.expected.ID, missionResource.ID)
			}
			if missionResource.Mission != test.expected.Mission {
				t.Errorf("Expected Mission %s, got %s", test.expected.Mission, missionResource.Mission)
			}
			if missionResource.Description != test.expected.Description {
				t.Errorf("Expected Description %s, got %s", test.expected.Description, missionResource.Description)
			}
			if missionResource.Coins != test.expected.Coins {
				t.Errorf("Expected Coins %v, got %v", test.expected.Coins, missionResource.Coins)
			}
			if missionResource.Experience != test.expected.Experience {
				t.Errorf("Expected Experience %v, got %v", test.expected.Experience, missionResource.Experience)
			}
			if missionResource.Status != test.expected.Status {
				t.Errorf("Expected Status %s, got %s", test.expected.Status, missionResource.Status)
			}
			if missionResource.Frequency != test.expected.Frequency {
				t.Errorf("Expected Frequency %s, got %s", test.expected.Frequency, missionResource.Frequency)
			}
			if missionResource.ResetTime != test.expected.ResetTime {
				t.Errorf("Expected ResetTime %s, got %s", test.expected.ResetTime, missionResource.ResetTime)
			}
			if missionResource.CreatedAt != test.expected.CreatedAt {
				t.Errorf("Expected CreatedAt %s, got %s", test.expected.CreatedAt, missionResource.CreatedAt)
			}

			for _, tr := range test.inputMission.MissionRequirements {
				missionRequirementResource := resources.TransformMissionRequirement(tr)

				if missionRequirementResource.ID != tr.ID {
					t.Errorf("Expected ID %d, got %d", tr.ID, missionRequirementResource.ID)
				}
				if missionRequirementResource.Task != tr.Task {
					t.Errorf("Expected Task %s, got %s", tr.Task, missionRequirementResource.Task)
				}
				if missionRequirementResource.Description != tr.Description {
					t.Errorf("Expected Description %s, got %s", tr.Description, missionRequirementResource.Description)
				}
				if missionRequirementResource.Goal != tr.Goal {
					t.Errorf("Expected Goal %d, got %d", tr.Goal, missionRequirementResource.Goal)
				}
				if missionRequirementResource.CreatedAt != utils.FormatTimestamp(tr.CreatedAt) {
					t.Errorf("Expected CreatedAt %s, got %s", tr.CreatedAt, missionRequirementResource.CreatedAt)
				}
			}

			for _, rw := range test.inputMission.Rewards {
				missionReward := resources.TransformReward(rw)

				if missionReward.ID != rw.ID {
					t.Errorf("Expected ID %d, got %d", rw.ID, missionReward.ID)
				}
				if missionReward.RewardableType != rw.RewardableType {
					t.Errorf("Expected RewardableType %s, got %s", rw.RewardableType, missionReward.RewardableType)
				}
				if missionReward.SourceableType != rw.SourceableType {
					t.Errorf("Expected SourceableType %s, got %s", rw.SourceableType, missionReward.SourceableType)
				}
			}
		})
	}
}

func TestTransformMissions(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputMissions []*domain.Mission
		expected      []resources.MissionResource
	}{
		"multiple titles": {
			inputMissions: []*domain.Mission{
				{
					ID:          1,
					Mission:     "Mission 1",
					Description: "Mission 1",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  100,
					Frequency:   "one-time",
					ResetTime:   fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
					MissionRequirements: []domain.MissionRequirement{
						{
							ID:          1,
							Task:        "Do something",
							Key:         "do_something",
							Goal:        10,
							Description: "Do something.",
							CreatedAt:   fixedTime,
							UpdatedAt:   fixedTime,
						},
					},
				},
				{
					ID:          2,
					Mission:     "Mission 2",
					Description: "Mission 2",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  100,
					Frequency:   "one-time",
					ResetTime:   fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				},
			},
			expected: []resources.MissionResource{
				{
					ID:          1,
					Mission:     "Mission 1",
					Description: "Mission 1",
					Status:      "available",
					Coins:       10,
					Experience:  100,
					Frequency:   "one-time",
					ResetTime:   utils.FormatTimestamp(fixedTime),
					CreatedAt:   utils.FormatTimestamp(fixedTime),
					MissionRequirements: resources.TransformMissionRequirements([]domain.MissionRequirement{
						{
							ID:          1,
							Task:        "Do something",
							Key:         "do_something",
							Goal:        10,
							Description: "Do something.",
							CreatedAt:   fixedTime,
							UpdatedAt:   fixedTime,
						},
					}),
				},
				{
					ID:          2,
					Mission:     "Mission 2",
					Description: "Mission 2",
					Status:      "available",
					Coins:       10,
					Experience:  100,
					Frequency:   "one-time",
					ResetTime:   utils.FormatTimestamp(fixedTime),
					CreatedAt:   utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			missionResources := resources.TransformMissions(test.inputMissions)

			if len(missionResources) != len(test.expected) {
				t.Errorf("Expected %d titles, got %d", len(test.expected), len(missionResources))
			}

			for i := range missionResources {
				if missionResources[i].ID != test.expected[i].ID {
					t.Errorf("Expected ID %d, got %d", test.expected[i].ID, missionResources[i].ID)
				}
				if missionResources[i].Coins != test.expected[i].Coins {
					t.Errorf("Expected Coins %d, got %d", test.expected[i].Coins, missionResources[i].Coins)
				}
				if missionResources[i].Description != test.expected[i].Description {
					t.Errorf("Expected Description %s, got %s", test.expected[i].Description, missionResources[i].Description)
				}
				if missionResources[i].Experience != test.expected[i].Experience {
					t.Errorf("Expected Experience %d, got %d", test.expected[i].Experience, missionResources[i].Experience)
				}
				if missionResources[i].Status != test.expected[i].Status {
					t.Errorf("Expected Status %s, got %s", test.expected[i].Status, missionResources[i].Status)
				}
				if missionResources[i].CreatedAt != test.expected[i].CreatedAt {
					t.Errorf("Expected CreatedAt %s, got %s", test.expected[i].CreatedAt, missionResources[i].CreatedAt)
				}
			}
		})
	}
}
