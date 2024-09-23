package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock := tests.Setup(t)
	fixedTime := time.Date(2024, 9, 22, 12, 57, 51, 0, time.UTC)

	testCases := map[string]struct {
		user     domain.User
		mockFunc func()
		wantErr  bool
	}{
		"Valid user creation": {
			user: domain.User{
				Name:       "John Doe",
				Email:      "johndoe@example.com",
				Nickname:   "johnny",
				Blocked:    false,
				Experience: 500,
				Birthdate:  fixedTime,
				Password:   "supersecretpassword",
				LevelID:    1,
				Wallet:     domain.Wallet{Amount: 0},
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					"John Doe",
					"johndoe@example.com",
					"johnny",
					500,
					false,
					fixedTime,
					sqlmock.AnyArg(),
					1,
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Duplicate email error": {
			user: domain.User{
				Name:      "John Doe",
				Email:     "johndoe@example.com",
				Nickname:  "johnny",
				Birthdate: fixedTime,
				Password:  "password123",
				LevelID:   1,
				Wallet:    domain.Wallet{Amount: 0},
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					"John Doe",
					"johndoe@example.com",
					"johnny",
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					fixedTime,
					sqlmock.AnyArg(),
					1,
				).WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			err := db.Create(&tc.user).Error

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

func TestGetUserByID(t *testing.T) {
	db, mock := tests.Setup(t)
	fixedTime := time.Date(2024, 9, 22, 12, 57, 51, 0, time.UTC)

	testCases := map[string]struct {
		userID    uint
		mockFunc  func()
		wantUser  domain.User
		wantError bool
	}{
		"Valid user fetch": {
			userID: 1,
			wantUser: domain.User{
				ID:        1,
				Name:      "John Doe",
				Email:     "johndoe@example.com",
				Nickname:  "johnny",
				Blocked:   false,
				Birthdate: fixedTime,
				Wallet:    domain.Wallet{Amount: 0},
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "blocked", "birthdate"}).
					AddRow(1, "John Doe", "johndoe@example.com", "johnny", false, fixedTime)
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"User not found": {
			userID:    2,
			wantUser:  domain.User{},
			wantError: true,
			mockFunc: func() {
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var user domain.User
			err := db.First(&user, tc.userID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantUser, user)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSoftDeleteUser(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		userID   uint
		mockFunc func()
		wantErr  bool
	}{
		"Valid soft delete": {
			userID: 1,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			userID: 2,
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete user"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			err := db.Delete(&domain.User{}, tc.userID).Error

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

func TestUpdateUser(t *testing.T) {
	db, mock := tests.Setup(t)
	fixedTime := time.Date(2024, 9, 22, 12, 57, 51, 0, time.UTC)

	testCases := map[string]struct {
		user     domain.User
		mockFunc func()
		wantErr  bool
	}{
		"Valid user update": {
			user: domain.User{
				ID:         1,
				Name:       "John Doe",
				Email:      "johndoe@example.com",
				Nickname:   "johnny",
				Blocked:    false,
				Experience: 500,
				Birthdate:  fixedTime,
				Password:   "supersecretpassword",
				LevelID:    1,
				Wallet:     domain.Wallet{Amount: 0},
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					"John Doe",
					"johndoe@example.com",
					"johnny",
					500,
					false,
					fixedTime,
					sqlmock.AnyArg(),
					1,
					1,
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Update fails": {
			user: domain.User{
				ID:         2,
				Name:       "Jane Doe",
				Email:      "janedoe@example.com",
				Nickname:   "jane",
				Blocked:    false,
				Experience: 600,
				Birthdate:  fixedTime,
				Password:   "anotherpassword",
				LevelID:    2,
				Wallet:     domain.Wallet{Amount: 0},
			},
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					"Jane Doe",
					"janedoe@example.com",
					"jane",
					600,
					false,
					fixedTime,
					sqlmock.AnyArg(),
					2,
					2,
				).WillReturnError(fmt.Errorf("failed to update user"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			err := db.Save(&tc.user).Error

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

func TestCreateUserWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		user    domain.User
		wantErr string
	}{
		"Missing required fields": {
			user:    domain.User{},
			wantErr: "Name is a required field, Email is a required field, Nickname is a required field, Birthdate is a required field, Password is a required field",
		},
		"Missing Name and Email": {
			user: domain.User{
				Nickname:  "nick",
				Birthdate: time.Now(),
				Password:  "password123",
				Wallet:    domain.Wallet{Amount: 0},
			},
			wantErr: "Name is a required field, Email is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.user.ValidateUser()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
