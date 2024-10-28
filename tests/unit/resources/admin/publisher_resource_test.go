package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformPublisher(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.Publisher
		expected resources_admin.PublisherResource
	}{
		"as nil": {
			input: domain.Publisher{},
			expected: resources_admin.PublisherResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"single Publisher": {
			input: domain.Publisher{
				ID:        1,
				Name:      "Publisher 1",
				Acting:    false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.PublisherResource{
				ID:        1,
				Name:      "Publisher 1",
				Acting:    false,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformPublisher(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformPublishers(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Publisher
		expected []resources_admin.PublisherResource
	}{
		"empty slice": {
			input:    []domain.Publisher{},
			expected: []resources_admin.PublisherResource{},
		},
		"multiple Publishers": {
			input: []domain.Publisher{
				{
					ID:        1,
					Name:      "Publisher 1",
					Acting:    true,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Publisher 2",
					Acting:    false,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources_admin.PublisherResource{
				{
					ID:        1,
					Name:      "Publisher 1",
					Acting:    true,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Publisher 2",
					Acting:    false,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformPublishers(tc.input)

			if result == nil {
				result = []resources_admin.PublisherResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
