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

func TestCreatePermissionable(t *testing.T) {
	testCases := map[string]struct {
		permissionable domain.Permissionable
		mockBehavior   func(mock sqlmock.Sqlmock, permissionable domain.Permissionable)
		expectError    bool
	}{
		"Success": {
			permissionable: domain.Permissionable{
				PermissionableID:   1,
				PermissionableType: "users",
				PermissionID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permissionable domain.Permissionable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `permissionables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permissionable.PermissionableID,
						permissionable.PermissionableType,
						permissionable.PermissionID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			permissionable: domain.Permissionable{
				PermissionableID:   1,
				PermissionableType: "users",
				PermissionID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permissionable domain.Permissionable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `permissionables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permissionable.PermissionableID,
						permissionable.PermissionableType,
						permissionable.PermissionID,
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

			tc.mockBehavior(mock, tc.permissionable)

			err := db.Create(&tc.permissionable).Error

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

func TestUpdatePermissionable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		permissionable domain.Permissionable
		mockBehavior   func(mock sqlmock.Sqlmock, permissionable domain.Permissionable)
		expectError    bool
	}{
		"Success": {
			permissionable: domain.Permissionable{
				ID:                 1,
				PermissionableID:   1,
				PermissionableType: "users",
				PermissionID:       1,
				CreatedAt:          fixedTime,
				UpdatedAt:          fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permissionable domain.Permissionable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissionables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permissionable.PermissionableID,
						permissionable.PermissionableType,
						permissionable.PermissionID,
						permissionable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			permissionable: domain.Permissionable{
				ID:                 1,
				PermissionableID:   1,
				PermissionableType: "users",
				PermissionID:       1,
				CreatedAt:          fixedTime,
				UpdatedAt:          fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permissionable domain.Permissionable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissionables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permissionable.PermissionableID,
						permissionable.PermissionableType,
						permissionable.PermissionID,
						permissionable.ID,
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

			tc.mockBehavior(mock, tc.permissionable)

			err := db.Save(&tc.permissionable).Error

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

func TestSoftDeletePermissionable(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		permissionableID uint
		mockBehavior     func(mock sqlmock.Sqlmock, permissionableID uint)
		wantErr          bool
	}{
		"Can soft delete a Permissionable": {
			permissionableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, permissionableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissionables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), permissionableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			permissionableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, permissionableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissionables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Permissionable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.permissionableID)

			err := db.Delete(&domain.Permissionable{}, tc.permissionableID).Error

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

func TestValidatePermissionable(t *testing.T) {
	testCases := map[string]struct {
		permissionable domain.Permissionable
	}{
		"Can empty validations errors": {
			permissionable: domain.Permissionable{
				PermissionableID:   1,
				PermissionableType: "users",
				Permission: domain.Permission{
					ID:    1,
					Scope: "test:permission",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.permissionable.ValidatePermissionable()
			assert.NoError(t, err)
		})
	}
}

func TestCreatePermissionableWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		permissionable domain.Permissionable
		wantErr        string
	}{
		"Missing required fields": {
			permissionable: domain.Permissionable{},
			wantErr: `
				Scope is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.permissionable.ValidatePermissionable()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
