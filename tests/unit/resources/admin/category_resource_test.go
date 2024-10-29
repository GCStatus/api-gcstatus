package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformCategory(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.Category
		expected resources_admin.CategoryResource
	}{
		"as null": {
			input: domain.Category{},
			expected: resources_admin.CategoryResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"valid category": {
			input: domain.Category{
				ID:        1,
				Name:      "Category 1",
				Slug:      "category-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.CategoryResource{
				ID:        1,
				Name:      "Category 1",
				Slug:      "category-1",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoryResource := resources_admin.TransformCategory(test.input)

			if !reflect.DeepEqual(categoryResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, categoryResource)
			}
		})
	}
}

func TestTransformCategories(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.Category
		expected []resources_admin.CategoryResource
	}{
		"as null": {
			input:    []domain.Category{},
			expected: []resources_admin.CategoryResource{},
		},
		"multiple categories": {
			input: []domain.Category{
				{
					ID:        1,
					Name:      "Category 1",
					Slug:      "category-1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Category 2",
					Slug:      "category-2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources_admin.CategoryResource{
				{
					ID:        1,
					Name:      "Category 1",
					Slug:      "category-1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Category 2",
					Slug:      "category-2",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoriesResource := resources_admin.TransformCategories(test.input)

			if categoriesResource == nil {
				categoriesResource = []resources_admin.CategoryResource{}
			}

			if !reflect.DeepEqual(categoriesResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, categoriesResource)
			}
		})
	}
}
