package user

import (
	"encoding/json"
	"github.com/aifedorov/gophermart/internal/domain/auth"
	"net/http"

	"github.com/aifedorov/gophermart/internal/api"
	userDomain "github.com/aifedorov/gophermart/internal/domain/user"
)

type RegisterHandler struct {
	userService userDomain.Service
	authService auth.Service
}

func NewRegisterHandler(userService userDomain.Service, authService auth.Service) *RegisterHandler {
	return &RegisterHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req api.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
