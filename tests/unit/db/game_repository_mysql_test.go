package tests

import (
	"errors"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

				crackRows := mock.NewRows([]string{"id", "status", "cracked_at", "cracker_id", "protection_id", "game_id"}).
					AddRow(1, "uncracked", fixedTime, 1, 1, 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `cracks` WHERE `cracks`.`game_id` = ? AND `cracks`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(crackRows)

				crackerRows := mock.NewRows([]string{"id", "name", "acting"}).
					AddRow(1, "GOLDBERG", true)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `crackers` WHERE `crackers`.`id` = ? AND `crackers`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(crackerRows)

				protectionRows := mock.NewRows([]string{"id", "name"}).
					AddRow(1, "Denuvo")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `protections` WHERE `protections`.`id` = ? AND `protections`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(protectionRows)

				gameDevelopersRows := mock.NewRows([]string{"id", "developer_id", "game_id"}).
					AddRow(1, 1, 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `game_developers` WHERE `game_developers`.`game_id` = ? AND `game_developers`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(gameDevelopersRows)

				developersRows := mock.NewRows([]string{"id", "name", "acting"}).
					AddRow(1, "Game Science", true)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `developers` WHERE `developers`.`id` = ? AND `developers`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(developersRows)

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

				gamePublishersRows := mock.NewRows([]string{"id", "publisher_id", "game_id"}).
					AddRow(1, 1, 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `game_publishers` WHERE `game_publishers`.`game_id` = ? AND `game_publishers`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(gamePublishersRows)

				publishersRows := mock.NewRows([]string{"id", "name", "acting"}).
					AddRow(1, "Game Science", true)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `publishers` WHERE `publishers`.`id` = ? AND `publishers`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(publishersRows)

				requirementRows := mock.NewRows([]string{"id", "os", "dx", "cpu", "ram", "gpu", "rom", "obs", "network", "requirement_type_id", "game_id"}).
					AddRow(1, "Windows 11 64-bit", "DirectX 12", "Ryzen 5 3600", "16GB", "GeForce RTX 3090 Ti", "90GB", "Test", "Non necessary", 1, 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `requirements` WHERE `requirements`.`game_id` = ? AND `requirements`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(requirementRows)

				requirementTypeRows := mock.NewRows([]string{"id", "os", "potential"}).
					AddRow(1, "windows", "minimum")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `requirement_types` WHERE `requirement_types`.`id` = ? AND `requirement_types`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(requirementTypeRows)

				reviewablesRows := mock.NewRows([]string{"id", "rate", "review", "reviewable_id", "reviewable_type", "user_id"}).
					AddRow(1, 5, "Good game!", 1, "games", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reviewables` WHERE `reviewable_type` = ? AND `reviewables`.`reviewable_id` = ? AND `reviewables`.`deleted_at` IS NULL")).
					WithArgs("games", 1).
					WillReturnRows(reviewablesRows)

				usersRows := mock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"}).
					AddRow(1, "Fake", "fake@gmail.com", "fake", fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(usersRows)

				profilesRows := mock.NewRows([]string{"id", "share", "photo", "user_id"}).
					AddRow(1, true, "https://photo.co", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `profiles` WHERE `profiles`.`user_id` = ? AND `profiles`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(profilesRows)

				supportsRows := mock.NewRows([]string{"id", "url", "email", "contact", "game_id"}).
					AddRow(1, "https://google.com", "email@example.com", "fakeContact", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `game_supports` WHERE `game_supports`.`game_id` = ? AND `game_supports`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(supportsRows)

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

				torrentsRows := mock.NewRows([]string{"id", "url", "posted_at", "torrent_provider_id", "game_id"}).
					AddRow(1, "https://google.com", fixedTime, 1, 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `torrents` WHERE `torrents`.`game_id` = ? AND `torrents`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(torrentsRows)

				torrentProvidersRows := mock.NewRows([]string{"id", "url", "name"}).
					AddRow(1, "https://google.com", "Google")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `torrent_providers` WHERE `torrent_providers`.`id` = ? AND `torrent_providers`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(torrentProvidersRows)

				viewRows := mock.NewRows([]string{"id", "viewable_id", "viewable_type", "count"}).
					AddRow(1, 1, "games", 10)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `viewables` WHERE `viewable_type` = ? AND `viewables`.`viewable_id` = ? AND `viewables`.`deleted_at` IS NULL")).
					WithArgs("games", 1).
					WillReturnRows(viewRows)

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `viewables` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`count`=?,`viewable_id`=?,`viewable_type`=? WHERE `viewables`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 11, 1, "games", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
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
