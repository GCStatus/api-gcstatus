package tests

import (
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLevelRepositoryMySQL_GetAll(t *testing.T) {
	gormDB, mock := testutils.Setup(t)

	repo := db.NewLevelRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		mockBehavior func()
		expectedLen  int
		expectedErr  error
	}{
		"success case": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "experience", "level", "coins", "created_at", "updated_at"}).
					AddRow(1, 0, 1, 0, time.Now(), time.Now()).
					AddRow(2, 500, 2, 100, time.Now(), time.Now())
				mock.ExpectQuery("^SELECT \\* FROM `levels`").WillReturnRows(rows)

				rewardRows := sqlmock.NewRows([]string{"id", "sourceable_type", "sourceable_id", "rewardable_type", "rewardable_id"}).
					AddRow(1, "levels", 1, "titles", 10).
					AddRow(2, "levels", 2, "titles", 20)
				mock.ExpectQuery("^SELECT \\* FROM `rewards`").
					WithArgs("levels", 1, 2).WillReturnRows(rewardRows)

				titlesRows := sqlmock.NewRows([]string{"id", "title", "description", "purchasable", "status"}).
					AddRow(1, "Title 1", "Title 1", false, "available").
					AddRow(2, "Title 2", "Title 2", false, "available")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE id IN (?,?) AND `titles`.`deleted_at` IS NULL")).
					WillReturnRows(titlesRows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		"no records found": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "experience", "level", "coins", "created_at", "updated_at"})
				mock.ExpectQuery("^SELECT \\* FROM `levels`").WillReturnRows(rows)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior()

			levels, err := repo.GetAll()

			assert.Equal(t, tc.expectedErr, err)
			assert.Len(t, levels, tc.expectedLen)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLevelRepositoryMySQL_FindById(t *testing.T) {
	gormDB, mock := testutils.Setup(t)

	repo := db.NewLevelRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		id            uint
		mockBehavior  func()
		expectedLevel *domain.Level
		expectedErr   error
	}{
		"valid ID": {
			id: 1,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "experience", "level", "coins", "created_at", "updated_at"}).
					AddRow(1, 0, 1, 0, time.Now(), time.Now())
				mock.ExpectQuery("^SELECT \\* FROM `levels` WHERE `levels`.`id` = \\? AND `levels`.`deleted_at` IS NULL ORDER BY `levels`.`id` LIMIT \\?").
					WithArgs(1, 1).WillReturnRows(rows)

				rewardRows := sqlmock.NewRows([]string{"id", "sourceable_id", "sourceable_type", "rewardable_id", "rewardable_type"}).
					AddRow(1, 1, "levels", 1, "titles")
				mock.ExpectQuery("^SELECT \\* FROM `rewards` WHERE `sourceable_type` = \\? AND `rewards`.`sourceable_id` = \\? AND `rewards`.`deleted_at` IS NULL").
					WithArgs("levels", 1).WillReturnRows(rewardRows)
			},
			expectedLevel: &domain.Level{ID: 1, Experience: 0, Level: 1, Coins: 0},
			expectedErr:   nil,
		},
		"not found ID": {
			id: 999,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "experience", "level", "coins", "created_at", "updated_at"})
				mock.ExpectQuery("^SELECT \\* FROM `levels` WHERE `levels`.`id` = \\? AND `levels`.`deleted_at` IS NULL ORDER BY `levels`.`id` LIMIT \\?").
					WithArgs(999, 1).WillReturnRows(rows)
			},
			expectedLevel: &domain.Level{ID: 0},
			expectedErr:   gorm.ErrRecordNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior()

			level, err := repo.FindById(tc.id)

			assert.Equal(t, tc.expectedErr, err)
			if err == gorm.ErrRecordNotFound {
				assert.Equal(t, uint(0), level.ID)
			} else {
				assert.Equal(t, tc.expectedLevel.ID, level.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
