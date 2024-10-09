package resources

import "gcstatus/internal/domain"

type LanguageResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	ISO  string `json:"iso"`
}

func TransformLanguage(language domain.Language) LanguageResource {
	return LanguageResource{
		ID:   language.ID,
		Name: language.Name,
		ISO:  language.ISO,
	}
}

func TransformLanguages(languages []domain.Language) []LanguageResource {
	var resources []LanguageResource

	for _, language := range languages {
		resources = append(resources, TransformLanguage(language))
	}

	return resources
}
