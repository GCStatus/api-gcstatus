package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformTag(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.Tag
		expected resources_admin.TagResource
	}{
		"as null": {
			input: domain.Tag{},
			expected: resources_admin.TagResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"multiple categories": {
			input: domain.Tag{
				ID:        1,
				Name:      "Tag 1",
				Slug:      "tag-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.TagResource{
				ID:        1,
				Name:      "Tag 1",
				Slug:      "tag-1",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoryResource := resources_admin.TransformTag(test.input)

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
		expected []resources_admin.TagResource
	}{
		"as null": {
			input:    []domain.Tag{},
			expected: []resources_admin.TagResource{},
		},
		"multiple categories": {
			input: []domain.Tag{
				{
					ID:        1,
					Name:      "Tag 1",
					Slug:      "tag-1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Tag 2",
					Slug:      "tag-2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources_admin.TagResource{
				{
					ID:        1,
					Name:      "Tag 1",
					Slug:      "tag-1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Tag 2",
					Slug:      "tag-2",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			genresResource := resources_admin.TransformTags(test.input)

			if genresResource == nil {
				genresResource = []resources_admin.TagResource{}
			}

			if !reflect.DeepEqual(genresResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, genresResource)
			}
		})
	}
}
