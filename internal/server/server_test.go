package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/repository"
	"github.com/aifedorov/gophermart/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServer_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := newMockServer(newMockStorageForRegister(ctrl))

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
			s.router.ServeHTTP(res, req)

			assert.Equal(t, tt.want.statusCode, res.Code)

			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header().Get("Content-Type"))
			}
		})
	}
}

func TestServer_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := newMockServer(newMockStorageFoLogin(ctrl))

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
			s.router.ServeHTTP(res, req)

			assert.Equal(t, tt.want.statusCode, res.Code)

			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header().Get("Content-Type"))
			}
		})
	}
}

func newMockServer(repo repository.Repository) *Server {
	s := NewServer(newMockConfig(), repo)
	s.mountHandlers()
	return s
}

func newMockConfig() *config.Config {
	return &config.Config{
		ListenAddress:        "localhost:8080",
		StorageDSN:           "postgres://test",
		AccrualSystemAddress: "localhost:8081",
		LogLevel:             "debug",
	}
}

func newMockStorageForRegister(ctrl *gomock.Controller) repository.Repository {
	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		StoreUser("loginExists", gomock.Any()).
		Return(repository.ErrAlreadyExists).
		AnyTimes()

	mockRepo.EXPECT().
		StoreUser("newLogin", "test").
		Return(nil).
		AnyTimes()

	mockRepo.EXPECT().
		StoreUser("", gomock.Any()).
		Return(repository.ErrNotFound).AnyTimes()

	mockRepo.EXPECT().
		StoreUser(gomock.Any(), "").
		Return(repository.ErrNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		StoreUser("test", "test").
		Return(errors.New("internal error")).
		AnyTimes()

	return mockRepo
}

func newMockStorageFoLogin(ctrl *gomock.Controller) repository.Repository {
	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		FetchUser("loginExists", "test").
		Return(repository.User{ID: "1", Login: "loginExists", Password: "test"}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		FetchUser("loginNotExists", "test").
		Return(repository.User{}, repository.ErrNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		FetchUser("test", "wrongPass").
		Return(repository.User{}, repository.ErrInvalidateCredentials).
		AnyTimes()

	mockRepo.EXPECT().
		FetchUser("test", "test").
		Return(repository.User{}, errors.New("internal error")).
		AnyTimes()

	return mockRepo
}
