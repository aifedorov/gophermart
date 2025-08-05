package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func encodeResponse(rw http.ResponseWriter, orders []OrderResponse) error {
	encoder := json.NewEncoder(rw)

	if err := encoder.Encode(orders); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
		return errors.New("failed to encode response")
	}
	return nil
}

func decodeWithdraw(r *http.Request) (WithdrawRequest, error) {
	var body WithdrawRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		return WithdrawRequest{}, errors.New("request body is empty")
	}
	if err != nil {
		return WithdrawRequest{}, fmt.Errorf("failed to decode request: %w", err)
	}
	return body, nil
}

func encodeJSONResponse(rw http.ResponseWriter, data interface{}) error {
	encoder := json.NewEncoder(rw)

	if err := encoder.Encode(data); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
		return errors.New("failed to encode response")
	}
	return nil
}
