package handler

import (
	"encoding/json"
	"errors"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"go.uber.org/zap"
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
