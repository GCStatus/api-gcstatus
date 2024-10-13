package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGameStore(t *testing.T) {
	testCases := map[string]struct {
		gameStore    domain.GameStore
		mockBehavior func(mock sqlmock.Sqlmock, gameStore domain.GameStore)
		expectError  bool
	}{
		"Success": {
			gameStore: domain.GameStore{
				Price:       22999,
				URL:         "https://google.com",
				GameID:      1,
				StoreID:     1,
				StoreGameID: "1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameStore domain.GameStore) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameStore.Price,
						gameStore.URL,
						gameStore.GameID,
						gameStore.StoreID,
						gameStore.StoreGameID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			gameStore: domain.GameStore{
				Price:       22999,
				URL:         "https://google.com",
				GameID:      1,
				StoreID:     1,
				StoreGameID: "1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameStore domain.GameStore) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameStore.Price,
						gameStore.URL,
						gameStore.GameID,
						gameStore.StoreID,
						gameStore.StoreGameID,
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

			tc.mockBehavior(mock, tc.gameStore)

			err := db.Create(&tc.gameStore).Error

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

func TestUpdateGameStore(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameStore    domain.GameStore
		mockBehavior func(mock sqlmock.Sqlmock, gameStore domain.GameStore)
		expectError  bool
	}{
		"Success": {
			gameStore: domain.GameStore{
				ID:          1,
				Price:       22999,
				URL:         "https://google.com",
				GameID:      1,
				StoreID:     1,
				StoreGameID: "1",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameStore domain.GameStore) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameStore.Price,
						gameStore.URL,
						gameStore.GameID,
						gameStore.StoreID,
						gameStore.StoreGameID,
						gameStore.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			gameStore: domain.GameStore{
				ID:          1,
				Price:       22999,
				URL:         "https://google.com",
				GameID:      1,
				StoreID:     1,
				StoreGameID: "1",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameStore domain.GameStore) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameStore.Price,
						gameStore.URL,
						gameStore.GameID,
						gameStore.StoreID,
						gameStore.StoreGameID,
						gameStore.ID,
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

			tc.mockBehavior(mock, tc.gameStore)

			err := db.Save(&tc.gameStore).Error

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

func TestSoftDeleteGameStore(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		gameStoreID  uint
		mockBehavior func(mock sqlmock.Sqlmock, gameStoreID uint)
		wantErr      bool
	}{
		"Can soft delete a GameStore": {
			gameStoreID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, gameStoreID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_stores` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), gameStoreID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			gameStoreID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, gameStoreID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_stores` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete GameStore"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.gameStoreID)

			err := db.Delete(&domain.GameStore{}, tc.gameStoreID).Error

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

func TestValidateGameStoreValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameStore domain.GameStore
	}{
		"Can empty validations errors": {
			gameStore: domain.GameStore{
				Price:       22999,
				URL:         "https://google.com",
				StoreGameID: "1",
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
				Store: domain.Store{
					Name: "Store 1",
					URL:  "https://google.com",
					Slug: "store-1",
					Logo: "https://placehold.co/600x400/EEE/31343C",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.gameStore.ValidateGameStore()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGameStoreWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		gameStore domain.GameStore
		wantErr   string
	}{
		"Missing required fields": {
			gameStore: domain.GameStore{},
			wantErr: `
				Age is a required field,
				Slug is a required field,
				Title is a required field,
				Condition is a required field,
				Cover is a required field,
				About is a required field,
				Description is a required field,
				ShortDescription is a required field,
				ReleaseDate is a required field,
				Name is a required field,
				URL is a required field,
				Slug is a required field,
				Logo is a required field,
				StoreGameID is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.gameStore.ValidateGameStore()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
