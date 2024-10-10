package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformCategory(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    domain.Category
		expected resources.CategoryResource
	}{
		"as null": {
			input:    domain.Category{},
			expected: resources.CategoryResource{},
		},
		"valid category": {
			input: domain.Category{
				ID:        1,
				Name:      "Category 1",
				Slug:      "category-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.CategoryResource{
				ID:   1,
				Name: "Category 1",
				Slug: "category-1",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoryResource := resources.TransformCategory(test.input)

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
		expected []resources.CategoryResource
	}{
		"as null": {
			input:    []domain.Category{},
			expected: []resources.CategoryResource{},
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
			expected: []resources.CategoryResource{
				{
					ID:   1,
					Name: "Category 1",
					Slug: "category-1",
				},
				{
					ID:   2,
					Name: "Category 2",
					Slug: "category-2",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoriesResource := resources.TransformCategories(test.input)

			if categoriesResource == nil {
				categoriesResource = []resources.CategoryResource{}
			}

			if !reflect.DeepEqual(categoriesResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, categoriesResource)
			}
		})
	}
}
