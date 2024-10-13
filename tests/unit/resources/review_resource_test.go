package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformReview(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Reviewable
		expected resources.ReviewResource
	}{
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
					CreatedAt: fixedTime,
					Profile: domain.Profile{
						ID:    1,
						Share: true,
						Photo: "profile-photo-key",
					},
				},
			},
			expected: resources.ReviewResource{
				ID:        1,
				Rate:      5,
				Review:    "Good game!",
				Played:    true,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				User: resources.MinimalUserResource{
					ID:        1,
					Email:     "fake@gmail.com",
					Name:      "Fake",
					Nickname:  "fake",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/profile-photo-key"),
					CreatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockS3Client := &MockS3Client{}
			result := resources.TransformReview(tc.input, mockS3Client)

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
		expected []resources.ReviewResource
	}{
		"empty slice": {
			input:    []domain.Reviewable{},
			expected: []resources.ReviewResource{},
		},
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
						Profile: domain.Profile{
							ID:    1,
							Share: true,
							Photo: "photo-1-key",
						},
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
						Profile: domain.Profile{
							ID:    1,
							Share: true,
							Photo: "photo-2-key",
						},
					},
				},
			},
			expected: []resources.ReviewResource{
				{
					ID:        1,
					Rate:      1,
					Review:    "Bad game!",
					Played:    true,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
					User: resources.MinimalUserResource{
						ID:        1,
						Email:     "fake@gmail.com",
						Name:      "Fake",
						Nickname:  "fake",
						Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-1-key"),
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
					User: resources.MinimalUserResource{
						ID:        1,
						Email:     "fake@gmail.com",
						Name:      "Fake",
						Nickname:  "fake",
						Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-2-key"),
						CreatedAt: utils.FormatTimestamp(fixedTime),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockS3Client := &MockS3Client{}
			result := resources.TransformReviews(tc.input, mockS3Client)

			if result == nil {
				result = []resources.ReviewResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
