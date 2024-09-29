package ports

import "gcstatus/internal/domain"

type TransactionRepository interface {
	GetAllForUser(userID uint) ([]domain.Transaction, error)
	CreateTransaction(transaction *domain.Transaction) error
}
