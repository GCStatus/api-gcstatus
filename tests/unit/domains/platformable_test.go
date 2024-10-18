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

func TestCreatePlatformable(t *testing.T) {
	testCases := map[string]struct {
		platformable domain.Platformable
		mockBehavior func(mock sqlmock.Sqlmock, platformable domain.Platformable)
		expectError  bool
	}{
		"Success": {
			platformable: domain.Platformable{
				Platformable:     1,
				PlatformableType: "games",
				PlatformID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platformable domain.Platformable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `platformables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platformable.PlatformableID,
						platformable.PlatformableType,
						platformable.PlatformID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			platformable: domain.Platformable{
				Platformable:     1,
				PlatformableType: "games",
				PlatformID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platformable domain.Platformable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `platformables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platformable.PlatformableID,
						platformable.PlatformableType,
						platformable.PlatformID,
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

			tc.mockBehavior(mock, tc.platformable)

			err := db.Create(&tc.platformable).Error

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

func TestUpdatePlatformable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		platformable domain.Platformable
		mockBehavior func(mock sqlmock.Sqlmock, platformable domain.Platformable)
		expectError  bool
	}{
		"Success": {
			platformable: domain.Platformable{
				ID:               1,
				Platformable:     1,
				PlatformableType: "games",
				PlatformID:       1,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platformable domain.Platformable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platformables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platformable.PlatformableID,
						platformable.PlatformableType,
						platformable.PlatformID,
						platformable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			platformable: domain.Platformable{
				ID:               1,
				Platformable:     1,
				PlatformableType: "games",
				PlatformID:       1,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platformable domain.Platformable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platformables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platformable.PlatformableID,
						platformable.PlatformableType,
						platformable.PlatformID,
						platformable.ID,
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

			tc.mockBehavior(mock, tc.platformable)

			err := db.Save(&tc.platformable).Error

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

func TestSoftDeletePlatformable(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		platformable uint
		mockBehavior func(mock sqlmock.Sqlmock, platformable uint)
		wantErr      bool
	}{
		"Can soft delete a Platformable": {
			platformable: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, platformable uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platformables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), platformable).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			platformable: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, platformable uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `platformables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Platformable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.platformable)

			err := db.Delete(&domain.Platformable{}, tc.platformable).Error

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
