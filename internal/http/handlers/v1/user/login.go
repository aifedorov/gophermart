package user

import (
	"encoding/json"
	"net/http"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/domain/auth"
	userDomain "github.com/aifedorov/gophermart/internal/domain/user"
)

type LoginHandler struct {
	userService userDomain.Service
	authService auth.Service
}

func NewLoginHandler(userService userDomain.Service, authService auth.Service) *LoginHandler {
	return &LoginHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
