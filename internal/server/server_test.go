package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/domain/order"
	orderMocks "github.com/aifedorov/gophermart/internal/domain/order/mocks"
	"github.com/aifedorov/gophermart/internal/domain/user"
	userMocks "github.com/aifedorov/gophermart/internal/domain/user/mocks"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
)

type ServerTestSuite struct {
	suite.Suite
	server    *httptest.Server
	client    *http.Client
	ctrl      *gomock.Controller
	userRepo  *userMocks.MockRepository
	orderRepo *orderMocks.MockRepository
}

func (suite *ServerTestSuite) SetupSuite() {
	jar, _ := cookiejar.New(nil)
	suite.client = &http.Client{
		Jar: jar,
	}
}

func (suite *ServerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.userRepo = userMocks.NewMockRepository(suite.ctrl)
	suite.orderRepo = orderMocks.NewMockRepository(suite.ctrl)

	// Clear cookies between tests
	jar, _ := cookiejar.New(nil)
	suite.client.Jar = jar

	s := NewServer(newMockConfig(), suite.userRepo, suite.orderRepo)
	s.mountHandlers()

	suite.server = httptest.NewServer(s.router)
}

func (suite *ServerTestSuite) TearDownTest() {
	suite.server.Close()
	suite.ctrl.Finish()
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) TestUserRegistrationThenLogin() {
	login := "test"
	pass := "pass"
	userID := "user-id-1"

	// Mock expectations for registration
	suite.userRepo.EXPECT().
		CreateUser(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)

	// Mock expectations for login
	suite.userRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)

	// 1. Register user
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 2. Login user
	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestCreateOrderThenGetOrders() {
	login := "testuser"
	pass := "testpass"
	userID := "user-id-2"
	orderNumber := "4532015112830366"

	// Mock expectations for registration and login
	suite.userRepo.EXPECT().
		CreateUser(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)
	suite.userRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)

	// Mock expectations for order creation
	suite.orderRepo.EXPECT().
		GetOrderByNumber(orderNumber).
		Return(order.Order{}, order.ErrOrderNotFound)
	suite.orderRepo.EXPECT().
		CreateOrder(userID, orderNumber).
		Return(order.Order{ID: "order-id-1", UserID: userID, Number: orderNumber, Status: order.StatusNew}, nil)

	// Mock expectations for getting orders
	suite.orderRepo.EXPECT().
		GetOrdersByUserID(userID).
		Return([]order.Order{
			{ID: "order-id-1", UserID: userID, Number: orderNumber, Status: order.StatusNew},
		}, nil)

	// 1. Register and login to get auth
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 2. Create an order (authenticated)
	resp = suite.createOrder(orderNumber)
	suite.Equal(http.StatusAccepted, resp.StatusCode)
	_ = resp.Body.Close()

	// 3. Get orders (authenticated)
	resp = suite.getOrders()
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 4. Check returned orders
	var orders []api.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&orders)
	suite.Require().NoError(err)
	suite.Equal(1, len(orders))
	suite.Equal(orderNumber, orders[0].Number)
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestProtectedEndpointWithoutAuth() {
	// Test get orders without authentication
	resp := suite.getOrders()
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
	_ = resp.Body.Close()

	// Test creating order without authentication
	resp = suite.createOrder("4532015112830366")
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestProtectedEndpointWithAuth() {
	login := "testuser"
	pass := "testpass"
	userID := "user-auth-1"
	orderNumber := "4532015112830366"

	// Mock expectations for registration and login
	suite.userRepo.EXPECT().
		CreateUser(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)
	suite.userRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)

	// Mock expectations for order creation and retrieval
	suite.orderRepo.EXPECT().
		GetOrderByNumber(orderNumber).
		Return(order.Order{}, order.ErrOrderNotFound)
	suite.orderRepo.EXPECT().
		CreateOrder(userID, orderNumber).
		Return(order.Order{ID: "order-auth-1", UserID: userID, Number: orderNumber, Status: order.StatusNew}, nil)
	suite.orderRepo.EXPECT().
		GetOrdersByUserID(userID).
		Return([]order.Order{
			{ID: "order-auth-1", UserID: userID, Number: orderNumber, Status: order.StatusNew},
		}, nil)

	// 1. Register and login to get auth cookie
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var authCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == auth.CookieName {
			authCookie = cookie
			break
		}
	}
	suite.Require().NotNil(authCookie, "JWT cookie should be set after login")
	_ = resp.Body.Close()

	// 2. Test accessing protected endpoints with auth cookie
	// Create authenticated request for creating order
	req, err := http.NewRequest("POST", suite.server.URL+"/api/user/orders", strings.NewReader(orderNumber))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "text/plain")
	req.AddCookie(authCookie)

	resp, err = suite.client.Do(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusAccepted, resp.StatusCode)
	_ = resp.Body.Close()

	// Create an authenticated request for getting orders
	req, err = http.NewRequest("GET", suite.server.URL+"/api/user/orders", nil)
	suite.Require().NoError(err)
	req.AddCookie(authCookie)

	resp, err = suite.client.Do(req)
	suite.Require().NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// Verify response contains the created order
	var orders []api.OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orders)
	suite.Require().NoError(err)
	suite.Equal(1, len(orders))
	suite.Equal(orderNumber, orders[0].Number)
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestGetUserBalance() {
	login := "balanceuser"
	pass := "balancepass"
	userID := "user-balance-1"

	// Mock expectations for registration and login
	suite.userRepo.EXPECT().
		CreateUser(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)
	suite.userRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)

	// Mock expectations for getting balance
	suite.userRepo.EXPECT().
		GetUserByID(userID).
		Return(user.User{ID: userID, Login: login, Balance: 150.50}, nil)

	// 1. Register and login to get auth
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 2. Get user balance (authenticated)
	resp = suite.getBalance()
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 3. Check returned balance
	var balance api.BalanceResponse
	err := json.NewDecoder(resp.Body).Decode(&balance)
	suite.Require().NoError(err)
	suite.Equal(150.50, balance.Current)
	suite.Equal(0.0, balance.Withdrawn) // TODO: implement withdrawn calculation
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestSuccessfulWithdrawal() {
	login := "withdrawuser"
	pass := "withdrawpass"
	userID := "user-withdraw-1"
	orderNumber := "2377225624"
	withdrawAmount := 75.0

	// Mock expectations for registration and login
	suite.userRepo.EXPECT().
		CreateUser(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)
	suite.userRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)

	// Mock expectations for withdrawal
	suite.userRepo.EXPECT().
		Withdrawal(userID, orderNumber, withdrawAmount).
		Return(nil)

	// 1. Register and login to get auth
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 2. Perform withdrawal (authenticated)
	resp = suite.withdraw(orderNumber, withdrawAmount)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestWithdrawalInsufficientFunds() {
	login := "pooruser"
	pass := "poorpass"
	userID := "user-poor-1"
	orderNumber := "2377225624"
	withdrawAmount := 200.0

	// Mock expectations for registration and login
	suite.userRepo.EXPECT().
		CreateUser(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)
	suite.userRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(user.User{ID: userID, Login: login}, nil)

	// Mock expectations for withdrawal with insufficient funds
	suite.userRepo.EXPECT().
		Withdrawal(userID, orderNumber, withdrawAmount).
		Return(user.ErrWithdrawInsufficientFunds)

	// 1. Register and login to get auth
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 2. Attempt withdrawal with insufficient funds (authenticated)
	resp = suite.withdraw(orderNumber, withdrawAmount)
	suite.Equal(http.StatusPaymentRequired, resp.StatusCode)
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestWithdrawalWithoutAuth() {
	orderNumber := "2377225624"
	withdrawAmount := 50.0

	// No mock expectations needed since we're testing unauthenticated access

	// Attempt withdrawal without authentication
	resp := suite.withdraw(orderNumber, withdrawAmount)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
	_ = resp.Body.Close()
}

func (suite *ServerTestSuite) TestBalanceWithoutAuth() {
	// No mock expectations needed since we're testing unauthenticated access

	// Attempt to get balance without authentication
	resp := suite.getBalance()
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
	_ = resp.Body.Close()
}

// Helper methods

func (suite *ServerTestSuite) registerUser(login, password string) *http.Response {
	body := fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password)
	resp, err := suite.client.Post(
		suite.server.URL+"/api/user/register",
		"application/json",
		strings.NewReader(body),
	)

	suite.Require().NoError(err)
	return resp
}

