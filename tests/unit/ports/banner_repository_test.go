package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockBannerRepository struct {
	banners map[uint]*domain.Banner
}

func NewMockBannerRepository() *MockBannerRepository {
	return &MockBannerRepository{
		banners: make(map[uint]*domain.Banner),
	}
}

func (m *MockBannerRepository) CreateBanner(banner *domain.Banner) error {
	if banner == nil {
		return errors.New("invalid banner data")
	}
	m.banners[banner.ID] = banner
	return nil
}

func (m *MockBannerRepository) GetBannersForHome() ([]domain.Banner, error) {
	var banners []domain.Banner
	for _, banner := range m.banners {
		banners = append(banners, *banner)
	}
	return banners, nil
}

func TestMockBannerRepository_GetBannersForHome(t *testing.T) {
	fixedTime := time.Now()
	mockRepo := NewMockBannerRepository()

	testCases := map[string]struct {
		setupFunc   func()
		expected    []domain.Banner
		expectedErr error
	}{
		"no banners": {
			setupFunc:   func() {},
			expected:    []domain.Banner{},
			expectedErr: nil,
		},
		"single game banner": {
			setupFunc: func() {
				if err := mockRepo.CreateBanner(&domain.Banner{
					ID:             1,
					Component:      domain.HomeHeaderCarouselBannersComponent,
					BannerableType: "games",
					BannerableID:   1,
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
				}); err != nil {
					t.Fatalf("Failed to create banner: %v", err)
				}
			},
			expected: []domain.Banner{
				{
					ID:             1,
					Component:      domain.HomeHeaderCarouselBannersComponent,
					BannerableType: "games",
					BannerableID:   1,
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
			},
			expectedErr: nil,
		},
		"multiple banners with mixed types": {
			setupFunc: func() {
				if err := mockRepo.CreateBanner(&domain.Banner{
					ID:             1,
					Component:      domain.HomeHeaderCarouselBannersComponent,
					BannerableType: "games",
					BannerableID:   1,
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
				}); err != nil {
					t.Fatalf("Failed to create banner: %v", err)
				}
				if err := mockRepo.CreateBanner(&domain.Banner{
					ID:             2,
					Component:      domain.HomeHeaderCarouselBannersComponent,
					BannerableType: "categories",
					BannerableID:   20,
				}); err != nil {
					t.Fatalf("Failed to create banner: %v", err)
				}
			},
			expected: []domain.Banner{
				{
					ID:             1,
					Component:      domain.HomeHeaderCarouselBannersComponent,
					BannerableType: "games",
					BannerableID:   1,
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
					Component:      domain.HomeHeaderCarouselBannersComponent,
					BannerableType: "categories",
					BannerableID:   20,
					Bannerable:     nil,
				},
			},
			expectedErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.setupFunc()

			actual, err := mockRepo.GetBannersForHome()
			sort.Slice(actual, func(i, j int) bool {
				return actual[i].ID < actual[j].ID
			})
			sort.Slice(tc.expected, func(i, j int) bool {
				return tc.expected[i].ID < tc.expected[j].ID
			})

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if actual == nil {
				actual = []domain.Banner{}
			}

			assert.Equal(t, tc.expected, actual)
			mockRepo.banners = make(map[uint]*domain.Banner)
		})
	}
}
