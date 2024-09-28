package ports

type WalletRepository interface {
	Add(userID uint, amount uint) error
	Subtract(userID uint, amount uint) error
}
