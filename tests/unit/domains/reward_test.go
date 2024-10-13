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

func CreateRewardTest(t *testing.T) {
	testCases := map[string]struct {
		reward       domain.Reward
		mockBehavior func(mock sqlmock.Sqlmock, reward domain.Reward)
		expectErr    bool
	}{
		"Successfully created": {
			reward: domain.Reward{
				SourceableID:   1,
				SourceableType: "levels",
				RewardableID:   1,
				RewardableType: "titles",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reward domain.Reward) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `rewards`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reward.SourceableID,
						reward.SourceableType,
						reward.RewardableID,
						reward.RewardableType,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			reward: domain.Reward{
				SourceableID:   1,
				SourceableType: "levels",
				RewardableID:   1,
				RewardableType: "titles",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reward domain.Reward) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `rewards`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reward.SourceableID,
						reward.SourceableType,
						reward.RewardableID,
						reward.RewardableType,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.reward)

			err := db.Create(&tc.reward).Error

			if tc.expectErr {
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

func TestSoftDeleteReward(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		rewardID     uint
		mockBehavior func(mock sqlmock.Sqlmock, rewardID uint)
		wantErr      bool
	}{
		"Can soft delete a title progress": {
			rewardID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, rewardID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `rewards` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), rewardID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			rewardID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, rewardID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `rewards` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete reward"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.rewardID)

			err := db.Delete(&domain.Reward{}, tc.rewardID).Error

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

func TestUpdateReward(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		reward       domain.Reward
		mockBehavior func(mock sqlmock.Sqlmock, reward domain.Reward)
		expectError  bool
	}{
		"Success": {
			reward: domain.Reward{
				ID:             1,
				SourceableID:   1,
				SourceableType: "levels",
				RewardableID:   1,
				RewardableType: "titles",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reward domain.Reward) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `rewards`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reward.SourceableID,
						reward.SourceableType,
						reward.RewardableID,
						reward.RewardableType,
						reward.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			reward: domain.Reward{
				ID:             1,
				SourceableID:   1,
				SourceableType: "levels",
				RewardableID:   1,
				RewardableType: "titles",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reward domain.Reward) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `rewards`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reward.SourceableID,
						reward.SourceableType,
						reward.RewardableID,
						reward.RewardableType,
						reward.ID,
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

			tc.mockBehavior(mock, tc.reward)

			err := db.Save(&tc.reward).Error

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

func TestValidateRewardValidData(t *testing.T) {
	testCases := map[string]struct {
		reward domain.Reward
	}{
		"Can empty validations errors": {
			reward: domain.Reward{
				ID:             1,
				SourceableID:   1,
				SourceableType: "levels",
				RewardableID:   1,
				RewardableType: "titles",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.reward.ValidateReward()

			assert.NoError(t, err)
		})
	}
}

func TestCreateRewardWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		reward  domain.Reward
		wantErr string
	}{
		"Missing required fields": {
			reward: domain.Reward{},
			wantErr: `
				SourceableID is a required field,
				SourceableType is a required field,
				RewardableID is a required field,
				RewardableType is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.reward.ValidateReward()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
