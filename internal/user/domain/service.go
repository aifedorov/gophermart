package domain

import (
	"errors"
	"fmt"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	repository "github.com/aifedorov/gophermart/internal/user/repository/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req RegisterRequest) (*User, error)
	Login(req LoginRequest) (*User, error)
}
type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(req RegisterRequest) (*User, error) {
	if !s.isValidCredentials(req.Login, req.Password) {
		logger.Log.Info("userservice: invalid credentials")
		return nil, ErrEmptyCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("userservice: failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("userservice: failed to hash password: %w", err)
	}

	dbUser, err := s.repo.CreateUser(req.Login, string(hashedPassword))
	if errors.Is(err, repository.ErrUserAlreadyExists) {
		logger.Log.Info("userservice: user already exists", zap.Error(err))
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		logger.Log.Error("userservice: failed to create user", zap.Error(err))
		return nil, err
	}

	domainUser := s.convertUserToDomain(dbUser)
	return &domainUser, nil
}

func (s *service) Login(req LoginRequest) (*User, error) {
	if !s.isValidCredentials(req.Login, req.Password) {
		logger.Log.Info("userservice: invalid credentials")
		return nil, ErrEmptyCredentials
	}

	dbUser, err := s.repo.GetUserByUsername(req.Login)
	if errors.Is(err, repository.ErrUserNotFound) {
		logger.Log.Info("userservice: user not found", zap.String("login", req.Login))
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		logger.Log.Error("userservice: failed to get user", zap.Error(err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(req.Password))
	if err != nil {
		logger.Log.Info("userservice: invalid password", zap.String("login", req.Login))
		return nil, ErrInvalidCredentials
	}

	domainUser := s.convertUserToDomain(dbUser)
	return &domainUser, nil
}

func (s *service) isValidCredentials(login, password string) bool {
	return login != "" && password != ""
}

func (s *service) convertUserToDomain(dbUser repository.User) User {
	return User{
		ID:       dbUser.ID.String(),
		Login:    dbUser.Username,
		Password: dbUser.PasswordHash,
	}
}
