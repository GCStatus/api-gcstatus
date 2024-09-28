package usecases

import "gcstatus/internal/ports"

type WalletService struct {
	repo ports.WalletRepository
}

func NewWalletService(repo ports.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (r *WalletService) Add(userID uint, amount uint) error {
	return r.repo.Add(userID, amount)
}

func (r *WalletService) Subtract(userID uint, amount uint) error {
	return r.repo.Subtract(userID, amount)
}
