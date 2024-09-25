package tests

import (
	"errors"
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepositoryMySQL_GetAllUsers(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewUserRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		mockBehavior func()
		expectedLen  int
		expectedErr  error
	}{
		"success case": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"}).
					AddRow(1, "Fake", "fake@gmail.com", "fake", time.Now(), time.Now()).
					AddRow(2, "Fake 2", "fake2@gmail.com", "fake2", time.Now(), time.Now())
				mock.ExpectQuery("^SELECT \\* FROM `users`").WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		"no records found": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"})
				mock.ExpectQuery("^SELECT \\* FROM `users`").WillReturnRows(rows)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior()

			users, err := repo.GetAllUsers()

			assert.Equal(t, tc.expectedErr, err)
			assert.Len(t, users, tc.expectedLen)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepositoryMySQL_GetUserByID(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewUserRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		id           uint
		mockBehavior func()
		expectedUser *domain.User
		expectedErr  error
	}{
		"valid ID": {
			id: 1,
			mockBehavior: func() {
				userRows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"}).
					AddRow(1, "Fake", "fake@gmail.com", "fake", time.Now(), time.Now())
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(1, 1).WillReturnRows(userRows)

				profileRows := sqlmock.NewRows([]string{"id", "user_id"}).
					AddRow(1, 1)
				mock.ExpectQuery("^SELECT \\* FROM `profiles` WHERE `profiles`.`user_id` = \\? AND `profiles`.`deleted_at` IS NULL").
					WithArgs(1).WillReturnRows(profileRows)

				walletRows := sqlmock.NewRows([]string{"id", "user_id"}).
					AddRow(1, 1)
				mock.ExpectQuery("^SELECT \\* FROM `wallets` WHERE `wallets`.`user_id` = \\? AND `wallets`.`deleted_at` IS NULL").
					WithArgs(1).WillReturnRows(walletRows)
			},
			expectedUser: &domain.User{ID: 1, Profile: domain.Profile{ID: 1}, Wallet: domain.Wallet{ID: 1}},
			expectedErr:  nil,
		},
		"not found ID": {
			id: 999,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"})
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(999, 1).WillReturnRows(rows)
			},
			expectedUser: &domain.User{ID: 0},
			expectedErr:  gorm.ErrRecordNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior()

			user, err := repo.GetUserByID(tc.id)

			assert.Equal(t, tc.expectedErr, err)
			if err == gorm.ErrRecordNotFound {
				assert.Equal(t, uint(0), user.ID)
			} else {
				assert.Equal(t, tc.expectedUser.ID, user.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepositoryMySQL_CreateUser(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		user         *domain.User
		mockBehavior func(mock sqlmock.Sqlmock, user *domain.User)
		expectedErr  error
	}{
		"success case": {
			user: &domain.User{
				ID:         1,
				Email:      "fake@gmail.com",
				Name:       "Fake",
				Nickname:   "fake",
				Experience: 0,
				Birthdate:  fixedTime,
				Password:   "fake1234",
				Blocked:    false,
				LevelID:    1,
				Wallet: domain.Wallet{
					ID:     1,
					Amount: 100,
					UserID: 1,
				},
				Profile: domain.Profile{
					ID:     1,
					UserID: 1,
				},
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, user *domain.User) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						user.Name,
						user.Email,
						user.Nickname,
						user.Experience,
						user.Blocked,
						user.Birthdate,
						user.Password,
						user.LevelID,
						user.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `profiles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						false,
						"",
						"",
						"",
						"",
						"",
						"",
						"",
						"",
						user.ID,
						user.Profile.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `wallets`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						100,
						user.ID,
						user.Profile.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"Failure - Insert Error": {
			user: &domain.User{
				ID:         1,
				Email:      "fake@gmail.com",
				Name:       "Fake",
				Nickname:   "fake",
				Experience: 0,
				Birthdate:  fixedTime,
				Password:   "fake1234",
				Blocked:    false,
				LevelID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, user *domain.User) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						user.Name,
						user.Email,
						user.Nickname,
						user.Experience,
						user.Blocked,
						user.Birthdate,
						user.Password,
						user.LevelID,
						user.ID,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := tests.Setup(t)

			repo := db.NewUserRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.user)

			err := repo.CreateUser(tc.user)

			assert.Equal(t, tc.expectedErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepositoryMySQL_FindUserByEmailOrNickname(t *testing.T) {
	testCases := map[string]struct {
		searchable   string
		mockBehavior func(mock sqlmock.Sqlmock)
		expectedUser *domain.User
		expectedErr  error
	}{
		"find by username": {
			searchable: "fake",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"}).
					AddRow(1, "Fake", "fake@gmail.com", "fake", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (nickname = ? OR email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
					WithArgs("fake", "fake", 1).WillReturnRows(rows)
			},
			expectedUser: &domain.User{ID: 1, Email: "fake@gmail.com", Nickname: "fake"},
			expectedErr:  nil,
		},
		"find by email": {
			searchable: "fake@gmail.com",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"}).
					AddRow(1, "Fake", "fake@gmail.com", "fake", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (nickname = ? OR email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
					WithArgs("fake@gmail.com", "fake@gmail.com", 1).WillReturnRows(rows)
			},
			expectedUser: &domain.User{ID: 1, Email: "fake@gmail.com", Nickname: "fake"},
			expectedErr:  nil,
		},
		"not found user": {
			searchable: "wrong",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (nickname = ? OR email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
					WithArgs("wrong", "wrong", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser: &domain.User{ID: 0},
			expectedErr:  gorm.ErrRecordNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := tests.Setup(t)

			repo := db.NewUserRepositoryMySQL(gormDB)

			tc.mockBehavior(mock)

			user, err := repo.FindUserByEmailOrNickname(tc.searchable)

			assert.Equal(t, tc.expectedErr, err)
			if err == gorm.ErrRecordNotFound {
				assert.Equal(t, uint(0), user.ID)
			} else {
				assert.Equal(t, tc.expectedUser.ID, user.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepositoryMySQL_UpdateUserPassword(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewUserRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID         uint
		hashedPassword string
		mockBehavior   func()
		expectedErr    error
	}{
		"successful password update": {
			userID:         1,
			hashedPassword: "newHashedPassword123",
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `password`=?,`updated_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs("newHashedPassword123", sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"failed password update (db error)": {
			userID:         1,
			hashedPassword: "newHashedPassword123",
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `password`=?,`updated_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs("newHashedPassword123", sqlmock.AnyArg(), 1).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior()

			err := repo.UpdateUserPassword(tc.userID, tc.hashedPassword)

			assert.Equal(t, tc.expectedErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func UserRepositoryMySQL_UpdateUserNickAndEmail(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewUserRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID        uint
		newNickname   string
		newEmail      string
		password      string
		mockBehavior  func(newNickname string, newEmail string, userID uint)
		expectedError error
	}{
		"success": {
			userID:      1,
			newNickname: "user2",
			newEmail:    "user2@example.com",
			password:    "validpass1234",
			mockBehavior: func(newNickname string, newEmail string, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `nickname`=?,`email`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(newNickname, newEmail, userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		"duplicated nickname": {
			userID:      1,
			newNickname: "user1",
			newEmail:    "user2@example.com",
			password:    "validpass1234",
			mockBehavior: func(newNickname string, newEmail string, userID uint) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE nickname = ? AND id != ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(newNickname, userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
			},
			expectedError: errors.New("nickname already in use"),
		},
		"duplicated email": {
			userID:      1,
			newNickname: "user2",
			newEmail:    "user@example.com",
			password:    "validpass1234",
			mockBehavior: func(newNickname string, newEmail string, userID uint) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND id != ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(newEmail, userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
			},
			expectedError: errors.New("email already in use"),
		},
		"db error": {
			userID:      999,
			newNickname: "user2",
			newEmail:    "user2@example.com",
			password:    "validpass1234",
			mockBehavior: func(newNickname string, newEmail string, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `nickname`=?,`email`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(newNickname, newEmail, userID).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("db error"),
		},
		"wrongpass": {
			userID:      1,
			newNickname: "user2",
			newEmail:    "user2@example.com",
			password:    "wrongpass",
			mockBehavior: func(newNickname string, newEmail string, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `nickname`=?,`email`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(newNickname, newEmail, userID).
					WillReturnError(errors.New("password does not match"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("password does not match"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.mockBehavior(tc.newNickname, tc.newEmail, tc.userID)

			request := ports.UpdateNickAndEmailRequest{
				Password: tc.password,
				Nickname: tc.newNickname,
				Email:    tc.newEmail,
			}

			err := repo.UpdateUserNickAndEmail(tc.userID, request)

			assert.Equal(t, tc.expectedError, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepositoryMySQL_UpdateUserBasics(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := tests.Setup(t)

	repo := db.NewUserRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID        uint
		newName       string
		newBirthdate  string
		mockBehavior  func(newName string, newBirthdate string, userID uint)
		expectedError error
	}{
		"success": {
			userID:       1,
			newName:      "user2",
			newBirthdate: fixedTime.Format("2006-01-02"),
			mockBehavior: func(newName string, newBirthdate string, userID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `birthdate`=?,`name`=?,`updated_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(newBirthdate, newName, sqlmock.AnyArg(), userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.mockBehavior(tc.newName, tc.newBirthdate, tc.userID)

			request := ports.UpdateUserBasicsRequest{
				Name:      tc.newName,
				Birthdate: tc.newBirthdate,
			}

			err := repo.UpdateUserBasics(tc.userID, request)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
