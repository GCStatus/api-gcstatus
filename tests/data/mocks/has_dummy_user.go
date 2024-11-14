package test_mocks

import (
	"fmt"
	"gcstatus/config"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	testingutils "gcstatus/tests/utils"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"gorm.io/gorm"
)

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CreateDummyUser(t *testing.T, dbConn *gorm.DB, overrides *domain.User) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword("admin1234")
	if err != nil {
		t.Fatalf("failed to hash dummy user password: %+v", err)
	}

	randomSuffix := randomString(8)

	defaultUser := domain.User{
		Name:       fmt.Sprintf("User_%s", randomSuffix),
		Email:      fmt.Sprintf("user_%s@example.com", randomSuffix),
		Nickname:   fmt.Sprintf("nickname_%s", randomSuffix),
		Experience: 0,
		Blocked:    false,
		Birthdate:  time.Now(),
		Password:   hashedPassword,
		LevelID:    1,
		Profile:    domain.Profile{Share: false},
		Wallet:     domain.Wallet{Amount: 0},
	}

	if overrides != nil {
		if overrides.Name != "" {
			defaultUser.Name = overrides.Name
		}
		if overrides.Email != "" {
			defaultUser.Email = overrides.Email
		}
		if overrides.Nickname != "" {
			defaultUser.Nickname = overrides.Nickname
		}
		if overrides.Experience != 0 {
			defaultUser.Experience = overrides.Experience
		}
		if overrides.Password != "" {
			hashedPassword, err := utils.HashPassword(overrides.Password)
			if err != nil {
				t.Fatalf("failed to hash dummy user password: %+v", err)
			}
			defaultUser.Password = hashedPassword
		}
		if !overrides.Birthdate.IsZero() {
			defaultUser.Birthdate = overrides.Birthdate
		}
	}

	if err := dbConn.Create(&defaultUser).Error; err != nil {
		return nil, err
	}

	return &defaultUser, nil
}

func ActingAsDummyUser(
	t *testing.T,
	dbConn *gorm.DB,
	overrides *domain.User,
	req *http.Request,
	env *config.Config,
) (*domain.User, error) {
	user, err := CreateDummyUser(t, dbConn, overrides)
	if err != nil {
		t.Fatalf("failed to create user on acting method: %+v", err)
	}

	token := testingutils.GenerateAuthTokenForUser(t, user)
	if token != "" {
		req.AddCookie(&http.Cookie{
			Name:     env.AccessTokenKey,
			Value:    token,
			Path:     "/",
			Domain:   env.Domain,
			HttpOnly: true,
			Secure:   false,
			MaxAge:   86400,
		})
	}

	return user, err
}
