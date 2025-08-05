package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func decodeRegister(r *http.Request) (RegisterRequest, error) {
	var body RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		return RegisterRequest{}, errors.New("request body is empty")
	}
	if err != nil {
		return RegisterRequest{}, fmt.Errorf("failed to decode request: %w", err)
	}
	return body, nil
}

func decodeLogin(r *http.Request) (LoginRequest, error) {
	var body LoginRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		return LoginRequest{}, errors.New("request body is empty")
	}
	if err != nil {
		return LoginRequest{}, fmt.Errorf("failed to decode request: %w", err)
	}
	return body, nil
}
