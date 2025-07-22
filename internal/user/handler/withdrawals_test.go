package handler

import (
	"context"
	"encoding/json"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/aifedorov/gophermart/internal/user/domain"
	"github.com/aifedorov/gophermart/internal/user/repository"
	userMocks "github.com/aifedorov/gophermart/internal/user/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithdrawalsHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type want struct {
		statusCode int
		body       []WithdrawalResponse
	}

	testWithdrawals := []repository.Withdrawal{
		{
			ID:          "1",
			UserID:      "1",
			OrderNumber: "2377225624",
			Sum:         500.0,
		},
		{
			ID:          "2",
			UserID:      "1",
			OrderNumber: "1234567890",
			Sum:         250.0,
		},
	}

	expectedResponse := []WithdrawalResponse{
		{
			Order: "2377225624",
			Sum:   500.0,
		},
		{
			Order: "1234567890",
			Sum:   250.0,
		},
	}

	tests := []struct {
		name        string
		method      string
		path        string
		userID      string
		mock        func(mockRepo *userMocks.MockRepository)
		want        want
		expectError bool
		errorType   error
	}{
		{
			name:   "successful retrieval of withdrawals",
			method: http.MethodGet,
			path:   "//user/withdrawals",
			userID: "1",
			mock: func(mockRepo *userMocks.MockRepository) {
				mockRepo.EXPECT().
					GetWithdrawalsByUserID("1").
					Return(testWithdrawals, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusOK,
				body:       expectedResponse,
			},
		},
		{
			name:   "no withdrawals found - returns 204",
			method: http.MethodGet,
			path:   "/api/user/withdrawals",
			userID: "1",
			mock: func(mockRepo *userMocks.MockRepository) {
				mockRepo.EXPECT().
					GetWithdrawalsByUserID("1").
					Return([]repository.Withdrawal{}, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name:   "unauthorized - no user id in context",
			method: http.MethodGet,
			mock:   func(mockRepo *userMocks.MockRepository) {},
			path:   "/api/user/withdrawals",
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "internal server error - repository error",
			method: http.MethodGet,
			path:   "/api/user/withdrawals",
			userID: "1",
			mock: func(mockRepo *userMocks.MockRepository) {
				mockRepo.EXPECT().
					GetWithdrawalsByUserID("1").
					Return(nil, assert.AnError).
					Times(1)
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUserRepo := userMocks.NewMockRepository(ctrl)
			tt.mock(mockUserRepo)

			userService := domain.NewService(mockUserRepo)
			handlerFunc := NewWithdrawalsHandler(userService)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			handlerFunc(res, req)

			assert.Equal(t, tt.want.statusCode, res.Code)

			if tt.want.body != nil {
				var actualResponse []WithdrawalResponse
				err := json.Unmarshal(res.Body.Bytes(), &actualResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.body, actualResponse)
			}
		})
	}
}
