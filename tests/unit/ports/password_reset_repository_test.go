package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"
)

type MockPasswordResetRepository struct {
	passwordResets map[uint]*domain.PasswordReset
}

func NewMockPasswordResetRepository() *MockPasswordResetRepository {
	return &MockPasswordResetRepository{
		passwordResets: make(map[uint]*domain.PasswordReset),
	}
}

func (m *MockPasswordResetRepository) CreatePasswordReset(passwordReset *domain.PasswordReset) error {
	if passwordReset == nil {
		return errors.New("invalid password reset data")
	}
	m.passwordResets[passwordReset.ID] = passwordReset
	return nil
}

func (m *MockPasswordResetRepository) FindPasswordResetByToken(token string) (*domain.PasswordReset, error) {
	for _, pr := range m.passwordResets {
		if pr.Token == token {
			return pr, nil
		}
	}
	return nil, errors.New("password reset not found")
}

func (m *MockPasswordResetRepository) DeletePasswordResetByID(id uint) error {
	if _, exists := m.passwordResets[id]; !exists {
		return errors.New("password reset not found")
	}
	delete(m.passwordResets, id)
	return nil
}

func TestMockPasswordResetRepository_CreatePasswordReset(t *testing.T) {
	mockRepo := NewMockPasswordResetRepository()
	fixedTime := time.Now()

	testCases := map[string]struct {
		input         *domain.PasswordReset
		expectedError bool
	}{
		"valid input": {
			input: &domain.PasswordReset{
				ID:        1,
				Token:     "validToken123",
				Email:     "valid@gmail.com",
				ExpiresAt: fixedTime,
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

			err := mockRepo.CreatePasswordReset(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.passwordResets[tc.input.ID] == nil {
					t.Fatalf("expected password reset to be created, but it wasn't")
				}
			}
		})
	}
}

func TestMockPasswordResetRepository_FindPasswordResetByToken(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockPasswordResetRepository()

	err := mockRepo.CreatePasswordReset(&domain.PasswordReset{
		ID:        1,
		Token:     "validToken123",
		Email:     "valid@gmail.com",
		ExpiresAt: fixedTime,
	})

	if err != nil {
		t.Fatalf("failed to create the password reset: %s", err.Error())
	}

	testCases := map[string]struct {
		token         string
		expectedError bool
	}{
		"valid token": {
			token:         "validToken123",
			expectedError: false,
		},
		"invalid token": {
			token:         "invalidToken",
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := mockRepo.FindPasswordResetByToken(tc.token)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if result != nil {
					t.Fatalf("expected result to be nil, got %v", result)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if result == nil || result.Token != tc.token {
					t.Fatalf("expected password reset with token %s, got %v", tc.token, result)
				}
			}
		})
	}
}

func TestMockPasswordResetRepository_DeletePasswordResetByID(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockPasswordResetRepository()

	err := mockRepo.CreatePasswordReset(&domain.PasswordReset{
		ID:        1,
		Token:     "validToken123",
		Email:     "valid@gmail.com",
		ExpiresAt: fixedTime,
	})

	if err != nil {
		t.Fatalf("failed to create the password reset: %s", err.Error())
	}

	testCases := map[string]struct {
		id            uint
		expectedError bool
	}{
		"valid ID": {
			id:            1,
			expectedError: false,
		},
		"invalid ID": {
			id:            999,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.DeletePasswordResetByID(tc.id)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.passwordResets[tc.id] != nil {
					t.Fatalf("expected password reset to be deleted, but it wasn't")
				}
			}
		})
	}
}
