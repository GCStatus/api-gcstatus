package tests

import (
	"errors"
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
	"testing"
	"time"
)

type MockAdminTagRepository struct {
	tags map[uint]*domain.Tag
}

func NewMockAdminTagRepository() *MockAdminTagRepository {
	return &MockAdminTagRepository{
		tags: make(map[uint]*domain.Tag),
	}
}

func (m *MockAdminTagRepository) GetAll() ([]domain.Tag, error) {
	var tags []domain.Tag
	for _, tag := range m.tags {
		tags = append(tags, *tag)
	}
	return tags, nil
}

func (m *MockAdminTagRepository) CreateTag(tag *domain.Tag) error {
	if tag == nil {
		return errors.New("invalid tag data")
	}
	m.tags[tag.ID] = tag
	return nil
}

func (m *MockAdminTagRepository) Update(id uint, request ports_admin.UpdateTagInterface) error {
	if request.Name == "" || request.Slug == "" {
		return errors.New("invalid payload data")
	}
	if _, exists := m.tags[id]; !exists {
		return errors.New("tag not found")
	}
	for _, tag := range m.tags {
		if tag.ID == id {
			tag.Name = request.Name
			tag.Slug = utils.Slugify(request.Name)
		}
	}

	return nil
}

func (m *MockAdminTagRepository) Delete(id uint) error {
	if _, exists := m.tags[id]; !exists {
		return errors.New("tag not found")
	}
	delete(m.tags, id)
	return nil
}

func TestMockAdminTagRepository_GetAll(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		expectedTagsCount int
		mockCreatetags    func(repo *MockAdminTagRepository)
	}{
		"multiple tags": {
			expectedTagsCount: 2,
			mockCreatetags: func(repo *MockAdminTagRepository) {
				if err := repo.CreateTag(&domain.Tag{
					ID:        1,
					Name:      "Tag 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the tag: %s", err.Error())
				}
				if err := repo.CreateTag(&domain.Tag{
					ID:        2,
					Name:      "tag 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the tag: %s", err.Error())
				}
			},
		},
		"no tags": {
			expectedTagsCount: 0,
			mockCreatetags:    func(repo *MockAdminTagRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminTagRepository()

			tc.mockCreatetags(mockRepo)

			tags, err := mockRepo.GetAll()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(tags) != tc.expectedTagsCount {
				t.Fatalf("expected %d tags, got %d", tc.expectedTagsCount, len(tags))
			}
		})
	}
}

func TestMockAdminTagRepository_Update(t *testing.T) {
	testCases := map[string]struct {
		tagID               uint
		updateRequest       ports_admin.UpdateTagInterface
		setupTags           func(repo *MockAdminTagRepository)
		expectedError       error
		expectedUpdatedName string
	}{
		"successful update": {
			tagID: 1,
			updateRequest: ports_admin.UpdateTagInterface{
				Name: "Updated Tag 1",
				Slug: "updated-tag-1",
			},
			setupTags: func(repo *MockAdminTagRepository) {
				if err := repo.CreateTag(&domain.Tag{ID: 1, Name: "Tag 1"}); err != nil {
					t.Fatalf("failed to create tag: %+v", err)
				}
			},
			expectedError:       nil,
			expectedUpdatedName: "Updated Tag 1",
		},
		"invalid payload - empty name": {
			tagID: 1,
			updateRequest: ports_admin.UpdateTagInterface{
				Name: "",
				Slug: "some-slug",
			},
			setupTags: func(repo *MockAdminTagRepository) {
				if err := repo.CreateTag(&domain.Tag{ID: 1, Name: "Tag 1"}); err != nil {
					t.Fatalf("failed to create tag: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Tag 1",
		},
		"invalid payload - empty slug": {
			tagID: 1,
			updateRequest: ports_admin.UpdateTagInterface{
				Name: "Tag 1",
				Slug: "",
			},
			setupTags: func(repo *MockAdminTagRepository) {
				if err := repo.CreateTag(&domain.Tag{ID: 1, Name: "Tag 1"}); err != nil {
					t.Fatalf("failed to create tag: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Tag 1",
		},
		"tag not found": {
			tagID: 99,
			updateRequest: ports_admin.UpdateTagInterface{
				Name: "Nonexistent tag",
				Slug: "nonexistent-tag",
			},
			setupTags:           func(repo *MockAdminTagRepository) {},
			expectedError:       errors.New("tag not found"),
			expectedUpdatedName: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminTagRepository()

			tc.setupTags(mockRepo)

			err := mockRepo.Update(tc.tagID, tc.updateRequest)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tag, exists := mockRepo.tags[tc.tagID]; exists {
				if tag.Name != tc.expectedUpdatedName {
					t.Fatalf("expected tag name to be %s, got %s", tc.expectedUpdatedName, tag.Name)
				}
			} else if tc.expectedUpdatedName != "" {
				t.Fatalf("expected tag %d to exist, but it does not", tc.tagID)
			}
		})
	}
}

func TestMockAdminTagRepository_Delete(t *testing.T) {
	testCases := map[string]struct {
		tagToDelete   uint
		expectedError error
		setupTags     func(repo *MockAdminTagRepository)
	}{
		"successful deletion": {
			tagToDelete:   1,
			expectedError: nil,
			setupTags: func(repo *MockAdminTagRepository) {
				if err := repo.CreateTag(&domain.Tag{ID: 1, Name: "Tag 1"}); err != nil {
					t.Fatalf("failed to create tag: %+v", err)
				}
			},
		},
		"tag does not exist": {
			tagToDelete:   99,
			expectedError: errors.New("tag not found"),
			setupTags:     func(repo *MockAdminTagRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminTagRepository()

			tc.setupTags(mockRepo)

			err := mockRepo.Delete(tc.tagToDelete)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if _, exists := mockRepo.tags[tc.tagToDelete]; exists {
					t.Fatalf("expected tag %d to be deleted, but it still exists", tc.tagToDelete)
				}
			}
		})
	}
}
