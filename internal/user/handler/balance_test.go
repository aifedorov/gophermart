package handler

import (
	"context"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/aifedorov/gophermart/internal/user/domain"
	"github.com/aifedorov/gophermart/internal/user/repository"
	userMocks "github.com/aifedorov/gophermart/internal/user/repository/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBalanceHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageBalanceHandler(ctrl)
	userService := domain.NewService(repo)
	handlerFunc := NewBalanceHandler(userService)

	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		userID string
		want   want
	}{
		{
			name:   "return user balance",
			method: http.MethodGet,
			path:   "/api/user/balance",
			userID: "1",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				body:        `{"current":100,"withdrawn":0}` + "\n",
			},
		},
		{
			name:   "zero balance",
			method: http.MethodGet,
			path:   "/api/user/balance",
			userID: "2",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				body:        `{"current":0,"withdrawn":0}` + "\n",
			},
		},
		{
			name:   "unauthorized - no user id in context",
			method: http.MethodGet,
			path:   "/api/user/balance",
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "unauthorized - empty user id",
			method: http.MethodGet,
			path:   "/api/user/balance",
			userID: "",
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			res := httptest.NewRecorder()

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			handlerFunc(res, req)

			assert.Equal(t, tt.want.statusCode, res.Code)
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header().Get("Content-Type"))
			}
			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, res.Body.String())
			}
		})
	}
}

func newMockStorageBalanceHandler(ctrl *gomock.Controller) repository.Repository {
	mockRepo := userMocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetUserByID("1").
		Return(repository.User{ID: "1", Login: "test", Password: "test", Balance: 100}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByID("2").
		Return(repository.User{ID: "2", Login: "zero", Password: "zero"}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByID("4").
		Return(repository.User{}, domain.ErrNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByID("5").
		Return(repository.User{}, assert.AnError).
		AnyTimes()

	return mockRepo
}
