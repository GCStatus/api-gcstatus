package tests

import (
	"context"
	"fmt"
	"gcstatus/internal/middlewares"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var MockThrottlePasswordCache = struct {
	GetPasswordThrottleCache func(ctx context.Context, key string) (string, error)
	SetPasswordThrottleCache func(ctx context.Context, key string, duration time.Duration) error
}{}

func GetPasswordThrottleCache(ctx context.Context, key string) (string, error) {
	return MockThrottlePasswordCache.GetPasswordThrottleCache(ctx, key)
}

func SetPasswordThrottleCache(ctx context.Context, key string, duration time.Duration) error {
	return MockThrottlePasswordCache.SetPasswordThrottleCache(ctx, key, duration)
}

func TestLimitResetRequestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		cacheError         error
		cacheHit           bool
		expectedStatusCode int
		expectedMessage    string
	}{
		// "First request (no existing cache), should pass": {
		// 	cacheError:         redis.Nil,
		// 	cacheHit:           false,
		// 	expectedStatusCode: http.StatusOK,
		// 	expectedMessage:    "ok",
		// },
		"Subsequent request within throttle, should block": {
			cacheError:         nil,
			cacheHit:           true,
			expectedStatusCode: http.StatusTooManyRequests,
			expectedMessage:    "You must wait for 60 seconds before sending the email again!",
		},
	}

	mockMail := "test@example.com"

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			MockThrottlePasswordCache.GetPasswordThrottleCache = func(_ context.Context, key string) (string, error) {
				if tt.cacheError != nil {
					return "", tt.cacheError
				}

				if tt.cacheHit {
					return "password-reset:" + mockMail, nil
				}

				return "", redis.Nil
			}

			MockThrottlePasswordCache.SetPasswordThrottleCache = func(_ context.Context, key string, duration time.Duration) error {
				assert.Equal(t, time.Minute, duration)
				return nil
			}

			r := gin.New()
			r.Use(middlewares.LimitResetRequestMiddleware())
			r.POST("/reset-password", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			w := httptest.NewRecorder()
			reqBody := fmt.Sprintf(`{"email": "%s"}`, mockMail)
			req, _ := http.NewRequest(http.MethodPost, "/reset-password", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedMessage)
		})
	}
}
