package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformPublisher(t *testing.T) {
	testCases := map[string]struct {
		input    domain.Publisher
		expected resources.PublisherResource
	}{
		"single Publisher": {
			input: domain.Publisher{
				ID:     1,
				Name:   "Publisher 1",
				Acting: false,
			},
			expected: resources.PublisherResource{
				ID:     1,
				Name:   "Publisher 1",
				Acting: false,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformPublisher(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformPublishers(t *testing.T) {
	testCases := map[string]struct {
		input    []domain.Publisher
		expected []resources.PublisherResource
	}{
		"empty slice": {
			input:    []domain.Publisher{},
			expected: []resources.PublisherResource{},
		},
		"multiple Publishers": {
			input: []domain.Publisher{
				{
					ID:     1,
					Name:   "Publisher 1",
					Acting: true,
				},
				{
					ID:     2,
					Name:   "Publisher 2",
					Acting: false,
				},
			},
			expected: []resources.PublisherResource{
				{
					ID:     1,
					Name:   "Publisher 1",
					Acting: true,
				},
				{
					ID:     2,
					Name:   "Publisher 2",
					Acting: false,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformPublishers(tc.input)

			if result == nil {
				result = []resources.PublisherResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
