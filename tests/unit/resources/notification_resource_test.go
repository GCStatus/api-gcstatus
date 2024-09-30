package tests

import (
	"encoding/json"
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"log"
	"testing"
	"time"
)

func TestTransformNotification(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)
	formattedReadAt := utils.FormatTimestamp(fixedTime)

	notificationContent := &domain.NotificationData{
		Title:     "You have some test notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Fatalf("failed to marshal notification content: %+v", err)
	}

	tests := map[string]struct {
		input    domain.Notification
		expected resources.NotificationResource
	}{
		"normal notification": {
			input: domain.Notification{
				ID:        1,
				Data:      string(dataJson),
				ReadAt:    &fixedTime,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.NotificationResource{
				ID: 1,
				Data: domain.NotificationData{
					Title:     notificationContent.Title,
					Icon:      notificationContent.Icon,
					ActionUrl: notificationContent.ActionUrl,
				},
				ReadAt:    &formattedReadAt,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
		"notification without read_at": {
			input: domain.Notification{
				ID:        2,
				Data:      string(dataJson),
				ReadAt:    nil,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.NotificationResource{
				ID: 2,
				Data: domain.NotificationData{
					Title:     notificationContent.Title,
					Icon:      notificationContent.Icon,
					ActionUrl: notificationContent.ActionUrl,
				},
				ReadAt:    nil,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformNotification(test.input)

			if test.expected.ReadAt != nil {
				if result.ReadAt == nil || *result.ReadAt != *test.expected.ReadAt {
					t.Errorf("Expected ReadAt %+v, got %+v", test.expected.ReadAt, result.ReadAt)
				}
			} else if result.ReadAt != nil {
				t.Errorf("Expected ReadAt nil, got %+v", result.ReadAt)
			}

			if result.ID != test.expected.ID ||
				result.Data.ActionUrl != test.expected.Data.ActionUrl ||
				result.Data.Icon != test.expected.Data.Icon ||
				result.Data.Title != test.expected.Data.Title ||
				result.CreatedAt != test.expected.CreatedAt ||
				result.UpdatedAt != test.expected.UpdatedAt {
				t.Errorf("Expected %+v, got %+v", test.expected, result)
			}
		})
	}
}

func TestTransformNotifications(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)
	formattedReadAt := utils.FormatTimestamp(fixedTime)

	notificationContent := &domain.NotificationData{
		Title:     "You have some test notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Fatalf("failed to marshal notification content: %+v", err)
	}

	tests := map[string]struct {
		input    []domain.Notification
		expected []resources.NotificationResource
	}{
		"multiple notifications": {
			input: []domain.Notification{
				{
					ID:        1,
					Type:      "NewTestType",
					Data:      string(dataJson),
					ReadAt:    &fixedTime,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Type:      "NewTestType",
					Data:      string(dataJson),
					ReadAt:    nil,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources.NotificationResource{
				{
					ID: 1,
					Data: domain.NotificationData{
						Title:     notificationContent.Title,
						Icon:      notificationContent.Icon,
						ActionUrl: notificationContent.ActionUrl,
					},
					ReadAt:    &formattedReadAt,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID: 2,
					Data: domain.NotificationData{
						Title:     notificationContent.Title,
						Icon:      notificationContent.Icon,
						ActionUrl: notificationContent.ActionUrl,
					},
					ReadAt:    nil,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"empty notifications": {
			input:    []domain.Notification{},
			expected: []resources.NotificationResource{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformNotifications(test.input)

			if len(result) != len(test.expected) {
				t.Errorf("Expected %d requirements, got %d", len(test.expected), len(result))
				return
			}

			for i := range result {
				if test.expected[i].ReadAt != nil {
					if result[i].ReadAt == nil || *result[i].ReadAt != *test.expected[i].ReadAt {
						t.Errorf("Expected ReadAt %+v, got %+v", test.expected[i].ReadAt, result[i].ReadAt)
					}
				} else if result[i].ReadAt != nil {
					t.Errorf("Expected ReadAt nil, got %+v", result[i].ReadAt)
				} else {
					if result[i] != test.expected[i] {
						t.Errorf("Expected %+v, got %+v", test.expected[i], result[i])
					}
				}
			}
		})
	}
}
