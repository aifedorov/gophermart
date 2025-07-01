package repository

import (
	"sync"
	"time"

	"github.com/aifedorov/gophermart/internal/domain/order"
	"github.com/aifedorov/gophermart/internal/domain/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type InMemoryStorage struct {
	mu          sync.RWMutex
	users       map[string]user.User
	orders      map[string]order.Order
	withdrawals map[string]user.Withdrawal
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:  make(map[string]user.User),
		orders: make(map[string]order.Order),
	}
}

func (ms *InMemoryStorage) CreateUser(login, password string) (user.User, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.users[login]
	if ok {
		return user.User{}, user.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, err
	}

	newUser := user.User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: string(hashedPassword),
	}

	ms.users[login] = newUser
	return newUser, nil
}

func (ms *InMemoryStorage) GetUserByCredentials(login, password string) (user.User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	userObj, ok := ms.users[login]
	if !ok {
		return user.User{}, user.ErrNotFound
	}

	err := bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(password))
	if err != nil {
		return user.User{}, user.ErrInvalidateCredentials
	}
	return userObj, nil
}

func (ms *InMemoryStorage) GetUserByID(userID string) (user.User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for _, userObj := range ms.users {
		if userObj.ID == userID {
			return userObj, nil
		}
	}
	return user.User{}, user.ErrNotFound
}

func (ms *InMemoryStorage) CreateOrderByUserID(userID, orderNumber string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, ok := ms.orders[orderNumber]
	if ok {
		return order.ErrAlreadyExists
	}

	ms.orders[orderNumber] = order.Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Number:    orderNumber,
		Status:    order.StatusNew,
		CreatedAt: time.Now(),
	}
	return nil
}

func (ms *InMemoryStorage) GetOrdersByUserID(userID string) ([]order.Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	res := make([]order.Order, 0)
	for _, orderObj := range ms.orders {
		if orderObj.UserID == userID {
			res = append(res, orderObj)
		}
	}
	return res, nil
}

func (ms *InMemoryStorage) CreateOrder(userID, number string) (order.Order, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	orderObj := order.Order{
		ID:        uuid.NewString(),
		UserID:    userID,
		Number:    number,
		Status:    order.StatusNew,
		CreatedAt: time.Now(),
	}

	ms.orders[number] = orderObj
	return orderObj, nil
}

func (ms *InMemoryStorage) GetOrderByNumber(number string) (order.Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	orderObj, ok := ms.orders[number]
	if !ok {
		return order.Order{}, order.ErrOrderNotFound
	}
	return orderObj, nil
}

func (ms *InMemoryStorage) UpdateOrder(orderObj order.Order) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.orders[orderObj.Number]; !ok {
		return order.ErrOrderNotFound
	}

	ms.orders[orderObj.Number] = orderObj
	return nil
}

func (ms *InMemoryStorage) Withdrawal(userID, orderNumber string, amount float64) error {
	if amount <= 0 {
		return user.ErrWithdrawNegativeAmount
	}

	userObj, err := ms.GetUserByID(userID)
	if err != nil {
		return err
	}
	if userObj.Balance < amount {
		return user.ErrWithdrawInsufficientFunds
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	userObj.Balance -= amount
	ms.users[userObj.Login] = userObj

	withdrawal := user.Withdrawal{
		ID:    uuid.NewString(),
		Order: orderNumber,
		Sum:   amount,
	}
	ms.withdrawals[userObj.ID] = withdrawal

	return nil
}

func (ms *InMemoryStorage) UpdateOrderStatus(number string, status order.Status) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	orderObj, ok := ms.orders[number]
	if !ok {
		return order.ErrOrderNotFound
	}

	orderObj.Status = status
	ms.orders[number] = orderObj
	return nil
}
