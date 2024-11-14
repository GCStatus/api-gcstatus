package testutils

import (
	"errors"
	"fmt"
	"gcstatus/config"
	"gcstatus/di"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	data_test "gcstatus/tests/data"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	envPath := filepath.Join(basepath, "../..", ".env.testing")
	err := godotenv.Load(envPath)

	return err
}

func Seed(db *gorm.DB) {
	levels := []domain.Level{
		{ID: 1, Experience: 0, Level: 1, Coins: 0, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Experience: 500, Level: 2, Coins: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	titles := []domain.Title{
		{ID: 1, Title: "Title 1", Description: "Title 1", Purchasable: false, Status: "available"},
	}

	db.Create(&titles).Create(&levels)
}

func SetupGinTestContext(method, url string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(method, url, strings.NewReader(body))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

func SetupTestDB(t *testing.T) (*gorm.DB, *config.Config) {
	err := LoadEnv()
	if err != nil {
		t.Fatalf("failed to load env file: %+v", err)
	}

	env := config.LoadConfig()
	dsn := GetDBConnectionURL(env)
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	models := di.GetModels()
	for _, model := range models {
		if err := dbConn.AutoMigrate(model); err != nil {
			t.Fatalf("Failed to migrate table for model %T: %v", model, err)
		}
	}

	data_test.Seed(t, dbConn)

	return dbConn, env
}

func GetDBConnectionURL(config *config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
	)
}

func RefreshDatabase(t *testing.T, dbConn *gorm.DB, models []any) {
	dbConn.Exec("SET FOREIGN_KEY_CHECKS=0;")

	for _, model := range models {
		stmt := &gorm.Statement{DB: dbConn}
		err := stmt.Parse(model)
		if err != nil {
			t.Fatalf("failed to parse model table: %+v", err)
		}

		tableName := stmt.Schema.Table

		if err := dbConn.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
			t.Fatalf("Failed to truncate table for model %T: %v", model, err)
		}
	}

	dbConn.Exec("SET FOREIGN_KEY_CHECKS=1;")
}

func GenerateAuthTokenForUser(t *testing.T, user *domain.User) string {
	env := config.LoadConfig()

	secret := []byte(env.JwtSecret)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}).SignedString(secret)
	if err != nil {
		t.Fatalf("failed to generate user token: %+v", err)
	}

	encryptedToken, err := utils.Encrypt(token, env.JwtSecret)
	if err != nil {
		t.Fatalf("encryption error: %+v", err)
	}

	return encryptedToken
}
