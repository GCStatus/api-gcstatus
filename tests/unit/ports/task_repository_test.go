package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockTaskRepository simulates TaskRepository behavior
type MockTaskRepository struct {
	Requirements []domain.TitleRequirement
	Progress     map[uint]*domain.TitleProgress
	UserTitles   map[uint]map[uint]bool // userID -> titleID -> hasTitle
}

// Ensure MockTaskRepository implements the TaskRepository interface
var _ ports.TaskRepository = &MockTaskRepository{}

// NewMockTaskRepository creates a new mock repository
func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		Requirements: []domain.TitleRequirement{},
		Progress:     make(map[uint]*domain.TitleProgress),
		UserTitles:   make(map[uint]map[uint]bool),
	}
}

func (m *MockTaskRepository) GetTitleRequirementsByKey(actionKey string) ([]domain.TitleRequirement, error) {
	if len(m.Requirements) == 0 {
		return nil, errors.New("no requirements found")
	}
	return m.Requirements, nil
}

func (m *MockTaskRepository) GetOrCreateTitleProgress(userID, requirementID uint) (*domain.TitleProgress, error) {
	progress, exists := m.Progress[requirementID]
	if !exists {
		// Simulate creating new progress
		progress = &domain.TitleProgress{
			UserID:             userID,
			TitleRequirementID: requirementID,
			Progress:           0,
			Completed:          false,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		m.Progress[requirementID] = progress
	}
	return progress, nil
}

func (m *MockTaskRepository) UpdateTitleProgress(progress *domain.TitleProgress) error {
	m.Progress[progress.TitleRequirementID] = progress
	return nil
}

func (m *MockTaskRepository) UserHasTitle(userID uint, titleID uint) (bool, error) {
	userTitles, exists := m.UserTitles[userID]
	if !exists {
		return false, nil
	}
	hasTitle := userTitles[titleID]
	return hasTitle, nil
}

func (m *MockTaskRepository) AwardTitleToUser(userID uint, titleID uint) error {
	if _, exists := m.UserTitles[userID]; !exists {
		m.UserTitles[userID] = make(map[uint]bool)
	}
	m.UserTitles[userID][titleID] = true
	return nil
}

func MockTaskRepository_TestGetTitleRequirementsByKey(t *testing.T) {
	mockRepo := NewMockTaskRepository()

	testCases := []struct {
		name          string
		setup         func()
		expectedError string
	}{
		{
			name: "no requirements",
			setup: func() {
				mockRepo.Requirements = []domain.TitleRequirement{}
			},
			expectedError: "no requirements found",
		},
		{
			name: "requirements found",
			setup: func() {
				mockRepo.Requirements = []domain.TitleRequirement{
					{ID: 1, TitleID: 1, Key: "test_action", Goal: 10},
				}
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			_, err := mockRepo.GetTitleRequirementsByKey("test_action")
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func MockTaskRepository_TestGetOrCreateTitleProgress(t *testing.T) {
	mockRepo := NewMockTaskRepository()

	testCases := []struct {
		name           string
		userID         uint
		requirementID  uint
		setup          func()
		expectedExists bool
	}{
		{
			name:          "create new progress",
			userID:        1,
			requirementID: 100,
			setup: func() {
				mockRepo.Progress = make(map[uint]*domain.TitleProgress)
			},
			expectedExists: false,
		},
		{
			name:          "return existing progress",
			userID:        1,
			requirementID: 100,
			setup: func() {
				mockRepo.Progress[100] = &domain.TitleProgress{
					UserID:             1,
					TitleRequirementID: 100,
					Progress:           5,
					Completed:          false,
				}
			},
			expectedExists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			progress, err := mockRepo.GetOrCreateTitleProgress(tc.userID, tc.requirementID)
			assert.NoError(t, err)
			if tc.expectedExists {
				assert.Equal(t, 5, progress.Progress)
			} else {
				assert.Equal(t, 0, progress.Progress)
			}
		})
	}
}

func MockTaskRepository_TestUpdateTitleProgress(t *testing.T) {
	mockRepo := NewMockTaskRepository()

	testCases := []struct {
		name           string
		progress       domain.TitleProgress
		setup          func()
		expectedAmount int
	}{
		{
			name: "update existing progress",
			progress: domain.TitleProgress{
				UserID:             1,
				TitleRequirementID: 100,
				Progress:           10,
				Completed:          true,
			},
			setup: func() {
				mockRepo.Progress[100] = &domain.TitleProgress{
					UserID:             1,
					TitleRequirementID: 100,
					Progress:           5,
					Completed:          false,
				}
			},
			expectedAmount: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			err := mockRepo.UpdateTitleProgress(&tc.progress)
			assert.NoError(t, err)

			updatedProgress, _ := mockRepo.GetOrCreateTitleProgress(1, 100)
			assert.Equal(t, tc.expectedAmount, updatedProgress.Progress)
			assert.True(t, updatedProgress.Completed)
		})
	}
}

func MockTaskRepository_TestUserHasTitle(t *testing.T) {
	mockRepo := NewMockTaskRepository()

	testCases := []struct {
		name       string
		userID     uint
		titleID    uint
		setup      func()
		expectTrue bool
	}{
		{
			name:    "user has no title",
			userID:  1,
			titleID: 100,
			setup: func() {
				mockRepo.UserTitles = make(map[uint]map[uint]bool)
			},
			expectTrue: false,
		},
		{
			name:    "user has title",
			userID:  1,
			titleID: 100,
			setup: func() {
				mockRepo.UserTitles[1] = map[uint]bool{
					100: true,
				}
			},
			expectTrue: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			hasTitle, err := mockRepo.UserHasTitle(tc.userID, tc.titleID)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectTrue, hasTitle)
		})
	}
}

func MockTaskRepository_TestAwardTitleToUser(t *testing.T) {
	mockRepo := NewMockTaskRepository()

	testCases := []struct {
		name       string
		userID     uint
		titleID    uint
		setup      func()
		expectTrue bool
	}{
		{
			name:    "award title to user",
			userID:  1,
			titleID: 100,
			setup: func() {
				mockRepo.UserTitles = make(map[uint]map[uint]bool)
			},
			expectTrue: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			err := mockRepo.AwardTitleToUser(tc.userID, tc.titleID)
			assert.NoError(t, err)

			hasTitle, err := mockRepo.UserHasTitle(tc.userID, tc.titleID)
			assert.NoError(t, err)
			assert.True(t, hasTitle)
		})
	}
}
