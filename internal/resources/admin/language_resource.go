package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type LanguageResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ISO       string `json:"iso"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformLanguage(language domain.Language) LanguageResource {
	return LanguageResource{
		ID:        language.ID,
		Name:      language.Name,
		ISO:       language.ISO,
		CreatedAt: utils.FormatTimestamp(language.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(language.UpdatedAt),
	}
}

func TransformLanguages(languages []domain.Language) []LanguageResource {
	var resources []LanguageResource

	for _, language := range languages {
		resources = append(resources, TransformLanguage(language))
	}

	return resources
}
