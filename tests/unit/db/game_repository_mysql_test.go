package tests

import (
	"errors"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGameRepositoryMySQL_FindBySlug(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := tests.Setup(t)
	mockRepo := db.NewGameRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		slug         string
		wantErr      bool
		expectedErr  error
		wantGame     domain.Game
		mockBehavior func(slug string)
	}{
		"game found": {
			slug:    "valid",
			wantErr: false,
			wantGame: domain.Game{
				ID:               1,
				Slug:             "valid",
				Age:              18,
				Title:            "Game Test",
				Condition:        domain.CommomCondition,
				Cover:            "https://placehold.co/600x400/EEE/31343C",
				About:            "About game",
				Description:      "Description",
				ShortDescription: "Short description",
				Free:             false,
				ReleaseDate:      fixedTime,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
			},
			mockBehavior: func(slug string) {
				rows := mock.NewRows([]string{"id", "age", "slug", "title", "condition", "cover", "about", "description", "short_description", "free", "release_date", "created_at", "updated_at"}).
					AddRow(1, 18, "valid", "Game Test", domain.CommomCondition, "https://placehold.co/600x400/EEE/31343C", "About game", "Description", "Short description", false, fixedTime, fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `games` WHERE slug = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT ?")).
					WithArgs(slug, 1).
					WillReturnRows(rows)

				categoriableRows := mock.NewRows([]string{"id", "categoriable_id", "categoriable_type", "category_id"}).
					AddRow(1, 1, "games", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categoriables` WHERE `categoriable_type` = ? AND `categoriables`.`categoriable_id` = ? AND `categoriables`.`deleted_at` IS NULL")).
					WithArgs("games", 1).
					WillReturnRows(categoriableRows)

				categoriesRows := mock.NewRows([]string{"id", "name"}).
					AddRow(1, "FPS")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(categoriesRows)

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

				gameLanguageRows := mock.NewRows([]string{"id", "menu", "dubs", "subtitles", "game_id", "language_id"}).
					AddRow(1, false, true, false, 1, 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `game_languages` WHERE `game_languages`.`game_id` = ? AND `game_languages`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(gameLanguageRows)

				languageRows := mock.NewRows([]string{"id", "name", "iso"}).
					AddRow(1, "Portuguese", "pt_BR")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `languages` WHERE `languages`.`id` = ? AND `languages`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(languageRows)

				platformableRows := mock.NewRows([]string{"id", "platformable_id", "platformable_type", "platform_id"}).
					AddRow(1, 1, "games", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platformables` WHERE `platformable_type` = ? AND `platformables`.`platformable_id` = ? AND `platformables`.`deleted_at` IS NULL")).
					WithArgs("games", 1).
					WillReturnRows(platformableRows)

				platformsRows := mock.NewRows([]string{"id", "name"}).
					AddRow(1, "PC")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platforms` WHERE `platforms`.`id` = ? AND `platforms`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(platformsRows)

				taggablesRows := mock.NewRows([]string{"id", "taggable_id", "taggable_type", "tag_id"}).
					AddRow(1, 1, "games", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `taggables` WHERE `taggable_type` = ? AND `taggables`.`taggable_id` = ? AND `taggables`.`deleted_at` IS NULL")).
					WithArgs("games", 1).
					WillReturnRows(taggablesRows)

				tagsRows := mock.NewRows([]string{"id", "name"}).
					AddRow(1, "Adventure")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tags` WHERE `tags`.`id` = ? AND `tags`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(tagsRows)
			},
		},
		"game not found": {
			slug:        "invalid",
			wantErr:     true,
			expectedErr: errors.New("record not found"),
			wantGame:    domain.Game{},
			mockBehavior: func(slug string) {
				rows := mock.NewRows([]string{"id", "age", "slug", "title", "condition", "cover", "about", "description", "short_description", "free", "release_date", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `games` WHERE slug = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT ?")).
					WithArgs(slug, 1).
					WillReturnRows(rows)
			},
		},
		"db error": {
			slug:        "valid",
			wantErr:     true,
			expectedErr: errors.New("db error"),
			wantGame:    domain.Game{},
			mockBehavior: func(slug string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `games` WHERE slug = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT ?")).
					WithArgs(slug, 1).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(tc.slug)

			game, err := mockRepo.FindBySlug(tc.slug)

			assert.Equal(t, tc.expectedErr, err)
			if err == gorm.ErrRecordNotFound {
				assert.Equal(t, uint(0), game.ID)
			} else {
				assert.Equal(t, tc.wantGame.ID, game.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
