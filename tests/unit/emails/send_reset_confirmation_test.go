package tests

import (
	"gcstatus/pkg/ses"
	"strings"
	"testing"
)

func TestSendPasswordResetConfirmationEmail(t *testing.T) {
	tests := map[string]struct {
		userEmail    string
		Name         string
		sendFunc     ses.SendEmailFunc
		expectedBody string
		expectError  bool
	}{
		"successful email": {
			userEmail:    "test@example.com",
			Name:         "Test Name",
			sendFunc:     MockSendEmail,
			expectedBody: "Hello, Test!",
			expectError:  false,
		},
		"failed email sending": {
			userEmail:    "fail@example.com",
			Name:         "Test Name",
			sendFunc:     MockSendEmail,
			expectedBody: "Hello, Test!",
			expectError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ses.SendPasswordResetConfirmationEmail(tc.userEmail, tc.Name, tc.sendFunc)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError {
				if !strings.Contains(tc.expectedBody, "Test") {
					t.Errorf("Expected name %s in email body, but it was not found", "Test")
				}
			}
		})
	}
}
