package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireAuth(t *testing.T) {
	t.Parallel()

	middleware := NewJWTMiddleware("test-secret")

	tests := []struct {
		name       string
		userID     string
		expectCode int
	}{
		{
			name:       "with user ID in context",
			userID:     "test-user-id",
			expectCode: http.StatusOK,
		},
		{
			name:       "without user ID in context",
			userID:     "",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handlerCalled := false
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			requireAuthHandler := middleware.RequireAuth(testHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			res := httptest.NewRecorder()

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			requireAuthHandler.ServeHTTP(res, req)

			assert.Equal(t, tt.expectCode, res.Code)
			if tt.expectCode == http.StatusOK {
				assert.True(t, handlerCalled, "Handler should have been called")
			} else {
				assert.False(t, handlerCalled, "Handler should not have been called")
			}
		})
	}
}
