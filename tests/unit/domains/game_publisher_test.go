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

func TestCreateGamePublisher(t *testing.T) {
	testCases := map[string]struct {
		gamePublisher domain.GamePublisher
		mockBehavior  func(mock sqlmock.Sqlmock, gamePublisher domain.GamePublisher)
		expectError   bool
	}{
		"Success": {
			gamePublisher: domain.GamePublisher{
				GameID:      1,
				PublisherID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gamePublisher domain.GamePublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gamePublisher.GameID,
						gamePublisher.PublisherID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			gamePublisher: domain.GamePublisher{
				GameID:      1,
				PublisherID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gamePublisher domain.GamePublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gamePublisher.GameID,
						gamePublisher.PublisherID,
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

			tc.mockBehavior(mock, tc.gamePublisher)

			err := db.Create(&tc.gamePublisher).Error

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

func TestUpdateGamePublisher(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gamePublisher domain.GamePublisher
		mockBehavior  func(mock sqlmock.Sqlmock, gamePublisher domain.GamePublisher)
		expectError   bool
	}{
		"Success": {
			gamePublisher: domain.GamePublisher{
				ID:          1,
				GameID:      1,
				PublisherID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gamePublisher domain.GamePublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gamePublisher.GameID,
						gamePublisher.PublisherID,
						gamePublisher.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			gamePublisher: domain.GamePublisher{
				ID:          1,
				GameID:      1,
				PublisherID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gamePublisher domain.GamePublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gamePublisher.GameID,
						gamePublisher.PublisherID,
						gamePublisher.ID,
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

			tc.mockBehavior(mock, tc.gamePublisher)

			err := db.Save(&tc.gamePublisher).Error

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

func TestSoftDeleteGamePublisher(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		gamePublisherID uint
		mockBehavior    func(mock sqlmock.Sqlmock, gamePublisherID uint)
		wantErr         bool
	}{
		"Can soft delete a GamePublisher": {
			gamePublisherID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, gamePublisherID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_publishers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), gamePublisherID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			gamePublisherID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, GamePublisher uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_publishers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete GamePublisher"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.gamePublisherID)

			err := db.Delete(&domain.GamePublisher{}, tc.gamePublisherID).Error

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

func TestValidateGamePublisher(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gamePublisher domain.GamePublisher
	}{
		"Can empty validations errors": {
			gamePublisher: domain.GamePublisher{
				Publisher: domain.Publisher{
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
			err := tc.gamePublisher.ValidateGamePublisher()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGamePublisherWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		gamePublisher domain.GamePublisher
		wantErr       string
	}{
		"Missing required fields": {
			gamePublisher: domain.GamePublisher{},
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
			err := tc.gamePublisher.ValidateGamePublisher()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
