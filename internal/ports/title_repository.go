package ports

import "gcstatus/internal/domain"

type TitleRepository interface {
	GetAllForUser(userID uint) ([]domain.Title, error)
	FindById(titleID uint) (domain.Title, error)
	ToggleEnableTitle(userID uint, titleID uint) error
}
