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

func TestCreateProtection(t *testing.T) {
	testCases := map[string]struct {
		protection   domain.Protection
		mockBehavior func(mock sqlmock.Sqlmock, protection domain.Protection)
		expectError  bool
	}{
		"Success": {
			protection: domain.Protection{
				Name: "Denuvo",
				Slug: "denuvo",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, protection domain.Protection) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `protections`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						protection.Name,
						protection.Slug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			protection: domain.Protection{
				Name: "Denuvo",
				Slug: "denuvo",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, protection domain.Protection) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `protections`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						protection.Name,
						protection.Slug,
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

			tc.mockBehavior(mock, tc.protection)

			err := db.Create(&tc.protection).Error

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

func TestUpdateProtection(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		protection   domain.Protection
		mockBehavior func(mock sqlmock.Sqlmock, protection domain.Protection)
		expectError  bool
	}{
		"Success": {
			protection: domain.Protection{
				ID:        1,
				Name:      "Denuvo",
				Slug:      "denuvo",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, protection domain.Protection) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `protections`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						protection.Name,
						protection.Slug,
						protection.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			protection: domain.Protection{
				ID:        1,
				Name:      "Denuvo",
				Slug:      "denuvo",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, protection domain.Protection) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `protections`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						protection.Name,
						protection.Slug,
						protection.ID,
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

			tc.mockBehavior(mock, tc.protection)

			err := db.Save(&tc.protection).Error

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

func TestSoftDeleteProtection(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		protectionID uint
		mockBehavior func(mock sqlmock.Sqlmock, protectionID uint)
		wantErr      bool
	}{
		"Can soft delete a Protection": {
			protectionID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, protectionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `protections` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), protectionID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			protectionID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, protectionID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `protections` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Protection"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.protectionID)

			err := db.Delete(&domain.Protection{}, tc.protectionID).Error

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

func TestValidateProtection(t *testing.T) {
	testCases := map[string]struct {
		protection domain.Protection
	}{
		"Can empty validations errors": {
			protection: domain.Protection{
				Name: "Denuvo",
				Slug: "denuvo",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.protection.ValidateProtection()
			assert.NoError(t, err)
		})
	}
}

func TestCreateProtectionWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		protection domain.Protection
		wantErr    string
	}{
		"Missing required fields": {
			protection: domain.Protection{},
			wantErr: `
				Name is a required field,
				Slug is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.protection.ValidateProtection()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
