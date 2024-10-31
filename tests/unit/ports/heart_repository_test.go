package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
)

type MockHeartRepository struct {
	heartables map[uint]*domain.Heartable
}

func NewMockHeartRepository() *MockHeartRepository {
	return &MockHeartRepository{
		heartables: make(map[uint]*domain.Heartable),
	}
}

func (m *MockHeartRepository) Create(heartable *domain.Heartable) error {
	if heartable == nil {
		return errors.New("invalid heartable data")
	}
	m.heartables[heartable.ID] = heartable
	return nil
}

func (m *MockHeartRepository) FindForUser(heartableID uint, heartableType string, userID uint) (*domain.Heartable, error) {
	for _, heartable := range m.heartables {
		if heartable.HeartableID == heartableID && heartable.HeartableType == heartableType && heartable.UserID == userID {
			return heartable, nil
		}
	}

	return nil, errors.New("heartable not found")
}

func (m *MockHeartRepository) Delete(id uint) error {
	if _, exists := m.heartables[id]; !exists {
		return errors.New("heartable not found")
	}
	delete(m.heartables, id)
	return nil
}

func TestMockHeartRepository_Create(t *testing.T) {
	mockRepo := NewMockHeartRepository()

	testCases := map[string]struct {
		input         *domain.Heartable
		expectedError bool
	}{
		"valid input": {
			input: &domain.Heartable{
				HeartableID:   1,
				HeartableType: "games",
				UserID:        1,
			},
			expectedError: false,
		},
		"nil input": {
			input:         nil,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.Create(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.heartables[tc.input.ID] == nil {
					t.Fatalf("expected heartable to be created, but it wasn't")
				}
			}
		})
	}
}

func TestMockHeartRepository_FindForUser(t *testing.T) {
	mockRepo := NewMockHeartRepository()

	if err := mockRepo.Create(&domain.Heartable{
		ID:            1,
		HeartableID:   1,
		HeartableType: "games",
		UserID:        1,
	}); err != nil {
		t.Fatalf("failed to create the heartable: %s", err.Error())
	}

	testCases := map[string]struct {
		heartableID   uint
		heartableType string
		userID        uint
		expectedError bool
	}{
		"valid payload": {
			heartableID:   1,
			heartableType: "games",
			userID:        1,
			expectedError: false,
		},
		"invalid token": {
			heartableID:   2,
			heartableType: "games",
			userID:        1,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := mockRepo.FindForUser(tc.heartableID, tc.heartableType, tc.userID)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if result != nil {
					t.Fatalf("expected result to be nil, got %v", result)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == nil || result.HeartableID != tc.heartableID || result.HeartableType != tc.heartableType || result.UserID != tc.userID {
					t.Fatalf("expected heartable %v, got %v", tc, result)
				}
			}
		})
	}
}

func TestMockHeartableRepository_Delete(t *testing.T) {
	mockRepo := NewMockHeartRepository()

	if err := mockRepo.Create(&domain.Heartable{
		ID:            1,
		HeartableID:   1,
		HeartableType: "games",
		UserID:        1,
	}); err != nil {
		t.Fatalf("failed to create the heartable: %s", err.Error())
	}

	testCases := map[string]struct {
		id            uint
		expectedError bool
	}{
		"valid ID": {
			id:            1,
			expectedError: false,
		},
		"invalid ID": {
			id:            999,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.Delete(tc.id)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.heartables[tc.id] != nil {
					t.Fatalf("expected heartable to be deleted, but it wasn't")
				}
			}
		})
	}
}
