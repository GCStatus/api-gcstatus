package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type MorphsFormat struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type LanguageFormat struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ISO       string `json:"iso"`
	Menu      bool   `json:"menu"`
	Dubs      bool   `json:"dubs"`
	Subtitles bool   `json:"subtitles"`
}

type RequirementTypeFormat struct {
	ID        uint   `json:"id"`
	OS        string `json:"os"`
	Potential string `json:"potential"`
}

type RequirementFormat struct {
	ID              uint                  `json:"id"`
	OS              string                `json:"os"`
	DX              string                `json:"dx"`
	CPU             string                `json:"cpu"`
	RAM             string                `json:"ram"`
	GPU             string                `json:"gpu"`
	ROM             string                `json:"rom"`
	OBS             *string               `json:"obs,omitempty"`
	Network         string                `json:"network"`
	RequirementType RequirementTypeFormat `json:"requirement_type"`
}

type CrackByFormat struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Acting bool   `json:"acting"`
}

type ProtectionFormat struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CrackFormat struct {
	ID         uint              `json:"id"`
	Status     string            `json:"status"`
	CrackedAt  *string           `json:"cracked_at"`
	By         *CrackByFormat    `json:"by"`
	Protection *ProtectionFormat `json:"protection"`
}

type GameResource struct {
	ID               uint                `json:"id"`
	Age              uint                `json:"age"`
	Slug             string              `json:"slug"`
	Title            string              `json:"title"`
	Condition        string              `json:"condition"`
	Cover            string              `json:"cover"`
	About            string              `json:"about"`
	Description      string              `json:"description"`
	ShortDescription string              `json:"short_description"`
	Free             bool                `json:"is_free"`
	Legal            *string             `json:"legal"`
	Website          *string             `json:"website"`
	ReleaseDate      string              `json:"release_date"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
	Categories       []MorphsFormat      `json:"categories"`
	Platforms        []MorphsFormat      `json:"platforms"`
	Genres           []MorphsFormat      `json:"genres"`
	Tags             []MorphsFormat      `json:"tags"`
	Languages        []LanguageFormat    `json:"languages"`
	Requirements     []RequirementFormat `json:"requirements"`
	Crack            *CrackFormat        `json:"crack"`
}

func TransformGame(game domain.Game) GameResource {
	resource := GameResource{
		ID:               game.ID,
		Age:              uint(game.Age),
		Slug:             game.Slug,
		Title:            game.Title,
		Condition:        game.Condition,
		Cover:            game.Cover,
		About:            game.About,
		Description:      game.Description,
		ShortDescription: game.ShortDescription,
		Free:             game.Free,
		Legal:            game.Legal,
		Website:          game.Website,
		ReleaseDate:      utils.FormatTimestamp(game.ReleaseDate),
		CreatedAt:        utils.FormatTimestamp(game.CreatedAt),
		UpdatedAt:        utils.FormatTimestamp(game.UpdatedAt),
		Categories:       []MorphsFormat{},
		Platforms:        []MorphsFormat{},
		Genres:           []MorphsFormat{},
		Tags:             []MorphsFormat{},
		Languages:        []LanguageFormat{},
		Requirements:     []RequirementFormat{},
		Crack:            nil,
	}

	for _, categoriable := range game.Categories {
		if categoriable.Category.ID != 0 {
			category := MorphsFormat{
				ID:   categoriable.Category.ID,
				Name: categoriable.Category.Name,
			}

			resource.Categories = append(resource.Categories, category)
		}
	}

	for _, platformable := range game.Platforms {
		if platformable.Platform.ID != 0 {
			platform := MorphsFormat{
				ID:   platformable.Platform.ID,
				Name: platformable.Platform.Name,
			}

			resource.Platforms = append(resource.Platforms, platform)
		}
	}

	for _, genreable := range game.Genres {
		if genreable.Genre.ID != 0 {
			genre := MorphsFormat{
				ID:   genreable.Genre.ID,
				Name: genreable.Genre.Name,
			}

			resource.Genres = append(resource.Genres, genre)
		}
	}

	for _, taggable := range game.Tags {
		if taggable.Tag.ID != 0 {
			tag := MorphsFormat{
				ID:   taggable.Tag.ID,
				Name: taggable.Tag.Name,
			}

			resource.Tags = append(resource.Tags, tag)
		}
	}

	for _, gameLanguage := range game.Languages {
		if gameLanguage.Language.ID != 0 {
			language := LanguageFormat{
				ID:        gameLanguage.Language.ID,
				Name:      gameLanguage.Language.Name,
				ISO:       gameLanguage.Language.ISO,
				Menu:      gameLanguage.Menu,
				Dubs:      gameLanguage.Dubs,
				Subtitles: gameLanguage.Subtitles,
			}

			resource.Languages = append(resource.Languages, language)
		}
	}

	for _, gameRequirement := range game.Requirements {
		if gameRequirement.ID != 0 {
			requirement := RequirementFormat{
				ID:              gameRequirement.ID,
				OS:              gameRequirement.OS,
				DX:              gameRequirement.DX,
				CPU:             gameRequirement.CPU,
				RAM:             gameRequirement.RAM,
				GPU:             gameRequirement.GPU,
				ROM:             gameRequirement.ROM,
				OBS:             gameRequirement.OBS,
				Network:         gameRequirement.Network,
				RequirementType: RequirementTypeFormat{},
			}

			if gameRequirement.RequirementType.ID != 0 {
				requirement.RequirementType = RequirementTypeFormat{
					ID:        gameRequirement.RequirementType.ID,
					OS:        gameRequirement.RequirementType.OS,
					Potential: gameRequirement.RequirementType.Potential,
				}
			}

			resource.Requirements = append(resource.Requirements, requirement)
		}
	}

	if game.Crack != nil && game.Crack.ID != 0 {
		crack := CrackFormat{
			ID:     game.Crack.ID,
			Status: game.Crack.Status,
		}

		if game.Crack.CrackedAt != nil {
			formattedTime := utils.FormatTimestamp(*game.Crack.CrackedAt)
			crack.CrackedAt = &formattedTime
		}

		if game.Crack.Cracker.ID != 0 {
			crack.By = &CrackByFormat{
				ID:     game.Crack.Cracker.ID,
				Name:   game.Crack.Cracker.Name,
				Acting: game.Crack.Cracker.Acting,
			}
		} else {
			crack.By = nil
		}

		if game.Crack.Protection.ID != 0 {
			crack.Protection = &ProtectionFormat{
				ID:   game.Crack.Protection.ID,
				Name: game.Crack.Protection.Name,
			}
		} else {
			crack.Protection = nil
		}

		resource.Crack = &crack
	}

	return resource
}

func TransformGames(games []domain.Game) []GameResource {
	var resources []GameResource

	resources = make([]GameResource, 0, len(games))

	for _, game := range games {
		resources = append(resources, TransformGame(game))
	}

	return resources
}
