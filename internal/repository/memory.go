package repository

import (
	"sync"

	"github.com/google/uuid"
)

type InMemoryStorage struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users: make(map[string]User),
	}
}

func (ms *InMemoryStorage) StoreUser(login, password string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.users[login]
	if ok {
		return ErrAlreadyExists
	}
	ms.users[login] = User{
		ID:       uuid.New(),
		Login:    login,
		Password: password,
	}
	return nil
}

func (ms *InMemoryStorage) FetchUser(login string) (User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	user, ok := ms.users[login]
	if !ok {
		return User{}, ErrNotFound
	}
	return user, nil
}
