package feature_tests

import (
	"encoding/json"
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/internal/usecases"
	test_mocks "gcstatus/tests/data/mocks"
	testutils "gcstatus/tests/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var levelTruncateModels = []any{
	&domain.User{},
	&domain.Wallet{},
	&domain.Profile{},
	&domain.Level{},
}

func setupLevelHandler(dbConn *gorm.DB) *api.LevelHandler {
	levelService := usecases.NewLevelService(db.NewLevelRepositoryMySQL(dbConn))
	return api.NewLevelHandler(levelService)
}

func TestLevelHandler_GetAll(t *testing.T) {

	levelHandler := setupLevelHandler(dbConn)

	tests := map[string]struct {
		expectCode     int
		expectResponse []map[string]any
	}{
		"successful get all": {
			expectCode: http.StatusOK,
			expectResponse: []map[string]any{
				{
					"id":         float64(1),
					"level":      float64(1),
					"coins":      float64(0),
					"experience": float64(0),
				},
				{
					"id":         float64(2),
					"level":      float64(2),
					"coins":      float64(100),
					"experience": float64(50),
				},
				{
					"id":         float64(3),
					"level":      float64(3),
					"coins":      float64(200),
					"experience": float64(100),
				},
				{
					"id":         float64(4),
					"level":      float64(4),
					"coins":      float64(400),
					"experience": float64(200),
					"rewards": []map[string]any{
						{
							"id":              float64(1),
							"rewardable_type": "titles",
							"sourceable_type": "levels",
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/levels", nil)
			req.Header.Set("Content-Type", "application/json")

			_, err := test_mocks.CreateDummyLevel(t, dbConn, &domain.Level{
				Level:      4,
				Experience: 200,
				Coins:      400,
				Rewards: []domain.Reward{
					{
						ID:             1,
						SourceableID:   1,
						SourceableType: "levels",
						RewardableID:   1,
						RewardableType: "titles",
					},
				},
			})
			if err != nil {
				t.Fatalf("failed to create dummy level: %+v", err)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			levelHandler.GetAll(c)

			assert.Equal(t, tc.expectCode, w.Code)

			var responseBody map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("failed to parse JSON response: %+v", err)
			}

			data, ok := responseBody["data"].([]any)
			if assert.True(t, ok, "response should contain 'data' field as an array") {
				for i, expectedLevel := range tc.expectResponse {
					levelData, ok := data[i].(map[string]any)
					if assert.True(t, ok, "each level should be a map") {
						for key, expectedValue := range expectedLevel {
							if key == "rewards" {
								rewards, ok := levelData["rewards"].([]any)
								if assert.True(t, ok, "rewards should be an array") {
									for j, expectedReward := range expectedLevel["rewards"].([]map[string]any) {
										rewardData := rewards[j].(map[string]any)
										for rewardKey, rewardValue := range expectedReward {
											assert.Equal(t, rewardValue, rewardData[rewardKey], "unexpected value for '%s'", rewardKey)
										}
									}
								}
							} else {
								assert.Equal(t, expectedValue, levelData[key], "unexpected value for '%s'", key)
							}
						}
					}
				}
			}
		})
	}

	t.Cleanup(func() {
		testutils.RefreshDatabase(t, dbConn, levelTruncateModels)
	})
}
