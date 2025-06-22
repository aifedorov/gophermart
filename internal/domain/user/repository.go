package user

import "github.com/aifedorov/gophermart/internal/repository"

type Repository interface {
	CreateUser(login, password string) (repository.User, error)
	GetUserByCredentials(login, password string) (repository.User, error)
}
