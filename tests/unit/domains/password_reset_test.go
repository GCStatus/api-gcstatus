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

func TestCreatePasswordReset(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		passwordReset domain.PasswordReset
		mockBehavior  func(mock sqlmock.Sqlmock, passwordReset domain.PasswordReset)
		expectError   bool
	}{
		"Success": {
			passwordReset: domain.PasswordReset{
				Email:     "fake@gmail.com",
				Token:     "mVs0byFtjAoetlNnbk84vOh5BTDT8PTF",
				ExpiresAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, passwordReset domain.PasswordReset) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `password_resets`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						passwordReset.Email,
						passwordReset.Token,
						passwordReset.ExpiresAt,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			passwordReset: domain.PasswordReset{
				Email:     "fake@gmail.com",
				Token:     "mVs0byFtjAoetlNnbk84vOh5BTDT8PTF",
				ExpiresAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, passwordReset domain.PasswordReset) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `password_resets`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						passwordReset.Email,
						passwordReset.Token,
						passwordReset.ExpiresAt,
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

			tc.mockBehavior(mock, tc.passwordReset)

			err := db.Create(&tc.passwordReset).Error

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

func TestSoftDeletePasswordReset(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		passwordResetID uint
		mockFunc        func()
		wantErr         bool
	}{
		"Valid soft delete": {
			passwordResetID: 1,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `password_resets` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			passwordResetID: 2,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `password_resets` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete user"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			err := db.Delete(&domain.PasswordReset{}, tc.passwordResetID).Error

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

func TestValidatePasswordResetValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		passwordReset domain.PasswordReset
	}{
		"Can empty validations errors": {
			passwordReset: domain.PasswordReset{
				Email:     "fake@gmail.com",
				Token:     "mVs0byFtjAoetlNnbk84vOh5BTDT8PTF",
				ExpiresAt: fixedTime,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.passwordReset.ValidatePasswordReset()
			assert.NoError(t, err)
		})
	}
}

func TestCreatePasswordResetWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		passwordReset domain.PasswordReset
		wantErr       string
	}{
		"Missing required fields": {
			passwordReset: domain.PasswordReset{},
			wantErr:       "Email is a required field, Token is a required field, ExpiresAt is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.passwordReset.ValidatePasswordReset()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
