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

func CreateMissionRequirementTest(t *testing.T) {
	testCases := map[string]struct {
		missionRequirement domain.MissionRequirement
		mockBehavior       func(mock sqlmock.Sqlmock, missionRequirement domain.MissionRequirement)
		expectErr          bool
	}{
		"Successfully created": {
			missionRequirement: domain.MissionRequirement{
				Task:        "Do something",
				Key:         "do_something",
				Goal:        10,
				Description: "Mission 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionRequirement domain.MissionRequirement) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `mission_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionRequirement.Task,
						missionRequirement.Description,
						missionRequirement.Key,
						missionRequirement.Goal,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			missionRequirement: domain.MissionRequirement{
				Task:        "Do something",
				Key:         "do_something",
				Goal:        10,
				Description: "Mission 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionRequirement domain.MissionRequirement) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `mission_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionRequirement.Task,
						missionRequirement.Description,
						missionRequirement.Key,
						missionRequirement.Goal,
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

			tc.mockBehavior(mock, tc.missionRequirement)

			err := db.Create(&tc.missionRequirement).Error

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

func TestSoftDeleteMissionRequirement(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		missionRequirementID uint
		mockBehavior         func(mock sqlmock.Sqlmock, missionRequirementID uint)
		wantErr              bool
	}{
		"Can soft delete a mission requirement": {
			missionRequirementID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, missionRequirementID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_requirements` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), missionRequirementID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			missionRequirementID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, missionRequirementID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_requirements` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete mission requirement"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.missionRequirementID)

			err := db.Delete(&domain.MissionRequirement{}, tc.missionRequirementID).Error

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

func TestUpdateMissionRequirement(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		missionRequirement domain.MissionRequirement
		mockBehavior       func(mock sqlmock.Sqlmock, mission domain.MissionRequirement)
		expectError        bool
	}{
		"Success": {
			missionRequirement: domain.MissionRequirement{
				ID:          1,
				Task:        "Do something",
				Key:         "do_something",
				Description: "Mission 1",
				Goal:        10,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				MissionID:   1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionRequirement domain.MissionRequirement) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionRequirement.Task,
						missionRequirement.Key,
						missionRequirement.Goal,
						missionRequirement.Description,
						missionRequirement.MissionID,
						missionRequirement.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			missionRequirement: domain.MissionRequirement{
				ID:          1,
				Task:        "Do something",
				Key:         "do_something",
				Description: "Mission 1",
				Goal:        10,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				MissionID:   1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, missionRequirement domain.MissionRequirement) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `mission_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						missionRequirement.Task,
						missionRequirement.Key,
						missionRequirement.Goal,
						missionRequirement.Description,
						missionRequirement.MissionID,
						missionRequirement.ID,
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

			tc.mockBehavior(mock, tc.missionRequirement)

			err := db.Save(&tc.missionRequirement).Error

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

func TestValidateMissionRequirementValidData(t *testing.T) {
	testCases := map[string]struct {
		missionRequirement domain.MissionRequirement
	}{
		"Can empty validations errors": {
			missionRequirement: domain.MissionRequirement{
				ID:          1,
				Task:        "Do something",
				Key:         "do_something",
				Description: "Mission 1",
				Goal:        10,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				MissionID:   1,
				MissionProgress: domain.MissionProgress{
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
				Mission: domain.Mission{
					Mission:     "Mission 1",
					Description: "Mission 1",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  50,
					Frequency:   domain.OneTimeMission,
					ResetTime:   time.Now(),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.missionRequirement.ValidateMissionRequirement()

			assert.NoError(t, err)
		})
	}
}

func TestCreateMissionRequirementWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		missionRequirement domain.MissionRequirement
		wantErr            string
	}{
		"Missing required fields": {
			missionRequirement: domain.MissionRequirement{},
			wantErr: `
				Task is a required field,
				Key is a required field,
				Goal is a required field,
				Description is a required field,
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

			err := tc.missionRequirement.ValidateMissionRequirement()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
