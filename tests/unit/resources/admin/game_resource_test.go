package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransformGame(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Game
		expected resources_admin.GameResource
	}{
		"not liked game": {
			input: domain.Game{
				ID:               1,
				Age:              16,
				Slug:             "test-game",
				Title:            "Test Game",
				Condition:        "New",
				Cover:            "test-cover.jpg",
				About:            "About Test Game",
				Description:      "Detailed description of Test Game",
				ShortDescription: "Short description",
				Free:             true,
				Legal:            utils.StringPtr("Some legal info"),
				Website:          utils.StringPtr("http://testgame.com"),
				ReleaseDate:      fixedTime,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
				Categories:       []domain.Categoriable{},
				Platforms:        []domain.Platformable{},
				Genres:           []domain.Genreable{},
				Tags:             []domain.Taggable{},
				Hearts: []domain.Heartable{
					{
						ID:            1,
						HeartableID:   1,
						HeartableType: "games",
						UserID:        2,
					},
				},
			},
			expected: resources_admin.GameResource{
				ID:               1,
				Age:              16,
				Slug:             "test-game",
				Title:            "Test Game",
				Condition:        "New",
				Cover:            "test-cover.jpg",
				About:            "About Test Game",
				Description:      "Detailed description of Test Game",
				ShortDescription: "Short description",
				Free:             true,
				HeartsCount:      1,
				Legal:            utils.StringPtr("Some legal info"),
				Website:          utils.StringPtr("http://testgame.com"),
				ReleaseDate:      utils.FormatTimestamp(fixedTime),
				CreatedAt:        utils.FormatTimestamp(fixedTime),
				UpdatedAt:        utils.FormatTimestamp(fixedTime),
				Categories:       []resources_admin.CategoryResource{},
				Platforms:        []resources_admin.PlatformResource{},
				Genres:           []resources_admin.GenreResource{},
				Tags:             []resources_admin.TagResource{},
				Languages:        []resources_admin.GameLanguageResource{},
				Requirements:     []resources_admin.RequirementResource{},
				Torrents:         []resources_admin.TorrentResource{},
				Publishers:       []resources_admin.PublisherResource{},
				Developers:       []resources_admin.DeveloperResource{},
				Reviews:          []resources_admin.ReviewResource{},
				Critics:          []resources_admin.CriticableResource{},
				Stores:           []resources_admin.GameStoreResource{},
				Comments:         []resources_admin.CommentableResource{},
				Galleries:        []resources_admin.GalleriableResource{},
				DLCs:             []resources_admin.DLCResource{},
			},
		},
		"No Morph Relations": {
			input: domain.Game{
				ID:               1,
				Age:              16,
				Slug:             "test-game",
				Title:            "Test Game",
				Condition:        "New",
				Cover:            "test-cover.jpg",
				About:            "About Test Game",
				Description:      "Detailed description of Test Game",
				ShortDescription: "Short description",
				Free:             true,
				Legal:            utils.StringPtr("Some legal info"),
				Website:          utils.StringPtr("http://testgame.com"),
				ReleaseDate:      fixedTime,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
				Categories:       []domain.Categoriable{},
				Platforms:        []domain.Platformable{},
				Genres:           []domain.Genreable{},
				Tags:             []domain.Taggable{},
			},
			expected: resources_admin.GameResource{
				ID:               1,
				Age:              16,
				Slug:             "test-game",
				Title:            "Test Game",
				Condition:        "New",
				Cover:            "test-cover.jpg",
				About:            "About Test Game",
				HeartsCount:      0,
				Description:      "Detailed description of Test Game",
				ShortDescription: "Short description",
				Free:             true,
				Legal:            utils.StringPtr("Some legal info"),
				Website:          utils.StringPtr("http://testgame.com"),
				ReleaseDate:      utils.FormatTimestamp(fixedTime),
				CreatedAt:        utils.FormatTimestamp(fixedTime),
				UpdatedAt:        utils.FormatTimestamp(fixedTime),
				Categories:       []resources_admin.CategoryResource{},
				Platforms:        []resources_admin.PlatformResource{},
				Genres:           []resources_admin.GenreResource{},
				Tags:             []resources_admin.TagResource{},
				Languages:        []resources_admin.GameLanguageResource{},
				Requirements:     []resources_admin.RequirementResource{},
				Torrents:         []resources_admin.TorrentResource{},
				Publishers:       []resources_admin.PublisherResource{},
				Developers:       []resources_admin.DeveloperResource{},
				Reviews:          []resources_admin.ReviewResource{},
				Critics:          []resources_admin.CriticableResource{},
				Stores:           []resources_admin.GameStoreResource{},
				Comments:         []resources_admin.CommentableResource{},
				Galleries:        []resources_admin.GalleriableResource{},
				DLCs:             []resources_admin.DLCResource{},
			},
		},
		"With One Category": {
			input: domain.Game{
				ID:          2,
				Age:         18,
				Slug:        "fps-game",
				Title:       "FPS Game",
				Condition:   "Used",
				Cover:       "fps-cover.jpg",
				About:       "About FPS Game",
				Description: "Detailed description of FPS Game",
				Free:        false,
				Legal:       nil,
				Website:     nil,
				ReleaseDate: fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Categories: []domain.Categoriable{
					{
						Category: domain.Category{
							ID:        1,
							Name:      "FPS",
							CreatedAt: fixedTime,
							UpdatedAt: fixedTime,
						},
					},
				},
				Platforms: []domain.Platformable{},
				Genres:    []domain.Genreable{},
				Tags:      []domain.Taggable{},
			},
			expected: resources_admin.GameResource{
				ID:               2,
				Age:              18,
				Slug:             "fps-game",
				Title:            "FPS Game",
				Condition:        "Used",
				Cover:            "fps-cover.jpg",
				About:            "About FPS Game",
				Description:      "Detailed description of FPS Game",
				ShortDescription: "",
				Free:             false,
				Legal:            nil,
				Website:          nil,
				ReleaseDate:      utils.FormatTimestamp(fixedTime),
				CreatedAt:        utils.FormatTimestamp(fixedTime),
				UpdatedAt:        utils.FormatTimestamp(fixedTime),
				Categories: []resources_admin.CategoryResource{
					{
						ID:        1,
						Name:      "FPS",
						CreatedAt: utils.FormatTimestamp(fixedTime),
						UpdatedAt: utils.FormatTimestamp(fixedTime),
					},
				},
				Platforms:    []resources_admin.PlatformResource{},
				Genres:       []resources_admin.GenreResource{},
				Tags:         []resources_admin.TagResource{},
				Languages:    []resources_admin.GameLanguageResource{},
				Requirements: []resources_admin.RequirementResource{},
				Torrents:     []resources_admin.TorrentResource{},
				Publishers:   []resources_admin.PublisherResource{},
				Developers:   []resources_admin.DeveloperResource{},
				Reviews:      []resources_admin.ReviewResource{},
				Critics:      []resources_admin.CriticableResource{},
				Stores:       []resources_admin.GameStoreResource{},
				Comments:     []resources_admin.CommentableResource{},
				Galleries:    []resources_admin.GalleriableResource{},
				DLCs:         []resources_admin.DLCResource{},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockS3Client := &testutils.MockS3Client{}
			result := resources_admin.TransformGame(tc.input, mockS3Client)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTransformGames(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Game
		expected []resources_admin.GameResource
	}{
		"Empty Game List": {
			input:    []domain.Game{},
			expected: []resources_admin.GameResource{},
		},
		"Single Game With No Morph Relations": {
			input: []domain.Game{
				{
					ID:               1,
					Age:              16,
					Slug:             "test-game",
					Title:            "Test Game",
					Condition:        "New",
					Cover:            "test-cover.jpg",
					About:            "About Test Game",
					Description:      "Detailed description of Test Game",
					ShortDescription: "Short description",
					Free:             true,
					Legal:            utils.StringPtr("Some legal info"),
					Website:          utils.StringPtr("http://testgame.com"),
					ReleaseDate:      fixedTime,
					CreatedAt:        fixedTime,
					UpdatedAt:        fixedTime,
					Categories:       []domain.Categoriable{},
					Platforms:        []domain.Platformable{},
					Genres:           []domain.Genreable{},
					Tags:             []domain.Taggable{},
				},
			},
			expected: []resources_admin.GameResource{
				{
					ID:               1,
					Age:              16,
					Slug:             "test-game",
					Title:            "Test Game",
					Condition:        "New",
					Cover:            "test-cover.jpg",
					About:            "About Test Game",
					Description:      "Detailed description of Test Game",
					ShortDescription: "Short description",
					Free:             true,
					Legal:            utils.StringPtr("Some legal info"),
					Website:          utils.StringPtr("http://testgame.com"),
					ReleaseDate:      utils.FormatTimestamp(fixedTime),
					CreatedAt:        utils.FormatTimestamp(fixedTime),
					UpdatedAt:        utils.FormatTimestamp(fixedTime),
					Categories:       []resources_admin.CategoryResource{},
					Platforms:        []resources_admin.PlatformResource{},
					Genres:           []resources_admin.GenreResource{},
					Tags:             []resources_admin.TagResource{},
					Languages:        []resources_admin.GameLanguageResource{},
					Requirements:     []resources_admin.RequirementResource{},
					Torrents:         []resources_admin.TorrentResource{},
					Publishers:       []resources_admin.PublisherResource{},
					Developers:       []resources_admin.DeveloperResource{},
					Reviews:          []resources_admin.ReviewResource{},
					Critics:          []resources_admin.CriticableResource{},
					Stores:           []resources_admin.GameStoreResource{},
					Comments:         []resources_admin.CommentableResource{},
					Galleries:        []resources_admin.GalleriableResource{},
					DLCs:             []resources_admin.DLCResource{},
				},
			},
		},
		"Multiple Games With Mixed Morph Relations": {
			input: []domain.Game{
				{
					ID:          2,
					Age:         18,
					Slug:        "fps-game",
					Title:       "FPS Game",
					Condition:   "Used",
					Cover:       "fps-cover.jpg",
					About:       "About FPS Game",
					Description: "Detailed description of FPS Game",
					Free:        false,
					ReleaseDate: fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
					Categories: []domain.Categoriable{
						{
							Category: domain.Category{
								ID:        1,
								Name:      "FPS",
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
							},
						},
					},
					Platforms: []domain.Platformable{
						{
							Platform: domain.Platform{
								ID:        1,
								Name:      "PC",
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
							},
						},
					},
					Genres: []domain.Genreable{},
					Tags:   []domain.Taggable{},
				},
				{
					ID:          3,
					Age:         13,
					Slug:        "rpg-game",
					Title:       "RPG Game",
					Condition:   "New",
					Cover:       "rpg-cover.jpg",
					About:       "About RPG Game",
					Description: "Detailed description of RPG Game",
					Free:        true,
					ReleaseDate: fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
					Categories:  []domain.Categoriable{},
					Platforms:   []domain.Platformable{},
					Genres: []domain.Genreable{
						{
							Genre: domain.Genre{
								ID:        2,
								Name:      "Fantasy",
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
							},
						},
					},
					Tags: []domain.Taggable{
						{
							Tag: domain.Tag{
								ID:        1,
								Name:      "Adventure",
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
							},
						},
					},
				},
			},
			expected: []resources_admin.GameResource{
				{
					ID:          2,
					Age:         18,
					Slug:        "fps-game",
					Title:       "FPS Game",
					Condition:   "Used",
					Cover:       "fps-cover.jpg",
					About:       "About FPS Game",
					Description: "Detailed description of FPS Game",
					Free:        false,
					ReleaseDate: utils.FormatTimestamp(fixedTime),
					CreatedAt:   utils.FormatTimestamp(fixedTime),
					UpdatedAt:   utils.FormatTimestamp(fixedTime),
					Categories: []resources_admin.CategoryResource{
						{
							ID:        1,
							Name:      "FPS",
							CreatedAt: utils.FormatTimestamp(fixedTime),
							UpdatedAt: utils.FormatTimestamp(fixedTime),
						},
					},
					Platforms: []resources_admin.PlatformResource{
						{
							ID:        1,
							Name:      "PC",
							CreatedAt: utils.FormatTimestamp(fixedTime),
							UpdatedAt: utils.FormatTimestamp(fixedTime),
						},
					},
					Genres:       []resources_admin.GenreResource{},
					Tags:         []resources_admin.TagResource{},
					Languages:    []resources_admin.GameLanguageResource{},
					Requirements: []resources_admin.RequirementResource{},
					Torrents:     []resources_admin.TorrentResource{},
					Publishers:   []resources_admin.PublisherResource{},
					Developers:   []resources_admin.DeveloperResource{},
					Reviews:      []resources_admin.ReviewResource{},
					Critics:      []resources_admin.CriticableResource{},
					Stores:       []resources_admin.GameStoreResource{},
					Comments:     []resources_admin.CommentableResource{},
					Galleries:    []resources_admin.GalleriableResource{},
					DLCs:         []resources_admin.DLCResource{},
				},
				{
					ID:          3,
					Age:         13,
					Slug:        "rpg-game",
					Title:       "RPG Game",
					Condition:   "New",
					Cover:       "rpg-cover.jpg",
					About:       "About RPG Game",
					Description: "Detailed description of RPG Game",
					Free:        true,
					ReleaseDate: utils.FormatTimestamp(fixedTime),
					CreatedAt:   utils.FormatTimestamp(fixedTime),
					UpdatedAt:   utils.FormatTimestamp(fixedTime),
					Categories:  []resources_admin.CategoryResource{},
					Platforms:   []resources_admin.PlatformResource{},
					Genres: []resources_admin.GenreResource{
						{
							ID:        2,
							Name:      "Fantasy",
							CreatedAt: utils.FormatTimestamp(fixedTime),
							UpdatedAt: utils.FormatTimestamp(fixedTime),
						},
					},
					Tags: []resources_admin.TagResource{
						{
							ID:        1,
							Name:      "Adventure",
							CreatedAt: utils.FormatTimestamp(fixedTime),
							UpdatedAt: utils.FormatTimestamp(fixedTime),
						},
					},
					Languages:    []resources_admin.GameLanguageResource{},
					Requirements: []resources_admin.RequirementResource{},
					Torrents:     []resources_admin.TorrentResource{},
					Publishers:   []resources_admin.PublisherResource{},
					Developers:   []resources_admin.DeveloperResource{},
					Reviews:      []resources_admin.ReviewResource{},
					Critics:      []resources_admin.CriticableResource{},
					Stores:       []resources_admin.GameStoreResource{},
					Comments:     []resources_admin.CommentableResource{},
					Galleries:    []resources_admin.GalleriableResource{},
					DLCs:         []resources_admin.DLCResource{},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockS3Client := &testutils.MockS3Client{}
			result := resources_admin.TransformGames(tc.input, mockS3Client)
			assert.Equal(t, tc.expected, result)
		})
	}
}
