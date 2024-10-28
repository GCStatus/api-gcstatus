package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformGameLanguage(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.GameLanguage
		expected resources_admin.GameLanguageResource
	}{
		"basic transformation": {
			input: domain.GameLanguage{
				ID:        1,
				Menu:      true,
				Dubs:      false,
				Subtitles: true,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Language: domain.Language{
					ID:        1,
					Name:      "English",
					ISO:       "en",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources_admin.GameLanguageResource{
				ID:        1,
				Menu:      true,
				Dubs:      false,
				Subtitles: true,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				Language: resources_admin.LanguageResource{
					ID:        1,
					Name:      "English",
					ISO:       "en",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"language with nil values": {
			input: domain.GameLanguage{
				ID:        2,
				Menu:      false,
				Dubs:      false,
				Subtitles: false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Language: domain.Language{
					ID:   0,
					Name: "",
					ISO:  "",
				},
			},
			expected: resources_admin.GameLanguageResource{
				ID:        2,
				Menu:      false,
				Dubs:      false,
				Subtitles: false,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				Language: resources_admin.LanguageResource{
					ID:        0,
					Name:      "",
					ISO:       "",
					CreatedAt: utils.FormatTimestamp(zeroTime),
					UpdatedAt: utils.FormatTimestamp(zeroTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformGameLanguage(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
