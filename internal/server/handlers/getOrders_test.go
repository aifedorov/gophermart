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

func TestGetOrdersHandler(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := newMockStorageForGetOrders(ctrl)
	handlerFunc := NewGetOrdersHandler(repo)

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
			name:   "return order list",
			method: http.MethodGet,
			path:   "/api/user/orders",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				body:        `{"id":"1","number":"4532015112830366","status":"CREATED"}`,
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

func newMockStorageForGetOrders(ctrl *gomock.Controller) repository.Repository {
	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetOrders().
		Return([]repository.Order{
			{ID: "1", Number: "4532015112830366", Status: repository.New},
		}, nil).
		AnyTimes()

	mockRepo.EXPECT().
		GetOrders().
		Return(nil, repository.ErrNotFound).
		AnyTimes()

	return mockRepo
}
