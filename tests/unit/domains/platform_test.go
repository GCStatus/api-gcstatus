package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePlatform(t *testing.T) {
	testCases := map[string]struct {
		platform     domain.Platform
		mockBehavior func(mock sqlmock.Sqlmock, platform domain.Platform)
		expectError  bool
	}{
		"Success": {
			platform: domain.Platform{
				Name: "Platform 1",
				Slug: "platform-1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platform domain.Platform) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `platforms`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platform.Name,
						platform.Slug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			platform: domain.Platform{
				Name: "Failure",
				Slug: "failure",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platform domain.Platform) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `platforms`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platform.Name,
						platform.Slug,
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

			tc.mockBehavior(mock, tc.platform)

			err := db.Create(&tc.platform).Error

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

func TestUpdatePlatform(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		platform     domain.Platform
		mockBehavior func(mock sqlmock.Sqlmock, platform domain.Platform)
		expectError  bool
	}{
		"Success": {
			platform: domain.Platform{
				ID:        1,
				Name:      "Platform 1",
				Slug:      "platform-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platform domain.Platform) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platforms`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platform.Name,
						platform.Slug,
						platform.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			platform: domain.Platform{
				ID:        1,
				Name:      "Platform 1",
				Slug:      "platform-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platform domain.Platform) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platforms`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platform.Name,
						platform.Slug,
						platform.ID,
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

			tc.mockBehavior(mock, tc.platform)

			err := db.Save(&tc.platform).Error

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

func TestSoftDeletePlatform(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		platformID   uint
		mockBehavior func(mock sqlmock.Sqlmock, platformID uint)
		wantErr      bool
	}{
		"Can soft delete a Platform": {
			platformID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, platformID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platforms` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), platformID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			platformID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, platformID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platforms` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Platform"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.platformID)

			err := db.Delete(&domain.Platform{}, tc.platformID).Error

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

func TestGetPlatformByID(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		platformID   uint
		mockFunc     func()
		wantPlatform domain.Platform
		wantError    bool
	}{
		"Valid Platform fetch": {
			platformID: 1,
			wantPlatform: domain.Platform{
				ID:   1,
				Name: "Platform 1",
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Platform 1")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platforms` WHERE `platforms`.`id` = ? AND `platforms`.`deleted_at` IS NULL ORDER BY `platforms`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Platform not found": {
			platformID:   2,
			wantPlatform: domain.Platform{},
			wantError:    true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platforms` WHERE `platforms`.`id` = ? AND `platforms`.`deleted_at` IS NULL ORDER BY `platforms`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var platform domain.Platform
			err := db.First(&platform, tc.platformID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantPlatform, platform)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidatePlatformValidData(t *testing.T) {
	testCases := map[string]struct {
		platform domain.Platform
	}{
		"Can empty validations errors": {
			platform: domain.Platform{
				Name: "Platform 1",
				Slug: "platform-1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.platform.ValidatePlatform()
			assert.NoError(t, err)
		})
	}
}

func TestCreatePlatformWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		platform domain.Platform
		wantErr  string
	}{
		"Missing required fields": {
			platform: domain.Platform{},
			wantErr:  "Name is a required field, Slug is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.platform.ValidatePlatform()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
