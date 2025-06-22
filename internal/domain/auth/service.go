package auth

import (
	"net/http"

	"github.com/aifedorov/gophermart/internal/domain/user"
)

type Service interface {
	GenerateToken(userID user.ID) (string, error)
	ValidateToken(token string) (user.ID, error)
	SetAuthCookies(userID user.ID, w http.ResponseWriter) error
	GetUserIDFromRequest(r *http.Request) (user.ID, error)
}

type service struct {
	secretKey string
}

func (s service) GenerateToken(userID user.ID) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) ValidateToken(token string) (user.ID, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) SetAuthCookies(userID user.ID, w http.ResponseWriter) error {
	//TODO implement me
	panic("implement me")
}

func (s service) GetUserIDFromRequest(r *http.Request) (user.ID, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(secretKey string) Service {
	return &service{
		secretKey: secretKey,
	}
}
