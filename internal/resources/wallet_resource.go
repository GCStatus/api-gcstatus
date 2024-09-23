package resources

import "gcstatus/internal/domain"

type WalletResource struct {
	ID     uint `json:"id"`
	Amount uint `json:"amount"`
}

func TransformWallet(wallet *domain.Wallet) *WalletResource {
	return &WalletResource{
		ID:     wallet.ID,
		Amount: uint(wallet.Amount),
	}
}
