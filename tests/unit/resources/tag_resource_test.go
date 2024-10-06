package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformTag(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    domain.Tag
		expected resources.TagResource
	}{
		"as null": {
			input:    domain.Tag{},
			expected: resources.TagResource{},
		},
		"multiple categories": {
			input: domain.Tag{
				ID:        1,
				Name:      "Tag 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.TagResource{
				ID:   1,
				Name: "Tag 1",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoryResource := resources.TransformTag(test.input)

			if !reflect.DeepEqual(categoryResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, categoryResource)
			}
		})
	}
}

func TestTransformTags(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.Tag
		expected []resources.TagResource
	}{
		"as null": {
			input:    []domain.Tag{},
			expected: []resources.TagResource{},
		},
		"multiple categories": {
			input: []domain.Tag{
				{
					ID:        1,
					Name:      "Tag 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Tag 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources.TagResource{
				{
					ID:   1,
					Name: "Tag 1",
				},
				{
					ID:   2,
					Name: "Tag 2",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			genresResource := resources.TransformTags(test.input)

			if genresResource == nil {
				genresResource = []resources.TagResource{}
			}

			if !reflect.DeepEqual(genresResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, genresResource)
			}
		})
	}
}
