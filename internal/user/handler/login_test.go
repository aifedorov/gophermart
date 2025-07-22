package handler

import (
	"errors"
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

func TestLoginHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForLogin(ctrl)
	userService := domain.NewService(repo)
	handlerFunc := NewLoginHandler(newMockConfig(), userService)

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
		want   want
	}{
		{
			name:   "valid credentials",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: `{
				"login": "loginExists",
				"password": "test"
			}`,
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				body:        `{"id":"1","login":"loginExists","password":"pass"}`,
			},
		},
		{
			name:   "invalid login",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: `{
				"login": "loginNotExists",
				"password": "test"
			}`,
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "invalid password",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: `{
				"login": "test",
				"password": "wrongPass"
			}`,
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "bad request",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: `{
				"lg": "loginExists",
				"password": "test"
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "invalid json body",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: `{
				"login: "loginExists",
				"password": "test"
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: `{
				"login": "test",
				"password": "test"
			}`,
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			res := httptest.NewRecorder()
			handlerFunc(res, req)

			assert.Equal(t, tt.want.statusCode, res.Code)

			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header().Get("Content-Type"))
			}
		})
	}
}

func newMockStorageForLogin(ctrl *gomock.Controller) repository.Repository {
	mockRepo := userMocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetUserByCredentials("loginExists", "test").
		Return(repository.User{ID: "1", Login: "loginExists", Password: "test"}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByCredentials("loginNotExists", "test").
		Return(repository.User{}, domain.ErrNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByCredentials("test", "wrongPass").
		Return(repository.User{}, domain.ErrInvalidateCredentials).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByCredentials("test", "test").
		Return(repository.User{}, errors.New("internal error")).
		AnyTimes()

	return mockRepo
}
