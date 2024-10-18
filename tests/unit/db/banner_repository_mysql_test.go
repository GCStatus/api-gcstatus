package tests

import (
	"errors"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBannerRepositoryMySQL_GetBannersForHome(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)
	mockRepo := db.NewBannerRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		mockBehavior func()
		expected     []domain.Banner
		wantErr      bool
		expectedErr  error
	}{
		"no banners": {
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `banners` WHERE component = ? AND `banners`.`deleted_at` IS NULL")).
					WithArgs(domain.HomeHeaderCarouselBannersComponent).
					WillReturnRows(sqlmock.NewRows([]string{"id", "component", "bannerable_type", "bannerable_id"}))
			},
			expected:    []domain.Banner{},
			wantErr:     false,
			expectedErr: nil,
		},
		"single game banner": {
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `banners` WHERE component = ? AND `banners`.`deleted_at` IS NULL")).
					WithArgs(domain.HomeHeaderCarouselBannersComponent).
					WillReturnRows(sqlmock.NewRows([]string{"id", "component", "bannerable_type", "bannerable_id"}).
						AddRow(1, domain.HomeHeaderCarouselBannersComponent, "games", 1))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `games` WHERE id = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "release_date"}).
						AddRow(1, "Test Game", fixedTime))

				crackRows := mock.NewRows([]string{"id", "status", "cracked_at", "cracker_id", "protection_id", "game_id"}).
					AddRow(1, "uncracked", fixedTime, 1, 1, 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `cracks` WHERE `cracks`.`game_id` = ? AND `cracks`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(crackRows)

				genreableRows := mock.NewRows([]string{"id", "genreable_id", "genreable_type", "genre_id"}).
					AddRow(1, 1, "games", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `genreables` WHERE `genreable_type` = ? AND `genreables`.`genreable_id` = ? AND `genreables`.`deleted_at` IS NULL")).
					WithArgs("games", 1).
					WillReturnRows(genreableRows)

				genresRows := mock.NewRows([]string{"id", "name"}).
					AddRow(1, "Action")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `genres` WHERE `genres`.`id` = ? AND `genres`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(genresRows)

				platformableGamesRows := mock.NewRows([]string{"id", "platformable_id", "platformable_type", "platform_id"}).
					AddRow(1, 1, "games", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platformables` WHERE `platformable_type` = ? AND `platformables`.`platformable_id` = ? AND `platformables`.`deleted_at` IS NULL")).
					WithArgs("games", 1).
					WillReturnRows(platformableGamesRows)

				platformsRows := mock.NewRows([]string{"id", "name"}).
					AddRow(1, "PC")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platforms` WHERE `platforms`.`id` = ? AND `platforms`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(platformsRows)
			},
			expected: []domain.Banner{
				{
					ID:             1,
					Component:      domain.HomeHeaderCarouselBannersComponent,
					BannerableType: "games",
					BannerableID:   1,
					Bannerable: domain.Game{
						ID:          1,
						Title:       "Test Game",
						ReleaseDate: fixedTime,
						Genres: []domain.Genreable{
							{
								ID:            1,
								GenreableID:   1,
								GenreableType: "games",
								GenreID:       1,
								Genre: domain.Genre{
									ID:   1,
									Name: "Action",
								},
							},
						},
						Platforms: []domain.Platformable{
							{
								ID:               1,
								PlatformableID:   1,
								PlatformableType: "games",
								PlatformID:       1,
								Platform: domain.Platform{
									ID:   1,
									Name: "PC",
								},
							},
						},
						Crack: &domain.Crack{
							ID:           1,
							Status:       "uncracked",
							CrackedAt:    &fixedTime,
							CrackerID:    1,
							ProtectionID: 1,
							GameID:       1,
						},
					},
				},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		"database error on banners query": {
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `banners` WHERE component = ? AND `banners`.`deleted_at` IS NULL")).
					WithArgs(domain.HomeHeaderCarouselBannersComponent).
					WillReturnError(errors.New("database error"))
			},
			expected:    nil,
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior()

			actual, err := mockRepo.GetBannersForHome()

			if tc.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
