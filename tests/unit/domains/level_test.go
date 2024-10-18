package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateLevel(t *testing.T) {
	testCases := map[string]struct {
		level        domain.Level
		mockBehavior func(mock sqlmock.Sqlmock, level domain.Level)
		expectError  bool
	}{
		"Success": {
			level: domain.Level{
				Level:      1,
				Experience: 500,
				Coins:      1029,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, level domain.Level) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `levels`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						level.Level,
						level.Experience,
						level.Coins,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			level: domain.Level{
				Level:      1,
				Experience: 500,
				Coins:      1029,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, level domain.Level) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `levels`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						level.Level,
						level.Experience,
						level.Coins,
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

			tc.mockBehavior(mock, tc.level)

			err := db.Create(&tc.level).Error

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

func TestUpdateLevel(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		level        domain.Level
		mockBehavior func(mock sqlmock.Sqlmock, level domain.Level)
		expectError  bool
	}{
		"Success": {
			level: domain.Level{
				ID:         1,
				Level:      1,
				Experience: 500,
				Coins:      1029,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, level domain.Level) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `levels`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						level.Level,
						level.Experience,
						level.Coins,
						level.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			level: domain.Level{
				ID:         1,
				Level:      1,
				Experience: 500,
				Coins:      1029,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, level domain.Level) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `levels`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						level.Level,
						level.Experience,
						level.Coins,
						level.ID,
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

			tc.mockBehavior(mock, tc.level)

			err := db.Save(&tc.level).Error

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

func TestSoftDeleteLevel(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		levelID      uint
		mockBehavior func(mock sqlmock.Sqlmock, levelID uint)
		wantErr      bool
	}{
		"Can soft delete a level": {
			levelID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, levelID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `levels` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), levelID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			levelID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, levelID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `levels` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete level"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.levelID)

			err := db.Delete(&domain.Level{}, tc.levelID).Error

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

func TestGetLevelByID(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		levelID   uint
		mockFunc  func()
		wantLevel domain.Level
		wantError bool
	}{
		"Valid level fetch": {
			levelID: 1,
			wantLevel: domain.Level{
				ID:         1,
				Level:      1,
				Experience: 500,
				Coins:      1029,
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "level", "experience", "coins"}).
					AddRow(1, 1, 500, 1029)
				mock.ExpectQuery("SELECT \\* FROM `levels` WHERE `levels`.`id` = \\? AND `levels`.`deleted_at` IS NULL ORDER BY `levels`.`id` LIMIT \\?").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Level not found": {
			levelID:   2,
			wantLevel: domain.Level{},
			wantError: true,
			mockFunc: func() {
				mock.ExpectQuery("SELECT \\* FROM `levels` WHERE `levels`.`id` = \\? AND `levels`.`deleted_at` IS NULL ORDER BY `levels`.`id` LIMIT \\?").
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var level domain.Level
			err := db.First(&level, tc.levelID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantLevel, level)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateLevelValidData(t *testing.T) {
	testCases := map[string]struct {
		level domain.Level
	}{
		"Can empty validations errors": {
			level: domain.Level{
				Level:      1,
				Experience: 500,
				Coins:      1029,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.level.ValidateLevel()
			assert.NoError(t, err)
		})
	}
}

func TestCreateLevelWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		level   domain.Level
		wantErr string
	}{
		"Missing required fields": {
			level:   domain.Level{},
			wantErr: "Level is a required field, Experience is a required field, Coins is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.level.ValidateLevel()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
