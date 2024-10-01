package ses

import (
	"bytes"
	"fmt"
	"gcstatus/pkg/utils"
	"html/template"
)

// HTML Email Template for Password Reset Confirmation
const confirmationEmailTemplate = `
  <main style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; margin: 0;">
    <div style="max-width: 600px; background-color: #ffffff; padding: 20px; border-radius: 5px; margin: 0 auto; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
      <h2>Hello, {{.Name}}!</h2>
      <p>We wanted to let you know that your password has been successfully reset.</p>
      <p>If you didn't request this change or if you believe this was done in error, please contact our support team immediately.</p>
      <p>For security reasons, we recommend you check your recent activities and ensure everything is as expected.</p>
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

type ConfirmationEmailData struct {
	Name string
}

func SendPasswordResetConfirmationEmail(userEmail, name string, sendFunc SendEmailFunc) error {
	firstName, _ := utils.GetFirstAndLastName(name)

	data := ConfirmationEmailData{
		Name: firstName,
	}

	tmpl, err := template.New("confirmationEmail").Parse(confirmationEmailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	// Execute template with data and store it in a buffer
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	err = sendFunc(userEmail, body.String(), "Password Reset Confirmation")
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
