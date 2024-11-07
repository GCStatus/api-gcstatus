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

var commentTruncateModels = []any{
	&domain.User{},
	&domain.Wallet{},
	&domain.Profile{},
	&domain.Commentable{},
}

func setupCommentHandler(dbConn *gorm.DB) *api.CommentHandler {
	userService := usecases.NewUserService(db.NewUserRepositoryMySQL(dbConn))
	commentService := usecases.NewCommentService(db.NewCommentRepositoryMySQL(dbConn))
	return api.NewCommentHandler(userService, commentService)
}

func TestCommentHandler_Create(t *testing.T) {
	commentHandler := setupCommentHandler(dbConn)

	tests := map[string]struct {
		payload        string
		expectCode     int
		expectResponse map[string]any
	}{
		"valid comment payload": {
			payload: fmt.Sprintf(`{
				"commentable_id": %d,
				"commentable_type": "games",
				"comment": "Just testing comment"
    		}`, uint(1)),
			expectCode: http.StatusCreated,
			expectResponse: map[string]any{
				"comment":          "Just testing comment",
				"commentable_id":   float64(1),
				"commentable_type": "games",
			},
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

			req := httptest.NewRequest(http.MethodPost, "/comments", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			user, err := test_mocks.ActingAsDummyUser(t, dbConn, &domain.User{}, req, env)
			if err != nil {
				t.Fatalf("failed to create dummy user: %+v", err)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			commentHandler.Create(c)

			assert.Equal(t, tc.expectCode, w.Code)
			var responseBody map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("failed to parse JSON response: %+v", err)
			}

			if w.Code == http.StatusCreated {
				data, ok := responseBody["data"].(map[string]any)
				if assert.True(t, ok, "response should contain 'data' field") {
					for key, expectedValue := range tc.expectResponse {
						if key == "comment" { // the only value from response
							actualValue := data[key]

							assert.Equal(t, expectedValue, actualValue, "unexpected value for '%s'", key)
						}
					}
				}

				var createdComment domain.Commentable
				err := dbConn.First(&createdComment).Error
				assert.NoError(t, err, "Comment record should exist in the database")

				var payloadData map[string]any
				if err := json.Unmarshal([]byte(tc.payload), &payloadData); err != nil {
					t.Fatalf("failed to unmarshal payload body: %+v", err)
				}

				assert.Equal(t, payloadData["comment"], createdComment.Comment)
				assert.Equal(t, uint(payloadData["commentable_id"].(float64)), createdComment.CommentableID)
				assert.Equal(t, payloadData["commentable_type"], createdComment.CommentableType)
				assert.Equal(t, user.ID, createdComment.UserID)
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
		testutils.RefreshDatabase(t, dbConn, commentTruncateModels)
	})
}

func TestCommentHandler_Delete(t *testing.T) {
	commentHandler := setupCommentHandler(dbConn)

	tests := map[string]struct {
		expectCode     int
		expectResponse map[string]any
		setupComment   bool
		anotherUser    bool
	}{
		"can delete a comment": {
			expectCode:     http.StatusOK,
			expectResponse: map[string]any{"message": "Your comment was successfully removed!"},
			setupComment:   true,
			anotherUser:    false,
		},
		"comment not found": {
			expectCode:     http.StatusNotFound,
			expectResponse: map[string]any{"message": "Could not found the given comment!"},
			setupComment:   false,
			anotherUser:    false,
		},
		"cannot delete another user's comment": {
			expectCode:     http.StatusForbidden,
			expectResponse: map[string]any{"message": "This comment does not belongs to you user!"},
			setupComment:   true,
			anotherUser:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var comment *domain.Commentable
			var err error
			var req *http.Request

			user, err := test_mocks.CreateDummyUser(t, dbConn, &domain.User{})
			if err != nil {
				t.Fatalf("failed to create dummy user: %+v", err)
			}

			commentOwner := user
			if tc.anotherUser {
				commentOwner, err = test_mocks.CreateDummyUser(t, dbConn, &domain.User{})
				if err != nil {
					t.Fatalf("failed to create another dummy user: %+v", err)
				}
			}

			if tc.setupComment {
				comment, err = test_mocks.CreateDummyComment(t, dbConn, &domain.Commentable{User: *commentOwner})
				if err != nil {
					t.Fatalf("failed to create dummy comment: %+v", err)
				}
				req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/comments/%d", comment.ID), nil)
			} else {
				req = httptest.NewRequest(http.MethodDelete, "/comments/999999", nil)
			}

			req.Header.Set("Content-Type", "application/json")
			if token := testutils.GenerateAuthTokenForUser(t, user); token != "" {
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

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if tc.setupComment {
				c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", comment.ID)}}
			} else {
				c.Params = gin.Params{{Key: "id", Value: "999999"}}
			}

			commentHandler.Delete(c)

			assert.Equal(t, tc.expectCode, w.Code)

			var responseBody map[string]any
			if err = json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("failed to parse JSON response: %+v", err)
			}

			if data, exists := tc.expectResponse["data"]; exists {
				if message, exists := data.(map[string]any)["message"]; exists {
					assert.Equal(t, message, responseBody["message"], "unexpected response message")
				}
			}
		})
	}

	t.Cleanup(func() {
		testutils.RefreshDatabase(t, dbConn, commentTruncateModels)
	})
}
