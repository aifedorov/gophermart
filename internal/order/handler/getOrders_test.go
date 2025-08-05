package handler

import (
	"context"
	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/order/repository/db"
	orderMocks "github.com/aifedorov/gophermart/internal/order/repository/mocks"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetOrdersHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForGetOrders(ctrl)
	orderService := domain.NewService(repo)
	handlerFunc := NewGetOrdersHandler(orderService)

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
			name:   "return order list",
			method: http.MethodGet,
			path:   "/api/user/orders",
			userID: TestUserID1.String(),
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name:   "unauthorized - no user id in context",
			method: http.MethodGet,
			path:   "/api/user/orders",
			want: want{
				statusCode: http.StatusUnauthorized,
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

func newMockStorageForGetOrders(ctrl *gomock.Controller) repository.Repository {
	mockRepo := orderMocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetOrdersByUserID(TestUserID1.String()).
		Return([]repository.Order{
			{Number: "4532015112830366", Status: repository.OrderstatusNEW},
		}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetOrdersByUserID("550e8400-e29b-41d4-a716-446655440002").
		Return(nil, domain.ErrOrderNotFound).
		AnyTimes()

	return mockRepo
}
