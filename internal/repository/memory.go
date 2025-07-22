package repository

import (
	"github.com/aifedorov/gophermart/internal/order/domain"
	"github.com/aifedorov/gophermart/internal/order/repository"
	domain2 "github.com/aifedorov/gophermart/internal/user/domain"
	repository2 "github.com/aifedorov/gophermart/internal/user/repository"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type InMemoryStorage struct {
	mu          sync.RWMutex
	users       map[string]repository2.User
	orders      map[string]repository.Order
	withdrawals map[string]repository2.Withdrawal
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:       make(map[string]repository2.User),
		orders:      make(map[string]repository.Order),
		withdrawals: make(map[string]repository2.Withdrawal),
	}
}

func (ms *InMemoryStorage) CreateUser(login, password string) (repository2.User, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.users[login]
	if ok {
		return repository2.User{}, domain2.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return repository2.User{}, err
	}

	newUser := repository2.User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: string(hashedPassword),
	}

	ms.users[login] = newUser
	return newUser, nil
}

func (ms *InMemoryStorage) GetUserByCredentials(login, password string) (repository2.User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	userObj, ok := ms.users[login]
	if !ok {
		return repository2.User{}, domain2.ErrNotFound
	}

	err := bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(password))
	if err != nil {
		return repository2.User{}, domain2.ErrInvalidateCredentials
	}
	return userObj, nil
}

func (ms *InMemoryStorage) GetUserByID(userID string) (repository2.User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for _, userObj := range ms.users {
		if userObj.ID == userID {
			return userObj, nil
		}
	}
	return repository2.User{}, domain2.ErrNotFound
}

func (ms *InMemoryStorage) CreateOrderByUserID(userID, orderNumber string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.orders[orderNumber]
	if ok {
		return domain.ErrAlreadyExists
	}

	ms.orders[orderNumber] = repository.Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Number:    orderNumber,
		Status:    repository.StatusNew,
		CreatedAt: time.Now(),
	}
	return nil
}

func (ms *InMemoryStorage) GetOrdersByUserID(userID string) ([]repository.Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	res := make([]repository.Order, 0)
	for _, orderObj := range ms.orders {
		if orderObj.UserID == userID {
			res = append(res, orderObj)
		}
	}
	return res, nil
}

func (ms *InMemoryStorage) CreateOrder(userID, number string) (repository.Order, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	orderObj := repository.Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Number:    number,
		Status:    repository.StatusNew,
		CreatedAt: time.Now(),
	}

	ms.orders[number] = orderObj
	return orderObj, nil
}

func (ms *InMemoryStorage) GetOrderByNumber(number string) (repository.Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	orderObj, ok := ms.orders[number]
	if !ok {
		return repository.Order{}, domain.ErrOrderNotFound
	}
	return orderObj, nil
}

func (ms *InMemoryStorage) UpdateOrder(orderObj repository.Order) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.orders[orderObj.Number]; !ok {
		return domain.ErrOrderNotFound
	}

	ms.orders[orderObj.Number] = orderObj
	return nil
}

func (ms *InMemoryStorage) Withdrawal(userID, orderNumber string, amount float64) error {
	if amount <= 0 {
		return domain2.ErrWithdrawNegativeAmount
	}

	userObj, err := ms.GetUserByID(userID)
	if err != nil {
		return err
	}
	if userObj.Balance < amount {
		return domain2.ErrWithdrawInsufficientFunds
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	userObj.Balance -= amount
	ms.users[userObj.Login] = userObj

	withdrawal := repository2.Withdrawal{
		ID:          uuid.NewString(),
		OrderNumber: orderNumber,
		Sum:         amount,
	}
	ms.withdrawals[userObj.ID] = withdrawal

	return nil
}

func (ms *InMemoryStorage) GetWithdrawalsByUserID(userID string) ([]repository2.Withdrawal, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	withdrawals := make([]repository2.Withdrawal, 0)
	for _, w := range ms.withdrawals {
		if w.UserID == userID {
			withdrawals = append(withdrawals, w)
		}
	}

	return withdrawals, nil
}

func (ms *InMemoryStorage) UpdateOrderStatus(number string, status repository.Status) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	orderObj, ok := ms.orders[number]
	if !ok {
		return domain.ErrOrderNotFound
	}

	orderObj.Status = status
	ms.orders[number] = orderObj
	return nil
}
