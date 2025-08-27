package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/order/repository/db"
	orderMocks "github.com/aifedorov/gophermart/internal/order/repository/mocks"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateOrdersHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForCreateOrders(ctrl)
	orderService := domain.NewService(repo)
	handlerFunc := NewCreateOrdersHandler(orderService)

	type want struct {
		contentType string
		statusCode  int
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
			name:   "valid order number",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `4532015112830366`,
			userID: TestUserID1.String(),
			want: want{
				statusCode:  http.StatusAccepted,
				contentType: "text/plain",
			},
		},
		{
			name:   "invalid order number",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `1234567890`,
			userID: TestUserID1.String(),
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:   "order already exists",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `5555555555554444`,
			userID: TestUserID1.String(),
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "order uploaded by another user",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `4111111111111111`,
			userID: TestUserID1.String(),
			want: want{
				statusCode: http.StatusConflict,
			},
		},
		{
			name:   "unauthorized",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `4532015112830366`,
			want: want{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:   "empty body",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   "",
			userID: TestUserID1.String(),
			want: want{
				statusCode: http.StatusBadRequest,
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
		})
	}
}

func newMockStorageForCreateOrders(ctrl *gomock.Controller) repository.Repository {
	mockRepo := orderMocks.NewMockRepository(ctrl)

	// Success case - order created
	mockRepo.EXPECT().
		GetOrderByNumber("4532015112830366").
		Return(repository.Order{}, domain.ErrOrderNotFound).
		AnyTimes()
	mockRepo.EXPECT().
		CreateTopUpOrder(TestUserID1.String(), "4532015112830366").
		Return(repository.Order{}, true, nil).
		AnyTimes()

	// OrderNumber already uploaded case - same user
	mockRepo.EXPECT().
		GetOrderByNumber("5555555555554444").
		Return(repository.Order{UserID: TestUserID1, Number: "5555555555554444"}, nil).
		AnyTimes()
	mockRepo.EXPECT().
		CreateTopUpOrder(TestUserID1.String(), "5555555555554444").
		Return(repository.Order{UserID: TestUserID1, Number: "5555555555554444"}, false, nil).
		AnyTimes()

	// OrderNumber uploaded by another user
	mockRepo.EXPECT().
		GetOrderByNumber("4111111111111111").
		Return(repository.Order{UserID: TestUserID2, Number: "4111111111111111"}, nil).
		AnyTimes()
	mockRepo.EXPECT().
		CreateTopUpOrder(TestUserID1.String(), "4111111111111111").
		Return(repository.Order{UserID: TestUserID2, Number: "4111111111111111"}, false, nil).
		AnyTimes()

	return mockRepo
}
