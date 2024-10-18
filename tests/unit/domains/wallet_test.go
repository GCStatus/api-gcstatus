package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	testCases := map[string]struct {
		wallet       domain.Wallet
		mockBehavior func(mock sqlmock.Sqlmock, wallet domain.Wallet)
		expectError  bool
	}{
		"Success": {
			wallet: domain.Wallet{
				Amount: 0,
				UserID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, wallet domain.Wallet) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `wallets`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						wallet.Amount,
						wallet.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			wallet: domain.Wallet{
				Amount: 0,
				UserID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, wallet domain.Wallet) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `wallets`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						wallet.Amount,
						wallet.UserID,
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

			tc.mockBehavior(mock, tc.wallet)

			err := db.Create(&tc.wallet).Error

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

func TestUpdateWallet(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		wallet       domain.Wallet
		mockBehavior func(mock sqlmock.Sqlmock, wallet domain.Wallet)
		expectError  bool
	}{
		"Success": {
			wallet: domain.Wallet{
				ID:        1,
				Amount:    100,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, wallet domain.Wallet) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `wallets`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						wallet.Amount,
						wallet.UserID,
						wallet.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			wallet: domain.Wallet{
				ID:        1,
				Amount:    0,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, wallet domain.Wallet) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `wallets`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						wallet.Amount,
						wallet.UserID,
						wallet.ID,
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

			tc.mockBehavior(mock, tc.wallet)

			err := db.Save(&tc.wallet).Error

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

func TestSoftDeleteWallet(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		walletID     uint
		mockBehavior func(mock sqlmock.Sqlmock, walletID uint)
		wantErr      bool
	}{
		"Can soft delete a wallet": {
			walletID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, walletID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `wallets` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), walletID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			walletID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, walletID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `wallets` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete wallet"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.walletID)

			err := db.Delete(&domain.Wallet{}, tc.walletID).Error

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

func TestValidateWalletValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		wallet domain.Wallet
	}{
		"Valid wallet with zero amount": {
			wallet: domain.Wallet{
				ID:        1,
				Amount:    100,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.wallet.ValidateWallet()
			assert.NoError(t, err)
		})
	}
}

func TestCreateWalletWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		wallet  domain.Wallet
		wantErr string
	}{
		"Missing required fields": {
			wallet:  domain.Wallet{},
			wantErr: "Amount is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.wallet.ValidateWallet()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
