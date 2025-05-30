package server

import (
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

	s := newMockServer(ctrl)

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
			name:   "valid request",
			method: http.MethodPost,
			path:   "/api/user/register",
			body: `{
				"login": "test",
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

func newMockServer(ctrl *gomock.Controller) *Server {
	s := NewServer(newMockConfig(), newMockStorage(ctrl))
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

func newMockStorage(ctrl *gomock.Controller) repository.Repository {
	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		StoreUser("loginExists", gomock.Any()).
		Return(repository.ErrAlreadyExists).
		AnyTimes()

	mockRepo.EXPECT().
		StoreUser("test", "test").
		Return(nil).AnyTimes()

	mockRepo.EXPECT().
		StoreUser("", gomock.Any()).
		Return(repository.ErrNotFound).AnyTimes()

	mockRepo.EXPECT().
		StoreUser(gomock.Any(), "").
		Return(repository.ErrNotFound).
		AnyTimes()

	return mockRepo
}
