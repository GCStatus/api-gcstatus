package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransformGame(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Game
		expected resources.GameResource
	}{
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
			expected: resources.GameResource{
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
				Categories:       []resources.CategoryResource{},
				Platforms:        []resources.PlatformResource{},
				Genres:           []resources.GenreResource{},
				Tags:             []resources.TagResource{},
				Languages:        []resources.GameLanguageResource{},
				Requirements:     []resources.RequirementResource{},
				Torrents:         []resources.TorrentResource{},
				Publishers:       []resources.PublisherResource{},
				Developers:       []resources.DeveloperResource{},
				Reviews:          []resources.ReviewResource{},
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
				ReleaseDate: time.Date(2022, 5, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Categories: []domain.Categoriable{
					{
						Category: domain.Category{
							ID:   1,
							Name: "FPS",
						},
					},
				},
				Platforms: []domain.Platformable{},
				Genres:    []domain.Genreable{},
				Tags:      []domain.Taggable{},
			},
			expected: resources.GameResource{
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
				ReleaseDate:      utils.FormatTimestamp(time.Date(2022, 5, 15, 0, 0, 0, 0, time.UTC)),
				CreatedAt:        utils.FormatTimestamp(fixedTime),
				UpdatedAt:        utils.FormatTimestamp(fixedTime),
				Categories: []resources.CategoryResource{
					{ID: 1, Name: "FPS"},
				},
				Platforms:    []resources.PlatformResource{},
				Genres:       []resources.GenreResource{},
				Tags:         []resources.TagResource{},
				Languages:    []resources.GameLanguageResource{},
				Requirements: []resources.RequirementResource{},
				Torrents:     []resources.TorrentResource{},
				Publishers:   []resources.PublisherResource{},
				Developers:   []resources.DeveloperResource{},
				Reviews:      []resources.ReviewResource{},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockS3Client := &MockS3Client{}
			result := resources.TransformGame(tc.input, mockS3Client)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTransformGames(t *testing.T) {
	testCases := map[string]struct {
		input    []domain.Game
		expected []resources.GameResource
	}{
		"Empty Game List": {
			input:    []domain.Game{},
			expected: []resources.GameResource{},
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
					ReleaseDate:      time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC),
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
					Categories:       []domain.Categoriable{},
					Platforms:        []domain.Platformable{},
					Genres:           []domain.Genreable{},
					Tags:             []domain.Taggable{},
				},
			},
			expected: []resources.GameResource{
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
					ReleaseDate:      utils.FormatTimestamp(time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC)),
					CreatedAt:        utils.FormatTimestamp(time.Now()),
					UpdatedAt:        utils.FormatTimestamp(time.Now()),
					Categories:       []resources.CategoryResource{},
					Platforms:        []resources.PlatformResource{},
					Genres:           []resources.GenreResource{},
					Tags:             []resources.TagResource{},
					Languages:        []resources.GameLanguageResource{},
					Requirements:     []resources.RequirementResource{},
					Torrents:         []resources.TorrentResource{},
					Publishers:       []resources.PublisherResource{},
					Developers:       []resources.DeveloperResource{},
					Reviews:          []resources.ReviewResource{},
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
					ReleaseDate: time.Date(2022, 5, 15, 0, 0, 0, 0, time.UTC),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Categories: []domain.Categoriable{
						{
							Category: domain.Category{
								ID:   1,
								Name: "FPS",
							},
						},
					},
					Platforms: []domain.Platformable{
						{
							Platform: domain.Platform{
								ID:   1,
								Name: "PC",
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
					ReleaseDate: time.Date(2021, 7, 21, 0, 0, 0, 0, time.UTC),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Categories:  []domain.Categoriable{},
					Platforms:   []domain.Platformable{},
					Genres: []domain.Genreable{
						{
							Genre: domain.Genre{
								ID:   2,
								Name: "Fantasy",
							},
						},
					},
					Tags: []domain.Taggable{
						{
							Tag: domain.Tag{
								ID:   1,
								Name: "Adventure",
							},
						},
					},
				},
			},
			expected: []resources.GameResource{
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
					ReleaseDate: utils.FormatTimestamp(time.Date(2022, 5, 15, 0, 0, 0, 0, time.UTC)),
					CreatedAt:   utils.FormatTimestamp(time.Now()),
					UpdatedAt:   utils.FormatTimestamp(time.Now()),
					Categories: []resources.CategoryResource{
						{ID: 1, Name: "FPS"},
					},
					Platforms: []resources.PlatformResource{
						{ID: 1, Name: "PC"},
					},
					Genres:       []resources.GenreResource{},
					Tags:         []resources.TagResource{},
					Languages:    []resources.GameLanguageResource{},
					Requirements: []resources.RequirementResource{},
					Torrents:     []resources.TorrentResource{},
					Publishers:   []resources.PublisherResource{},
					Developers:   []resources.DeveloperResource{},
					Reviews:      []resources.ReviewResource{},
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
					ReleaseDate: utils.FormatTimestamp(time.Date(2021, 7, 21, 0, 0, 0, 0, time.UTC)),
					CreatedAt:   utils.FormatTimestamp(time.Now()),
					UpdatedAt:   utils.FormatTimestamp(time.Now()),
					Categories:  []resources.CategoryResource{},
					Platforms:   []resources.PlatformResource{},
					Genres: []resources.GenreResource{
						{ID: 2, Name: "Fantasy"},
					},
					Tags: []resources.TagResource{
						{ID: 1, Name: "Adventure"},
					},
					Languages:    []resources.GameLanguageResource{},
					Requirements: []resources.RequirementResource{},
					Torrents:     []resources.TorrentResource{},
					Publishers:   []resources.PublisherResource{},
					Developers:   []resources.DeveloperResource{},
					Reviews:      []resources.ReviewResource{},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockS3Client := &MockS3Client{}
			result := resources.TransformGames(tc.input, mockS3Client)
			assert.Equal(t, tc.expected, result)
		})
	}
}
