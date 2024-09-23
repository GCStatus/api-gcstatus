package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func MockSendEmail(recipient, body, subject string) error {
	if recipient == "fail@example.com" {
		return errors.New("failed to send email")
	}

	return nil
}

func Setup(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	err := LoadEnv()
	if err != nil {
		t.Fatalf("failed to load env variables: %v", err)
	}

	db, mock := SetupMockDB(t)

	Seed(db)

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

func LoadEnv() error {
	envPath := filepath.Join(basepath, "..", ".env.testing")
	err := godotenv.Load(envPath)

	return err
}

func Seed(db *gorm.DB) {
	levels := []domain.Level{
		{ID: 1, Experience: 0, Level: 1, Coins: 0, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Experience: 500, Level: 2, Coins: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	db.Create(&levels)
}