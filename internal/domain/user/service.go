package user

import (
	"context"

	"github.com/aifedorov/gophermart/internal/api"
)

type Service interface {
	RegisterUser(ctx context.Context, req *api.RegisterRequest) error
	AuthenticateUser(ctx context.Context, req *api.LoginRequest) (*User, error)
}

type userService struct {
	userRepo Repository
	userSvc  Service
}

func NewService(userRepo Repository, userSvc Service) Service {
	return &userService{
		userRepo: userRepo,
		userSvc:  userSvc,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req *api.RegisterRequest) error {
	return nil
}

func (s *userService) AuthenticateUser(ctx context.Context, req *api.LoginRequest) (*User, error) {
	return nil, nil
}
