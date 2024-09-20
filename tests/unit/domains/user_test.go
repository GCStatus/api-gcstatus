package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock := tests.Setup(t)

	defer func() {
		dbConn, _ := db.DB()
		dbConn.Close()
	}()

	user := domain.User{
		Name:      "John Doe",
		Email:     "johndoe@example.com",
		Nickname:  "johnny",
		Blocked:   false,
		Birthdate: time.Date(1990, 01, 01, 0, 0, 0, 0, time.UTC),
		Password:  "supersecretpassword",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		user.Name,
		user.Email,
		user.Nickname,
		user.Blocked,
		user.Birthdate,
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := db.Create(&user).Error
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock := tests.Setup(t)

	defer func() {
		dbConn, _ := db.DB()
		dbConn.Close()
	}()

	userID := uint(1)
	mockUser := domain.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "johndoe@example.com",
		Nickname:  "johnny",
		Blocked:   false,
		Birthdate: time.Date(1990, 01, 01, 0, 0, 0, 0, time.UTC),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "blocked", "birthdate"}).
		AddRow(mockUser.ID, mockUser.Name, mockUser.Email, mockUser.Nickname, mockUser.Blocked, mockUser.Birthdate)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(userID, 1).
		WillReturnRows(rows)

	var user domain.User
	err := db.First(&user, userID).Error
	assert.NoError(t, err)
	assert.Equal(t, mockUser.ID, user.ID)
	assert.Equal(t, mockUser.Email, user.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSoftDeleteUser(t *testing.T) {
	db, mock := tests.Setup(t)

	defer func() {
		dbConn, _ := db.DB()
		dbConn.Close()
	}()

	userID := uint(1)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), userID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := db.Delete(&domain.User{}, userID).Error
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
