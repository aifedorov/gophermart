package user

import (
	"errors"

	"github.com/aifedorov/gophermart/internal/repository"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(req RegisterRequest) (*User, error) {
	if !s.isValidCredentials(req.Login, req.Password) {
		return nil, ErrEmptyCredentials
	}

	repoUser, err := s.repo.CreateUser(req.Login, req.Password)
	if errors.Is(err, repository.ErrAlreadyExists) {
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       repoUser.ID,
		Login:    repoUser.Login,
		Password: repoUser.Password,
	}, nil
}

func (s *Service) Login(req LoginRequest) (*User, error) {
	if !s.isValidCredentials(req.Login, req.Password) {
		return nil, ErrEmptyCredentials
	}

	repoUser, err := s.repo.GetUserByCredentials(req.Login, req.Password)
	if errors.Is(err, repository.ErrInvalidateCredentials) || errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       repoUser.ID,
		Login:    repoUser.Login,
		Password: repoUser.Password,
	}, nil
}

func (s *Service) isValidCredentials(login, password string) bool {
	return login != "" && password != ""
}

func (s *Service) GetUserBalance(userID string) (*Balance, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Current:   user.Balance,
		Withdrawn: 0,
	}, nil
}
