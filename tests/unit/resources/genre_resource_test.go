package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformGenre(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    domain.Genre
		expected resources.GenreResource
	}{
		"as null": {
			input:    domain.Genre{},
			expected: resources.GenreResource{},
		},
		"valid category": {
			input: domain.Genre{
				ID:        1,
				Name:      "Genre 1",
				Slug:      "genre-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.GenreResource{
				ID:   1,
				Name: "Genre 1",
				Slug: "genre-1",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			genreResource := resources.TransformGenre(test.input)

			if !reflect.DeepEqual(genreResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, genreResource)
			}
		})
	}
}

func TestTransformGenres(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.Genre
		expected []resources.GenreResource
	}{
		"as null": {
			input:    []domain.Genre{},
			expected: []resources.GenreResource{},
		},
		"multiple genres": {
			input: []domain.Genre{
				{
					ID:        1,
					Name:      "Genre 1",
					Slug:      "genre-1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Genre 2",
					Slug:      "genre-2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources.GenreResource{
				{
					ID:   1,
					Name: "Genre 1",
					Slug: "genre-1",
				},
				{
					ID:   2,
					Name: "Genre 2",
					Slug: "genre-2",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			genresResource := resources.TransformGenres(test.input)

			if genresResource == nil {
				genresResource = []resources.GenreResource{}
			}

			if !reflect.DeepEqual(genresResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, genresResource)
			}
		})
	}
}
