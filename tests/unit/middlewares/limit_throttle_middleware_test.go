package tests

import (
	"context"
	"gcstatus/internal/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var MockThrottleCache = struct {
	AddThrottleCache    func(ctx context.Context, key string) error
	ExpireThrottleCache func(ctx context.Context, key string, timeWindow time.Duration) error
}{}

func AddThrottleCache(ctx context.Context, key string) error {
	return MockThrottleCache.AddThrottleCache(ctx, key)
}

func ExpireThrottleCache(ctx context.Context, key string, timeWindow time.Duration) error {
	return MockThrottleCache.ExpireThrottleCache(ctx, key, timeWindow)
}

func TestLimitThrottleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalRateLimit := middlewares.RateLimit

	tests := map[string]struct {
		cacheError           error
		expectedStatusCode   int
		expectedCacheExpires bool
		tempRateLimit        int
	}{
		"First request, should pass": {
			cacheError:           nil,
			expectedStatusCode:   http.StatusOK,
			expectedCacheExpires: true,
			tempRateLimit:        200,
		},
		"Within rate limit, should pass": {
			cacheError:         nil,
			expectedStatusCode: http.StatusOK,
			tempRateLimit:      200,
		},
		"Exceeded rate limit, should block": {
			cacheError:         nil,
			expectedStatusCode: http.StatusTooManyRequests,
			tempRateLimit:      1,
		},
		"Exceeded rate limit with lower threshold, should block": {
			cacheError:         nil,
			expectedStatusCode: http.StatusTooManyRequests,
			tempRateLimit:      1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			middlewares.RateLimit = tt.tempRateLimit

			MockThrottleCache.AddThrottleCache = func(ctx context.Context, key string) error {
				return tt.cacheError
			}

			MockThrottleCache.ExpireThrottleCache = func(ctx context.Context, key string, timeWindow time.Duration) error {
				if tt.expectedCacheExpires {
					assert.Equal(t, middlewares.TimeWindow, timeWindow)
				}

				return nil
			}

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
