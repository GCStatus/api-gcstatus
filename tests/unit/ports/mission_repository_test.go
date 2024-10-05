package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockMissionRepository struct {
	missions        map[uint]*domain.Mission
	userMissions    map[uint][]uint
	userMissionData map[uint]map[uint]domain.UserMission
	requirements    map[uint][]domain.MissionRequirement
	userProgress    map[uint]map[uint]domain.MissionProgress
}

func NewMockMissionRepository() *MockMissionRepository {
	return &MockMissionRepository{
		missions:        make(map[uint]*domain.Mission),
		userMissions:    make(map[uint][]uint),
		userMissionData: make(map[uint]map[uint]domain.UserMission),
		requirements:    make(map[uint][]domain.MissionRequirement),
		userProgress:    make(map[uint]map[uint]domain.MissionProgress),
	}
}

func (m *MockMissionRepository) GetAllForUser(userID uint) ([]domain.Mission, error) {
	var missions []domain.Mission
	missionIDs := m.userMissions[userID]

	for _, missionID := range missionIDs {
		if mission, exists := m.missions[missionID]; exists {
			missions = append(missions, *mission)
		}
	}

	return missions, nil
}

func (m *MockMissionRepository) CreateMission(mission *domain.Mission) error {
	if mission == nil {
		return errors.New("invalid mission data")
	}
	m.missions[mission.ID] = mission
	return nil
}

func (m *MockMissionRepository) AddMissionRequirements(missionID uint, requirements []domain.MissionRequirement) {
	m.requirements[missionID] = requirements
}

func (m *MockMissionRepository) AddUserProgress(userID uint, requirementID uint, progress domain.MissionProgress) {
	if _, ok := m.userProgress[userID]; !ok {
		m.userProgress[userID] = make(map[uint]domain.MissionProgress)
	}
	m.userProgress[userID][requirementID] = progress
}

func (m *MockMissionRepository) AddUserMission(userID, missionID uint) {
	m.userMissions[userID] = append(m.userMissions[userID], missionID)
}

func (m *MockMissionRepository) FindById(id uint) (*domain.Mission, error) {
	for _, mission := range m.missions {
		if mission.ID == id {
			return mission, nil
		}
	}

	return nil, errors.New("mission not found")
}

func (m *MockMissionRepository) CompleteMission(userID uint, missionID uint) error {
	mission, exists := m.missions[missionID]
	if !exists || mission.Status == domain.MissionUnavailable || mission.Status == domain.MissionCanceled {
		return errors.New("mission not found or unavailable")
	}

	if !mission.ForAll {
		userMissions, ok := m.userMissions[userID]
		if !ok || !contains(userMissions, missionID) {
			return errors.New("user is not assigned to this mission")
		}
	}

	if _, ok := m.userMissionData[userID]; !ok {
		m.userMissionData[userID] = make(map[uint]domain.UserMission)
	}

	userMission, ok := m.userMissionData[userID][missionID]
	if !ok {
		userMission = domain.UserMission{
			UserID:    userID,
			MissionID: missionID,
			Completed: false,
		}
		m.userMissionData[userID][missionID] = userMission
	}

	if userMission.Completed {
		return errors.New("mission already completed by user")
	}

	requirements, exists := m.requirements[missionID]
	if !exists {
		return errors.New("no requirements found for this mission")
	}

	for _, req := range requirements {
		userProgress, ok := m.userProgress[userID][req.ID]
		if !ok || !userProgress.Completed {
			return errors.New("mission requirements not yet fully completed")
		}
	}

	userMission.Completed = true
	userMission.LastCompletedAt = time.Now()
	m.userMissionData[userID][missionID] = userMission

	return nil
}

