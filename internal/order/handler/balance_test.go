package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	orderDomain "github.com/aifedorov/gophermart/internal/order/domain"
	repository "github.com/aifedorov/gophermart/internal/order/repository/db"
	orderMocks "github.com/aifedorov/gophermart/internal/order/repository/mocks"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBalanceHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageBalanceHandler(ctrl)
	orderService := orderDomain.NewService(repo)
	handlerFunc := NewBalanceHandler(orderService)

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
			userID: TestUserID1.String(),
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
			userID: "550e8400-e29b-41d4-a716-446655440002",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				body:        `{"current":0,"withdrawn":0}` + "\n",
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
	mockRepo := orderMocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetUserBalanceByUserID(TestUserID1.String()).
		Return(decimal.NewFromInt(100), nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserWithdrawByUserID(TestUserID1.String()).
		Return(decimal.NewFromInt(0), nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserBalanceByUserID("550e8400-e29b-41d4-a716-446655440002").
		Return(decimal.NewFromInt(0), nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserWithdrawByUserID("550e8400-e29b-41d4-a716-446655440002").
		Return(decimal.NewFromInt(0), nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserBalanceByUserID("4").
		Return(decimal.Decimal{}, orderDomain.ErrOrderNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserWithdrawByUserID("4").
		Return(decimal.Decimal{}, orderDomain.ErrOrderNotFound).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserBalanceByUserID("5").
		Return(decimal.Decimal{}, assert.AnError).
		AnyTimes()

	mockRepo.EXPECT().
		GetUserWithdrawByUserID("5").
		Return(decimal.Decimal{}, assert.AnError).
		AnyTimes()

	return mockRepo
}