func (suite *ServerTestSuite) loginUser(login, password string) *http.Response {
	body := fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password)
	resp, err := suite.client.Post(
		suite.server.URL+"/api/user/login",
		"application/json",
		strings.NewReader(body),
	)

	suite.Require().NoError(err)
	return resp
}

func (suite *ServerTestSuite) createOrder(orderNumber string) *http.Response {
	resp, err := suite.client.Post(
		suite.server.URL+"/api/user/orders",
		"text/plain",
		strings.NewReader(orderNumber),
	)

	suite.Require().NoError(err)
	return resp
}

func (suite *ServerTestSuite) getOrders() *http.Response {
	resp, err := suite.client.Get(suite.server.URL + "/api/user/orders")
	suite.Require().NoError(err)
	return resp
}

func (suite *ServerTestSuite) getBalance() *http.Response {
	resp, err := suite.client.Get(suite.server.URL + "/api/user/balance")
	suite.Require().NoError(err)
	return resp
}

func (suite *ServerTestSuite) withdraw(orderNumber string, amount float64) *http.Response {
	body := fmt.Sprintf(`{"order":"%s","sum":%f}`, orderNumber, amount)
	resp, err := suite.client.Post(
		suite.server.URL+"/api/user/balance/withdraw",
		"application/json",
		strings.NewReader(body),
	)
	suite.Require().NoError(err)
	return resp
}

func newMockConfig() config.Config {
	return config.Config{
		ListenAddress:        "localhost:8080",
		StorageDSN:           "postgres://test",
		AccrualSystemAddress: "localhost:8081",
		LogLevel:             "debug",
		SecretKey:            "secret",
	}
}
