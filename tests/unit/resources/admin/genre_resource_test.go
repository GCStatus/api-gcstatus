package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformGenre(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.Genre
		expected resources_admin.GenreResource
	}{
		"as null": {
			input: domain.Genre{},
			expected: resources_admin.GenreResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"multiple categories": {
			input: domain.Genre{
				ID:        1,
				Name:      "Genre 1",
				Slug:      "genre-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.GenreResource{
				ID:        1,
				Name:      "Genre 1",
				Slug:      "genre-1",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			GenreResource := resources_admin.TransformGenre(test.input)

			if !reflect.DeepEqual(GenreResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, GenreResource)
			}
		})
	}
}

func TestTransformGenres(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.Genre
		expected []resources_admin.GenreResource
	}{
		"as null": {
			input:    []domain.Genre{},
			expected: []resources_admin.GenreResource{},
		},
		"multiple categories": {
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
			expected: []resources_admin.GenreResource{
				{
					ID:        1,
					Name:      "Genre 1",
					Slug:      "genre-1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Genre 2",
					Slug:      "genre-2",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			genresResources_admin := resources_admin.TransformGenres(test.input)

			if genresResources_admin == nil {
				genresResources_admin = []resources_admin.GenreResource{}
			}

			if !reflect.DeepEqual(genresResources_admin, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, genresResources_admin)
			}
		})
	}
}
