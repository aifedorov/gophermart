package handler

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/user/domain"
	repository "github.com/aifedorov/gophermart/internal/user/repository/db"
	userMocks "github.com/aifedorov/gophermart/internal/user/repository/mocks"
	"golang.org/x/crypto/bcrypt"
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
				"login": "loginExists",
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

	// Generate bcrypt hash for password "test"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	mockRepo.EXPECT().
		GetUserByUsername("loginExists").
		Return(repository.User{Username: "loginExists", PasswordHash: string(hashedPassword)}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByUsername("loginNotExists").
		Return(repository.User{}, repository.ErrUserNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserByUsername("test").
		Return(repository.User{}, errors.New("internal error")).
		AnyTimes()

	return mockRepo
}
