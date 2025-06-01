package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/aifedorov/gophermart/internal/api"
	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/repository"
	"github.com/aifedorov/gophermart/internal/repository/mocks"
)

type ServerTestSuite struct {
	suite.Suite
	server   *httptest.Server
	client   *http.Client
	ctrl     *gomock.Controller
	mockRepo *mocks.MockRepository
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.client = &http.Client{}
}

func (suite *ServerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.ctrl)

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
		Return(nil)

	suite.mockRepo.EXPECT().
		GetUserByCredentials(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

	// 1. Register user
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 2. Login user
	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *ServerTestSuite) TestCreateOrderThenGetOrders() {
	orderNumber := "4532015112830366"

	suite.mockRepo.EXPECT().
		CreateOrder(orderNumber).
		Return(nil).
		AnyTimes()

	suite.mockRepo.EXPECT().
		GetOrders().
		Return([]repository.Order{
			{
				ID:        "1",
				Number:    orderNumber,
				Status:    repository.New,
				CreatedAt: time.Time{},
			},
		}, nil).
		AnyTimes()

	// 1. Create order
	resp := suite.createOrder(orderNumber)
	suite.Equal(http.StatusAccepted, resp.StatusCode)

	// 2. Get orders
	resp = suite.getOrders()
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 3. Check returned orders
	var orders []api.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&orders)
	suite.Require().NoError(err)
	suite.Equal(1, len(orders))
	suite.Equal(orderNumber, orders[0].Number)
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

func newMockConfig() *config.Config {
	return &config.Config{
		ListenAddress:        "localhost:8080",
		StorageDSN:           "postgres://test",
		AccrualSystemAddress: "localhost:8081",
		LogLevel:             "debug",
	}
}
