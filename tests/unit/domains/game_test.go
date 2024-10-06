package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		game         domain.Game
		mockBehavior func(mock sqlmock.Sqlmock, game domain.Game)
		expectError  bool
	}{
		"Success": {
			game: domain.Game{
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
			mockBehavior: func(mock sqlmock.Sqlmock, game domain.Game) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `games`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						game.Age,
						game.Slug,
						game.Title,
						game.Condition,
						game.Cover,
						game.About,
						game.Description,
						game.ShortDescription,
						game.Free,
						game.Legal,
						game.Website,
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			game: domain.Game{
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
			mockBehavior: func(mock sqlmock.Sqlmock, game domain.Game) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `games`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						game.Age,
						game.Slug,
						game.Title,
						game.Condition,
						game.Cover,
						game.About,
						game.Description,
						game.ShortDescription,
						game.Free,
						game.Legal,
						game.Website,
						sqlmock.AnyArg(),
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.game)

			err := db.Create(&tc.game).Error

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateGame(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		game         domain.Game
		mockBehavior func(mock sqlmock.Sqlmock, game domain.Game)
		expectError  bool
	}{
		"Success": {
			game: domain.Game{
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
			mockBehavior: func(mock sqlmock.Sqlmock, game domain.Game) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `games`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						game.Age,
						game.Slug,
						game.Title,
						game.Condition,
						game.Cover,
						game.About,
						game.Description,
						game.ShortDescription,
						game.Free,
						game.Legal,
						game.Website,
						sqlmock.AnyArg(),
						game.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			game: domain.Game{
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
			mockBehavior: func(mock sqlmock.Sqlmock, game domain.Game) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `games`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						game.Age,
						game.Slug,
						game.Title,
						game.Condition,
						game.Cover,
						game.About,
						game.Description,
						game.ShortDescription,
						game.Free,
						game.Legal,
						game.Website,
						sqlmock.AnyArg(),
						game.ID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.game)

			err := db.Save(&tc.game).Error

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSoftDeleteGame(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		gameID       uint
		mockBehavior func(mock sqlmock.Sqlmock, gameID uint)
		wantErr      bool
	}{
		"Can soft delete a Game": {
			gameID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, gameID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `games` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), gameID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			gameID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, gameID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `games` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Game"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.gameID)

			err := db.Delete(&domain.Game{}, tc.gameID).Error

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetGameBySlug(t *testing.T) {
	fixedTime := time.Now()
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		gameSlug  string
		mockFunc  func()
		wantGame  domain.Game
		wantError bool
	}{
		"Valid game fetch": {
			gameSlug: "valid",
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
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "slug", "age", "cover", "about", "description", "short_description", "free", "release_date", "condition", "created_at", "updated_at"}).
					AddRow(1, "Game Test", "valid", 18, "https://placehold.co/600x400/EEE/31343C", "About game", "Description", "Short description", false, fixedTime, domain.CommomCondition, fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `games` WHERE slug = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT ?")).
					WithArgs("valid", 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Game not found": {
			gameSlug:  "invalid",
			wantGame:  domain.Game{},
			wantError: true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `games` WHERE slug = ? AND `games`.`deleted_at` IS NULL ORDER BY `games`.`id` LIMIT ?")).
					WithArgs("invalid", 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var game domain.Game
			err := db.Where("slug = ?", tc.gameSlug).First(&game).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantGame, game)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateGameValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		game domain.Game
	}{
		"Can empty validations errors": {
			game: domain.Game{
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
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.game.ValidateGame()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGameWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		game    domain.Game
		wantErr string
	}{
		"Missing required fields": {
			game:    domain.Game{},
			wantErr: "Age is a required field, Slug is a required field, Title is a required field, Condition is a required field, Cover is a required field, About is a required field, Description is a required field, ShortDescription is a required field, ReleaseDate is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.game.ValidateGame()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
