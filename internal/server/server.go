package server

import (
	orderDomain "github.com/aifedorov/gophermart/internal/order/domain"
	orderHandler "github.com/aifedorov/gophermart/internal/order/handler"
	orderRepository "github.com/aifedorov/gophermart/internal/order/repository"
	"github.com/aifedorov/gophermart/internal/pkg/config"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	userDomain "github.com/aifedorov/gophermart/internal/user/domain"
	userHandler "github.com/aifedorov/gophermart/internal/user/handler"
	userRepository "github.com/aifedorov/gophermart/internal/user/repository"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	router       *chi.Mux
	config       config.Config
	userService  *userDomain.Service
	orderService *orderDomain.Service
}

func NewServer(cfg config.Config, userRepo userRepository.Repository, orderRepo orderRepository.Repository) *Server {
	userService := userDomain.NewService(userRepo)
	orderService := orderDomain.NewService(orderRepo)
	return &Server{
		router:       chi.NewRouter(),
		config:       cfg,
		userService:  userService,
		orderService: orderService,
	}
}

func (s *Server) Run() {
	s.mountHandlers()

	logger.Log.Info("server: running on", zap.String("address", s.config.ListenAddress))
	err := http.ListenAndServe(s.config.ListenAddress, s.router)
	if err != nil {
		logger.Log.Fatal("server: failed to run", zap.Error(err))
	}
}

func (s *Server) mountHandlers() {

	jwtMiddleware := middleware.NewJWTMiddleware(s.config.SecretKey)

	s.router.Use(chimiddleware.Compress(6, "application/json", "text/plain", "text/html"))
	s.router.Use(middleware.RequestLogger)
	s.router.Use(middleware.ResponseLogger)

	s.router.Post("/api/user/register", userHandler.NewRegisterHandler(s.config, s.userService))
	s.router.Post("/api/user/login", userHandler.NewLoginHandler(s.config, s.userService))

	s.router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.CheckJWT)
		r.Post("/api/user/orders", orderHandler.NewCreateOrdersHandler(s.orderService))
		r.Get("/api/user/orders", orderHandler.NewGetOrdersHandler(s.orderService))
		r.Get("/api/user/balance", userHandler.NewBalanceHandler(s.userService))
		r.Post("/api/user/balance/withdraw", userHandler.NewWithdrawHandler(s.userService))
		r.Get("/api/user/withdrawals", userHandler.NewWithdrawalsHandler(s.userService))
	})
}
