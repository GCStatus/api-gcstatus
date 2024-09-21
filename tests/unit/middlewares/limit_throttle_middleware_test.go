package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/middlewares"
	"gcstatus/pkg/cache"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockCache struct {
	AddFunc                    func(key string) (int64, error)
	ExpireFunc                 func(key string, timeWindow time.Duration)
	GetPasswordThrottleFunc    func(key string) (string, error)
	SetPasswordThrottleFunc    func(key string, duration time.Duration) error
	RemovePasswordThrottleFunc func(email string) error
	GetUserFromCacheFunc       func(userID uint) (*domain.User, bool)
	SetUserInCacheFunc         func(user *domain.User)
	RemoveUserFromCacheFunc    func(userID uint)
}

func (m *MockCache) AddThrottleCache(key string) (int64, error) {
	return m.AddFunc(key)
}

func (m *MockCache) ExpireThrottleCache(key string, timeWindow time.Duration) {
	m.ExpireFunc(key, timeWindow)
}

func (m *MockCache) GetPasswordThrottleCache(key string) (string, error) {
	return "", nil
}

func (m *MockCache) SetPasswordThrottleCache(key string, duration time.Duration) error {
	return nil
}

func (m *MockCache) RemovePasswordThrottleCache(email string) error {
	return nil
}

func (m *MockCache) GetUserFromCache(userID uint) (*domain.User, bool) {
	return nil, false
}

func (m *MockCache) SetUserInCache(user *domain.User) {}

func (m *MockCache) RemoveUserFromCache(userID uint) {}

func TestLimitThrottleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalRateLimit := middlewares.RateLimit

	tests := map[string]struct {
		cacheError           error
		expectedStatusCode   int
		expectedCacheExpires bool
		tempRateLimit        int
		requestCount         int64
	}{
		"First request, should pass": {
			cacheError:           nil,
			expectedStatusCode:   http.StatusOK,
			expectedCacheExpires: true,
			tempRateLimit:        200,
			requestCount:         1,
		},
		"Within rate limit, should pass": {
			cacheError:         nil,
			expectedStatusCode: http.StatusOK,
			tempRateLimit:      200,
			requestCount:       50,
		},
		"Exceeded rate limit, should block": {
			cacheError:         nil,
			expectedStatusCode: http.StatusTooManyRequests,
			tempRateLimit:      1,
			requestCount:       201,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			middlewares.RateLimit = tt.tempRateLimit

			mockCache := &MockCache{
				AddFunc: func(key string) (int64, error) {
					return tt.requestCount, tt.cacheError
				},
				ExpireFunc: func(key string, timeWindow time.Duration) {
					if tt.expectedCacheExpires {
						assert.Equal(t, middlewares.TimeWindow, timeWindow)
					}
				},
			}

			cache.GlobalCache = mockCache

			r := gin.New()
			r.Use(middlewares.LimitThrottleMiddleware())

			r.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "192.168.1.1:1234"

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}

	middlewares.RateLimit = originalRateLimit
}
