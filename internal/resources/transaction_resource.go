package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type TransactionResource struct {
	ID          uint   `json:"id"`
	Amount      uint   `json:"amount"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
}

func TransformTransaction(transaction domain.Transaction) TransactionResource {
	transactionResource := TransactionResource{
		ID:          transaction.ID,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		CreatedAt:   utils.FormatTimestamp(transaction.CreatedAt),
	}

	if transaction.TransactionType.ID != 0 {
		transactionResource.Type = transaction.TransactionType.Type
	}

	return transactionResource
}

func TransformTransactions(transactions []domain.Transaction) []TransactionResource {
	var resources []TransactionResource

	for _, transaction := range transactions {
		resources = append(resources, TransformTransaction(transaction))
	}

	return resources
}
