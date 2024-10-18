package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateMissionProgressTest(t *testing.T) {
	testCases := map[string]struct {
		missionProgress domain.MissionProgress
		mockBehavior    func(mock sqlmock.Sqlmock, missionProgress domain.MissionProgress)
		expectErr       bool
	}{
		"Successfully created": {
			missionProgress: domain.MissionProgress{
				Progress:  5,
				Completed: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionProgress domain.MissionProgress) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `mission_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionProgress.Progress,
						missionProgress.Completed,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			missionProgress: domain.MissionProgress{
				Progress:  5,
				Completed: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionProgress domain.MissionProgress) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `mission_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionProgress.Progress,
						missionProgress.Completed,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := testutils.Setup(t)

			tc.mockBehavior(mock, tc.missionProgress)

			err := db.Create(&tc.missionProgress).Error

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

func TestSoftDeleteMissionProgress(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		missionProgressID uint
		mockBehavior      func(mock sqlmock.Sqlmock, missionProgressID uint)
		wantErr           bool
	}{
		"Can soft delete a title progress": {
			missionProgressID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, missionProgressID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_progresses` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), missionProgressID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			missionProgressID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, missionProgressID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_progresses` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete title requirement"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.missionProgressID)

			err := db.Delete(&domain.MissionProgress{}, tc.missionProgressID).Error

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

func TestUpdateMissionProgress(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		missionProgress domain.MissionProgress
		mockBehavior    func(mock sqlmock.Sqlmock, missionProgress domain.MissionProgress)
		expectError     bool
	}{
		"Success": {
			missionProgress: domain.MissionProgress{
				ID:                   1,
				Progress:             5,
				Completed:            false,
				CreatedAt:            fixedTime,
				UpdatedAt:            fixedTime,
				UserID:               1,
				MissionRequirementID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionProgress domain.MissionProgress) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionProgress.Progress,
						missionProgress.Completed,
						missionProgress.UserID,
						missionProgress.MissionRequirementID,
						missionProgress.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			missionProgress: domain.MissionProgress{
				ID:                   1,
				Progress:             5,
				Completed:            false,
				CreatedAt:            fixedTime,
				UpdatedAt:            fixedTime,
				UserID:               1,
				MissionRequirementID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionProgress domain.MissionProgress) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionProgress.Progress,
						missionProgress.Completed,
						missionProgress.UserID,
						missionProgress.MissionRequirementID,
						missionProgress.ID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := testutils.Setup(t)

			tc.mockBehavior(mock, tc.missionProgress)

			err := db.Save(&tc.missionProgress).Error

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

func TestValidateMissionProgressValidData(t *testing.T) {
	testCases := map[string]struct {
		missionProgress domain.MissionProgress
	}{
		"Can empty validations errors": {
			missionProgress: domain.MissionProgress{
				Progress:  5,
				Completed: false,
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

			err := tc.missionProgress.ValidateMissionProgress()

			assert.NoError(t, err)
		})
	}
}

func TestCreateMissionProgressWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		missionProgress domain.MissionProgress
		wantErr         string
	}{
		"Missing required fields": {
			missionProgress: domain.MissionProgress{},
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
				Amount is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.missionProgress.ValidateMissionProgress()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
