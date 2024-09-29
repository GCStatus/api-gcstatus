package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"
)

type MockTransactionRepository struct {
	transactions     map[uint]*domain.Transaction
	userTransactions map[uint][]uint
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions:     make(map[uint]*domain.Transaction),
		userTransactions: make(map[uint][]uint),
	}
}

func (m *MockTransactionRepository) GetAllForUser(userID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	transactionIDs := m.userTransactions[userID]

	for _, transactionID := range transactionIDs {
		if transaction, exists := m.transactions[transactionID]; exists {
			transactions = append(transactions, *transaction)
		}
	}

	return transactions, nil
}

func (m *MockTransactionRepository) AddUserTransaction(userID uint, transactionID uint) {
	m.userTransactions[userID] = append(m.userTransactions[userID], transactionID)
}

func (m *MockTransactionRepository) CreateTransaction(transaction *domain.Transaction) error {
	if transaction == nil {
		return errors.New("invalid transaction data")
	}
	m.transactions[transaction.ID] = transaction
	return nil
}

func (m *MockTransactionRepository) MockTransactionRepository_GetAllForUser(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		userID                    uint
		expectedTransactionsCount int
		mockCreateTransactions    func(repo *MockTransactionRepository)
	}{
		"multiple transactions for user 1": {
			userID:                    1,
			expectedTransactionsCount: 2,
			mockCreateTransactions: func(repo *MockTransactionRepository) {
				err := repo.CreateTransaction(&domain.Transaction{
					ID:          1,
					Amount:      200,
					Description: "Transaction 1",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the transaction: %s", err.Error())
				}
				err = repo.CreateTransaction(&domain.Transaction{
					ID:          2,
					Amount:      200,
					Description: "Transaction 2",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the transaction: %s", err.Error())
				}

				repo.AddUserTransaction(1, 1)
				repo.AddUserTransaction(1, 2)
			},
		},
		"no transactions for user 1": {
			userID:                    1,
			expectedTransactionsCount: 0,
			mockCreateTransactions:    func(repo *MockTransactionRepository) {},
		},
		"transactions for user 2": {
			userID:                    2,
			expectedTransactionsCount: 1,
			mockCreateTransactions: func(repo *MockTransactionRepository) {
				err := repo.CreateTransaction(&domain.Transaction{
					ID:          3,
					Amount:      200,
					Description: "Transaction 3",
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				})
				if err != nil {
					t.Fatalf("failed to create the transaction: %s", err.Error())
				}

				repo.AddUserTransaction(2, 3)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockTransactionRepository()

			tc.mockCreateTransactions(mockRepo)

			transactions, err := mockRepo.GetAllForUser(tc.userID)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(transactions) != tc.expectedTransactionsCount {
				t.Fatalf("expected %d transactions, got %d", tc.expectedTransactionsCount, len(transactions))
			}
		})
	}
}

func TestMockTransactionRepository_CreateTransaction(t *testing.T) {
	mockRepo := NewMockTransactionRepository()
	fixedTime := time.Now()

	testCases := map[string]struct {
		input         *domain.Transaction
		expectedError bool
	}{
		"valid input": {
			input: &domain.Transaction{
				ID:          1,
				Amount:      200,
				Description: "transaction 1",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			expectedError: false,
		},
		"nil input": {
			input:         nil,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.CreateTransaction(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.transactions[tc.input.ID] == nil {
					t.Fatalf("expected password reset to be created, but it wasn't")
				}
			}
		})
	}
}
