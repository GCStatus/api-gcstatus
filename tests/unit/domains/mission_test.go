package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateMissionTest(t *testing.T) {
	testCases := map[string]struct {
		mission      domain.Mission
		mockBehavior func(mock sqlmock.Sqlmock, mission domain.Mission)
		expectErr    bool
	}{
		"Successfully created": {
			mission: domain.Mission{
				Mission:     "Mission 1",
				Description: "Mission 1",
				Status:      "available",
				ForAll:      true,
				Coins:       100,
				Experience:  100,
				Frequency:   domain.OneTimeMission,
				ResetTime:   time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mission domain.Mission) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mission.Mission,
						mission.Description,
						mission.Status,
						mission.ForAll,
						mission.Coins,
						mission.Experience,
						mission.Frequency,
						mission.ResetTime,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			mission: domain.Mission{
				Mission:     "Mission 1",
				Description: "Mission 1",
				Status:      "available",
				ForAll:      true,
				Coins:       100,
				Experience:  100,
				Frequency:   domain.OneTimeMission,
				ResetTime:   time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mission domain.Mission) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mission.Mission,
						mission.Description,
						mission.Status,
						mission.ForAll,
						mission.Coins,
						mission.Experience,
						mission.Frequency,
						mission.ResetTime,
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

			tc.mockBehavior(mock, tc.mission)

			err := db.Create(&tc.mission).Error

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

func TestSoftDeleteMission(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		MissionID    uint
		mockBehavior func(mock sqlmock.Sqlmock, MissionID uint)
		wantErr      bool
	}{
		"Can soft delete a Mission": {
			MissionID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, MissionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `missions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), MissionID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			MissionID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, MissionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `missions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete mission"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.MissionID)

			err := db.Delete(&domain.Mission{}, tc.MissionID).Error

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

func TestUpdateMission(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		Mission      domain.Mission
		mockBehavior func(mock sqlmock.Sqlmock, Mission domain.Mission)
		expectError  bool
	}{
		"Success": {
			Mission: domain.Mission{
				ID:          1,
				Mission:     "Mission 1",
				Description: "Mission 1",
				Status:      "available",
				ForAll:      true,
				Coins:       100,
				Experience:  100,
				Frequency:   domain.OneTimeMission,
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mission domain.Mission) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mission.Mission,
						mission.Description,
						mission.Status,
						mission.ForAll,
						mission.Coins,
						mission.Experience,
						mission.Frequency,
						mission.ResetTime,
						mission.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			Mission: domain.Mission{
				ID:          1,
				Mission:     "Mission 1",
				Description: "Mission 1",
				Status:      "available",
				ForAll:      true,
				Coins:       100,
				Experience:  100,
				Frequency:   domain.OneTimeMission,
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mission domain.Mission) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `missions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mission.Mission,
						mission.Description,
						mission.Status,
						mission.ForAll,
						mission.Coins,
						mission.Experience,
						mission.Frequency,
						mission.ResetTime,
						mission.ID,
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

			tc.mockBehavior(mock, tc.Mission)

			err := db.Save(&tc.Mission).Error

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

func TestValidateMissionValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		mission domain.Mission
	}{
		"Can empty validations errors": {
			mission: domain.Mission{
				ID:          1,
				Mission:     "Mission 1",
				Description: "Mission 1",
				Status:      "available",
				ForAll:      true,
				Coins:       100,
				Experience:  100,
				Frequency:   domain.OneTimeMission,
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.mission.ValidateMission()

			assert.NoError(t, err)
		})
	}
}

func TestCreateMissionWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		mission domain.Mission
		wantErr string
	}{
		"Missing required fields": {
			mission: domain.Mission{},
			wantErr: "Mission is a required field, Description is a required field, Status is a required field, Coins is a required field, Experience is a required field, Frequency is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.mission.ValidateMission()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
