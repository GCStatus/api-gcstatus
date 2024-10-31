package tests

import (
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHeartRepositoryMySQL_CreateHeart(t *testing.T) {
	testCases := map[string]struct {
		heartable    *domain.Heartable
		mockBehavior func(mock sqlmock.Sqlmock, heartable *domain.Heartable)
		expectedErr  error
	}{
		"success case": {
			heartable: &domain.Heartable{
				HeartableID:   1,
				HeartableType: "games",
				UserID:        1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, heartable *domain.Heartable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `heartables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						heartable.HeartableID,
						heartable.HeartableType,
						heartable.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"Failure - Insert Error": {
			heartable: &domain.Heartable{
				HeartableID:   1,
				HeartableType: "games",
				UserID:        1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, heartable *domain.Heartable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `heartables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						heartable.HeartableID,
						heartable.HeartableType,
						heartable.UserID,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewHeartRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.heartable)

			err := repo.Create(tc.heartable)

			assert.Equal(t, tc.expectedErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHeartRepositoryMySQL_FindForUser(t *testing.T) {
	testCases := map[string]struct {
		heartableID       uint
		heartableType     string
		userID            uint
		mockBehavior      func(mock sqlmock.Sqlmock)
		expectedheartable *domain.Heartable
		expectedErr       error
	}{
		"valid payload": {
			heartableID:   1,
			heartableType: "games",
			userID:        1,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "heartable_id", "heartable_type", "user_id"}).
					AddRow(1, 1, "games", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `heartables` WHERE (heartable_id = ? AND heartable_type = ? AND user_id = ?) AND `heartables`.`deleted_at` IS NULL ORDER BY `heartables`.`id` LIMIT ?")).
					WithArgs(1, "games", 1, 1).WillReturnRows(rows)
			},
			expectedheartable: &domain.Heartable{ID: 1, HeartableID: 1, HeartableType: "games", UserID: 1},
			expectedErr:       nil,
		},
		"not found payload": {
			heartableID:   999,
			heartableType: "games",
			userID:        999,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `heartables` WHERE (heartable_id = ? AND heartable_type = ? AND user_id = ?) AND `heartables`.`deleted_at` IS NULL ORDER BY `heartables`.`id` LIMIT ?")).
					WithArgs(999, "games", 999, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedheartable: nil,
			expectedErr:       gorm.ErrRecordNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)
			repo := db.NewHeartRepositoryMySQL(gormDB)

			tc.mockBehavior(mock)

			heartable, err := repo.FindForUser(tc.heartableID, tc.heartableType, tc.userID)

			assert.Equal(t, tc.expectedErr, err)

			if tc.expectedheartable == nil {
				assert.Nil(t, heartable)
			} else {
				assert.NotNil(t, heartable)
				assert.Equal(t, tc.expectedheartable.ID, heartable.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHeartRepositoryMySQL_DeleteHeartByID(t *testing.T) {
	testCases := map[string]struct {
		heartableID  uint
		mockBehavior func(mock sqlmock.Sqlmock, heartableID uint)
		wantErr      bool
	}{
		"can delete a heartable": {
			heartableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, heartableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `heartables` WHERE `heartables`.`id` = ?")).
					WithArgs(heartableID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"delete fails": {
			heartableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, heartableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `heartables` WHERE `heartables`.`id` = ?")).
					WithArgs(2).
					WillReturnError(fmt.Errorf("failed to delete heartable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewHeartRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.heartableID)

			err := repo.Delete(tc.heartableID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete heartable")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
