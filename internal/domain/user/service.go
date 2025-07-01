package user

import (
	"errors"
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

	user, err := s.repo.CreateUser(req.Login, req.Password)
	if errors.Is(err, ErrAlreadyExists) {
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Login(req LoginRequest) (*User, error) {
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

func (s *Service) GetUserBalance(userID string) (Balance, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return Balance{}, err
	}

	return Balance{
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
