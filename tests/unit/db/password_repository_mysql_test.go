package tests

import (
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPasswordResetRepositoryMySQL_CreatePasswordReset(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		passwordReset *domain.PasswordReset
		mockBehavior  func(mock sqlmock.Sqlmock, passwordReset *domain.PasswordReset)
		expectedErr   error
	}{
		"success case": {
			passwordReset: &domain.PasswordReset{
				Email:     "fake@gmail.com",
				Token:     "asjkdasjdkajskdajsd",
				ExpiresAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, passwordReset *domain.PasswordReset) {
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
			expectedErr: nil,
		},
		"Failure - Insert Error": {
			passwordReset: &domain.PasswordReset{
				Email:     "fake@gmail.com",
				Token:     "asjkdasjdkajskdajsd",
				ExpiresAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, passwordReset *domain.PasswordReset) {
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
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewPasswordResetRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.passwordReset)

			err := repo.CreatePasswordReset(tc.passwordReset)

			assert.Equal(t, tc.expectedErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPasswordResetRepositoryMySQL_FindPasswordResetByToken(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		token                 string
		mockBehavior          func(mock sqlmock.Sqlmock)
		expectedPasswordReset *domain.PasswordReset
		expectedErr           error
	}{
		"valid token": {
			token: "123456789",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "token", "expires_at"}).
					AddRow(1, "fake@gmail.com", "123456789", fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `password_resets` WHERE token = ? AND `password_resets`.`deleted_at` IS NULL ORDER BY `password_resets`.`id` LIMIT ?")).
					WithArgs("123456789", 1).WillReturnRows(rows)
			},
			expectedPasswordReset: &domain.PasswordReset{ID: 1, Email: "fake@gmail.com", Token: "123456789", ExpiresAt: fixedTime},
			expectedErr:           nil,
		},
		"not found token": {
			token: "987654321",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `password_resets` WHERE token = ? AND `password_resets`.`deleted_at` IS NULL ORDER BY `password_resets`.`id` LIMIT ?")).
					WithArgs("987654321", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedPasswordReset: &domain.PasswordReset{Token: ""},
			expectedErr:           gorm.ErrRecordNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewPasswordResetRepositoryMySQL(gormDB)

			tc.mockBehavior(mock)

			passwordReset, err := repo.FindPasswordResetByToken(tc.token)

			assert.Equal(t, tc.expectedErr, err)
			if err == gorm.ErrRecordNotFound {
				assert.Equal(t, "", passwordReset.Token)
			} else {
				assert.Equal(t, tc.expectedPasswordReset.Token, passwordReset.Token)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPasswordResetRepositoryMySQL_DeletePasswordResetByID(t *testing.T) {
	testCases := map[string]struct {
		passwordResetID uint
		mockBehavior    func(mock sqlmock.Sqlmock, passwordResetID uint)
		wantErr         bool
	}{
		"Can soft delete a password reset": {
			passwordResetID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, passwordResetID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `password_resets` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), passwordResetID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			passwordResetID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, passwordResetID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `password_resets` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete password reset"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewPasswordResetRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.passwordResetID)

			err := repo.DeletePasswordResetByID(tc.passwordResetID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete password reset")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
