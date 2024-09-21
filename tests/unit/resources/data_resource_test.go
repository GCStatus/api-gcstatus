package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestDataResponse(t *testing.T) {
	staticTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    any
		expected resources.Response
	}{
		"as null": {
			input: nil,
			expected: resources.Response{
				Data: nil,
			},
		},
		"as domain": {
			input: domain.User{
				ID:        1,
				Name:      "John Doe",
				Email:     "john@example.com",
				Nickname:  "Johnny",
				Blocked:   false,
				Birthdate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt: staticTime,
				UpdatedAt: staticTime,
			},
			expected: resources.Response{
				Data: domain.User{
					ID:        1,
					Name:      "John Doe",
					Email:     "john@example.com",
					Nickname:  "Johnny",
					Blocked:   false,
					Birthdate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					CreatedAt: staticTime,
					UpdatedAt: staticTime,
				},
			},
		},
		"as collection": {
			input: []string{
				"hello",
				"world",
			},
			expected: resources.Response{
				Data: []string{
					"hello",
					"world",
				},
			},
		},
		"as domain collection": {
			input: []domain.User{
				{
					ID:        1,
					Name:      "John Doe",
					Email:     "john@example.com",
					Nickname:  "Johnny",
					Blocked:   false,
					Birthdate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					CreatedAt: staticTime,
					UpdatedAt: staticTime,
				},
				{
					ID:        2,
					Name:      "John Doe 2",
					Email:     "john2@example.com",
					Nickname:  "Johnny2",
					Blocked:   false,
					Birthdate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					CreatedAt: staticTime,
					UpdatedAt: staticTime,
				},
			},
			expected: resources.Response{
				Data: []domain.User{
					{
						ID:        1,
						Name:      "John Doe",
						Email:     "john@example.com",
						Nickname:  "Johnny",
						Blocked:   false,
						Birthdate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
						CreatedAt: staticTime,
						UpdatedAt: staticTime,
					},
					{
						ID:        2,
						Name:      "John Doe 2",
						Email:     "john2@example.com",
						Nickname:  "Johnny2",
						Blocked:   false,
						Birthdate: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
						CreatedAt: staticTime,
						UpdatedAt: staticTime,
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dataResource := resources.Response{Data: test.input}

			if !reflect.DeepEqual(dataResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, dataResource)
			}
		})
	}
}
