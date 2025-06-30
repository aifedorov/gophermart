package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/logger"
)

func decodeRegister(r *http.Request) (api.RegisterRequest, error) {
	var body api.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		return api.RegisterRequest{}, errors.New("request body is empty")
	}
	if err != nil {
		return api.RegisterRequest{}, fmt.Errorf("failed to decode request: %w", err)
	}
	return body, nil
}

func decodeLogin(r *http.Request) (api.LoginRequest, error) {
	var body api.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		return api.LoginRequest{}, errors.New("request body is empty")
	}
	if err != nil {
		return api.LoginRequest{}, fmt.Errorf("failed to decode request: %w", err)
	}
	return body, nil
}

func encodeResponse(rw http.ResponseWriter, orders []api.OrderResponse) error {
	encoder := json.NewEncoder(rw)

	if err := encoder.Encode(orders); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
		return errors.New("failed to encode response")
	}
	return nil
}

func encodeJSONResponse(rw http.ResponseWriter, data interface{}) error {
	encoder := json.NewEncoder(rw)

	if err := encoder.Encode(data); err != nil {
		logger.Log.Error("failed to encode response", zap.Error(err))
		return errors.New("failed to encode response")
	}
	return nil
}
