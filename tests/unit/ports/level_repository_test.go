package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"
)

type MockLevelRepository struct {
	levels map[uint]*domain.Level
}

func NewMockLevelRepository() *MockLevelRepository {
	return &MockLevelRepository{
		levels: make(map[uint]*domain.Level),
	}
}

func (m *MockLevelRepository) GetAll() ([]domain.Level, error) {
	var levels []domain.Level
	for _, level := range m.levels {
		levels = append(levels, *level)
	}

	return levels, nil
}

func (m *MockLevelRepository) CreateLevel(level *domain.Level) error {
	if level == nil {
		return errors.New("invalid level data")
	}
	m.levels[level.ID] = level
	return nil
}

func (m *MockLevelRepository) FindById(id uint) (*domain.Level, error) {
	for _, level := range m.levels {
		if level.ID == id {
			return level, nil
		}
	}
	return nil, errors.New("level not found")
}

func (m *MockLevelRepository) FindByLevel(level uint) (*domain.Level, error) {
	for _, lvl := range m.levels {
		if lvl.Level == level {
			return lvl, nil
		}
	}
	return nil, errors.New("level not found")
}

func TestMockLevelRepository_GetAll(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		expectedLevelCount int
		mockCreateLevels   func(repo *MockLevelRepository)
	}{
		"multiple levels": {
			expectedLevelCount: 2,
			mockCreateLevels: func(repo *MockLevelRepository) {
				err := repo.CreateLevel(&domain.Level{
					ID:         1,
					Level:      1,
					Experience: 500,
					Coins:      100,
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the level: %s", err.Error())
				}
				err = repo.CreateLevel(&domain.Level{
					ID:         2,
					Level:      2,
					Experience: 1000,
					Coins:      150,
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the level: %s", err.Error())
				}
			},
		},
		"no levels": {
			expectedLevelCount: 0,
			mockCreateLevels:   func(repo *MockLevelRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockLevelRepository()

			tc.mockCreateLevels(mockRepo)

			levels, err := mockRepo.GetAll()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(levels) != tc.expectedLevelCount {
				t.Fatalf("expected %d levels, got %d", tc.expectedLevelCount, len(levels))
			}
		})
	}
}

func TestMockLevelRepository_FindById(t *testing.T) {
	fixedTime := time.Now()

	mockRepo := NewMockLevelRepository()
	err := mockRepo.CreateLevel(&domain.Level{
		ID:         1,
		Level:      1,
		Experience: 500,
		Coins:      100,
		CreatedAt:  fixedTime,
		UpdatedAt:  fixedTime,
	})

	if err != nil {
		t.Fatalf("failed to create the level: %s", err.Error())
	}

	testCases := map[string]struct {
		levelID     uint
		expectError bool
	}{
		"valid level ID": {
			levelID:     1,
			expectError: false,
		},
		"invalid level ID": {
			levelID:     999,
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			level, err := mockRepo.FindById(tc.levelID)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if level != nil {
					t.Fatalf("expected nil level, got %v", level)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if level == nil || level.ID != tc.levelID {
					t.Fatalf("expected level ID %d, got %v", tc.levelID, level)
				}
			}
		})
	}
}

func TestMockLevelRepository_FindByLevel(t *testing.T) {
	fixedTime := time.Now()

	mockRepo := NewMockLevelRepository()
	err := mockRepo.CreateLevel(&domain.Level{
		ID:         1,
		Level:      1,
		Experience: 500,
		Coins:      100,
		CreatedAt:  fixedTime,
		UpdatedAt:  fixedTime,
	})

	if err != nil {
		t.Fatalf("failed to create the level: %s", err.Error())
	}

	testCases := map[string]struct {
		levelLevel  uint
		expectError bool
	}{
		"valid level ID": {
			levelLevel:  1,
			expectError: false,
		},
		"invalid level ID": {
			levelLevel:  999,
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			level, err := mockRepo.FindByLevel(tc.levelLevel)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if level != nil {
					t.Fatalf("expected nil level, got %v", level)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if level == nil || level.ID != tc.levelLevel {
					t.Fatalf("expected level ID %d, got %v", tc.levelLevel, level)
				}
			}
		})
	}
}
