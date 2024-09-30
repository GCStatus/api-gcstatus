package tests

import (
	"encoding/json"
	"errors"
	"gcstatus/internal/domain"
	"log"
	"testing"
	"time"
)

type MockNotificationRepository struct {
	notifications     map[uint]*domain.Notification
	userNotifications map[uint][]uint
}

func NewMockNotificationRepository() *MockNotificationRepository {
	return &MockNotificationRepository{
		notifications:     make(map[uint]*domain.Notification),
		userNotifications: make(map[uint][]uint),
	}
}

func (m *MockNotificationRepository) GetAllForUser(userID uint) ([]domain.Notification, error) {
	var notifications []domain.Notification
	notificationIDs := m.userNotifications[userID]

	for _, notificationID := range notificationIDs {
		if notification, exists := m.notifications[notificationID]; exists {
			notifications = append(notifications, *notification)
		}
	}

	return notifications, nil
}

func (m *MockNotificationRepository) AddUserNotification(userID uint, notificationID uint) {
	m.userNotifications[userID] = append(m.userNotifications[userID], notificationID)
}

func (m *MockNotificationRepository) CreateNotification(notification *domain.Notification) error {
	if notification == nil {
		return errors.New("invalid notification data")
	}
	m.notifications[notification.ID] = notification
	return nil
}

func (m *MockNotificationRepository) GetNotificationByID(notificationID uint) (*domain.Notification, error) {
	if notification, exists := m.notifications[notificationID]; exists {
		return notification, nil
	}
	return nil, errors.New("notification not found")
}

func (m *MockNotificationRepository) DeleteNotification(notificationID uint) error {
	if _, exists := m.notifications[notificationID]; !exists {
		return errors.New("notification not found")
	}
	delete(m.notifications, notificationID)
	return nil
}

func (m *MockNotificationRepository) MarkAsRead(notificationID uint) error {
	now := time.Now()
	notification, err := m.GetNotificationByID(notificationID)
	if err != nil {
		return err
	}
	notification.ReadAt = &now
	return nil
}

func (m *MockNotificationRepository) MarkAsUnread(notificationID uint) error {
	notification, err := m.GetNotificationByID(notificationID)
	if err != nil {
		return err
	}
	notification.ReadAt = nil
	return nil
}

func (m *MockNotificationRepository) MarkAllAsRead(userID uint) error {
	now := time.Now()
	for _, notification := range m.notifications {
		if notification.UserID == userID {
			notification.ReadAt = &now
		}
	}
	return nil
}

func (m *MockNotificationRepository) MarkAllAsUnread(userID uint) error {
	for _, notification := range m.notifications {
		if notification.UserID == userID {
			notification.ReadAt = nil
		}
	}
	return nil
}

func (m *MockNotificationRepository) DeleteAllNotifications(userID uint) error {
	for notificationID, notification := range m.notifications {
		if notification.UserID == userID {
			delete(m.notifications, notificationID)
		}
	}
	return nil
}

