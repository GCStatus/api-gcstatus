package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGameDeveloper(t *testing.T) {
	testCases := map[string]struct {
		gameDeveloper domain.GameDeveloper
		mockBehavior  func(mock sqlmock.Sqlmock, gameDeveloper domain.GameDeveloper)
		expectError   bool
	}{
		"Success": {
			gameDeveloper: domain.GameDeveloper{
				GameID:      1,
				DeveloperID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameDeveloper domain.GameDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameDeveloper.GameID,
						gameDeveloper.DeveloperID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			gameDeveloper: domain.GameDeveloper{
				GameID:      1,
				DeveloperID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameDeveloper domain.GameDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameDeveloper.GameID,
						gameDeveloper.DeveloperID,
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

			tc.mockBehavior(mock, tc.gameDeveloper)

			err := db.Create(&tc.gameDeveloper).Error

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

func TestUpdateGameDeveloper(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameDeveloper domain.GameDeveloper
		mockBehavior  func(mock sqlmock.Sqlmock, gameDeveloper domain.GameDeveloper)
		expectError   bool
	}{
		"Success": {
			gameDeveloper: domain.GameDeveloper{
				ID:          1,
				GameID:      1,
				DeveloperID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameDeveloper domain.GameDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameDeveloper.GameID,
						gameDeveloper.DeveloperID,
						gameDeveloper.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			gameDeveloper: domain.GameDeveloper{
				ID:          1,
				GameID:      1,
				DeveloperID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameDeveloper domain.GameDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameDeveloper.GameID,
						gameDeveloper.DeveloperID,
						gameDeveloper.ID,
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

			tc.mockBehavior(mock, tc.gameDeveloper)

			err := db.Save(&tc.gameDeveloper).Error

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

func TestSoftDeleteGameDeveloper(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		gameDeveloperID uint
		mockBehavior    func(mock sqlmock.Sqlmock, gameDeveloperID uint)
		wantErr         bool
	}{
		"Can soft delete a GameDeveloper": {
			gameDeveloperID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, gameDeveloperID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_developers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), gameDeveloperID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			gameDeveloperID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, GameDeveloper uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_developers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete GameDeveloper"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.gameDeveloperID)

			err := db.Delete(&domain.GameDeveloper{}, tc.gameDeveloperID).Error

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

func TestValidateGameDeveloper(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameDeveloper domain.GameDeveloper
	}{
		"Can empty validations errors": {
			gameDeveloper: domain.GameDeveloper{
				Developer: domain.Developer{
					Name:      "Game Science",
					Acting:    true,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				Game: domain.Game{
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
					Views: []domain.Viewable{
						{
							UserID:       10,
							ViewableID:   1,
							ViewableType: "games",
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.gameDeveloper.ValidateGameDeveloper()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGameDeveloperWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		gameDeveloper domain.GameDeveloper
		wantErr       string
	}{
		"Missing required fields": {
			gameDeveloper: domain.GameDeveloper{},
			wantErr: `
				Age is a required field,
				Slug is a required field,
				Title is a required field,
				Condition is a required field,
				Cover is a required field,
				About is a required field,
				Description is a required field,
				ShortDescription is a required field,
				ReleaseDate is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.gameDeveloper.ValidateGameDeveloper()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
