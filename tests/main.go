package tests

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Setup(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	godotenv.Load(".env.testing")

	db, mock := SetupMockDB(t)

	return db, mock
}

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize mock database: %+v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open mock sql db, got error: %+v", err)
	}

	return gormDB, mock
}
