package order

import (
	"io"
	"net/http"
	"strings"

	authDomain "github.com/aifedorov/gophermart/internal/domain/auth"
	orderDomain "github.com/aifedorov/gophermart/internal/domain/order"
)

type CreateHandler struct {
	orderService orderDomain.Service
	authService  authDomain.Service
}

func NewCreateHandler(orderService orderDomain.Service, authService authDomain.Service) *CreateHandler {
	return &CreateHandler{
		orderService: orderService,
		authService:  authService,
	}
}

func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	_, err := h.authService.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		http.Error(w, "Empty order number", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
