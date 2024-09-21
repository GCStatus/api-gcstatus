package email

import (
	"bytes"
	"fmt"
	"html/template"
)

// HTML Email Template
const emailTemplate = `
  <main style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; margin: 0;">
    <div style="max-width: 600px; background-color: #ffffff; padding: 20px; border-radius: 5px; margin: 0 auto; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
      <h2>Hello,</h2>
      <p>You're receiving this email because you requested a password change. If you didn't request this, you can safely discard this email.</p>
      <p>To reset your password, please click the button below:</p>
      <p>
        <a href="{{.ResetURL}}" style="background-color: #28a745; color: white; padding: 10px 20px; border-radius: 5px; text-align: center; display: inline-block; text-decoration: none; font-size: 16px;">
          Reset Password
        </a>
      </p>
      <p>If the button doesn't work, please use the following link to reset your password. You can try to click it or just copy and paste on your browser.</p>
      <p><a href="{{.ResetURL}}">{{.ResetURL}}</a></p>
      <div class="footer">
        <p>If you have any questions, please contact support.</p>
      </div>
      <div style="margin-top: 20px; color: #888; text-align: center;">
        <p style="font-size: 1rem;">Graciously,</p>
        <p style="font-size: 1rem; font-weight: 900;">Team GCStatus</p>
      </div>
    </div>
  </main>
`

type EmailData struct {
	ResetURL string
}

func SendPasswordResetEmail(userEmail, resetToken string, sendFunc SendEmailFunc) error {
	resetURL := fmt.Sprintf("https://gcstatus.cloud/password/reset/%s/?email=%s", resetToken, userEmail)

	data := EmailData{
		ResetURL: resetURL,
	}

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	// Execute template with data and store it in a buffer
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	err = sendFunc(userEmail, body.String(), "Password Reset Request")
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
