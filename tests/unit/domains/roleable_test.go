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

func TestCreateRoleable(t *testing.T) {
	testCases := map[string]struct {
		roleable     domain.Roleable
		mockBehavior func(mock sqlmock.Sqlmock, roleable domain.Roleable)
		expectError  bool
	}{
		"Success": {
			roleable: domain.Roleable{
				RoleableID:   1,
				RoleableType: "users",
				RoleID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, roleable domain.Roleable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `roleables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						roleable.RoleableID,
						roleable.RoleableType,
						roleable.RoleID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			roleable: domain.Roleable{
				RoleableID:   1,
				RoleableType: "users",
				RoleID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, roleable domain.Roleable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `roleables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						roleable.RoleableID,
						roleable.RoleableType,
						roleable.RoleID,
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

			tc.mockBehavior(mock, tc.roleable)

			err := db.Create(&tc.roleable).Error

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

func TestUpdateRoleable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		roleable     domain.Roleable
		mockBehavior func(mock sqlmock.Sqlmock, roleable domain.Roleable)
		expectError  bool
	}{
		"Success": {
			roleable: domain.Roleable{
				ID:           1,
				RoleableID:   1,
				RoleableType: "users",
				RoleID:       1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, roleable domain.Roleable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roleables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						roleable.RoleableID,
						roleable.RoleableType,
						roleable.RoleID,
						roleable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			roleable: domain.Roleable{
				ID:           1,
				RoleableID:   1,
				RoleableType: "users",
				RoleID:       1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, roleable domain.Roleable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roleables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						roleable.RoleableID,
						roleable.RoleableType,
						roleable.RoleID,
						roleable.ID,
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

			tc.mockBehavior(mock, tc.roleable)

			err := db.Save(&tc.roleable).Error

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

func TestSoftDeleteRoleable(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		roleableID   uint
		mockBehavior func(mock sqlmock.Sqlmock, roleableID uint)
		wantErr      bool
	}{
		"Can soft delete a Roleable": {
			roleableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, roleableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roleables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), roleableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			roleableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, roleableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roleables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Roleable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.roleableID)

			err := db.Delete(&domain.Roleable{}, tc.roleableID).Error

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

func TestValidateRoleable(t *testing.T) {
	testCases := map[string]struct {
		roleable domain.Roleable
	}{
		"Can empty validations errors": {
			roleable: domain.Roleable{
				RoleableID:   1,
				RoleableType: "users",
				Role: domain.Role{
					ID:   1,
					Name: "Role 1",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.roleable.ValidateRoleable()
			assert.NoError(t, err)
		})
	}
}

func TestCreateRoleableWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		roleable domain.Roleable
		wantErr  string
	}{
		"Missing required fields": {
			roleable: domain.Roleable{},
			wantErr: `
				Name is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.roleable.ValidateRoleable()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
