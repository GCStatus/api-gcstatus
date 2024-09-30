package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/email"
	"strings"
	"testing"
	"time"
)

func TestSendNewTransactionEmail(t *testing.T) {
	user := &domain.User{Name: "Test", Email: "test@example.com"}
	transaction := &domain.Transaction{
		Amount:            100,
		Description:       "Test transaction",
		TransactionTypeID: domain.AdditionTransactionTypeID,
		CreatedAt:         time.Now(),
	}

	tests := map[string]struct {
		user         *domain.User
		transaction  *domain.Transaction
		sendFunc     email.SendEmailFunc
		expectedBody string
		expectError  bool
	}{
		"successful email": {
			user:         user,
			transaction:  transaction,
			sendFunc:     MockSendEmail,
			expectedBody: "Hello, Test!",
			expectError:  false,
		},
		"failed email sending": {
			user:         &domain.User{Name: "Test", Email: "fail@example.com"},
			transaction:  transaction,
			sendFunc:     MockSendEmail,
			expectedBody: "Hello, Test!",
			expectError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := email.SendTransactionEmail(tc.user, tc.transaction, tc.sendFunc)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError {
				if !strings.Contains(tc.expectedBody, "Hello, Test!") {
					t.Errorf("Expected greeting 'Hello, Test!' in email body, but it was not found")
				}
			}
		})
	}
}
