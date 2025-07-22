package handler

import (
	"context"
	"encoding/json"
	"github.com/aifedorov/gophermart/internal/order/repository"
	orderMocks "github.com/aifedorov/gophermart/internal/order/repository/mocks"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"github.com/aifedorov/gophermart/internal/user/domain"
	repository2 "github.com/aifedorov/gophermart/internal/user/repository"
	userMocks "github.com/aifedorov/gophermart/internal/user/repository/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const testOrderNumber = "2377225624"

func TestWithdrawHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo, _ := newMockStorageForWithdraw(ctrl)
	userService := domain.NewService(userRepo)
	handlerFunc := NewWithdrawHandler(userService)

	type want struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name    string
		method  string
		path    string
		userID  string
		request WithdrawRequest
		want    want
	}{
		{
			name:   "successful withdrawal - amount less then balance",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: "1",
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   50,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "successful withdrawal - amount is equal to balance",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: "1",
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   100,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "unauthorized - no user id in context",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   100,
			},
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "insufficient funds	- amount more then balance",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: "1",
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   200,
			},
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
		},
		{
			name:   "insufficient funds	- amount is zero",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: "1",
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   0,
			},
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
		},
		{
			name:   "insufficient funds - negative amount",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: "1",
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   -100,
			},
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
		},
		{
			name:   "order number - invalid value",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: "1",
			request: WithdrawRequest{
				Order: "invalid_order",
				Sum:   751,
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:   "order number - empty value",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: "1",
			request: WithdrawRequest{
				Order: "",
				Sum:   751,
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:    "empty request body",
			method:  http.MethodPost,
			path:    "/api/user/balance/withdraw",
			userID:  "1",
			request: WithdrawRequest{},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqJSON, _ := json.Marshal(tt.request)
			body := strings.NewReader(string(reqJSON))
			req := httptest.NewRequest(tt.method, tt.path, body)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.userID)
				req = req.WithContext(ctx)
			}

			handlerFunc(res, req)

			assert.Equal(t, tt.want.statusCode, res.Code)
			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, res.Body.String())
			}
		})
	}
}

func newMockStorageForWithdraw(ctrl *gomock.Controller) (repository2.Repository, repository.Repository) {
	mockUserRepo := userMocks.NewMockRepository(ctrl)
	mockOrderRepo := orderMocks.NewMockRepository(ctrl)

	// Mock successful withdrawal - amount less than balance (50.0)
	mockUserRepo.EXPECT().
		Withdrawal("1", testOrderNumber, 50.0).
		Return(nil).
		AnyTimes()

	// Mock successful withdrawal - amount equal to balance (100.0)
	mockUserRepo.EXPECT().
		Withdrawal("1", testOrderNumber, 100.0).
		Return(nil).
		AnyTimes()

	// Mock insufficient funds - amount more than balance (200.0)
	mockUserRepo.EXPECT().
		Withdrawal("1", testOrderNumber, 200.0).
		Return(domain.ErrWithdrawInsufficientFunds).
		AnyTimes()

	return mockUserRepo, mockOrderRepo
}
