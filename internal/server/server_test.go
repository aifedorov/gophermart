package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/repository"
	"github.com/aifedorov/gophermart/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
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

func (suite *ServerTestSuite) TestUserRegistrationAndLogin() {
	login := "test"
	pass := "pass"

	suite.mockRepo.EXPECT().
		StoreUser(login, pass).
		Return(nil)

	suite.mockRepo.EXPECT().
		FetchUser(login, pass).
		Return(repository.User{ID: "1", Login: login}, nil)

	// 1. Register user
	resp := suite.registerUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 2. Login user
	resp = suite.loginUser(login, pass)
	suite.Equal(http.StatusOK, resp.StatusCode)
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

func newMockConfig() *config.Config {
	return &config.Config{
		ListenAddress:        "localhost:8080",
		StorageDSN:           "postgres://test",
		AccrualSystemAddress: "localhost:8081",
		LogLevel:             "debug",
	}
}