func (m *MockNotificationRepository) MockNotificationRepository_GetAllForUser(t *testing.T) {
	fixedTime := time.Now()

	notificationContent := &domain.NotificationData{
		Title:     "You have some test notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Fatalf("failed to marshal notification content: %+v", err)
	}

	testCases := map[string]struct {
		userID                     uint
		expectedNotificationsCount int
		mockCreateNotifications    func(repo *MockNotificationRepository)
	}{
		"multiple notifications for user 1": {
			userID:                     1,
			expectedNotificationsCount: 2,
			mockCreateNotifications: func(repo *MockNotificationRepository) {
				err := repo.CreateNotification(&domain.Notification{
					ID:        1,
					Type:      "NewTestType",
					Data:      string(dataJson),
					ReadAt:    &fixedTime,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the notification: %s", err.Error())
				}
				err = repo.CreateNotification(&domain.Notification{
					ID:        2,
					Type:      "NewTestType",
					Data:      string(dataJson),
					ReadAt:    nil,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the notification: %s", err.Error())
				}

				repo.AddUserNotification(1, 1)
				repo.AddUserNotification(1, 2)
			},
		},
		"no notifications for user 1": {
			userID:                     1,
			expectedNotificationsCount: 0,
			mockCreateNotifications:    func(repo *MockNotificationRepository) {},
		},
		"notifications for user 2": {
			userID:                     2,
			expectedNotificationsCount: 1,
			mockCreateNotifications: func(repo *MockNotificationRepository) {
				err := repo.CreateNotification(&domain.Notification{
					ID:        3,
					Type:      "NewTestType",
					Data:      string(dataJson),
					ReadAt:    nil,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the notification: %s", err.Error())
				}

				repo.AddUserNotification(2, 3)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockNotificationRepository()

			tc.mockCreateNotifications(mockRepo)

			notifications, err := mockRepo.GetAllForUser(tc.userID)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(notifications) != tc.expectedNotificationsCount {
				t.Fatalf("expected %d notifications, got %d", tc.expectedNotificationsCount, len(notifications))
			}
		})
	}
}

func TestMockNotificationRepository_CreateNotification(t *testing.T) {
	mockRepo := NewMockNotificationRepository()
	fixedTime := time.Now()

	notificationContent := &domain.NotificationData{
		Title:     "You have some test notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Fatalf("failed to marshal notification content: %+v", err)
	}

	testCases := map[string]struct {
		input         *domain.Notification
		expectedError bool
	}{
		"valid input": {
			input: &domain.Notification{
				ID:        1,
				Type:      "NewTestType",
				Data:      string(dataJson),
				ReadAt:    &fixedTime,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expectedError: false,
		},
		"nil input": {
			input:         nil,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.CreateNotification(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.notifications[tc.input.ID] == nil {
					t.Fatalf("expected notification to be created, but it wasn't")
				}
			}
		})
	}
}

func TestMockNotificationRepository_GetNotificationByID(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockNotificationRepository()

	notificationContent := &domain.NotificationData{
		Title:     "You have some test notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Fatalf("failed to marshal notification content: %+v", err)
	}

	err = mockRepo.CreateNotification(&domain.Notification{
		ID:        1,
		Type:      "NewTestType",
		Data:      string(dataJson),
		ReadAt:    &fixedTime,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	})
	if err != nil {
		t.Fatalf("failed to create the notification: %s", err.Error())
	}

	testCases := map[string]struct {
		notificationID uint
		expectedError  bool
	}{
		"valid notification ID": {
			notificationID: 1,
			expectedError:  false,
		},
		"invalid notification ID": {
			notificationID: 999,
			expectedError:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			notification, err := mockRepo.GetNotificationByID(tc.notificationID)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if notification != nil {
					t.Fatalf("expected nil user, got %v", notification)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if notification == nil || notification.ID != tc.notificationID {
					t.Fatalf("expected notification ID %d, got %v", tc.notificationID, notification)
				}
			}
		})
	}
}

func TestMockNotificationRepository_DeleteNotification(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockNotificationRepository()

	notificationContent := &domain.NotificationData{
		Title:     "You have some test notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Fatalf("failed to marshal notification content: %+v", err)
	}

	err = mockRepo.CreateNotification(&domain.Notification{
		ID:        1,
		Type:      "NewTestType",
		Data:      string(dataJson),
		ReadAt:    &fixedTime,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	})
	if err != nil {
		t.Fatalf("failed to create the notification: %s", err.Error())
	}

	testCases := map[string]struct {
		id            uint
		expectedError bool
	}{
		"valid ID": {
			id:            1,
			expectedError: false,
		},
		"invalid ID": {
			id:            999,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.DeleteNotification(tc.id)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.notifications[tc.id] != nil {
					t.Fatalf("expected password reset to be deleted, but it wasn't")
				}
			}
		})
	}
}

func TestMockNotificationRepository_MarkAsRead(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockNotificationRepository()

	notificationContent := &domain.NotificationData{
		Title:     "Test Notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		t.Fatalf("failed to marshal notification content: %+v", err)
	}

	err = mockRepo.CreateNotification(&domain.Notification{
		ID:        1,
		Type:      "NewTestType",
		Data:      string(dataJson),
		ReadAt:    nil,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	})
	if err != nil {
		t.Fatalf("failed to create notification: %s", err.Error())
	}

	testCases := map[string]struct {
		id            uint
		expectedError bool
	}{
		"valid ID": {
			id:            1,
			expectedError: false,
		},
		"invalid ID": {
			id:            999,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.MarkAsRead(tc.id)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.notifications[tc.id].ReadAt == nil {
					t.Fatalf("expected notification to be marked as read, but it wasn't")
				}
			}
		})
	}
}

