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

func TestCreateRole(t *testing.T) {
	testCases := map[string]struct {
		role         domain.Role
		mockBehavior func(mock sqlmock.Sqlmock, role domain.Role)
		expectError  bool
	}{
		"Success": {
			role: domain.Role{
				Name: "Role 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, role domain.Role) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `roles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						role.Name,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			role: domain.Role{
				Name: "Role 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, role domain.Role) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `roles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						role.Name,
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

			tc.mockBehavior(mock, tc.role)

			err := db.Create(&tc.role).Error

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

func TestUpdateRole(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		role         domain.Role
		mockBehavior func(mock sqlmock.Sqlmock, role domain.Role)
		expectError  bool
	}{
		"Success": {
			role: domain.Role{
				ID:        1,
				Name:      "Role 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, role domain.Role) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						role.Name,
						role.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			role: domain.Role{
				ID:        1,
				Name:      "Role 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, role domain.Role) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						role.Name,
						role.ID,
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

			tc.mockBehavior(mock, tc.role)

			err := db.Save(&tc.role).Error

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

func TestSoftDeleteRole(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		roleID       uint
		mockBehavior func(mock sqlmock.Sqlmock, roleID uint)
		wantErr      bool
	}{
		"Can soft delete a Role": {
			roleID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, roleID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), roleID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			roleID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, roleID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `roles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Role"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.roleID)

			err := db.Delete(&domain.Role{}, tc.roleID).Error

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

func TestValidateRole(t *testing.T) {
	testCases := map[string]struct {
		role domain.Role
	}{
		"Can empty validations errors": {
			role: domain.Role{
				Name: "Role 1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.role.ValidateRole()
			assert.NoError(t, err)
		})
	}
}

func TestCreateRoleWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		role    domain.Role
		wantErr string
	}{
		"Missing required fields": {
			role: domain.Role{},
			wantErr: `
				Name is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.role.ValidateRole()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
