package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type GameLanguageResource struct {
	ID        uint             `json:"id"`
	Menu      bool             `json:"menu"`
	Dubs      bool             `json:"dubs"`
	Subtitles bool             `json:"subtitles"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	Language  LanguageResource `json:"language"`
}

func TransformGameLanguage(gameLanguage domain.GameLanguage) GameLanguageResource {
	return GameLanguageResource{
		ID:        gameLanguage.ID,
		Menu:      gameLanguage.Menu,
		Dubs:      gameLanguage.Dubs,
		Subtitles: gameLanguage.Subtitles,
		CreatedAt: utils.FormatTimestamp(gameLanguage.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(gameLanguage.UpdatedAt),
		Language:  TransformLanguage(gameLanguage.Language),
	}
}
