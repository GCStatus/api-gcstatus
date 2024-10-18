package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCritic(t *testing.T) {
	testCases := map[string]struct {
		critic       domain.Critic
		mockBehavior func(mock sqlmock.Sqlmock, critic domain.Critic)
		expectError  bool
	}{
		"Success": {
			critic: domain.Critic{
				Name: "Critic 1",
				URL:  "https://google.com",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, critic domain.Critic) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `critics`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						critic.Name,
						critic.URL,
						critic.Logo,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			critic: domain.Critic{
				Name: "Critic 1",
				URL:  "https://google.com",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, critic domain.Critic) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `critics`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						critic.Name,
						critic.URL,
						critic.Logo,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := testutils.Setup(t)

			tc.mockBehavior(mock, tc.critic)

			err := db.Create(&tc.critic).Error

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

func TestUpdateCritic(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		critic       domain.Critic
		mockBehavior func(mock sqlmock.Sqlmock, critic domain.Critic)
		expectError  bool
	}{
		"Success": {
			critic: domain.Critic{
				ID:        1,
				Name:      "Critic 1",
				URL:       "https://google.com",
				Logo:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, critic domain.Critic) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `critics`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						critic.Name,
						critic.URL,
						critic.Logo,
						critic.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			critic: domain.Critic{
				ID:        1,
				Name:      "Critic 1",
				URL:       "https://google.com",
				Logo:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, critic domain.Critic) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `critics`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						critic.Name,
						critic.URL,
						critic.Logo,
						critic.ID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := testutils.Setup(t)

			tc.mockBehavior(mock, tc.critic)

			err := db.Save(&tc.critic).Error

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

func TestSoftDeleteCritic(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		criticID     uint
		mockBehavior func(mock sqlmock.Sqlmock, criticID uint)
		wantErr      bool
	}{
		"Can soft delete a Critic": {
			criticID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, criticID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `critics` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), criticID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			criticID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, criticID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `critics` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Critic"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.criticID)

			err := db.Delete(&domain.Critic{}, tc.criticID).Error

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

func TestGetCriticByID(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		criticID   uint
		mockFunc   func()
		wantCritic domain.Critic
		wantError  bool
	}{
		"Valid Critic fetch": {
			criticID: 1,
			wantCritic: domain.Critic{
				ID:   1,
				Name: "Critic 1",
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Critic 1")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `critics` WHERE `critics`.`id` = ? AND `critics`.`deleted_at` IS NULL ORDER BY `critics`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Critic not found": {
			criticID:   2,
			wantCritic: domain.Critic{},
			wantError:  true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `critics` WHERE `critics`.`id` = ? AND `critics`.`deleted_at` IS NULL ORDER BY `critics`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var critic domain.Critic
			err := db.First(&critic, tc.criticID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantCritic, critic)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateCriticValidData(t *testing.T) {
	testCases := map[string]struct {
		critic domain.Critic
	}{
		"Can empty validations errors": {
			critic: domain.Critic{
				Name: "Critic 1",
				URL:  "https://google.com",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.critic.ValidateCritic()
			assert.NoError(t, err)
		})
	}
}

func TestCreateCriticWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		critic  domain.Critic
		wantErr string
	}{
		"Missing required fields": {
			critic:  domain.Critic{},
			wantErr: "Name is a required field, URL is a required field, Logo is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.critic.ValidateCritic()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
