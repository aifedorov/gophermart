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

func TestCreateOrdersHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForCreateOrders(ctrl)
	orderService := order.NewService(repo)
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
			userID: "1",
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
			userID: "1",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:   "order already exists",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `5555555555554444`,
			userID: "1",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "order uploaded by another user",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `4111111111111111`,
			userID: "1",
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
			userID: "1",
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

func newMockStorageForCreateOrders(ctrl *gomock.Controller) repository.Repository {
	mockRepo := mock_repository.NewMockRepository(ctrl)

	// Success case - order created
	mockRepo.EXPECT().
		GetOrderByNumber("4532015112830366").
		Return(repository.Order{}, repository.ErrOrderNotFound).
		AnyTimes()
	mockRepo.EXPECT().
		CreateOrder("1", "4532015112830366").
		Return(repository.Order{ID: "1", UserID: "1", Number: "4532015112830366"}, nil).
		AnyTimes()

	// Order already uploaded case - same user
	mockRepo.EXPECT().
		GetOrderByNumber("5555555555554444").
		Return(repository.Order{ID: "2", UserID: "1", Number: "5555555555554444"}, nil).
		AnyTimes()

	// Order uploaded by another user
	mockRepo.EXPECT().
		GetOrderByNumber("4111111111111111").
		Return(repository.Order{ID: "3", UserID: "2", Number: "4111111111111111"}, nil).
		AnyTimes()

	return mockRepo
}
