package tests

import "errors"

func MockSendEmail(recipient, body, subject string) error {
	if recipient == "fail@example.com" {
		return errors.New("failed to send email")
	}

	return nil
}
