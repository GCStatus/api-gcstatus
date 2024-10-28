package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformSupport(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.GameSupport
		expected *resources_admin.SupportResource
	}{
		"as nil": {
			input: domain.GameSupport{},
			expected: &resources_admin.SupportResource{
				URL:       nil,
				Email:     nil,
				Contact:   nil,
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"basic transformation": {
			input: domain.GameSupport{
				ID:        1,
				URL:       utils.StringPtr("fakeurl.com.br"),
				Email:     utils.StringPtr("fakemail@gmail.com"),
				Contact:   utils.StringPtr("+23 12312394"),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Game: domain.Game{
					Slug:             "valid",
					Age:              18,
					Title:            "Game Test",
					Condition:        domain.CommomCondition,
					Cover:            "https://placehold.co/600x400/EEE/31343C",
					About:            "About game",
					Description:      "Description",
					ShortDescription: "Short description",
					Free:             false,
					ReleaseDate:      fixedTime,
					CreatedAt:        fixedTime,
					UpdatedAt:        fixedTime,
					Views: []domain.Viewable{
						{
							UserID:       10,
							ViewableID:   1,
							ViewableType: "games",
						},
					},
				},
			},
			expected: &resources_admin.SupportResource{
				ID:        1,
				URL:       utils.StringPtr("fakeurl.com.br"),
				Email:     utils.StringPtr("fakemail@gmail.com"),
				Contact:   utils.StringPtr("+23 12312394"),
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformSupport(&tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
