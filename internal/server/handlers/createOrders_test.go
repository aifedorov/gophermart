package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/aifedorov/gophermart/internal/repository"
	"github.com/aifedorov/gophermart/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrdersHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForCreateOrders(ctrl)
	handlerFunc := NewCreateOrdersHandler(repo)

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
			name:   "valid order number",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `4532015112830366`,
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
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:   "order already exists",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   `5555555555554444`,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "empty body",
			method: http.MethodPost,
			path:   "/api/user/orders",
			body:   "",
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
			handlerFunc(res, req)

			assert.Equal(t, tt.want.statusCode, res.Code)

			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, res.Header().Get("Content-Type"))
			}
		})
	}
}

func newMockStorageForCreateOrders(ctrl *gomock.Controller) repository.Repository {
	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		CreateOrder("4532015112830366").
		Return(nil).
		AnyTimes()

	mockRepo.EXPECT().
		CreateOrder("5555555555554444").
		Return(repository.ErrAlreadyExists).
		AnyTimes()

	return mockRepo
}
