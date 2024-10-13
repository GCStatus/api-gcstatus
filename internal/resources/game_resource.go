package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
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
	ViewsCount       uint                   `json:"views_count"`
	HeartsCount      uint                   `json:"hearts_count"`
	IsHearted        bool                   `json:"is_hearted"`
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
	Developers       []DeveloperResource    `json:"developers"`
	Reviews          []ReviewResource       `json:"reviews"`
	Critics          []CriticableResource   `json:"critics"`
	Stores           []GameStoreResource    `json:"stores"`
	Comments         []CommentableResource  `json:"comments"`
	Galleries        []GalleriableResource  `json:"galleries"`
	DLCs             []DLCResource          `json:"dlcs"`
	Crack            *CrackResource         `json:"crack"`
	Support          *SupportResource       `json:"support"`
}

func TransformGame(game domain.Game, s3Client s3.S3ClientInterface, userID uint) GameResource {
	resource := GameResource{
		ID:               game.ID,
		Age:              uint(game.Age),
		Slug:             game.Slug,
		Title:            game.Title,
		Condition:        game.Condition,
		Cover:            game.Cover,
		About:            game.About,
		Description:      game.Description,
		IsHearted:        false,
		ShortDescription: game.ShortDescription,
		Free:             game.Free,
		Legal:            game.Legal,
		Website:          game.Website,
		ViewsCount:       uint(len(game.Views)),
		HeartsCount:      uint(len(game.Hearts)),
		ReleaseDate:      utils.FormatTimestamp(game.ReleaseDate),
		CreatedAt:        utils.FormatTimestamp(game.CreatedAt),
		UpdatedAt:        utils.FormatTimestamp(game.UpdatedAt),
	}

	resource.Categories = transformCategories(game.Categories)
	resource.Platforms = transformPlatforms(game.Platforms)
	resource.Genres = transformGenres(game.Genres)
	resource.Tags = transformTags(game.Tags)
	resource.Languages = transformLanguages(game.Languages)
	resource.Requirements = transformRequirements(game.Requirements)
	resource.Torrents = transformTorrents(game.Torrents)
	resource.Publishers = transformPublishers(game.Publishers)
	resource.Developers = transformDevelopers(game.Developers)
	resource.Reviews = transformReviews(game.Reviews, s3Client)
	resource.Critics = transformCritics(game.Critics)
	resource.Stores = transformStores(game.Stores)
	resource.Comments = transformComments(game.Comments, s3Client, userID)
	resource.Galleries = transformGalleries(game.Galleries, s3Client)
	resource.DLCs = transformDLCs(game.DLCs, s3Client)

	if game.Crack != nil && game.Crack.ID != 0 {
		resource.Crack = TransformCrack(game.Crack)
	}

	if game.Support != nil && game.Support.ID != 0 {
		resource.Support = TransformSupport(game.Support)
	}

	heartsMap := make(map[uint]bool)
	for _, heart := range game.Hearts {
		heartsMap[heart.UserID] = true
	}

	resource.IsHearted = heartsMap[userID]

	return resource
}

func TransformGames(games []domain.Game, s3Client s3.S3ClientInterface, userID uint) []GameResource {
	var resources []GameResource

	resources = make([]GameResource, 0, len(games))

	for _, game := range games {
		resources = append(resources, TransformGame(game, s3Client, userID))
	}

	return resources
}

func transformCategories(categories []domain.Categoriable) []CategoryResource {
	categoryResources := make([]CategoryResource, 0)
	for _, c := range categories {
		if c.Category.ID != 0 {
			categoryResources = append(categoryResources, TransformCategory(c.Category))
		}
	}

	return categoryResources
}

func transformPlatforms(platforms []domain.Platformable) []PlatformResource {
	platformResources := make([]PlatformResource, 0)
	for _, p := range platforms {
		if p.Platform.ID != 0 {
			platformResources = append(platformResources, TransformPlatform(p.Platform))
		}
	}

	return platformResources
}

