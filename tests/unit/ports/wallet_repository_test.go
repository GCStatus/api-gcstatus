package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockWalletRepository struct {
	wallet domain.Wallet
}

func NewMockWalletRepository() *MockWalletRepository {
	return &MockWalletRepository{
		wallet: domain.Wallet{},
	}
}

func (m *MockWalletRepository) CreateWallet(user_id uint) error {
	if user_id == 999 {
		return errors.New("user not found")
	}
	m.wallet = domain.Wallet{
		ID:        1,
		Amount:    1000,
		UserID:    user_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return nil
}

func (m *MockWalletRepository) Add(user_id uint, amount uint) error {
	if m.wallet.UserID != user_id {
		return errors.New("no wallet found for given user")
	}

	m.wallet.Amount = m.wallet.Amount + int(amount)

	return nil
}

func (m *MockWalletRepository) Subtract(user_id uint, amount uint) error {
	if m.wallet.UserID != user_id {
		return errors.New("no wallet found for given user")
	}

	if int(amount) > m.wallet.Amount {
		return errors.New("insufficient funds")
	}

	m.wallet.Amount = m.wallet.Amount - int(amount)

	return nil
}

func (m *MockWalletRepository) TestAdd(t *testing.T) {
	mock := NewMockWalletRepository()

	err := mock.CreateWallet(1)
	if err != nil {
		t.Errorf("failed to create wallet for user: %+v", err.Error())
	}

	tests := map[string]struct {
		user_id      uint
		amount       uint
		expect_err   bool
		expect_msg   string
		expected_amt int
	}{
		"success add on user wallet": {
			user_id:      1,
			amount:       100,
			expect_err:   false,
			expected_amt: 1100,
		},
		"no wallet for user": {
			user_id:    999,
			amount:     100,
			expect_err: true,
			expect_msg: "no wallet found for given user",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mock.Add(tc.user_id, tc.amount)

			if tc.expect_err {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.expect_msg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected_amt, mock.wallet.Amount)
			}
		})
	}
}

func (m *MockWalletRepository) TestSubtract(t *testing.T) {
	mock := NewMockWalletRepository()

	err := mock.CreateWallet(1)
	if err != nil {
		t.Errorf("failed to create wallet for user: %+v", err.Error())
	}

	tests := map[string]struct {
		user_id      uint
		amount       uint
		expect_err   bool
		expect_msg   string
		expected_amt int
	}{
		"success subtract on user wallet": {
			user_id:      1,
			amount:       200,
			expect_err:   false,
			expected_amt: 800,
		},
		"insufficient funds": {
			user_id:    1,
			amount:     2000,
			expect_err: true,
			expect_msg: "insufficient funds",
		},
		"no wallet for user": {
			user_id:    999,
			amount:     100,
			expect_err: true,
			expect_msg: "no wallet found for given user",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mock.Subtract(tc.user_id, tc.amount)

			if tc.expect_err {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.expect_msg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected_amt, mock.wallet.Amount)
			}
		})
	}
}
