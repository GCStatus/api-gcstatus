package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"net/http"
)

type NotificationService struct {
	repo ports.NotificationRepository
}

func NewNotificationService(repo ports.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) GetAllForUser(userID uint) ([]domain.Notification, error) {
	return s.repo.GetAllForUser(userID)
}

func (s *NotificationService) CreateNotification(notification *domain.Notification) error {
	return s.repo.CreateNotification(notification)
}

func (s *NotificationService) DeleteNotification(userID, id uint) error {
	notification, err := s.repo.GetNotificationByID(id)
	if err != nil {
		return err
	}

	if err = s.hasAccessToNotification(userID, notification); err != nil {
		return err
	}

	return s.repo.DeleteNotification(id)
}

func (s *NotificationService) MarkAsRead(userID, id uint) error {
	notification, err := s.repo.GetNotificationByID(id)
	if err != nil {
		return err
	}

	if err = s.hasAccessToNotification(userID, notification); err != nil {
		return err
	}

	return s.repo.MarkAsRead(id)
}

func (s *NotificationService) MarkAsUnread(userID, id uint) error {
	notification, err := s.repo.GetNotificationByID(id)
	if err != nil {
		return err
	}

	if err = s.hasAccessToNotification(userID, notification); err != nil {
		return err
	}

	return s.repo.MarkAsUnread(id)
}

func (s *NotificationService) GetNotificationByID(userID, id uint) (domain.Notification, error) {
	notification, err := s.repo.GetNotificationByID(id)
	if err != nil {
		return notification, err
	}

	if err = s.hasAccessToNotification(userID, notification); err != nil {
		return notification, err
	}

	return notification, nil
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *NotificationService) MarkAllAsUnread(userID uint) error {
	return s.repo.MarkAllAsUnread(userID)
}

func (s *NotificationService) DeleteAllNotifications(userID uint) error {
	return s.repo.DeleteAllNotifications(userID)
}

func (s *NotificationService) hasAccessToNotification(userID uint, notification domain.Notification) error {
	if notification.UserID != userID {
		return errors.NewHttpError(http.StatusForbidden, "User does not have access to this notification.")
	}

	return nil
}
