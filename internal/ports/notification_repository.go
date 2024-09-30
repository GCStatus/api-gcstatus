package ports

import "gcstatus/internal/domain"

type NotificationRepository interface {
	GetAllForUser(userID uint) ([]domain.Notification, error)
	CreateNotification(*domain.Notification) error
	MarkAsRead(id uint) error
	MarkAsUnread(id uint) error
	DeleteNotification(id uint) error
	GetNotificationByID(id uint) (domain.Notification, error)
	MarkAllAsRead(userID uint) error
	MarkAllAsUnread(userID uint) error
	DeleteAllNotifications(userID uint) error
}
