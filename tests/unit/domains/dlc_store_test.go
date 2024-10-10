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

func TestCreateDLCStore(t *testing.T) {
	testCases := map[string]struct {
		DLCStore     domain.DLCStore
		mockBehavior func(mock sqlmock.Sqlmock, DLCStore domain.DLCStore)
		expectError  bool
	}{
		"Success": {
			DLCStore: domain.DLCStore{
				Price:     2200,
				URL:       "https://google.com",
				DLCID:     1,
				StoreID:   1,
				StorDLCID: "1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCStore domain.DLCStore) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCStore.Price,
						DLCStore.URL,
						DLCStore.DLCID,
						DLCStore.StoreID,
						DLCStore.StorDLCID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			DLCStore: domain.DLCStore{
				Price:     2200,
				URL:       "https://google.com",
				DLCID:     1,
				StoreID:   1,
				StorDLCID: "1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCStore domain.DLCStore) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCStore.Price,
						DLCStore.URL,
						DLCStore.DLCID,
						DLCStore.StoreID,
						DLCStore.StorDLCID,
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

			tc.mockBehavior(mock, tc.DLCStore)

			err := db.Create(&tc.DLCStore).Error

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

func TestUpdateDLCStore(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCStore     domain.DLCStore
		mockBehavior func(mock sqlmock.Sqlmock, DLCStore domain.DLCStore)
		expectError  bool
	}{
		"Success": {
			DLCStore: domain.DLCStore{
				ID:        1,
				Price:     2200,
				URL:       "https://google.com",
				DLCID:     1,
				StoreID:   1,
				StorDLCID: "1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCStore domain.DLCStore) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCStore.Price,
						DLCStore.URL,
						DLCStore.DLCID,
						DLCStore.StoreID,
						DLCStore.StorDLCID,
						DLCStore.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			DLCStore: domain.DLCStore{
				ID:        1,
				Price:     2200,
				URL:       "https://google.com",
				DLCID:     1,
				StoreID:   1,
				StorDLCID: "1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCStore domain.DLCStore) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCStore.Price,
						DLCStore.URL,
						DLCStore.DLCID,
						DLCStore.StoreID,
						DLCStore.StorDLCID,
						DLCStore.ID,
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

			tc.mockBehavior(mock, tc.DLCStore)

			err := db.Save(&tc.DLCStore).Error

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

func TestSoftDeleteDLCStore(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		DLCStoreID   uint
		mockBehavior func(mock sqlmock.Sqlmock, DLCStoreID uint)
		wantErr      bool
	}{
		"Can soft delete a DLCStore": {
			DLCStoreID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCStoreID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_stores` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), DLCStoreID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			DLCStoreID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCStoreID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_stores` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete DLCStore"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.DLCStoreID)

			err := db.Delete(&domain.DLCStore{}, tc.DLCStoreID).Error

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

func TestValidateDLCStore(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCStore domain.DLCStore
	}{
		"Can empty validations errors": {
			DLCStore: domain.DLCStore{
				Price:     2200,
				URL:       "https://google.com",
				StorDLCID: "1",
				DLC: domain.DLC{
					Name:        "Game Science",
					Cover:       "https://google.com",
					ReleaseDate: fixedTime,
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
			err := tc.DLCStore.ValidateDLCStore()
			assert.NoError(t, err)
		})
	}
}

func TestCreateDLCStoreWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		DLCStore domain.DLCStore
		wantErr  string
	}{
		"Missing required fields": {
			DLCStore: domain.DLCStore{},
			wantErr: `
				Price is a required field,
				URL is a required field,
				Name is a required field,
				Cover is a required field,
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
				StorDLCID is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.DLCStore.ValidateDLCStore()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
