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

func TestCreatePermission(t *testing.T) {
	testCases := map[string]struct {
		permission   domain.Permission
		mockBehavior func(mock sqlmock.Sqlmock, permission domain.Permission)
		expectError  bool
	}{
		"Success": {
			permission: domain.Permission{
				Scope: "test:permission",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permission domain.Permission) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `permissions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permission.Scope,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			permission: domain.Permission{
				Scope: "test:permission",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permission domain.Permission) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `permissions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permission.Scope,
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

			tc.mockBehavior(mock, tc.permission)

			err := db.Create(&tc.permission).Error

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

func TestUpdatePermission(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		permission   domain.Permission
		mockBehavior func(mock sqlmock.Sqlmock, permission domain.Permission)
		expectError  bool
	}{
		"Success": {
			permission: domain.Permission{
				ID:        1,
				Scope:     "test:permission",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permission domain.Permission) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permission.Scope,
						permission.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			permission: domain.Permission{
				ID:        1,
				Scope:     "test:permission",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, permission domain.Permission) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						permission.Scope,
						permission.ID,
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

			tc.mockBehavior(mock, tc.permission)

			err := db.Save(&tc.permission).Error

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

func TestSoftDeletePermission(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		permissionID uint
		mockBehavior func(mock sqlmock.Sqlmock, permissionID uint)
		wantErr      bool
	}{
		"Can soft delete a Permission": {
			permissionID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, permissionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), permissionID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			permissionID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, permissionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `permissions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Permission"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.permissionID)

			err := db.Delete(&domain.Permission{}, tc.permissionID).Error

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

func TestValidatePermission(t *testing.T) {
	testCases := map[string]struct {
		permission domain.Permission
	}{
		"Can empty validations errors": {
			permission: domain.Permission{
				Scope: "test:permission",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.permission.ValidatePermission()
			assert.NoError(t, err)
		})
	}
}

func TestCreatePermissionWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		permission domain.Permission
		wantErr    string
	}{
		"Missing required fields": {
			permission: domain.Permission{},
			wantErr: `
				Scope is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.permission.ValidatePermission()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
