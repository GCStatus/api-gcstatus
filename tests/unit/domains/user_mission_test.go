package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateUserMissionTest(t *testing.T) {
	testCases := map[string]struct {
		userMission  domain.UserMission
		mockBehavior func(mock sqlmock.Sqlmock, userMission domain.UserMission)
		expectErr    bool
	}{
		"Successfully created": {
			userMission: domain.UserMission{
				Completed:       false,
				LastCompletedAt: time.Now(),
				UserID:          1,
				MissionID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMission domain.UserMission) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `user_missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMission.Completed,
						userMission.LastCompletedAt,
						userMission.UserID,
						userMission.MissionID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			userMission: domain.UserMission{
				Completed:       false,
				LastCompletedAt: time.Now(),
				UserID:          1,
				MissionID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMission domain.UserMission) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `user_missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMission.Completed,
						userMission.LastCompletedAt,
						userMission.UserID,
						userMission.MissionID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.userMission)

			err := db.Create(&tc.userMission).Error

			if tc.expectErr {
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

func TestSoftDeleteUserMission(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		userMissionID uint
		mockBehavior  func(mock sqlmock.Sqlmock, userMissionID uint)
		wantErr       bool
	}{
		"Can soft delete a title": {
			userMissionID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_missions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), userMissionID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			userMissionID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_missions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete title"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.userMissionID)

			err := db.Delete(&domain.UserMission{}, tc.userMissionID).Error

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

func TestUpdateUserMission(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		userMission  domain.UserMission
		mockBehavior func(mock sqlmock.Sqlmock, userMission domain.UserMission)
		expectError  bool
	}{
		"Success": {
			userMission: domain.UserMission{
				ID:              1,
				Completed:       false,
				LastCompletedAt: fixedTime,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
				UserID:          1,
				MissionID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMission domain.UserMission) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMission.Completed,
						userMission.LastCompletedAt,
						userMission.UserID,
						userMission.MissionID,
						userMission.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			userMission: domain.UserMission{
				ID:              1,
				Completed:       false,
				LastCompletedAt: fixedTime,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
				UserID:          1,
				MissionID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMission domain.UserMission) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMission.Completed,
						userMission.LastCompletedAt,
						userMission.UserID,
						userMission.MissionID,
						userMission.ID,
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

			tc.mockBehavior(mock, tc.userMission)

			err := db.Save(&tc.userMission).Error

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

func TestValidateUserMissionValidData(t *testing.T) {
	testCases := map[string]struct {
		userMission domain.UserMission
	}{
		"Can empty validations errors": {
			userMission: domain.UserMission{
				ID:              1,
				Completed:       false,
				LastCompletedAt: time.Now(),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Mission: domain.Mission{
					Mission:     "Mission 1",
					Description: "Description 1",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  50,
					Frequency:   domain.OneTimeMission,
					ResetTime:   time.Now(),
				},
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
			t.Parallel()

			err := tc.userMission.ValidateUserMission()

			assert.NoError(t, err)
		})
	}
}

func TestCreateUserMissionWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		title   domain.UserMission
		wantErr string
	}{
		"Missing required fields": {
			title: domain.UserMission{},
			wantErr: `
				Name is a required field,
				Email is a required field,
				Nickname is a required field,
				Birthdate is a required field,
				Password is a required field,
				Share is a required field,
				Level is a required field,
				Experience is a required field,
				Coins is a required field,
				Amount is a required field,
				Mission is a required field,
				Description is a required field,
				Status is a required field,
				Coins is a required field,
				Experience is a required field,
				Frequency is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.title.ValidateUserMission()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