func TestMockNotificationRepository_MarkAsUnread(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockNotificationRepository()

	notificationContent := &domain.NotificationData{
		Title:     "Test Notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		t.Fatalf("failed to marshal notification content: %+v", err)
	}

	err = mockRepo.CreateNotification(&domain.Notification{
		ID:        1,
		Type:      "NewTestType",
		Data:      string(dataJson),
		ReadAt:    &fixedTime,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	})
	if err != nil {
		t.Fatalf("failed to create notification: %s", err.Error())
	}

	testCases := map[string]struct {
		id            uint
		expectedError bool
	}{
		"valid ID": {
			id:            1,
			expectedError: false,
		},
		"invalid ID": {
			id:            999,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.MarkAsUnread(tc.id)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.notifications[tc.id].ReadAt != nil {
					t.Fatalf("expected notification to be marked as unread, but it wasn't")
				}
			}
		})
	}
}

func TestMockNotificationRepository_MarkAllAsRead(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockNotificationRepository()

	userID := uint(1)
	notificationContent := &domain.NotificationData{
		Title:     "Test Notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		t.Fatalf("failed to marshal notification content: %+v", err)
	}

	for i := 0; i < 3; i++ {
		err = mockRepo.CreateNotification(&domain.Notification{
			ID:        uint(i + 1),
			Type:      "NewTestType",
			Data:      string(dataJson),
			ReadAt:    nil,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
			UserID:    userID,
		})
		if err != nil {
			t.Fatalf("failed to create notification: %s", err.Error())
		}
	}

	err = mockRepo.MarkAllAsRead(userID)
	if err != nil {
		t.Fatalf("unexpected error while marking all as read: %v", err)
	}

	for i := 1; i <= 3; i++ {
		if mockRepo.notifications[uint(i)].ReadAt == nil {
			t.Fatalf("expected notification %d to be marked as read, but it wasn't", i)
		}
	}
}

func TestMockNotificationRepository_MarkAllAsUnread(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockNotificationRepository()

	userID := uint(1)
	notificationContent := &domain.NotificationData{
		Title:     "Test Notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		t.Fatalf("failed to marshal notification content: %+v", err)
	}

	for i := 0; i < 3; i++ {
		err = mockRepo.CreateNotification(&domain.Notification{
			ID:        uint(i + 1),
			Type:      "NewTestType",
			Data:      string(dataJson),
			ReadAt:    &fixedTime,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
			UserID:    userID,
		})
		if err != nil {
			t.Fatalf("failed to create notification: %s", err.Error())
		}
	}

	err = mockRepo.MarkAllAsUnread(userID)
	if err != nil {
		t.Fatalf("unexpected error while marking all as unread: %v", err)
	}

	for i := 1; i <= 3; i++ {
		if mockRepo.notifications[uint(i)].ReadAt != nil {
			t.Fatalf("expected notification %d to be marked as unread, but it wasn't", i)
		}
	}
}

func TestMockNotificationRepository_DeleteAllNotifications(t *testing.T) {
	mockRepo := NewMockNotificationRepository()

	userID := uint(1)
	notificationContent := &domain.NotificationData{
		Title:     "Test Notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		t.Fatalf("failed to marshal notification content: %+v", err)
	}

	for i := 0; i < 3; i++ {
		err = mockRepo.CreateNotification(&domain.Notification{
			ID:        uint(i + 1),
			Type:      "NewTestType",
			Data:      string(dataJson),
			ReadAt:    nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userID,
		})
		if err != nil {
			t.Fatalf("failed to create notification: %s", err.Error())
		}
	}

	err = mockRepo.DeleteAllNotifications(userID)
	if err != nil {
		t.Fatalf("unexpected error while deleting all notifications: %v", err)
	}

	for i := 1; i <= 3; i++ {
		if _, exists := mockRepo.notifications[uint(i)]; exists {
			t.Fatalf("expected notification %d to be deleted, but it still exists", i)
		}
	}
}
