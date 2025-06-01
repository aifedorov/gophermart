package repository

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type InMemoryStorage struct {
	mu     sync.RWMutex
	users  map[string]User
	orders map[string]Order
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:  make(map[string]User),
		orders: make(map[string]Order),
	}
}

func (ms *InMemoryStorage) CreateUser(login, password string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.users[login]
	if ok {
		return ErrAlreadyExists
	}
	ms.users[login] = User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: password,
	}
	return nil
}

func (ms *InMemoryStorage) GetUserByCredentials(login, password string) (User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	user, ok := ms.users[login]
	if !ok {
		return User{}, ErrNotFound
	}
	// TODO: Use hash function instead of comparing passwords directly.
	if user.Password != password {
		return User{}, ErrInvalidateCredentials
	}
	return user, nil
}

func (ms *InMemoryStorage) CreateOrder(orderNumber string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.orders[orderNumber]
	if ok {
		return ErrAlreadyExists
	}

	ms.orders[orderNumber] = Order{
		ID:        uuid.NewString(),
		Number:    orderNumber,
		Status:    New,
		CreatedAt: time.Now(),
	}
	return nil
}

func (ms *InMemoryStorage) GetOrders() ([]Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	res := make([]Order, 0)
	for _, order := range ms.orders {
		res = append(res, order)
	}
	return res, nil
}
