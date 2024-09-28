package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"
)

type MockTitleRepository struct {
	titles     map[uint]*domain.Title
	userTitles map[uint][]uint
}

func NewMockTitleRepository() *MockTitleRepository {
	return &MockTitleRepository{
		titles:     make(map[uint]*domain.Title),
		userTitles: make(map[uint][]uint),
	}
}

func (m *MockTitleRepository) GetAll(userID uint) ([]domain.Title, error) {
	var titles []domain.Title
	titleIDs := m.userTitles[userID]

	for _, titleID := range titleIDs {
		if title, exists := m.titles[titleID]; exists {
			titles = append(titles, *title)
		}
	}

	return titles, nil
}

func (m *MockTitleRepository) FindById(id uint) (*domain.Title, error) {
	for _, title := range m.titles {
		if title.ID == id {
			return title, nil
		}
	}
	return nil, errors.New("title not found")
}

func (m *MockTitleRepository) AddUserTitle(userID, titleID uint) {
	m.userTitles[userID] = append(m.userTitles[userID], titleID)
}

func (m *MockTitleRepository) CreateTitle(title *domain.Title) error {
	if title == nil {
		return errors.New("invalid title data")
	}
	m.titles[title.ID] = title
	return nil
}

func (m *MockTitleRepository) ToggleEnableTitle(userID uint, titleID uint) error {
	for _, title := range m.titles {
		if title.ID == titleID {
			for i, userTitle := range title.Users {
				if userTitle.UserID == userID {
					userTitle.Enabled = !userTitle.Enabled
					title.Users[i] = userTitle
					return nil
				}
			}
			return errors.New("user not found for given title")
		}
	}
	return errors.New("title not found")
}

func (m *MockTitleRepository) MockTitleRepository_GetAll(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		userID             uint
		expectedTitleCount int
		mockCreateTitles   func(repo *MockTitleRepository)
	}{
		"multiple titles for user 1": {
			userID:             1,
			expectedTitleCount: 2,
			mockCreateTitles: func(repo *MockTitleRepository) {
				err := repo.CreateTitle(&domain.Title{
					ID:          1,
					Title:       "Title 1",
					Description: "Title 1",
					Cost:        func() *int { i := 200; return &i }(),
					Purchasable: true,
					Status:      "available",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the title: %s", err.Error())
				}
				err = repo.CreateTitle(&domain.Title{
					ID:          2,
					Title:       "Title 2",
					Description: "Title 2",
					Cost:        nil,
					Purchasable: false,
					Status:      "available",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the title: %s", err.Error())
				}

				repo.AddUserTitle(1, 1)
				repo.AddUserTitle(1, 2)
			},
		},
		"no titles for user 1": {
			userID:             1,
			expectedTitleCount: 0,
			mockCreateTitles:   func(repo *MockTitleRepository) {},
		},
		"titles for user 2": {
			userID:             2,
			expectedTitleCount: 1,
			mockCreateTitles: func(repo *MockTitleRepository) {
				err := repo.CreateTitle(&domain.Title{
					ID:          3,
					Title:       "Title 3",
					Description: "Title 3",
					Cost:        func() *int { i := 100; return &i }(),
					Purchasable: true,
					Status:      "available",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the title: %s", err.Error())
				}

				repo.AddUserTitle(2, 3)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockTitleRepository()

			tc.mockCreateTitles(mockRepo)

			titles, err := mockRepo.GetAll(tc.userID)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(titles) != tc.expectedTitleCount {
				t.Fatalf("expected %d titles, got %d", tc.expectedTitleCount, len(titles))
			}
		})
	}
}

func TestMockTitleRepository_FindById(t *testing.T) {
	fixedTime := time.Now()

	mockRepo := NewMockTitleRepository()
	err := mockRepo.CreateTitle(&domain.Title{
		ID:          1,
		Title:       "Title 1",
		Description: "Title 1",
		Cost:        func() *int { i := 200; return &i }(),
		Purchasable: true,
		Status:      "available",
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	})

	if err != nil {
		t.Fatalf("failed to create the title: %s", err.Error())
	}

	testCases := map[string]struct {
		titleID     uint
		expectError bool
	}{
		"valid title ID": {
			titleID:     1,
			expectError: false,
		},
		"invalid title ID": {
			titleID:     999,
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			title, err := mockRepo.FindById(tc.titleID)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if title != nil {
					t.Fatalf("expected nil title, got %v", title)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if title == nil || title.ID != tc.titleID {
					t.Fatalf("expected title ID %d, got %v", tc.titleID, title)
				}
			}
		})
	}
}

func TestMockTitleRepository_ToggleEnableTitle(t *testing.T) {
	tests := map[string]struct {
		userID           uint
		titleID          uint
		mockCreateTitles func(repo *MockTitleRepository)
		expectedEnabled  bool
		expectError      bool
	}{
		"enable title for user": {
			userID:  1,
			titleID: 1,
			mockCreateTitles: func(repo *MockTitleRepository) {
				repo.titles[1] = &domain.Title{
					ID: 1,
					Users: []domain.UserTitle{
						{UserID: 1, Enabled: false},
					},
				}
			},
			expectedEnabled: true,
			expectError:     false,
		},
		"user not found": {
			userID:  2,
			titleID: 1,
			mockCreateTitles: func(repo *MockTitleRepository) {
				repo.titles[1] = &domain.Title{
					ID: 1,
					Users: []domain.UserTitle{
						{UserID: 1, Enabled: false},
					},
				}
			},
			expectedEnabled: false,
			expectError:     true,
		},
		"title not found": {
			userID:  1,
			titleID: 2,
			mockCreateTitles: func(repo *MockTitleRepository) {
				repo.titles[1] = &domain.Title{
					ID: 1,
					Users: []domain.UserTitle{
						{UserID: 1, Enabled: false},
					},
				}
			},
			expectedEnabled: false,
			expectError:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockRepo := NewMockTitleRepository()
			tc.mockCreateTitles(mockRepo)

			err := mockRepo.ToggleEnableTitle(tc.userID, tc.titleID)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected an error, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			for _, title := range mockRepo.titles {
				if title.ID == tc.titleID {
					for _, userTitle := range title.Users {
						if userTitle.UserID == tc.userID {
							if userTitle.Enabled != tc.expectedEnabled {
								t.Errorf("expected enabled to be %v, got %v", tc.expectedEnabled, userTitle.Enabled)
							}
						}
					}
				}
			}
		})
	}
}
