package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/aifedorov/gophermart/internal/domain/user"
	userMocks "github.com/aifedorov/gophermart/internal/domain/user/mocks"
)

func TestServer_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForRegister(ctrl)
	userService := user.NewService(repo)
	handlerFunc := NewRegisterHandler(newMockConfig(), userService)

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   want
	}{
		{
			name:   "success register",
			method: http.MethodPost,
			path:   "/api/user/register",
			body: `{
				"login": "newLogin",
				"password": "test"
			}`,
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name:   "missing body",
			method: http.MethodPost,
			path:   "/api/user/register",
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
			name:   "empty login",
			method: http.MethodPost,
			path:   "/api/user/register",
			body: `{
				"login": "",
				"password": "test"
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "empty password",
			method: http.MethodPost,
			path:   "/api/user/register",
			body: `{
				"login": "test",
				"password": ""
			}`,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "login already exists",
			method: http.MethodPost,
			path:   "/api/user/register",
			body: `{
				"login": "loginExists",
				"password": "test"
			}`,
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name:   "internal server error",
			method: http.MethodPost,
			path:   "/api/user/register",
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

func newMockStorageForRegister(ctrl *gomock.Controller) user.Repository {
	mockRepo := userMocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		CreateUser("loginExists", gomock.Any()).
		Return(user.User{}, user.ErrAlreadyExists).
		AnyTimes()

	mockRepo.EXPECT().
		CreateUser("newLogin", "test").
		Return(user.User{ID: "1", Login: "newLogin"}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		CreateUser("", gomock.Any()).
		Return(user.User{}, user.ErrNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), "").
		Return(user.User{}, user.ErrNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		CreateUser("test", "test").
		Return(user.User{}, errors.New("internal error")).
		AnyTimes()

	return mockRepo
}
