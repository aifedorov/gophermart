package middleware

import (
	"github.com/aifedorov/gophermart/internal/domain/auth"
	"net/http"
)

type AuthMiddleware struct {
	authService auth.Service
}

func NewAuthMiddleware(authService auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
