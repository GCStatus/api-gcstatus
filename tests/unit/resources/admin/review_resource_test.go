package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformReview(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.Reviewable
		expected resources_admin.ReviewResource
	}{
		"as nil": {
			input: domain.Reviewable{},
			expected: resources_admin.ReviewResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"single Review": {
			input: domain.Reviewable{
				ID:             1,
				Rate:           5,
				Review:         "Good game!",
				Played:         true,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
				ReviewableID:   1,
				ReviewableType: "games",
				User: domain.User{
					ID:        1,
					Email:     "fake@gmail.com",
					Name:      "Fake",
					Nickname:  "fake",
					Birthdate: fixedTime,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Profile: domain.Profile{
						ID:    1,
						Share: true,
						Photo: "profile-photo-key",
					},
					Roles:       []domain.Roleable{},
					Permissions: []domain.Permissionable{},
				},
			},
			expected: resources_admin.ReviewResource{
				ID:        1,
				Rate:      5,
				Review:    "Good game!",
				Played:    true,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				User: resources_admin.MinimalUserResource{
					ID:        1,
					Email:     "fake@gmail.com",
					Name:      "Fake",
					Nickname:  "fake",
					CreatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformReview(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformReviews(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Reviewable
		expected []resources_admin.ReviewResource
	}{
		"multiple Reviews": {
			input: []domain.Reviewable{
				{
					ID:             1,
					Rate:           1,
					Review:         "Bad game!",
					Played:         true,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
					ReviewableID:   1,
					ReviewableType: "games",
					User: domain.User{
						ID:         1,
						Email:      "fake@gmail.com",
						Name:       "Fake",
						Nickname:   "fake",
						Experience: 0,
						Birthdate:  fixedTime,
						Password:   "fake1234",
						Blocked:    false,
						LevelID:    1,
						CreatedAt:  fixedTime,
						UpdatedAt:  fixedTime,
						Profile: domain.Profile{
							ID:    1,
							Share: true,
							Photo: "photo-1-key",
						},
						Roles:       []domain.Roleable{},
						Permissions: []domain.Permissionable{},
					},
				},
				{
					ID:             2,
					Rate:           5,
					Review:         "Good game!",
					Played:         false,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
					ReviewableID:   1,
					ReviewableType: "games",
					User: domain.User{
						ID:         1,
						Email:      "fake@gmail.com",
						Name:       "Fake",
						Nickname:   "fake",
						Experience: 0,
						Birthdate:  fixedTime,
						Password:   "fake1234",
						Blocked:    false,
						LevelID:    1,
						CreatedAt:  fixedTime,
						UpdatedAt:  fixedTime,
						Profile: domain.Profile{
							ID:    1,
							Share: true,
							Photo: "photo-2-key",
						},
						Roles:       []domain.Roleable{},
						Permissions: []domain.Permissionable{},
					},
				},
			},
			expected: []resources_admin.ReviewResource{
				{
					ID:        1,
					Rate:      1,
					Review:    "Bad game!",
					Played:    true,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
					User: resources_admin.MinimalUserResource{
						ID:        1,
						Email:     "fake@gmail.com",
						Name:      "Fake",
						Nickname:  "fake",
						CreatedAt: utils.FormatTimestamp(fixedTime),
					},
				},
				{
					ID:        2,
					Rate:      5,
					Review:    "Good game!",
					Played:    false,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
					User: resources_admin.MinimalUserResource{
						ID:        1,
						Email:     "fake@gmail.com",
						Name:      "Fake",
						Nickname:  "fake",
						CreatedAt: utils.FormatTimestamp(fixedTime),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformReviews(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
