package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"testing"
	"time"
)

func TestTransformTitle(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputTitle domain.Title
		expected   resources.TitleResource
	}{
		"normal title": {
			inputTitle: domain.Title{
				ID:          1,
				Title:       "Title 1",
				Description: "Title 1",
				Cost:        func() *int { i := 200; return &i }(),
				Purchasable: true,
				Status:      "available",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				TitleRequirements: []domain.TitleRequirement{
					{
						ID:          1,
						Task:        "Do something",
						Key:         "do_something",
						Goal:        10,
						Description: "Do something.",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				},
			},
			expected: resources.TitleResource{
				ID:          1,
				Title:       "Title 1",
				Description: "Title 1",
				Cost:        func() *int { i := 200; return &i }(),
				Purchasable: true,
				Status:      "available",
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				TitleRequirements: resources.TransformTitleRequirements([]domain.TitleRequirement{
					{
						ID:          1,
						Task:        "Do something",
						Key:         "do_something",
						Goal:        10,
						Description: "Do something.",
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				}),
			},
		},
		"missing cost": {
			inputTitle: domain.Title{
				ID:          1,
				Title:       "Title 1",
				Description: "Title 1",
				Purchasable: false,
				Status:      "available",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			expected: resources.TitleResource{
				ID:          1,
				Title:       "Title 1",
				Description: "Title 1",
				Purchasable: false,
				Status:      "available",
				CreatedAt:   utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			titleResource := resources.TransformTitle(test.inputTitle)

			if titleResource.ID != test.expected.ID {
				t.Errorf("Expected ID %d, got %d", test.expected.ID, titleResource.ID)
			}
			if titleResource.Title != test.expected.Title {
				t.Errorf("Expected Title %s, got %s", test.expected.Title, titleResource.Title)
			}
			if titleResource.Description != test.expected.Description {
				t.Errorf("Expected Description %s, got %s", test.expected.Description, titleResource.Description)
			}
			if (titleResource.Cost == nil && test.expected.Cost != nil) || (titleResource.Cost != nil && test.expected.Cost == nil) || (titleResource.Cost != nil && test.expected.Cost != nil && *titleResource.Cost != *test.expected.Cost) {
				t.Errorf("Expected Cost %v, got %v", test.expected.Cost, titleResource.Cost)
			}
			if titleResource.Purchasable != test.expected.Purchasable {
				t.Errorf("Expected Purchasable %v, got %v", test.expected.Purchasable, titleResource.Purchasable)
			}
			if titleResource.Status != test.expected.Status {
				t.Errorf("Expected Status %s, got %s", test.expected.Status, titleResource.Status)
			}
			if titleResource.CreatedAt != test.expected.CreatedAt {
				t.Errorf("Expected CreatedAt %s, got %s", test.expected.CreatedAt, titleResource.CreatedAt)
			}

			for _, tr := range test.inputTitle.TitleRequirements {
				titleRequirementResource := resources.TransformTitleRequirement(tr)

				if titleRequirementResource.ID != tr.ID {
					t.Errorf("Expected ID %d, got %d", tr.ID, titleRequirementResource.ID)
				}
				if titleRequirementResource.Task != tr.Task {
					t.Errorf("Expected Task %s, got %s", tr.Task, titleRequirementResource.Task)
				}
				if titleRequirementResource.Description != tr.Description {
					t.Errorf("Expected Description %s, got %s", tr.Description, titleRequirementResource.Description)
				}
				if titleRequirementResource.Goal != tr.Goal {
					t.Errorf("Expected Goal %d, got %d", tr.Goal, titleRequirementResource.Goal)
				}
				if titleRequirementResource.CreatedAt != utils.FormatTimestamp(tr.CreatedAt) {
					t.Errorf("Expected CreatedAt %s, got %s", tr.CreatedAt, titleRequirementResource.CreatedAt)
				}
			}
		})
	}
}

func TestTransformTitles(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputTitles []domain.Title
		expected    []resources.TitleResource
	}{
		"multiple titles": {
			inputTitles: []domain.Title{
				{
					ID:          1,
					Title:       "Title 1",
					Description: "Description 1",
					Cost:        func() *int { i := 200; return &i }(),
					Purchasable: true,
					Status:      "available",
					CreatedAt:   fixedTime,
					TitleRequirements: []domain.TitleRequirement{
						{
							ID:          1,
							Task:        "Do something",
							Key:         "do_something",
							Goal:        10,
							Description: "Do something.",
							CreatedAt:   fixedTime,
							UpdatedAt:   fixedTime,
						},
					},
				},
				{
					ID:          2,
					Title:       "Title 2",
					Description: "Description 2",
					Cost:        nil,
					Purchasable: false,
					Status:      "unavailable",
					CreatedAt:   fixedTime,
				},
			},
			expected: []resources.TitleResource{
				{
					ID:          1,
					Title:       "Title 1",
					Description: "Description 1",
					Cost:        func() *int { i := 200; return &i }(),
					Purchasable: true,
					Status:      "available",
					CreatedAt:   utils.FormatTimestamp(fixedTime),
					TitleRequirements: resources.TransformTitleRequirements([]domain.TitleRequirement{
						{
							ID:          1,
							Task:        "Do something",
							Key:         "do_something",
							Goal:        10,
							Description: "Do something.",
							CreatedAt:   fixedTime,
							UpdatedAt:   fixedTime,
						},
					}),
				},
				{
					ID:          2,
					Title:       "Title 2",
					Description: "Description 2",
					Cost:        nil,
					Purchasable: false,
					Status:      "unavailable",
					CreatedAt:   utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			titleResources := resources.TransformTitles(test.inputTitles)

			if len(titleResources) != len(test.expected) {
				t.Errorf("Expected %d titles, got %d", len(test.expected), len(titleResources))
			}

			for i := range titleResources {
				if titleResources[i].ID != test.expected[i].ID {
					t.Errorf("Expected ID %d, got %d", test.expected[i].ID, titleResources[i].ID)
				}
				if titleResources[i].Title != test.expected[i].Title {
					t.Errorf("Expected Title %s, got %s", test.expected[i].Title, titleResources[i].Title)
				}
				if titleResources[i].Description != test.expected[i].Description {
					t.Errorf("Expected Description %s, got %s", test.expected[i].Description, titleResources[i].Description)
				}
				if (titleResources[i].Cost == nil && test.expected[i].Cost != nil) || (titleResources[i].Cost != nil && test.expected[i].Cost == nil) || (titleResources[i].Cost != nil && test.expected[i].Cost != nil && *titleResources[i].Cost != *test.expected[i].Cost) {
					t.Errorf("Expected Cost %v, got %v", test.expected[i].Cost, titleResources[i].Cost)
				}
				if titleResources[i].Purchasable != test.expected[i].Purchasable {
					t.Errorf("Expected Purchasable %v, got %v", test.expected[i].Purchasable, titleResources[i].Purchasable)
				}
				if titleResources[i].Status != test.expected[i].Status {
					t.Errorf("Expected Status %s, got %s", test.expected[i].Status, titleResources[i].Status)
				}
				if titleResources[i].CreatedAt != test.expected[i].CreatedAt {
					t.Errorf("Expected CreatedAt %s, got %s", test.expected[i].CreatedAt, titleResources[i].CreatedAt)
				}
			}
		})
	}
}
