package tests

import (
	"errors"
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
	"testing"
	"time"
)

type MockAdminGenreRepository struct {
	genres map[uint]*domain.Genre
}

func NewMockAdminGenreRepository() *MockAdminGenreRepository {
	return &MockAdminGenreRepository{
		genres: make(map[uint]*domain.Genre),
	}
}

func (m *MockAdminGenreRepository) GetAll() ([]domain.Genre, error) {
	var genres []domain.Genre
	for _, genre := range m.genres {
		genres = append(genres, *genre)
	}
	return genres, nil
}

func (m *MockAdminGenreRepository) CreateGenre(genre *domain.Genre) error {
	if genre == nil {
		return errors.New("invalid genre data")
	}
	m.genres[genre.ID] = genre
	return nil
}

func (m *MockAdminGenreRepository) Update(id uint, request ports_admin.UpdateGenreInterface) error {
	if request.Name == "" || request.Slug == "" {
		return errors.New("invalid payload data")
	}
	if _, exists := m.genres[id]; !exists {
		return errors.New("genre not found")
	}
	for _, genre := range m.genres {
		if genre.ID == id {
			genre.Name = request.Name
			genre.Slug = utils.Slugify(request.Name)
		}
	}

	return nil
}

func (m *MockAdminGenreRepository) Delete(id uint) error {
	if _, exists := m.genres[id]; !exists {
		return errors.New("genre not found")
	}
	delete(m.genres, id)
	return nil
}

func TestMockAdminGenreRepository_GetAll(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		expectedgenresCount int
		mockCreateGenres    func(repo *MockAdminGenreRepository)
	}{
		"multiple genres": {
			expectedgenresCount: 2,
			mockCreateGenres: func(repo *MockAdminGenreRepository) {
				if err := repo.CreateGenre(&domain.Genre{
					ID:        1,
					Name:      "Genre 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the Genre: %s", err.Error())
				}
				if err := repo.CreateGenre(&domain.Genre{
					ID:        2,
					Name:      "Genre 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the Genre: %s", err.Error())
				}
			},
		},
		"no genres": {
			expectedgenresCount: 0,
			mockCreateGenres:    func(repo *MockAdminGenreRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminGenreRepository()

			tc.mockCreateGenres(mockRepo)

			genres, err := mockRepo.GetAll()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(genres) != tc.expectedgenresCount {
				t.Fatalf("expected %d genres, got %d", tc.expectedgenresCount, len(genres))
			}
		})
	}
}

func TestMockAdminGenreRepository_Update(t *testing.T) {
	testCases := map[string]struct {
		GenreID             uint
		updateRequest       ports_admin.UpdateGenreInterface
		setupGenres         func(repo *MockAdminGenreRepository)
		expectedError       error
		expectedUpdatedName string
	}{
		"successful update": {
			GenreID: 1,
			updateRequest: ports_admin.UpdateGenreInterface{
				Name: "Updated genre 1",
				Slug: "updated-genre-1",
			},
			setupGenres: func(repo *MockAdminGenreRepository) {
				if err := repo.CreateGenre(&domain.Genre{ID: 1, Name: "Genre 1"}); err != nil {
					t.Fatalf("failed to create genre: %+v", err)
				}
			},
			expectedError:       nil,
			expectedUpdatedName: "Updated genre 1",
		},
		"invalid payload - empty name": {
			GenreID: 1,
			updateRequest: ports_admin.UpdateGenreInterface{
				Name: "",
				Slug: "some-slug",
			},
			setupGenres: func(repo *MockAdminGenreRepository) {
				if err := repo.CreateGenre(&domain.Genre{ID: 1, Name: "Genre 1"}); err != nil {
					t.Fatalf("failed to create genre: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Genre 1",
		},
		"invalid payload - empty slug": {
			GenreID: 1,
			updateRequest: ports_admin.UpdateGenreInterface{
				Name: "Genre 1",
				Slug: "",
			},
			setupGenres: func(repo *MockAdminGenreRepository) {
				if err := repo.CreateGenre(&domain.Genre{ID: 1, Name: "Genre 1"}); err != nil {
					t.Fatalf("failed to create genre: %+v", err)
				}
			},
			expectedError:       errors.New("invalid payload data"),
			expectedUpdatedName: "Genre 1",
		},
		"Genre not found": {
			GenreID: 99,
			updateRequest: ports_admin.UpdateGenreInterface{
				Name: "Nonexistent Genre",
				Slug: "nonexistent-Genre",
			},
			setupGenres:         func(repo *MockAdminGenreRepository) {},
			expectedError:       errors.New("genre not found"),
			expectedUpdatedName: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminGenreRepository()

			tc.setupGenres(mockRepo)

			err := mockRepo.Update(tc.GenreID, tc.updateRequest)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if Genre, exists := mockRepo.genres[tc.GenreID]; exists {
				if Genre.Name != tc.expectedUpdatedName {
					t.Fatalf("expected genre name to be %s, got %s", tc.expectedUpdatedName, Genre.Name)
				}
			} else if tc.expectedUpdatedName != "" {
				t.Fatalf("expected genre %d to exist, but it does not", tc.GenreID)
			}
		})
	}
}

func TestMockAdminGenreRepository_Delete(t *testing.T) {
	testCases := map[string]struct {
		GenreToDelete uint
		expectedError error
		setupGenres   func(repo *MockAdminGenreRepository)
	}{
		"successful deletion": {
			GenreToDelete: 1,
			expectedError: nil,
			setupGenres: func(repo *MockAdminGenreRepository) {
				if err := repo.CreateGenre(&domain.Genre{ID: 1, Name: "Genre 1"}); err != nil {
					t.Fatalf("failed to create genre: %+v", err)
				}
			},
		},
		"Genre does not exist": {
			GenreToDelete: 99,
			expectedError: errors.New("genre not found"),
			setupGenres:   func(repo *MockAdminGenreRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockAdminGenreRepository()

			tc.setupGenres(mockRepo)

			err := mockRepo.Delete(tc.GenreToDelete)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if _, exists := mockRepo.genres[tc.GenreToDelete]; exists {
					t.Fatalf("expected genre %d to be deleted, but it still exists", tc.GenreToDelete)
				}
			}
		})
	}
}
