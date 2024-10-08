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

func TestCreateStore(t *testing.T) {
	testCases := map[string]struct {
		store        domain.Store
		mockBehavior func(mock sqlmock.Sqlmock, store domain.Store)
		expectError  bool
	}{
		"Success": {
			store: domain.Store{
				Name: "Store 1",
				URL:  "https://google.com",
				Slug: "store-1",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, store domain.Store) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						store.Name,
						store.URL,
						store.Slug,
						store.Logo,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			store: domain.Store{
				Name: "Store 1",
				URL:  "https://google.com",
				Slug: "store-1",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, store domain.Store) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						store.Name,
						store.URL,
						store.Slug,
						store.Logo,
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

			tc.mockBehavior(mock, tc.store)

			err := db.Create(&tc.store).Error

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

func TestUpdateStore(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		store        domain.Store
		mockBehavior func(mock sqlmock.Sqlmock, store domain.Store)
		expectError  bool
	}{
		"Success": {
			store: domain.Store{
				ID:        1,
				Name:      "Store 1",
				URL:       "https://google.com",
				Slug:      "store-1",
				Logo:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, store domain.Store) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						store.Name,
						store.URL,
						store.Slug,
						store.Logo,
						store.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			store: domain.Store{
				ID:        1,
				Name:      "Store 1",
				URL:       "https://google.com",
				Slug:      "store-1",
				Logo:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, store domain.Store) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `stores`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						store.Name,
						store.URL,
						store.Slug,
						store.Logo,
						store.ID,
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

			tc.mockBehavior(mock, tc.store)

			err := db.Save(&tc.store).Error

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

func TestSoftDeleteStore(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		storeID      uint
		mockBehavior func(mock sqlmock.Sqlmock, storeID uint)
		wantErr      bool
	}{
		"Can soft delete a Store": {
			storeID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, storeID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `stores` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), storeID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			storeID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, storeID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `stores` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Store"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.storeID)

			err := db.Delete(&domain.Store{}, tc.storeID).Error

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

func TestGetStoreByID(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		storeID   uint
		mockFunc  func()
		wantStore domain.Store
		wantError bool
	}{
		"Valid Store fetch": {
			storeID: 1,
			wantStore: domain.Store{
				ID:   1,
				Name: "Store 1",
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Store 1")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `stores` WHERE `stores`.`id` = ? AND `stores`.`deleted_at` IS NULL ORDER BY `stores`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Store not found": {
			storeID:   2,
			wantStore: domain.Store{},
			wantError: true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `stores` WHERE `stores`.`id` = ? AND `stores`.`deleted_at` IS NULL ORDER BY `stores`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var store domain.Store
			err := db.First(&store, tc.storeID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantStore, store)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateStoreValidData(t *testing.T) {
	testCases := map[string]struct {
		store domain.Store
	}{
		"Can empty validations errors": {
			store: domain.Store{
				Name: "Store 1",
				URL:  "https://google.com",
				Slug: "store-1",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.store.ValidateStore()
			assert.NoError(t, err)
		})
	}
}

func TestCreateStoreWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		store   domain.Store
		wantErr string
	}{
		"Missing required fields": {
			store:   domain.Store{},
			wantErr: "Name is a required field, URL is a required field, Slug is a required field, Logo is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.store.ValidateStore()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
