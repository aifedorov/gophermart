package handler

import (
	"context"
	"encoding/json"
	orderDomain "github.com/aifedorov/gophermart/internal/order/domain"
	repository "github.com/aifedorov/gophermart/internal/order/repository/db"
	orderMocks "github.com/aifedorov/gophermart/internal/order/repository/mocks"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWithdrawalsHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type want struct {
		statusCode int
		body       []WithdrawalResponse
	}

	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	testWithdrawals := []repository.Order{
		{
			Number:      "2377225624",
			Amount:      decimal.NewFromInt(500),
			ProcessedAt: pgtype.Timestamptz{Time: fixedTime, Valid: true},
		},
		{
			Number:      "1234567890",
			Amount:      decimal.NewFromInt(250),
			ProcessedAt: pgtype.Timestamptz{Time: fixedTime, Valid: true},
		},
	}

	expectedResponse := []WithdrawalResponse{
		{
			Order:       "2377225624",
			Sum:         decimal.NewFromInt(500),
			ProcessedAt: fixedTime,
		},
		{
			Order:       "1234567890",
			Sum:         decimal.NewFromInt(250),
			ProcessedAt: fixedTime,
		},
	}

	tests := []struct {
		name        string
		method      string
		path        string
		userID      string
		mock        func(mockRepo *orderMocks.MockRepository)
		want        want
		expectError bool
		errorType   error
	}{
		{
			name:   "successful retrieval of withdrawals",
			method: http.MethodGet,
			path:   "//user/withdrawals",
			userID: TestUserID1.String(),
			mock: func(mockRepo *orderMocks.MockRepository) {
				mockRepo.EXPECT().
					GetWithdrawalsByUserID(TestUserID1.String()).
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
			userID: TestUserID1.String(),
			mock: func(mockRepo *orderMocks.MockRepository) {
				mockRepo.EXPECT().
					GetWithdrawalsByUserID(TestUserID1.String()).
					Return([]repository.Order{}, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name:   "unauthorized - no user id in context",
			method: http.MethodGet,
			mock:   func(mockRepo *orderMocks.MockRepository) {},
			path:   "/api/user/withdrawals",
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "internal server error - repository error",
			method: http.MethodGet,
			path:   "/api/user/withdrawals",
			userID: TestUserID1.String(),
			mock: func(mockRepo *orderMocks.MockRepository) {
				mockRepo.EXPECT().
					GetWithdrawalsByUserID(TestUserID1.String()).
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

			mockOrderRepo := orderMocks.NewMockRepository(ctrl)
			tt.mock(mockOrderRepo)

			orderService := orderDomain.NewService(mockOrderRepo)
			handlerFunc := NewWithdrawalsHandler(orderService)

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
