package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type TransactionRepositoryMySQL struct {
	db *gorm.DB
}

func NewTransactionRepositoryMySQL(db *gorm.DB) ports.TransactionRepository {
	return &TransactionRepositoryMySQL{db: db}
}

func (h *TransactionRepositoryMySQL) GetAllForUser(userID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := h.db.Model(&domain.Transaction{}).Preload("TransactionType").Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

func (h *TransactionRepositoryMySQL) CreateTransaction(transaction *domain.Transaction) error {
	return h.db.Create(&transaction).Error
}
