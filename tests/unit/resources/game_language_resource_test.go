package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformGameLanguage(t *testing.T) {
	testCases := map[string]struct {
		input    domain.GameLanguage
		expected resources.GameLanguageResource
	}{
		"basic transformation": {
			input: domain.GameLanguage{
				ID:        1,
				Menu:      true,
				Dubs:      false,
				Subtitles: true,
				Language: domain.Language{
					ID:   1,
					Name: "English",
					ISO:  "en",
				},
			},
			expected: resources.GameLanguageResource{
				ID:        1,
				Menu:      true,
				Dubs:      false,
				Subtitles: true,
				Language: resources.LanguageResource{
					ID:   1,
					Name: "English",
					ISO:  "en",
				},
			},
		},
		"language with nil values": {
			input: domain.GameLanguage{
				ID:        2,
				Menu:      false,
				Dubs:      false,
				Subtitles: false,
				Language: domain.Language{
					ID:   0,
					Name: "",
					ISO:  "",
				},
			},
			expected: resources.GameLanguageResource{
				ID:        2,
				Menu:      false,
				Dubs:      false,
				Subtitles: false,
				Language: resources.LanguageResource{
					ID:   0,
					Name: "",
					ISO:  "",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformGameLanguage(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
