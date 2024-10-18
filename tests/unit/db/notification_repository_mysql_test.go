package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNotificationRepositoryMySQL_CreateNotification(t *testing.T) {
	fixedTime := time.Now()
	notificationData := getNotificationData(t)

	testCases := map[string]struct {
		notification *domain.Notification
		mockBehavior func(mock sqlmock.Sqlmock, notification *domain.Notification)
		expectedErr  error
	}{
		"success case": {
			notification: &domain.Notification{
				Type:      "NewTestType",
				Data:      notificationData,
				ReadAt:    &fixedTime,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, notification *domain.Notification) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `notifications`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						notification.Type,
						notification.Data,
						notification.ReadAt,
						notification.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"success case - null read_at": {
			notification: &domain.Notification{
				Type:      "NewTestType",
				Data:      notificationData,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, notification *domain.Notification) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `notifications`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						notification.Type,
						notification.Data,
						nil,
						notification.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"Failure - Insert Error": {
			notification: &domain.Notification{
				Type:      "NewTestType",
				Data:      notificationData,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, notification *domain.Notification) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `notifications`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						notification.Type,
						notification.Data,
						nil,
						notification.UserID,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewNotificationRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.notification)

			err := repo.CreateNotification(tc.notification)

			assert.Equal(t, tc.expectedErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationRepositoryMySQL_GetAllForUser(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)
	repo := db.NewNotificationRepositoryMySQL(gormDB)

	notificationData := getNotificationData(t)

	testCases := map[string]struct {
		userID                uint
		mockSetup             func()
		expectedError         error
		expectedNotifications []domain.Notification
	}{
		"success - notifications found": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE user_id = ? AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type", "data", "read_at", "user_id", "created_at", "updated_at"}).
						AddRow(1, "NewTestType", notificationData, nil, 1, fixedTime, fixedTime).
						AddRow(2, "NewTestType", notificationData, fixedTime, 1, fixedTime, fixedTime))
			},
			expectedNotifications: []domain.Notification{
				{
					ID:        1,
					Type:      "NewTestType",
					Data:      notificationData,
					ReadAt:    nil,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					UserID:    1,
				},
				{
					ID:        2,
					Type:      "NewTestType",
					Data:      notificationData,
					ReadAt:    &fixedTime,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					UserID:    1,
				},
			},
			expectedError: nil,
		},
		"no notifications found": {
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE user_id = ? AND `notifications`.`deleted_at` IS NULL")).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedNotifications: []domain.Notification{},
			expectedError:         nil,
		},
		"error - db failure": {
			userID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE user_id = ? AND `notifications`.`deleted_at` IS NULL")).
					WillReturnError(errors.New("db error"))
			},
			expectedNotifications: nil,
			expectedError:         errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			notifications, err := repo.GetAllForUser(tc.userID)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedNotifications, notifications)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationRepositoryMySQL_GetNotificationByID(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)
	repo := db.NewNotificationRepositoryMySQL(gormDB)

	notificationData := getNotificationData(t)

	testCases := map[string]struct {
		notificationID       uint
		mockSetup            func()
		expectedError        error
		expectedNotification domain.Notification
	}{
		"success - notification found": {
			notificationID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE `notifications`.`id` = ? AND `notifications`.`deleted_at` IS NULL ORDER BY `notifications`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type", "data", "read_at", "user_id", "created_at", "updated_at"}).
						AddRow(1, "NewTestType", notificationData, nil, 1, fixedTime, fixedTime))
			},
			expectedNotification: domain.Notification{
				ID:        1,
				Type:      "NewTestType",
				Data:      notificationData,
				ReadAt:    nil,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			expectedError: nil,
		},
		"error - notification not found": {
			notificationID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE `notifications`.`id` = ? AND `notifications`.`deleted_at` IS NULL ORDER BY `notifications`.`id` LIMIT ?")).
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedNotification: domain.Notification{},
			expectedError:        gorm.ErrRecordNotFound,
		},
		"error - db failure": {
			notificationID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE `notifications`.`id` = ? AND `notifications`.`deleted_at` IS NULL ORDER BY `notifications`.`id` LIMIT ?")).
					WithArgs(3, 1).
					WillReturnError(errors.New("db error"))
			},
			expectedNotification: domain.Notification{},
			expectedError:        errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			notification, err := repo.GetNotificationByID(tc.notificationID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedNotification, notification)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNotificationRepositoryMySQL_DeleteNotification(t *testing.T) {
	testCases := map[string]struct {
		notificationID uint
		mockBehavior   func(mock sqlmock.Sqlmock, notificationID uint)
		wantErr        bool
	}{
		"Can soft delete a notification": {
			notificationID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, notificationID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `notifications` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), notificationID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			notificationID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, notificationID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `notifications` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete notification"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewNotificationRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.notificationID)

			err := repo.DeleteNotification(tc.notificationID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete notification")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNotificationRepositoryMySQL_MarkAsRead(t *testing.T) {
	testCases := map[string]struct {
		notificationID uint
		mockBehavior   func(mock sqlmock.Sqlmock, notificationID uint)
		wantErr        bool
	}{
		"Can mark notification as read": {
			notificationID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, notificationID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE id = ? AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), notificationID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Failed to mark notification as read": {
			notificationID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, notificationID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE id = ? AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), notificationID).
					WillReturnError(fmt.Errorf("failed to mark as read"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewNotificationRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.notificationID)

			err := repo.MarkAsRead(tc.notificationID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to mark as read")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNotificationRepositoryMySQL_MarkAsUnread(t *testing.T) {
	testCases := map[string]struct {
		notificationID uint
		mockBehavior   func(mock sqlmock.Sqlmock, notificationID uint)
		wantErr        bool
	}{
		"Can mark notification as unread": {
			notificationID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, notificationID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE id = ? AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(nil, sqlmock.AnyArg(), notificationID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Failed to mark notification as unread": {
			notificationID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, notificationID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE id = ? AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(nil, sqlmock.AnyArg(), notificationID).
					WillReturnError(fmt.Errorf("failed to mark as unread"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewNotificationRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.notificationID)

			err := repo.MarkAsUnread(tc.notificationID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to mark as unread")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNotificationRepositoryMySQL_MarkAllAsRead(t *testing.T) {
	testCases := map[string]struct {
		userID       uint
		mockBehavior func(mock sqlmock.Sqlmock, userID uint)
		wantErr      bool
	}{
		"Can mark all notifications as read": {
			userID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE (user_id = ? AND read_at IS NULL) AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Failed to mark all notifications as read": {
			userID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE (user_id = ? AND read_at IS NULL) AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
					WillReturnError(fmt.Errorf("failed to mark all as read"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewNotificationRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.userID)

			err := repo.MarkAllAsRead(tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to mark all as read")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNotificationRepositoryMySQL_MarkAllAsUnread(t *testing.T) {
	testCases := map[string]struct {
		userID       uint
		mockBehavior func(mock sqlmock.Sqlmock, userID uint)
		wantErr      bool
	}{
		"Can mark all notifications as unread": {
			userID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE (user_id = ? AND read_at IS NOT NULL) AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(nil, sqlmock.AnyArg(), userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Failed to mark all notifications as unread": {
			userID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `read_at`=?,`updated_at`=? WHERE (user_id = ? AND read_at IS NOT NULL) AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(nil, sqlmock.AnyArg(), userID).
					WillReturnError(fmt.Errorf("failed to mark all as unread"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewNotificationRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.userID)

			err := repo.MarkAllAsUnread(tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to mark all as unread")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNotificationRepositoryMySQL_DeleteAllNotifications(t *testing.T) {
	testCases := map[string]struct {
		userID       uint
		mockBehavior func(mock sqlmock.Sqlmock, userID uint)
		wantErr      bool
	}{
		"Can delete all notifications": {
			userID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `deleted_at`=? WHERE user_id = ? AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Failed to delete all notifications": {
			userID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `notifications` SET `deleted_at`=? WHERE user_id = ? AND `notifications`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), userID).
					WillReturnError(fmt.Errorf("failed to delete notifications"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewNotificationRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.userID)

			err := repo.DeleteAllNotifications(tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete notifications")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func getNotificationData(t *testing.T) string {
	notificationContent := &domain.NotificationData{
		Title:     "Test Notification",
		ActionUrl: "/tests",
		Icon:      "TestIcon",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		t.Fatalf("failed to marshal notification content: %+v", err)
	}

	return string(dataJson)
}
