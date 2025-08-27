package accrual

import (
	"github.com/aifedorov/gophermart/internal/pkg/config"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type HTTPClient interface {
	GetAccrualByOrderNumber(orderNumber string) (OrderResponse, bool, error)
	Close() error
}

type httpClient struct {
	client        *resty.Client
	ListenAddress string
}

func NewHTTPClient(cfg config.Config) HTTPClient {
	return &httpClient{
		client:        resty.New(),
		ListenAddress: cfg.AccrualSystemAddress,
	}
}

func (c *httpClient) Close() error {
	return c.client.Close()
}

func (c *httpClient) GetAccrualByOrderNumber(orderNumber string) (OrderResponse, bool, error) {
	res, err := c.client.R().
		SetResult(&OrderResponse{}).
		SetHeader("Accept", "application/json").
		Get(c.ListenAddress + "/api/orders/" + orderNumber)

	if err != nil {
		logger.Log.Error("accrualclient: order processing failed", zap.Error(err))
		return OrderResponse{}, false, err
	}
	if res.StatusCode() < 200 || res.StatusCode() > 299 {
		logger.Log.Error("accrualclient: order processing failed", zap.Int("status", res.StatusCode()))
		return OrderResponse{}, false, nil
	}
	result := res.Result().(*OrderResponse)
	return *result, true, nil
}
