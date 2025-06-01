package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RegisterRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func decodeRequest(r *http.Request) (RegisterRequestBody, error) {
	var body RegisterRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		return RegisterRequestBody{}, errors.New("request body is empty")
	}
	if err != nil {
		return RegisterRequestBody{}, fmt.Errorf("failed to decode request: %w", err)
	}
	return body, nil
}

func isValidCredentials(login, password string) bool {
	return login != "" && password != ""
}
