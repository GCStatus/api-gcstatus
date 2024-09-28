package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type WalletRepositoryMySQL struct {
	db *gorm.DB
}

func NewWalletRepositoryMySQL(db *gorm.DB) ports.WalletRepository {
	return &WalletRepositoryMySQL{db: db}
}

func (repo *WalletRepositoryMySQL) Add(userID uint, amount uint) error {
	tx := repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err := tx.Model(&domain.Wallet{}).
		Where("user_id = ?", userID).
		UpdateColumn("amount", gorm.Expr("amount + ?", amount)).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (repo *WalletRepositoryMySQL) Subtract(userID uint, amount uint) error {
	tx := repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err := tx.Model(&domain.Wallet{}).
		Where("user_id = ?", userID).
		UpdateColumn("amount", gorm.Expr("amount - ?", amount)).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
