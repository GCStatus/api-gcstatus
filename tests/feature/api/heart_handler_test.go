package feature_tests

import (
	"encoding/json"
	"fmt"
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/internal/usecases"
	test_mocks "gcstatus/tests/data/mocks"
	testutils "gcstatus/tests/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var heartTruncateModels = []any{
	&domain.User{},
	&domain.Wallet{},
	&domain.Profile{},
	&domain.Heartable{},
}

func setupHeartHandler(dbConn *gorm.DB) *api.HeartHandler {
	userService := usecases.NewUserService(db.NewUserRepositoryMySQL(dbConn))
	heartService := usecases.NewHeartService(db.NewHeartRepositoryMySQL(dbConn))
	return api.NewHeartHandler(userService, heartService)
}

func TestHeartHandler_Create(t *testing.T) {
	heartableHandler := setupHeartHandler(dbConn)

	tests := map[string]struct {
		payload        string
		expectCode     int
		expectResponse map[string]any
	}{
		"valid heartable payload": {
			payload: fmt.Sprintf(`{
				"heartable_id": %d,
				"heartable_type": "games"
    		}`, uint(1)),
			expectCode:     http.StatusOK,
			expectResponse: map[string]any{"message": "The heart operation runned successfully!"},
		},
		"invalid payload": {
			payload:        `{}`,
			expectCode:     http.StatusUnprocessableEntity,
			expectResponse: map[string]any{"message": "Invalid request data"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/hearts", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			user, err := test_mocks.ActingAsDummyUser(t, dbConn, &domain.User{}, req, env)
			if err != nil {
				t.Fatalf("failed to create dummy user: %+v", err)
			}

			var payloadData map[string]any
			if err := json.Unmarshal([]byte(tc.payload), &payloadData); err != nil {
				t.Fatalf("failed to unmarshal payload body: %+v", err)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			heartableHandler.ToggleHeartable(c)

			assert.Equal(t, tc.expectCode, w.Code)

			var responseBody map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("failed to parse JSON response: %+v", err)
			}

			if w.Code == http.StatusOK {
				var heartable domain.Heartable
				heartableID := uint(payloadData["heartable_id"].(float64))

				err := dbConn.First(&heartable, heartableID).Error
				assert.NoError(t, err, "Heart record should exist in the database")
				assert.Equal(t, user.ID, heartable.UserID)
				assert.Equal(t, payloadData["heartable_type"], heartable.HeartableType)
				assert.Equal(t, uint(payloadData["heartable_id"].(float64)), heartable.HeartableID)
				assert.False(t, heartable.DeletedAt.Valid, "Expected DeletedAt to be nil for active record")
			} else {
				if data, exists := tc.expectResponse["data"]; exists {
					if message, exists := data.(map[string]any)["message"]; exists {
						assert.Equal(t, message, responseBody["message"], "unexpected response message")
					}
				}
			}
		})
	}

	t.Cleanup(func() {
		testutils.RefreshDatabase(t, dbConn, heartTruncateModels)
	})
}
