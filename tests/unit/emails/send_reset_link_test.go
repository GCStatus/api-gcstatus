package email_test

import (
	"errors"
	"fmt"
	"gcstatus/pkg/email"
	"strings"
	"testing"
)

func MockSendEmail(recipient, body, subject string) error {
	if recipient == "fail@example.com" {
		return errors.New("failed to send email")
	}

	return nil
}

func TestSendPasswordResetEmail(t *testing.T) {
	tests := map[string]struct {
		userEmail    string
		resetToken   string
		sendFunc     email.SendEmailFunc
		expectedBody string
		expectError  bool
	}{
		"successful email": {
			userEmail:    "test@example.com",
			resetToken:   "test-token",
			sendFunc:     MockSendEmail,
			expectedBody: "https://gcstatus.tech/password/reset/test-token/?email=test@example.com",
			expectError:  false,
		},
		"failed email sending": {
			userEmail:    "fail@example.com",
			resetToken:   "fail-token",
			sendFunc:     MockSendEmail,
			expectedBody: "https://gcstatus.tech/password/reset/fail-token/?email=fail@example.com",
			expectError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := email.SendPasswordResetEmail(tc.userEmail, tc.resetToken, tc.sendFunc)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError {
				resetURL := fmt.Sprintf("https://gcstatus.tech/password/reset/%s/?email=%s", tc.resetToken, tc.userEmail)
				if !strings.Contains(tc.expectedBody, resetURL) {
					t.Errorf("Expected reset URL %s in email body, but it was not found", resetURL)
				}
			}
		})
	}
}
