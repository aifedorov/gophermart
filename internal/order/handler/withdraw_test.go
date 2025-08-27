package handler

import (
	"context"
	"encoding/json"
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

const testOrderNumber = "2377225624"

func TestWithdrawHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
		mock    func(mockRepo *orderMocks.MockRepository)
	}{
		{
			name:   "successful withdrawal - amount less then balance",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: TestUserID1.String(),
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   decimal.NewFromInt(50),
			},
			want: want{
				statusCode: http.StatusOK,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {
				mockRepo.EXPECT().
					CreateWithdrawalOrder(TestUserID1.String(), testOrderNumber, decimal.NewFromInt(50)).
					Return(repository.Order{}, nil).
					Times(1)
			},
		},
		{
			name:   "successful withdrawal - amount is equal to balance",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: TestUserID1.String(),
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   decimal.NewFromInt(100),
			},
			want: want{
				statusCode: http.StatusOK,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {
				mockRepo.EXPECT().
					CreateWithdrawalOrder(TestUserID1.String(), testOrderNumber, decimal.NewFromInt(100)).
					Return(repository.Order{}, nil).
					Times(1)
			},
		},

		{
			name:   "insufficient funds	- amount more then balance",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: TestUserID1.String(),
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   decimal.NewFromInt(200),
			},
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {
				mockRepo.EXPECT().
					CreateWithdrawalOrder(TestUserID1.String(), testOrderNumber, decimal.NewFromInt(200)).
					Return(repository.Order{}, orderDomain.ErrWithdrawInsufficientFunds).
					Times(1)
			},
		},
		{
			name:   "insufficient funds	- amount is zero",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: TestUserID1.String(),
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   decimal.NewFromInt(0),
			},
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {},
		},
		{
			name:   "insufficient funds - negative amount",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: TestUserID1.String(),
			request: WithdrawRequest{
				Order: testOrderNumber,
				Sum:   decimal.NewFromInt(-100),
			},
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {},
		},
		{
			name:   "order number - invalid value",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: TestUserID1.String(),
			request: WithdrawRequest{
				Order: "invalid_order",
				Sum:   decimal.NewFromInt(751),
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {},
		},
		{
			name:   "order number - empty value",
			method: http.MethodPost,
			path:   "/api/user/balance/withdraw",
			userID: TestUserID1.String(),
			request: WithdrawRequest{
				Order: "",
				Sum:   decimal.NewFromInt(751),
			},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {},
		},
		{
			name:    "empty request body",
			method:  http.MethodPost,
			path:    "/api/user/balance/withdraw",
			userID:  TestUserID1.String(),
			request: WithdrawRequest{},
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			mock: func(mockRepo *orderMocks.MockRepository) {},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderRepo := orderMocks.NewMockRepository(ctrl)
			tt.mock(mockOrderRepo)

			orderService := orderDomain.NewService(mockOrderRepo)
			handlerFunc := NewWithdrawHandler(orderService)

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
