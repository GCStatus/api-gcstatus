package tests

import (
	"errors"
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
	"testing"
	"time"
)

type MockAdminPlatformRepository struct {
	platforms map[uint]*domain.Platform
}

func NewMockAdminPlatformRepository() *MockAdminPlatformRepository {
	return &MockAdminPlatformRepository{
		platforms: make(map[uint]*domain.Platform),
	}
}

func (m *MockAdminPlatformRepository) GetAll() ([]domain.Platform, error) {
	var platforms []domain.Platform
	for _, platform := range m.platforms {
		platforms = append(platforms, *platform)
	}
	return platforms, nil
}

func (m *MockAdminPlatformRepository) CreatePlatform(platform *domain.Platform) error {
	if platform == nil {
		return errors.New("invalid platform data")
	}
	m.platforms[platform.ID] = platform
	return nil
}

func (m *MockAdminPlatformRepository) Update(id uint, request ports_admin.UpdatePlatformInterface) error {
	if request.Name == "" || request.Slug == "" {
		return errors.New("invalid payload data")
	}
	if _, exists := m.platforms[id]; !exists {
		return errors.New("platform not found")
	}
	for _, platform := range m.platforms {
		if platform.ID == id {
			platform.Name = request.Name
			platform.Slug = utils.Slugify(request.Name)
		}
	}

	return nil
}

func (m *MockAdminPlatformRepository) Delete(id uint) error {
	if _, exists := m.platforms[id]; !exists {
		return errors.New("platform not found")
	}
	delete(m.platforms, id)
	return nil
}

func TestMockAdminPlatformRepository_GetAll(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		expectedPlatformsCount int
		mockCreatePlatforms    func(repo *MockAdminPlatformRepository)
	}{
		"multiple platforms": {
			expectedPlatformsCount: 2,
			mockCreatePlatforms: func(repo *MockAdminPlatformRepository) {
				if err := repo.CreatePlatform(&domain.Platform{
					ID:        1,
					Name:      "Platform 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the platform: %s", err.Error())
				}
				if err := repo.CreatePlatform(&domain.Platform{
					ID:        2,
					Name:      "platform 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the platform: %s", err.Error())
				}
			},
		},
		"no platforms": {
			expectedPlatformsCount: 0,
			mockCreatePlatforms:    func(repo *MockAdminPlatformRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminPlatformRepository()

			tc.mockCreatePlatforms(mockRepo)

			platforms, err := mockRepo.GetAll()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(platforms) != tc.expectedPlatformsCount {
				t.Fatalf("expected %d platforms, got %d", tc.expectedPlatformsCount, len(platforms))
			}
		})
	}
}

func TestMockAdminPlatformRepository_Update(t *testing.T) {
	testCases := map[string]struct {
		platformID          uint
		updateRequest       ports_admin.UpdatePlatformInterface
		setupPlatforms      func(repo *MockAdminPlatformRepository)
		expectedError       error
		expectedUpdatedName string
	}{
		"successful update": {
			platformID: 1,
			updateRequest: ports_admin.UpdatePlatformInterface{
				Name: "Updated platform 1",
				Slug: "updated-platform-1",
			},
			setupPlatforms: func(repo *MockAdminPlatformRepository) {
				if err := repo.CreatePlatform(&domain.Platform{ID: 1, Name: "Platform 1"}); err != nil {
					t.Fatalf("failed to create platform: %+v", err)
				}
			},
			expectedError:       nil,
			expectedUpdatedName: "Updated platform 1",
		},
		"invalid payload - empty name": {
			platformID: 1,
			updateRequest: ports_admin.UpdatePlatformInterface{
				Name: "",
				Slug: "some-slug",
			},
			setupPlatforms: func(repo *MockAdminPlatformRepository) {
				if err := repo.CreatePlatform(&domain.Platform{ID: 1, Name: "Platform 1"}); err != nil {
					t.Fatalf("failed to create platform: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Platform 1",
		},
		"invalid payload - empty slug": {
			platformID: 1,
			updateRequest: ports_admin.UpdatePlatformInterface{
				Name: "Platform 1",
				Slug: "",
			},
			setupPlatforms: func(repo *MockAdminPlatformRepository) {
				if err := repo.CreatePlatform(&domain.Platform{ID: 1, Name: "Platform 1"}); err != nil {
					t.Fatalf("failed to create platform: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Platform 1",
		},
		"platform not found": {
			platformID: 99,
			updateRequest: ports_admin.UpdatePlatformInterface{
				Name: "Nonexistent platform",
				Slug: "nonexistent-platform",
			},
			setupPlatforms:      func(repo *MockAdminPlatformRepository) {},
			expectedError:       errors.New("platform not found"),
			expectedUpdatedName: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminPlatformRepository()

			tc.setupPlatforms(mockRepo)

			err := mockRepo.Update(tc.platformID, tc.updateRequest)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if platform, exists := mockRepo.platforms[tc.platformID]; exists {
				if platform.Name != tc.expectedUpdatedName {
					t.Fatalf("expected platform name to be %s, got %s", tc.expectedUpdatedName, platform.Name)
				}
			} else if tc.expectedUpdatedName != "" {
				t.Fatalf("expected platform %d to exist, but it does not", tc.platformID)
			}
		})
	}
}

func TestMockAdminPlatformRepository_Delete(t *testing.T) {
	testCases := map[string]struct {
		platformToDelete uint
		expectedError    error
		setupPlatforms   func(repo *MockAdminPlatformRepository)
	}{
		"successful deletion": {
			platformToDelete: 1,
			expectedError:    nil,
			setupPlatforms: func(repo *MockAdminPlatformRepository) {
				if err := repo.CreatePlatform(&domain.Platform{ID: 1, Name: "Platform 1"}); err != nil {
					t.Fatalf("failed to create platform: %+v", err)
				}
			},
		},
		"platform does not exist": {
			platformToDelete: 99,
			expectedError:    errors.New("platform not found"),
			setupPlatforms:   func(repo *MockAdminPlatformRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminPlatformRepository()

			tc.setupPlatforms(mockRepo)

			err := mockRepo.Delete(tc.platformToDelete)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if _, exists := mockRepo.platforms[tc.platformToDelete]; exists {
					t.Fatalf("expected platform %d to be deleted, but it still exists", tc.platformToDelete)
				}
			}
		})
	}
}
