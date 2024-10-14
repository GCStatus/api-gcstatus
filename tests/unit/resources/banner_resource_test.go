package tests

import (
	"testing"
	"time"

	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"

	"github.com/stretchr/testify/assert"
)

func TestTransformBanner(t *testing.T) {
	fixedTime := time.Now()
	var mockS3Client s3.S3ClientInterface

	testCases := map[string]struct {
		inputBanner domain.Banner
		userID      uint
		expected    resources.BannerResource
	}{
		"game banner": {
			inputBanner: domain.Banner{
				ID:             1,
				BannerableType: "games",
				Bannerable: domain.Game{
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
					Platforms:        []domain.Platformable{},
					Genres:           []domain.Genreable{},
				},
			},
			userID: 1,
			expected: resources.BannerResource{
				ID:             1,
				BannerableType: "games",
				Game: &resources.GameResource{
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
					HeartsCount:      0,
					Legal:            utils.StringPtr("Some legal info"),
					Website:          utils.StringPtr("http://testgame.com"),
					IsHearted:        false,
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
					Critics:          []resources.CriticableResource{},
					Stores:           []resources.GameStoreResource{},
					Comments:         []resources.CommentableResource{},
					Galleries:        []resources.GalleriableResource{},
					DLCs:             []resources.DLCResource{},
				},
			},
		},
		"non-game banner": {
			inputBanner: domain.Banner{
				ID:             2,
				BannerableType: "categories",
				Bannerable:     nil,
			},
			userID: 1,
			expected: resources.BannerResource{
				ID:             2,
				BannerableType: "categories",
				Game:           nil,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := resources.TransformBanner(tc.inputBanner, mockS3Client, tc.userID)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestTransformBanners(t *testing.T) {
	fixedTime := time.Now()
	var mockS3Client s3.S3ClientInterface

	testCases := map[string]struct {
		inputBanners []domain.Banner
		userID       uint
		expected     []resources.BannerResource
	}{
		"multiple banners with games": {
			inputBanners: []domain.Banner{
				{
					ID:             1,
					BannerableType: "games",
					Bannerable: domain.Game{
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
						Platforms:        []domain.Platformable{},
						Genres:           []domain.Genreable{},
					},
				},
				{
					ID:             2,
					BannerableType: "games",
					Bannerable: domain.Game{
						ID:               2,
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
						Platforms:        []domain.Platformable{},
						Genres:           []domain.Genreable{},
					},
				},
			},
			userID: 1,
			expected: []resources.BannerResource{
				{
					ID:             1,
					BannerableType: "games",
					Game: &resources.GameResource{
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
						HeartsCount:      0,
						Legal:            utils.StringPtr("Some legal info"),
						Website:          utils.StringPtr("http://testgame.com"),
						IsHearted:        false,
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
						Critics:          []resources.CriticableResource{},
						Stores:           []resources.GameStoreResource{},
						Comments:         []resources.CommentableResource{},
						Galleries:        []resources.GalleriableResource{},
						DLCs:             []resources.DLCResource{},
					},
				},
				{
					ID:             2,
					BannerableType: "games",
					Game: &resources.GameResource{
						ID:               2,
						Age:              16,
						Slug:             "test-game",
						Title:            "Test Game",
						Condition:        "New",
						Cover:            "test-cover.jpg",
						About:            "About Test Game",
						Description:      "Detailed description of Test Game",
						ShortDescription: "Short description",
						Free:             true,
						HeartsCount:      0,
						Legal:            utils.StringPtr("Some legal info"),
						Website:          utils.StringPtr("http://testgame.com"),
						IsHearted:        false,
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
						Critics:          []resources.CriticableResource{},
						Stores:           []resources.GameStoreResource{},
						Comments:         []resources.CommentableResource{},
						Galleries:        []resources.GalleriableResource{},
						DLCs:             []resources.DLCResource{},
					},
				},
			},
		},
		"banners with mixed types": {
			inputBanners: []domain.Banner{
				{
					ID:             1,
					BannerableType: "games",
					Bannerable: domain.Game{
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
						Platforms:        []domain.Platformable{},
						Genres:           []domain.Genreable{},
					},
				},
				{
					ID:             2,
					BannerableType: "categories",
					Bannerable:     nil,
				},
			},
			userID: 1,
			expected: []resources.BannerResource{
				{
					ID:             1,
					BannerableType: "games",
					Game: &resources.GameResource{
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
						HeartsCount:      0,
						Legal:            utils.StringPtr("Some legal info"),
						Website:          utils.StringPtr("http://testgame.com"),
						IsHearted:        false,
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
						Critics:          []resources.CriticableResource{},
						Stores:           []resources.GameStoreResource{},
						Comments:         []resources.CommentableResource{},
						Galleries:        []resources.GalleriableResource{},
						DLCs:             []resources.DLCResource{},
					},
				},
				{
					ID:             2,
					BannerableType: "categories",
					Game:           nil,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := resources.TransformBanners(tc.inputBanners, mockS3Client, tc.userID)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
