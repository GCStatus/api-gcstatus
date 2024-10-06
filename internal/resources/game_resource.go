package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type GameResource struct {
	ID               uint                   `json:"id"`
	Age              uint                   `json:"age"`
	Slug             string                 `json:"slug"`
	Title            string                 `json:"title"`
	Condition        string                 `json:"condition"`
	Cover            string                 `json:"cover"`
	About            string                 `json:"about"`
	Description      string                 `json:"description"`
	ShortDescription string                 `json:"short_description"`
	Free             bool                   `json:"is_free"`
	Legal            *string                `json:"legal"`
	Website          *string                `json:"website"`
	ReleaseDate      string                 `json:"release_date"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
	Categories       []CategoryResource     `json:"categories"`
	Platforms        []PlatformResource     `json:"platforms"`
	Genres           []GenreResource        `json:"genres"`
	Tags             []TagResource          `json:"tags"`
	Languages        []GameLanguageResource `json:"languages"`
	Requirements     []RequirementResource  `json:"requirements"`
	Torrents         []TorrentResource      `json:"torrents"`
	Publishers       []PublisherResource    `json:"publishers"`
	Crack            *CrackResource         `json:"crack"`
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
		Categories:       []CategoryResource{},
		Platforms:        []PlatformResource{},
		Genres:           []GenreResource{},
		Tags:             []TagResource{},
		Languages:        []GameLanguageResource{},
		Requirements:     []RequirementResource{},
		Torrents:         []TorrentResource{},
		Publishers:       []PublisherResource{},
		Crack:            nil,
	}

	for _, categoriable := range game.Categories {
		if categoriable.Category.ID != 0 {
			resource.Categories = append(resource.Categories, TransformCategory(categoriable.Category))
		}
	}

	for _, platformable := range game.Platforms {
		if platformable.Platform.ID != 0 {
			resource.Platforms = append(resource.Platforms, TransformPlatform(platformable.Platform))
		}
	}

	for _, genreable := range game.Genres {
		if genreable.Genre.ID != 0 {
			resource.Genres = append(resource.Genres, TransformGenre(genreable.Genre))
		}
	}

	for _, taggable := range game.Tags {
		if taggable.Tag.ID != 0 {
			resource.Tags = append(resource.Tags, TransformTag(taggable.Tag))
		}
	}

	for _, gameLanguage := range game.Languages {
		if gameLanguage.Language.ID != 0 {
			resource.Languages = append(resource.Languages, TransformGameLanguage(gameLanguage))
		}
	}

	for _, gameRequirement := range game.Requirements {
		if gameRequirement.ID != 0 {
			resource.Requirements = append(resource.Requirements, TransformRequirement(gameRequirement))
		}
	}

	for _, torrent := range game.Torrents {
		if torrent.ID != 0 {
			resource.Torrents = append(resource.Torrents, TransformTorrent(torrent))
		}
	}

	for _, gamePublisher := range game.Publishers {
		if gamePublisher.ID != 0 {
			resource.Publishers = append(resource.Publishers, TransformPublisher(gamePublisher.Publisher))
		}
	}

	if game.Crack != nil && game.Crack.ID != 0 {
		resource.Crack = TransformCrack(game.Crack)
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
