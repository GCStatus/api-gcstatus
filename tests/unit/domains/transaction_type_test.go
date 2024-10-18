package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateTransactionTypeTest(t *testing.T) {
	testCases := map[string]struct {
		transactionType domain.TransactionType
		mockBehavior    func(mock sqlmock.Sqlmock, transaction domain.TransactionType)
		expectErr       bool
	}{
		"Successfully created": {
			transactionType: domain.TransactionType{
				Type: domain.AdditionTransactionType,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, transactionType domain.TransactionType) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `transaction_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						transactionType.Type,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			transactionType: domain.TransactionType{
				Type: domain.AdditionTransactionType,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, transactionType domain.TransactionType) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `transaction_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						transactionType.Type,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := testutils.Setup(t)

			tc.mockBehavior(mock, tc.transactionType)

			err := db.Create(&tc.transactionType).Error

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

func TestSoftDeleteTransactionType(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		transactionTypeID uint
		mockFunc          func()
		wantErr           bool
	}{
		"Valid soft delete": {
			transactionTypeID: 1,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `transaction_types` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			transactionTypeID: 2,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `transaction_types` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete transaction"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			err := db.Delete(&domain.TransactionType{}, tc.transactionTypeID).Error

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

func TestValidateTransactionTypeValidData(t *testing.T) {
	testCases := map[string]struct {
		transactionType domain.TransactionType
	}{
		"Can empty validations errors": {
			transactionType: domain.TransactionType{
				Type: domain.AdditionTransactionType,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.transactionType.ValidateTransactionType()
			assert.NoError(t, err)
		})
	}
}

func TestCreateTransactionTypeWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		transactionType domain.TransactionType
		wantErr         string
	}{
		"Missing required fields": {
			transactionType: domain.TransactionType{},
			wantErr:         "Type is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.transactionType.ValidateTransactionType()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
