package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformLanguage(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.Language
		expected resources_admin.LanguageResource
	}{
		"as null": {
			input: domain.Language{},
			expected: resources_admin.LanguageResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"multiple categories": {
			input: domain.Language{
				ID:        1,
				Name:      "Language 1",
				ISO:       "pt_BR",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.LanguageResource{
				ID:        1,
				Name:      "Language 1",
				ISO:       "pt_BR",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			LanguageResource := resources_admin.TransformLanguage(test.input)

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
		expected []resources_admin.LanguageResource
	}{
		"as null": {
			input:    []domain.Language{},
			expected: []resources_admin.LanguageResource{},
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
			expected: []resources_admin.LanguageResource{
				{
					ID:        1,
					Name:      "Language 1",
					ISO:       "pt_BR",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Language 2",
					ISO:       "en_US",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			LanguagesResources_admin := resources_admin.TransformLanguages(test.input)

			if LanguagesResources_admin == nil {
				LanguagesResources_admin = []resources_admin.LanguageResource{}
			}

			if !reflect.DeepEqual(LanguagesResources_admin, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, LanguagesResources_admin)
			}
		})
	}
}
