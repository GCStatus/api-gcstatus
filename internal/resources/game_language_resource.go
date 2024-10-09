package resources

import "gcstatus/internal/domain"

type GameLanguageResource struct {
	ID        uint             `json:"id"`
	Menu      bool             `json:"menu"`
	Dubs      bool             `json:"dubs"`
	Subtitles bool             `json:"subtitles"`
	Language  LanguageResource `json:"language"`
}

func TransformGameLanguage(gameLanguage domain.GameLanguage) GameLanguageResource {
	return GameLanguageResource{
		ID:        gameLanguage.ID,
		Menu:      gameLanguage.Menu,
		Dubs:      gameLanguage.Dubs,
		Subtitles: gameLanguage.Subtitles,
		Language:  TransformLanguage(gameLanguage.Language),
	}
}
