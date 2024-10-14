package resources

import (
	"encoding/json"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"log"
)

type NotificationResource struct {
	ID        uint                    `json:"id"`
	Data      domain.NotificationData `json:"content"`
	ReadAt    *string                 `json:"read_at"`
	CreatedAt string                  `json:"created_at"`
	UpdatedAt string                  `json:"updated_at"`
}

func TransformNotification(notification domain.Notification) NotificationResource {
	resource := NotificationResource{
		ID:        notification.ID,
		CreatedAt: utils.FormatTimestamp(notification.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(notification.UpdatedAt),
	}

	if notification.ReadAt != nil {
		formattedReadAt := utils.FormatTimestamp(*notification.ReadAt)
		resource.ReadAt = &formattedReadAt
	}

	var data domain.NotificationData
	err := json.Unmarshal([]byte(notification.Data), &data)
	if err != nil {
		log.Fatalf("failed to parse notification data: %+v", err)
		return resource
	}

	resource.Data = data

	return resource
}

func TransformNotifications(notifications []domain.Notification) []NotificationResource {
	var resources []NotificationResource

	for _, notification := range notifications {
		resources = append(resources, TransformNotification(notification))
	}

	return resources
}