func contains(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func (m *MockMissionRepository) MockMissionRepository_GetAllForUser(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		userID               uint
		expectedMissionCount int
		mockCreateMissions   func(repo *MockMissionRepository)
	}{
		"multiple missions for user 1": {
			userID:               1,
			expectedMissionCount: 2,
			mockCreateMissions: func(repo *MockMissionRepository) {
				if err := repo.CreateMission(&domain.Mission{
					ID:          1,
					Mission:     "Mission 1",
					Description: "Description",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  100,
					Frequency:   "daily",
					ResetTime:   fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the mission: %s", err.Error())
				}
				if err := repo.CreateMission(&domain.Mission{
					ID:          2,
					Mission:     "Mission 2",
					Description: "Description",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  100,
					Frequency:   "one-time",
					ResetTime:   fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the mission: %s", err.Error())
				}

				repo.AddUserMission(1, 1)
				repo.AddUserMission(1, 2)
			},
		},
		"no missions for user 1": {
			userID:               1,
			expectedMissionCount: 0,
			mockCreateMissions:   func(repo *MockMissionRepository) {},
		},
		"missions for user 2": {
			userID:               2,
			expectedMissionCount: 1,
			mockCreateMissions: func(repo *MockMissionRepository) {
				if err := repo.CreateMission(&domain.Mission{
					ID:          3,
					Mission:     "Mission 3",
					Description: "Description",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  100,
					Frequency:   "one-time",
					ResetTime:   fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the mission: %s", err.Error())
				}

				repo.AddUserMission(2, 3)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockMissionRepository()

			tc.mockCreateMissions(mockRepo)

			missions, err := mockRepo.GetAllForUser(tc.userID)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(missions) != tc.expectedMissionCount {
				t.Fatalf("expected %d missions, got %d", tc.expectedMissionCount, len(missions))
			}
		})
	}
}

func TestMockMissionRepository_FindById(t *testing.T) {
	fixedTime := time.Now()

	mockRepo := NewMockMissionRepository()
	if err := mockRepo.CreateMission(&domain.Mission{
		ID:          1,
		Mission:     "Mission 1",
		Description: "Description",
		Status:      "available",
		ForAll:      true,
		Coins:       10,
		Experience:  100,
		Frequency:   "daily",
		ResetTime:   fixedTime,
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}); err != nil {
		t.Fatalf("failed to create the mission: %s", err.Error())
	}

	testCases := map[string]struct {
		missionID   uint
		expectError bool
	}{
		"valid mission ID": {
			missionID:   1,
			expectError: false,
		},
		"invalid mission ID": {
			missionID:   999,
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mission, err := mockRepo.FindById(tc.missionID)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if mission != nil {
					t.Fatalf("expected nil mission, got %v", mission)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mission == nil || mission.ID != tc.missionID {
					t.Fatalf("expected mission ID %d, got %v", tc.missionID, mission)
				}
			}
		})
	}
}

func TestMockMissionRepository_CompleteMission(t *testing.T) {
	mockRepo := &MockMissionRepository{
		missions:        make(map[uint]*domain.Mission),
		userMissions:    make(map[uint][]uint),
		userMissionData: make(map[uint]map[uint]domain.UserMission),
		requirements:    make(map[uint][]domain.MissionRequirement),
		userProgress:    make(map[uint]map[uint]domain.MissionProgress),
	}

	testCases := map[string]struct {
		userID               uint
		missionID            uint
		setupMock            func()
		expectedError        error
		expectedMissionState domain.UserMission
	}{
		"mission not found or unavailable": {
			userID:    1,
			missionID: 100,
			setupMock: func() {
				mockRepo.missions[100] = &domain.Mission{ID: 100, Status: domain.MissionUnavailable}
			},
			expectedError: errors.New("mission not found or unavailable"),
		},
		"user not assigned to mission": {
			userID:    1,
			missionID: 200,
			setupMock: func() {
				mockRepo.missions[200] = &domain.Mission{ID: 200, Status: domain.MissionAvailable, ForAll: false}
			},
			expectedError: errors.New("user is not assigned to this mission"),
		},
		"requirements not completed": {
			userID:    1,
			missionID: 300,
			setupMock: func() {
				mockRepo.missions[300] = &domain.Mission{ID: 300, Status: domain.MissionAvailable, ForAll: true}
				mockRepo.requirements[300] = []domain.MissionRequirement{{ID: 1, Task: "Requirement 1"}}
				mockRepo.userProgress[1] = map[uint]domain.MissionProgress{
					1: {ID: 1, Completed: false},
				}
			},
			expectedError: errors.New("mission requirements not yet fully completed"),
		},
		"mission already completed": {
			userID:    1,
			missionID: 400,
			setupMock: func() {
				mockRepo.missions[400] = &domain.Mission{ID: 400, Status: domain.MissionAvailable, ForAll: true}
				mockRepo.userMissionData[1] = map[uint]domain.UserMission{
					400: {MissionID: 400, UserID: 1, Completed: true},
				}
			},
			expectedError: errors.New("mission already completed by user"),
		},
		"mission successfully completed": {
			userID:    1,
			missionID: 500,
			setupMock: func() {
				mockRepo.missions[500] = &domain.Mission{ID: 500, Status: domain.MissionAvailable, ForAll: true}
				mockRepo.requirements[500] = []domain.MissionRequirement{{ID: 1, Task: "Requirement 1"}}
				mockRepo.userProgress[1] = map[uint]domain.MissionProgress{
					1: {ID: 1, Completed: true},
				}
			},
			expectedError: nil,
			expectedMissionState: domain.UserMission{
				UserID:          1,
				MissionID:       500,
				Completed:       true,
				LastCompletedAt: time.Now(),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockRepo.missions = make(map[uint]*domain.Mission)
			mockRepo.userMissions = make(map[uint][]uint)
			mockRepo.userMissionData = make(map[uint]map[uint]domain.UserMission)
			mockRepo.requirements = make(map[uint][]domain.MissionRequirement)
			mockRepo.userProgress = make(map[uint]map[uint]domain.MissionProgress)

			tc.setupMock()

			err := mockRepo.CompleteMission(tc.userID, tc.missionID)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)

				userMission := mockRepo.userMissionData[tc.userID][tc.missionID]
				assert.True(t, userMission.Completed)
				assert.WithinDuration(t, time.Now(), userMission.LastCompletedAt, time.Second)
			}
		})
	}
}
