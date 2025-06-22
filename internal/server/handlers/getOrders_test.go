package handlers

import (
	"context"
	"github.com/aifedorov/gophermart/internal/domain/order"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/aifedorov/gophermart/internal/repository"
	mock_repository "github.com/aifedorov/gophermart/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetOrdersHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForGetOrders(ctrl)
	orderService := order.NewService(repo)
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
			userID: "1",
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
				ctx := context.WithValue(req.Context(), auth.UserIDKey, tt.userID)
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
	mockRepo := mock_repository.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetOrdersByUserID("1").
		Return([]repository.Order{
			{ID: "1", UserID: "1", Number: "4532015112830366", Status: repository.New},
		}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetOrdersByUserID("2").
		Return(nil, repository.ErrNotFound).
		AnyTimes()

	return mockRepo
}
