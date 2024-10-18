package tests

import (
	"errors"
	"gcstatus/internal/adapters/db"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepositoryMySQL_Add(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db.NewWalletRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID          uint
		amount          uint
		setupMock       func()
		expectedError   error
		expectedMessage string
	}{
		"user not found": {
			userID: 1,
			amount: 100,
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `wallets` SET `amount`=amount + ? WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
					WithArgs(100, 1).
					WillReturnError(errors.New("no rows affected"))
				mock.ExpectRollback()
			},
			expectedError:   errors.New("no rows affected"),
			expectedMessage: "no rows affected",
		},
		"successful addition": {
			userID: 2,
			amount: 200,
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `wallets` SET `amount`=amount + ? WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
					WithArgs(200, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		"commit error": {
			userID: 3,
			amount: 300,
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `wallets` SET `amount`=amount + ? WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
					WithArgs(300, 3).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			},
			expectedError:   errors.New("commit failed"),
			expectedMessage: "commit failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.setupMock()

			err := repo.Add(tc.userID, tc.amount)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestWalletRepositoryMySQL_Subtract(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db.NewWalletRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID          uint
		amount          uint
		setupMock       func()
		expectedError   error
		expectedMessage string
	}{
		"user not found": {
			userID: 1,
			amount: 100,
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `wallets` SET `amount`=amount - ? WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
					WithArgs(100, 1).
					WillReturnError(errors.New("no rows affected"))
				mock.ExpectRollback()
			},
			expectedError:   errors.New("no rows affected"),
			expectedMessage: "no rows affected",
		},
		"successful addition": {
			userID: 2,
			amount: 200,
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `wallets` SET `amount`=amount - ? WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
					WithArgs(200, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		"commit error": {
			userID: 3,
			amount: 300,
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `wallets` SET `amount`=amount - ? WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
					WithArgs(300, 3).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			},
			expectedError:   errors.New("commit failed"),
			expectedMessage: "commit failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.setupMock()

			err := repo.Subtract(tc.userID, tc.amount)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
