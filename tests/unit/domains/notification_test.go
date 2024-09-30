package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotification(t *testing.T) {
	testCases := map[string]struct {
		notification domain.Notification
		mockBehavior func(mock sqlmock.Sqlmock, notification domain.Notification)
		expectError  bool
	}{
		"Success": {
			notification: domain.Notification{
				Type:   "NewTestType",
				Data:   "New datatest",
				UserID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, notification domain.Notification) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `notifications`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						notification.Type,
						notification.Data,
						sqlmock.AnyArg(),
						notification.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			notification: domain.Notification{
				Type:   "NewTestType",
				Data:   "New datatest",
				UserID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, notification domain.Notification) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `notifications`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						notification.Type,
						notification.Data,
						sqlmock.AnyArg(),
						notification.UserID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.notification)

			err := db.Create(&tc.notification).Error

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateNotification(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		notification domain.Notification
		mockBehavior func(mock sqlmock.Sqlmock, level domain.Notification)
		expectError  bool
	}{
		"Success": {
			notification: domain.Notification{
				ID:        1,
				Type:      "NewTestType",
				Data:      "New datatest",
				ReadAt:    &fixedTime,
				UserID:    1,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, notification domain.Notification) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `notifications`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						notification.Type,
						notification.Data,
						notification.ReadAt,
						notification.UserID,
						notification.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			notification: domain.Notification{
				ID:        1,
				Type:      "NewTestType",
				Data:      "New datatest",
				ReadAt:    &fixedTime,
				UserID:    1,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, notification domain.Notification) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `notifications`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						notification.Type,
						notification.Data,
						notification.ReadAt,
						notification.UserID,
						notification.ID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.notification)

			err := db.Save(&tc.notification).Error

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSoftDeleteNotification(t *testing.T) {
	db, mock := tests.Setup(t)

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
			tc.mockBehavior(mock, tc.notificationID)

			err := db.Delete(&domain.Notification{}, tc.notificationID).Error

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetNotificationByID(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		notificationID   uint
		mockFunc         func()
		wantNotification domain.Notification
		wantError        bool
	}{
		"Valid level fetch": {
			notificationID: 1,
			wantNotification: domain.Notification{
				ID:     1,
				Type:   "NewTestType",
				Data:   "New datatest",
				UserID: 1,
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "type", "data", "user_id"}).
					AddRow(1, "NewTestType", "New datatest", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE `notifications`.`id` = ? AND `notifications`.`deleted_at` IS NULL ORDER BY `notifications`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Level not found": {
			notificationID:   2,
			wantNotification: domain.Notification{},
			wantError:        true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `notifications` WHERE `notifications`.`id` = ? AND `notifications`.`deleted_at` IS NULL ORDER BY `notifications`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var level domain.Notification
			err := db.First(&level, tc.notificationID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantNotification, level)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateNotificationValidData(t *testing.T) {
	testCases := map[string]struct {
		notification domain.Notification
	}{
		"Can empty validations errors": {
			notification: domain.Notification{
				Type: "NewTestType",
				Data: "New test data",
				User: domain.User{
					Name:       "Name",
					Email:      "test@example.com",
					Nickname:   "test1",
					Experience: 100,
					Birthdate:  time.Now(),
					Password:   "fakepass123",
					Profile: domain.Profile{
						Share: true,
					},
					Level: domain.Level{
						Level:      1,
						Coins:      100,
						Experience: 100,
					},
					Wallet: domain.Wallet{
						Amount: 100,
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.notification.ValidateNotification()
			assert.NoError(t, err)
		})
	}
}

func TestCreateNotificationWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		notification domain.Notification
		wantErr      string
	}{
		"Missing required fields": {
			notification: domain.Notification{},
			wantErr: `
				Type is a required field,
				Data is a required field,
				Name is a required field,
				Email is a required field,
				Nickname is a required field,
				Birthdate is a required field,
				Password is a required field,
				Share is a required field,
				Level is a required field,
				Experience is a required field,
				Coins is a required field,
				Amount is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.notification.ValidateNotification()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