func transformGenres(genres []domain.Genreable) []GenreResource {
	genreResources := make([]GenreResource, 0)
	for _, g := range genres {
		if g.Genre.ID != 0 {
			genreResources = append(genreResources, TransformGenre(g.Genre))
		}
	}

	return genreResources
}

func transformTags(tags []domain.Taggable) []TagResource {
	tagResources := make([]TagResource, 0)
	for _, t := range tags {
		if t.Tag.ID != 0 {
			tagResources = append(tagResources, TransformTag(t.Tag))
		}
	}

	return tagResources
}

func transformLanguages(languages []domain.GameLanguage) []GameLanguageResource {
	languageResources := make([]GameLanguageResource, 0)
	for _, l := range languages {
		if l.Language.ID != 0 {
			languageResources = append(languageResources, TransformGameLanguage(l))
		}
	}

	return languageResources
}

func transformRequirements(requirements []domain.Requirement) []RequirementResource {
	requirementResources := make([]RequirementResource, 0)
	for _, r := range requirements {
		if r.ID != 0 {
			requirementResources = append(requirementResources, TransformRequirement(r))
		}
	}

	return requirementResources
}

func transformTorrents(torrents []domain.Torrent) []TorrentResource {
	torrentResources := make([]TorrentResource, 0)
	for _, t := range torrents {
		if t.ID != 0 {
			torrentResources = append(torrentResources, TransformTorrent(t))
		}
	}

	return torrentResources
}

func transformPublishers(publishers []domain.GamePublisher) []PublisherResource {
	publisherResources := make([]PublisherResource, 0)
	for _, p := range publishers {
		if p.Publisher.ID != 0 {
			publisherResources = append(publisherResources, TransformPublisher(p.Publisher))
		}
	}

	return publisherResources
}

func transformDevelopers(developers []domain.GameDeveloper) []DeveloperResource {
	developerResources := make([]DeveloperResource, 0)
	for _, d := range developers {
		if d.Developer.ID != 0 {
			developerResources = append(developerResources, TransformDeveloper(d.Developer))
		}
	}

	return developerResources
}

func transformReviews(reviews []domain.Reviewable, s3Client s3.S3ClientInterface) []ReviewResource {
	reviewResources := make([]ReviewResource, 0)
	for _, r := range reviews {
		if r.ID != 0 {
			reviewResources = append(reviewResources, TransformReview(r, s3Client))
		}
	}

	return reviewResources
}

func transformCritics(critics []domain.Criticable) []CriticableResource {
	criticResources := make([]CriticableResource, 0)
	for _, c := range critics {
		if c.ID != 0 {
			criticResources = append(criticResources, TransformCriticable(c))
		}
	}

	return criticResources
}

func transformStores(stores []domain.GameStore) []GameStoreResource {
	storeResources := make([]GameStoreResource, 0)
	for _, s := range stores {
		if s.ID != 0 {
			storeResources = append(storeResources, TransformGameStore(s))
		}
	}

	return storeResources
}

func transformComments(comments []domain.Commentable, s3Client s3.S3ClientInterface, userID uint) []CommentableResource {
	commentResources := make([]CommentableResource, 0)
	for _, c := range comments {
		if c.ID != 0 {
			commentResources = append(commentResources, TransformCommentable(c, s3Client, userID))
		}
	}

	return commentResources
}

func transformGalleries(galleries []domain.Galleriable, s3Client s3.S3ClientInterface) []GalleriableResource {
	galleryResources := make([]GalleriableResource, 0)
	for _, g := range galleries {
		if g.ID != 0 {
			galleryResources = append(galleryResources, TransformGalleriable(g, s3Client))
		}
	}

	return galleryResources
}

func transformDLCs(DLCs []domain.DLC, s3Client s3.S3ClientInterface) []DLCResource {
	DLCResources := make([]DLCResource, 0)
	for _, d := range DLCs {
		if d.ID != 0 {
			DLCResources = append(DLCResources, TransformDLC(d, s3Client))
		}
	}

	return DLCResources
}
