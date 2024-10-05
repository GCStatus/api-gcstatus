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

func CreateUserMissionAssignmentTest(t *testing.T) {
	testCases := map[string]struct {
		userMissionAssignment domain.UserMissionAssignment
		mockBehavior          func(mock sqlmock.Sqlmock, userMissionAssignment domain.UserMissionAssignment)
		expectErr             bool
	}{
		"Successfully created": {
			userMissionAssignment: domain.UserMissionAssignment{
				UserID:    1,
				MissionID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionAssignment domain.UserMissionAssignment) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `user_mission_assignments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMissionAssignment.UserID,
						userMissionAssignment.MissionID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			userMissionAssignment: domain.UserMissionAssignment{
				UserID:    1,
				MissionID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionAssignment domain.UserMissionAssignment) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `user_mission_assignments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMissionAssignment.UserID,
						userMissionAssignment.MissionID,
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

			tc.mockBehavior(mock, tc.userMissionAssignment)

			err := db.Create(&tc.userMissionAssignment).Error

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

func TestSoftDeleteUserMissionAssignment(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		userMissionAssignmentID uint
		mockBehavior            func(mock sqlmock.Sqlmock, userMissionAssignmentID uint)
		wantErr                 bool
	}{
		"Can soft delete a title": {
			userMissionAssignmentID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionAssignmentID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_mission_assignments` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), userMissionAssignmentID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			userMissionAssignmentID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionAssignmentID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_mission_assignments` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete title"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.userMissionAssignmentID)

			err := db.Delete(&domain.UserMissionAssignment{}, tc.userMissionAssignmentID).Error

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

func TestUpdateUserMissionAssignment(t *testing.T) {
	testCases := map[string]struct {
		userMissionAssignment domain.UserMissionAssignment
		mockBehavior          func(mock sqlmock.Sqlmock, userMissionAssignment domain.UserMissionAssignment)
		expectError           bool
	}{
		"Success": {
			userMissionAssignment: domain.UserMissionAssignment{
				ID:        1,
				UserID:    1,
				MissionID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionAssignment domain.UserMissionAssignment) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_mission_assignments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMissionAssignment.UserID,
						userMissionAssignment.MissionID,
						userMissionAssignment.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			userMissionAssignment: domain.UserMissionAssignment{
				ID:        1,
				UserID:    1,
				MissionID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userMissionAssignment domain.UserMissionAssignment) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_mission_assignments`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userMissionAssignment.UserID,
						userMissionAssignment.MissionID,
						userMissionAssignment.ID,
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

			tc.mockBehavior(mock, tc.userMissionAssignment)

			err := db.Save(&tc.userMissionAssignment).Error

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

func TestValidateUserMissionAssignmentValidData(t *testing.T) {
	testCases := map[string]struct {
		userMissionAssignment domain.UserMissionAssignment
	}{
		"Can empty validations errors": {
			userMissionAssignment: domain.UserMissionAssignment{
				ID: 1,
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

			err := tc.userMissionAssignment.ValidateUserMissionAssignment()

			assert.NoError(t, err)
		})
	}
}

func TestCreateUserMissionAssignmentWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		title   domain.UserMissionAssignment
		wantErr string
	}{
		"Missing required fields": {
			title: domain.UserMissionAssignment{},
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

			err := tc.title.ValidateUserMissionAssignment()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
