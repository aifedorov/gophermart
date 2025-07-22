package accrual

import (
	"encoding/json"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Order struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

type Client interface {
	getAccrualByOrderNumber(orderNumber string) (Order, error)
}

type HTTPClient struct {
	client        *http.Client
	ListenAddress string
}

func NewHTTPClient() Client {
	return &HTTPClient{
		client:        &http.Client{},
		ListenAddress: ":8084",
	}
}

var _ Client = (*HTTPClient)(nil)

func (c *HTTPClient) getAccrualByOrderNumber(orderNumber string) (Order, error) {
	request, err := http.NewRequest(http.MethodPost, c.ListenAddress, strings.NewReader(orderNumber))
	if err != nil {
		logger.Log.Error("accrualclient: failed to create request", zap.Error(err))
		return Order{}, err
	}
	request.Header.Add("Content-Length", "0")
	response, err := c.client.Do(request)
	if err != nil {
		logger.Log.Error("accrualclient: failed to send request", zap.Error(err))
		return Order{}, err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Log.Error("accrualclient: failed to close response body", zap.Error(err))
			return
		}
	}()

	var order Order
	if err := json.NewDecoder(response.Body).Decode(&order); err != nil {
		logger.Log.Error("accrualclient: failed to decode response body", zap.Error(err))
		return Order{}, err
	}

	return order, nil
}
