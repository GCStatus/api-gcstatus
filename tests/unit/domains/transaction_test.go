package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateTransactionTest(t *testing.T) {
	testCases := map[string]struct {
		transaction  domain.Transaction
		mockBehavior func(mock sqlmock.Sqlmock, transaction domain.Transaction)
		expectErr    bool
	}{
		"Successfully created": {
			transaction: domain.Transaction{
				Amount:            2000,
				Description:       "Transaction test of 2000 coins",
				UserID:            1,
				TransactionTypeID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, transaction domain.Transaction) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `transactions`").
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
		},
		"Failure - Insert Error": {
			transaction: domain.Transaction{
				Amount:            2000,
				Description:       "Transaction test of 2000 coins",
				UserID:            1,
				TransactionTypeID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, transaction domain.Transaction) {
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
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.transaction)

			err := db.Create(&tc.transaction).Error

			if tc.expectErr {
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

func TestSoftDeleteTransaction(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		transactionID uint
		mockFunc      func()
		wantErr       bool
	}{
		"Valid soft delete": {
			transactionID: 1,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `transactions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			transactionID: 2,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `transactions` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete transaction"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			err := db.Delete(&domain.Transaction{}, tc.transactionID).Error

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

func TestValidateTransactionValidData(t *testing.T) {
	testCases := map[string]struct {
		transaction domain.Transaction
	}{
		"Can empty validations errors": {
			transaction: domain.Transaction{
				ID:          1,
				Amount:      200,
				Description: "Test transaction.",
				TransactionType: domain.TransactionType{
					Type: "addition",
				},
				User: domain.User{
					Name:       "Name",
					Email:      "test@example.com",
					Nickname:   "test1",
					Experience: 100,
					Birthdate:  time.Now(),
					Password:   "fakepass123",
					Profile: domain.Profile{
						Share: true,
					},
					Level: domain.Level{
						Level:      1,
						Coins:      100,
						Experience: 100,
					},
					Wallet: domain.Wallet{
						Amount: 100,
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.transaction.ValidateTransaction()
			assert.NoError(t, err)
		})
	}
}

func TestCreateTransactionWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		transaction domain.Transaction
		wantErr     string
	}{
		"Missing required fields": {
			transaction: domain.Transaction{},
			wantErr: `Name is a required field,
				Email is a required field,
				Nickname is a required field,
				Birthdate is a required field,
				Password is a required field,
				Share is a required field,
				Level is a required field,
				Experience is a required field,
				Coins is a required field,
				Amount is a required field,
				Type is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.transaction.ValidateTransaction()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
