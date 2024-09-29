package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type TransactionService struct {
	repo ports.TransactionRepository
}

func NewTransactionService(repo ports.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (r *TransactionService) GetAllForUser(userID uint) ([]domain.Transaction, error) {
	return r.repo.GetAllForUser(userID)
}

func (r *TransactionService) CreateTransaction(transaction *domain.Transaction) error {
	return r.repo.CreateTransaction(transaction)
}
