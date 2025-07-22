package domain

import (
	"errors"
	"github.com/aifedorov/gophermart/internal/user/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(req repository.RegisterRequest) (*repository.User, error) {
	if !s.isValidCredentials(req.Login, req.Password) {
		return nil, ErrEmptyCredentials
	}

	user, err := s.repo.CreateUser(req.Login, req.Password)
	if errors.Is(err, ErrAlreadyExists) {
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Login(req repository.LoginRequest) (*repository.User, error) {
	if !s.isValidCredentials(req.Login, req.Password) {
		return nil, ErrEmptyCredentials
	}

	user, err := s.repo.GetUserByCredentials(req.Login, req.Password)
	if errors.Is(err, ErrInvalidateCredentials) || errors.Is(err, ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) isValidCredentials(login, password string) bool {
	return login != "" && password != ""
}

func (s *Service) GetUserBalance(userID string) (repository.Balance, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return repository.Balance{}, err
	}

	return repository.Balance{
		Current:   user.Balance,
		Withdrawn: 0, // TODO: Calculate withdrawn
	}, nil
}

func (s *Service) Withdraw(userID, orderNumber string, amount float64) error {
	if amount <= 0 {
		return ErrWithdrawNegativeAmount
	}
	return s.repo.Withdrawal(userID, orderNumber, amount)
}

func (s *Service) GetWithdrawals(userID string) ([]repository.Withdrawal, error) {
	return s.repo.GetWithdrawalsByUserID(userID)
}
