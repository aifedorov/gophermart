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

func (ms *InMemoryStorage) CreateUser(login, password string) (User, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	user, ok := ms.users[login]
	if ok {
		return User{}, ErrAlreadyExists
	}

	ms.users[login] = User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: password,
	}

	return user, nil
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

func (ms *InMemoryStorage) CreateOrderByUserID(userID, orderNumber string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.orders[orderNumber]
	if ok {
		return ErrAlreadyExists
	}

	ms.orders[orderNumber] = Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Number:    orderNumber,
		Status:    New,
		CreatedAt: time.Now(),
	}
	return nil
}

func (ms *InMemoryStorage) GetOrdersByUserID(userID string) ([]Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	res := make([]Order, 0)
	for _, order := range ms.orders {
		if order.UserID == userID {
			res = append(res, order)
		}
	}
	return res, nil
}

func (ms *InMemoryStorage) CreateOrder(userID, number string) (Order, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	order := Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Number:    number,
		Status:    New,
		CreatedAt: time.Now(),
	}

	ms.orders[number] = order
	return order, nil
}

func (ms *InMemoryStorage) GetOrderByNumber(number string) (Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	order, ok := ms.orders[number]
	if !ok {
		return Order{}, ErrOrderNotFound
	}
	return order, nil
}

func (ms *InMemoryStorage) UpdateOrder(order Order) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.orders[order.Number]; !ok {
		return ErrOrderNotFound
	}

	ms.orders[order.Number] = order
	return nil
}
