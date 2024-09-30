package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"time"

	"gorm.io/gorm"
)

type NotificationRepositoryMySQL struct {
	db *gorm.DB
}

func NewNotificationRepositoryMySQL(db *gorm.DB) ports.NotificationRepository {
	return &NotificationRepositoryMySQL{db: db}
}

func (h *NotificationRepositoryMySQL) GetAllForUser(userID uint) ([]domain.Notification, error) {
	var notifications []domain.Notification
	err := h.db.Model(&domain.Notification{}).Where("user_id = ?", userID).Find(&notifications).Error
	return notifications, err
}

func (h *NotificationRepositoryMySQL) GetNotificationByID(id uint) (domain.Notification, error) {
	var notification domain.Notification
	err := h.db.First(&notification, id).Error
	return notification, err
}

func (h *NotificationRepositoryMySQL) CreateNotification(notification *domain.Notification) error {
	return h.db.Create(notification).Error
}

func (h *NotificationRepositoryMySQL) DeleteNotification(id uint) error {
	return h.db.Delete(&domain.Notification{}, id).Error
}

func (h *NotificationRepositoryMySQL) MarkAsRead(id uint) error {
	return h.db.Model(&domain.Notification{}).Where("id = ?", id).Update("read_at", time.Now()).Error
}

func (h *NotificationRepositoryMySQL) MarkAsUnread(id uint) error {
	return h.db.Model(&domain.Notification{}).Where("id = ?", id).Update("read_at", nil).Error
}

func (h *NotificationRepositoryMySQL) MarkAllAsRead(userID uint) error {
	return h.db.Model(&domain.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Update("read_at", time.Now()).Error
}

func (h *NotificationRepositoryMySQL) MarkAllAsUnread(userID uint) error {
	return h.db.Model(&domain.Notification{}).Where("user_id = ? AND read_at IS NOT NULL", userID).Update("read_at", nil).Error
}

func (h *NotificationRepositoryMySQL) DeleteAllNotifications(userID uint) error {
	return h.db.Model(&domain.Notification{}).Where("user_id = ?", userID).Delete(&domain.Notification{}).Error
}
