package tests

import (
	"errors"
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
	"testing"
	"time"
)

type MockAdminCategoryRepository struct {
	categories map[uint]*domain.Category
}

func NewMockAdminCategoryRepository() *MockAdminCategoryRepository {
	return &MockAdminCategoryRepository{
		categories: make(map[uint]*domain.Category),
	}
}

func (m *MockAdminCategoryRepository) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	for _, categorie := range m.categories {
		categories = append(categories, *categorie)
	}
	return categories, nil
}

func (m *MockAdminCategoryRepository) CreateCategory(category *domain.Category) error {
	if category == nil {
		return errors.New("invalid category data")
	}
	m.categories[category.ID] = category
	return nil
}

func (m *MockAdminCategoryRepository) Update(id uint, request ports_admin.UpdateCategoryInterface) error {
	if request.Name == "" || request.Slug == "" {
		return errors.New("invalid payload data")
	}
	if _, exists := m.categories[id]; !exists {
		return errors.New("category not found")
	}
	for _, category := range m.categories {
		if category.ID == id {
			category.Name = request.Name
			category.Slug = utils.Slugify(request.Name)
		}
	}

	return nil
}

func (m *MockAdminCategoryRepository) Delete(id uint) error {
	if _, exists := m.categories[id]; !exists {
		return errors.New("category not found")
	}
	delete(m.categories, id)
	return nil
}

func TestMockAdminCategoryRepository_GetAll(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		expectedCategoriesCount int
		mockCreateCategories    func(repo *MockAdminCategoryRepository)
	}{
		"multiple categories": {
			expectedCategoriesCount: 2,
			mockCreateCategories: func(repo *MockAdminCategoryRepository) {
				if err := repo.CreateCategory(&domain.Category{
					ID:        1,
					Name:      "Category 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the category: %s", err.Error())
				}
				if err := repo.CreateCategory(&domain.Category{
					ID:        2,
					Name:      "Category 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the category: %s", err.Error())
				}
			},
		},
		"no categories": {
			expectedCategoriesCount: 0,
			mockCreateCategories:    func(repo *MockAdminCategoryRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminCategoryRepository()

			tc.mockCreateCategories(mockRepo)

			categories, err := mockRepo.GetAll()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(categories) != tc.expectedCategoriesCount {
				t.Fatalf("expected %d categories, got %d", tc.expectedCategoriesCount, len(categories))
			}
		})
	}
}

func TestMockAdminCategoryRepository_Update(t *testing.T) {
	testCases := map[string]struct {
		categoryID          uint
		updateRequest       ports_admin.UpdateCategoryInterface
		setupCategories     func(repo *MockAdminCategoryRepository)
		expectedError       error
		expectedUpdatedName string
	}{
		"successful update": {
			categoryID: 1,
			updateRequest: ports_admin.UpdateCategoryInterface{
				Name: "Updated Category 1",
				Slug: "updated-category-1",
			},
			setupCategories: func(repo *MockAdminCategoryRepository) {
				if err := repo.CreateCategory(&domain.Category{ID: 1, Name: "Category 1"}); err != nil {
					t.Fatalf("failed to create category: %+v", err)
				}
			},
			expectedError:       nil,
			expectedUpdatedName: "Updated Category 1",
		},
		"invalid payload - empty name": {
			categoryID: 1,
			updateRequest: ports_admin.UpdateCategoryInterface{
				Name: "",
				Slug: "some-slug",
			},
			setupCategories: func(repo *MockAdminCategoryRepository) {
				if err := repo.CreateCategory(&domain.Category{ID: 1, Name: "Category 1"}); err != nil {
					t.Fatalf("failed to create category: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Category 1",
		},
		"invalid payload - empty slug": {
			categoryID: 1,
			updateRequest: ports_admin.UpdateCategoryInterface{
				Name: "Category 1",
				Slug: "",
			},
			setupCategories: func(repo *MockAdminCategoryRepository) {
				if err := repo.CreateCategory(&domain.Category{ID: 1, Name: "Category 1"}); err != nil {
					t.Fatalf("failed to create category: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Category 1",
		},
		"category not found": {
			categoryID: 99,
			updateRequest: ports_admin.UpdateCategoryInterface{
				Name: "Nonexistent Category",
				Slug: "nonexistent-category",
			},
			setupCategories:     func(repo *MockAdminCategoryRepository) {},
			expectedError:       errors.New("category not found"),
			expectedUpdatedName: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminCategoryRepository()

			tc.setupCategories(mockRepo)

			err := mockRepo.Update(tc.categoryID, tc.updateRequest)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if category, exists := mockRepo.categories[tc.categoryID]; exists {
				if category.Name != tc.expectedUpdatedName {
					t.Fatalf("expected category name to be %s, got %s", tc.expectedUpdatedName, category.Name)
				}
			} else if tc.expectedUpdatedName != "" {
				t.Fatalf("expected category %d to exist, but it does not", tc.categoryID)
			}
		})
	}
}

func TestMockAdminCategoryRepository_Delete(t *testing.T) {
	testCases := map[string]struct {
		categoryToDelete uint
		expectedError    error
		setupCategories  func(repo *MockAdminCategoryRepository)
	}{
		"successful deletion": {
			categoryToDelete: 1,
			expectedError:    nil,
			setupCategories: func(repo *MockAdminCategoryRepository) {
				if err := repo.CreateCategory(&domain.Category{ID: 1, Name: "Category 1"}); err != nil {
					t.Fatalf("failed to create category: %+v", err)
				}
			},
		},
		"category does not exist": {
			categoryToDelete: 99,
			expectedError:    errors.New("category not found"),
			setupCategories:  func(repo *MockAdminCategoryRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminCategoryRepository()

			tc.setupCategories(mockRepo)

			err := mockRepo.Delete(tc.categoryToDelete)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if _, exists := mockRepo.categories[tc.categoryToDelete]; exists {
					t.Fatalf("expected category %d to be deleted, but it still exists", tc.categoryToDelete)
				}
			}
		})
	}
}
