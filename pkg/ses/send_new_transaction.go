package ses

import (
	"bytes"
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
	"html/template"
)

// Generic Email Template for Transaction Notification
const transactionEmailTemplate = `
  <main style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; margin: 0;">
    <div style="max-width: 600px; background-color: #ffffff; padding: 20px; border-radius: 5px; margin: 0 auto; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
      <h2>Hello, {{.Name}}!</h2>
      <p>You have a new transaction on your account.</p>
      <p>Here are the details of your transaction:</p>
      <ul>
        <li>Amount: <strong>{{.Amount}} coins</strong></li>
        <li>Description: <strong>{{.Description}}</strong></li>
        <li>Type: <strong>{{.TransactionType}}</strong></li>
        <li>Date: <strong>{{.Date}}</strong></li>
      </ul>
      <p>If you have any questions or concerns, feel free to contact our support team.</p>
      <div style="margin-top: 20px; color: #888; text-align: center;">
        <p style="font-size: 1rem;">Graciously,</p>
        <p style="font-size: 1rem; font-weight: 900;">Team GCStatus</p>
      </div>
    </div>
  </main>
`

type TransactionEmailData struct {
	Amount          uint
	Description     string
	TransactionType string
	Date            string
	Name            string
}

func SendTransactionEmail(user *domain.User, transaction *domain.Transaction, sendFunc SendEmailFunc) error {
	transactionType := "Unknown"
	if transaction.TransactionTypeID == domain.AdditionTransactionTypeID {
		transactionType = domain.AdditionTransactionType
	} else if transaction.TransactionTypeID == domain.SubtractionTransactionTypeID {
		transactionType = domain.SubtractionTransactionType
	}

	date := transaction.CreatedAt.Format("2006-01-02 15:04")
	fName, _ := utils.GetFirstAndLastName(user.Name)

	data := TransactionEmailData{
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		TransactionType: transactionType,
		Date:            date,
		Name:            fName,
	}

	tmpl, err := template.New("transactionEmail").Parse(transactionEmailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	err = sendFunc(user.Email, body.String(), "You have a new transaction")
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
