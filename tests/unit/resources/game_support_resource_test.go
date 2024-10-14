package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"testing"
)

func TestTransformSupport(t *testing.T) {
	testCases := map[string]struct {
		input    *domain.GameSupport
		expected resources.SupportResource
	}{
		"as null": {
			input:    &domain.GameSupport{},
			expected: resources.SupportResource{},
		},
		"basic transformation": {
			input: &domain.GameSupport{
				ID:      1,
				URL:     utils.StringPtr("https://google.com"),
				Email:   utils.StringPtr("email@example.com"),
				Contact: utils.StringPtr("fakeContact"),
				GameID:  1,
			},
			expected: resources.SupportResource{
				ID:      1,
				URL:     utils.StringPtr("https://google.com"),
				Email:   utils.StringPtr("email@example.com"),
				Contact: utils.StringPtr("fakeContact"),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformSupport(tc.input)

			if result.ID != tc.expected.ID ||
				!CompareStringPtr(result.URL, tc.expected.URL) ||
				!CompareStringPtr(result.Email, tc.expected.Email) ||
				!CompareStringPtr(result.Contact, tc.expected.Contact) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func CompareStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil {
		return *a == *b
	}
	return false
}
