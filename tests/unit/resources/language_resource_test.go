package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformLanguage(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    domain.Language
		expected resources.LanguageResource
	}{
		"as null": {
			input:    domain.Language{},
			expected: resources.LanguageResource{},
		},
		"multiple categories": {
			input: domain.Language{
				ID:        1,
				Name:      "Language 1",
				ISO:       "pt_BR",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.LanguageResource{
				ID:   1,
				Name: "Language 1",
				ISO:  "pt_BR",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			LanguageResource := resources.TransformLanguage(test.input)

			if !reflect.DeepEqual(LanguageResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, LanguageResource)
			}
		})
	}
}

func TestTransformLanguages(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.Language
		expected []resources.LanguageResource
	}{
		"as null": {
			input:    []domain.Language{},
			expected: []resources.LanguageResource{},
		},
		"multiple categories": {
			input: []domain.Language{
				{
					ID:        1,
					Name:      "Language 1",
					ISO:       "pt_BR",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Language 2",
					ISO:       "en_US",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources.LanguageResource{
				{
					ID:   1,
					Name: "Language 1",
					ISO:  "pt_BR",
				},
				{
					ID:   2,
					Name: "Language 2",
					ISO:  "en_US",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			LanguagesResources := resources.TransformLanguages(test.input)

			if LanguagesResources == nil {
				LanguagesResources = []resources.LanguageResource{}
			}

			if !reflect.DeepEqual(LanguagesResources, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, LanguagesResources)
			}
		})
	}
}
