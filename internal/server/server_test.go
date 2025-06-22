package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/repository"
	mock_repository "github.com/aifedorov/gophermart/internal/repository/mocks"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
)

type ServerTestSuite struct {
	suite.Suite
	server   *httptest.Server
	client   *http.Client
	ctrl     *gomock.Controller
	mockRepo *mock_repository.MockRepository
}

func (suite *ServerTestSuite) SetupSuite() {
	jar, _ := cookiejar.New(nil)
	suite.client = &http.Client{
		Jar: jar,
	}
}

func (suite *ServerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRepo = mock_repository.NewMockRepository(suite.ctrl)

	// Clear cookies between tests
	jar, _ := cookiejar.New(nil)
	suite.client.Jar = jar

	s := NewServer(newMockConfig(), suite.mockRepo)
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

	suite.mockRepo.EXPECT().
		CreateUser(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

	suite.mockRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

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
	orderNumber := "4532015112830366"

	// Setup mocks for registration, login, and orders
	suite.mockRepo.EXPECT().
		CreateUser(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

	suite.mockRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

	suite.mockRepo.EXPECT().
		GetOrderByNumber(orderNumber).
		Return(repository.Order{}, repository.ErrOrderNotFound).
		AnyTimes()

	suite.mockRepo.EXPECT().
		CreateOrder("1", orderNumber).
		Return(repository.Order{
			ID:        "1",
			UserID:    "1",
			Number:    orderNumber,
			Status:    repository.New,
			CreatedAt: time.Time{},
		}, nil).
		AnyTimes()

	suite.mockRepo.EXPECT().
		GetOrdersByUserID("1").
		Return([]repository.Order{
			{
				ID:        "1",
				UserID:    "1",
				Number:    orderNumber,
				Status:    repository.New,
				CreatedAt: time.Time{},
			},
		}, nil).
		AnyTimes()

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
	orderNumber := "4532015112830366"

	// Setup mocks
	suite.mockRepo.EXPECT().
		CreateUser(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

	suite.mockRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

	suite.mockRepo.EXPECT().
		GetOrderByNumber(orderNumber).
		Return(repository.Order{}, repository.ErrOrderNotFound)

	suite.mockRepo.EXPECT().
		CreateOrder("1", orderNumber).
		Return(repository.Order{
			ID:        "1",
			UserID:    "1",
			Number:    orderNumber,
			Status:    repository.New,
			CreatedAt: time.Time{},
		}, nil)

	suite.mockRepo.EXPECT().
		GetOrdersByUserID("1").
		Return([]repository.Order{
			{
				ID:        "1",
				UserID:    "1",
				Number:    orderNumber,
				Status:    repository.New,
				CreatedAt: time.Time{},
			},
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
	// Creat authenticated request for creating order
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

func newMockConfig() config.Config {
	return config.Config{
		ListenAddress:        "localhost:8080",
		StorageDSN:           "postgres://test",
		AccrualSystemAddress: "localhost:8081",
		LogLevel:             "debug",
		SecretKey:            "secret",
	}
}
