package tests

import (
	"fmt"
	"gcstatus/pkg/ses"
	"strings"
	"testing"
)

func TestSendPasswordResetEmail(t *testing.T) {
	tests := map[string]struct {
		userEmail    string
		resetToken   string
		sendFunc     ses.SendEmailFunc
		expectedBody string
		expectError  bool
	}{
		"successful email": {
			userEmail:    "test@example.com",
			resetToken:   "test-token",
			sendFunc:     MockSendEmail,
			expectedBody: "https://gcstatus.cloud/password/reset/test-token/?email=test@example.com",
			expectError:  false,
		},
		"failed email sending": {
			userEmail:    "fail@example.com",
			resetToken:   "fail-token",
			sendFunc:     MockSendEmail,
			expectedBody: "https://gcstatus.cloud/password/reset/fail-token/?email=fail@example.com",
			expectError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ses.SendPasswordResetEmail(tc.userEmail, tc.resetToken, tc.sendFunc)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError {
				resetURL := fmt.Sprintf("https://gcstatus.cloud/password/reset/%s/?email=%s", tc.resetToken, tc.userEmail)
				if !strings.Contains(tc.expectedBody, resetURL) {
					t.Errorf("Expected reset URL %s in email body, but it was not found", resetURL)
				}
			}
		})
	}
}
