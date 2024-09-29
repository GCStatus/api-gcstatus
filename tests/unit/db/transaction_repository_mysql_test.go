package tests

import (
	"errors"
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepositoryMySQL_CreateTransaction(t *testing.T) {
	testCases := map[string]struct {
		transaction  *domain.Transaction
		mockBehavior func(mock sqlmock.Sqlmock, transaction *domain.Transaction)
		expectedErr  error
	}{
		"success case": {
			transaction: &domain.Transaction{
				Amount:      200,
				Description: "Transaction 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, transaction *domain.Transaction) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `transactions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						transaction.Amount,
						transaction.Description,
						transaction.UserID,
						transaction.TransactionTypeID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"Failure - Insert Error": {
			transaction: &domain.Transaction{
				Amount:      200,
				Description: "Transaction 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, transaction *domain.Transaction) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `transactions`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						transaction.Amount,
						transaction.Description,
						transaction.UserID,
						transaction.TransactionTypeID,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := tests.Setup(t)

			repo := db.NewTransactionRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.transaction)

			err := repo.CreateTransaction(tc.transaction)

			assert.Equal(t, tc.expectedErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTransactionRepositoryMySQL_GetAllForUser(t *testing.T) {
	gormDB, mock := tests.Setup(t)
	repo := db.NewTransactionRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID               uint
		mockSetup            func()
		expectedError        error
		expectedTransactions []domain.Transaction
	}{
		"success - transactions found": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? AND `transactions`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "amount", "description", "user_id", "transaction_type_id"}).
						AddRow(1, 200, "Transaction 1", 1, 1).
						AddRow(2, 200, "Transaction 2", 1, 1))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transaction_types` WHERE `transaction_types`.`id` = ? AND `transaction_types`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).
						AddRow(1, "addition"))
			},
			expectedTransactions: []domain.Transaction{
				{
					ID:                1,
					Amount:            200,
					Description:       "Transaction 1",
					UserID:            1,
					TransactionTypeID: 1,
					TransactionType: domain.TransactionType{
						ID:   1,
						Type: "addition",
					},
				},
				{
					ID:                2,
					Amount:            200,
					Description:       "Transaction 2",
					UserID:            1,
					TransactionTypeID: 1,
					TransactionType: domain.TransactionType{
						ID:   1,
						Type: "addition",
					},
				},
			},
			expectedError: nil,
		},
		"no transactions found": {
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? AND `transactions`.`deleted_at` IS NULL")).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedTransactions: []domain.Transaction{},
			expectedError:        nil,
		},
		"error - db failure": {
			userID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? AND `transactions`.`deleted_at` IS NULL")).
					WillReturnError(errors.New("db error"))
			},
			expectedTransactions: nil,
			expectedError:        errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			transactions, err := repo.GetAllForUser(tc.userID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedTransactions, transactions)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
