package server

import (
	orderDomain "github.com/aifedorov/gophermart/internal/order/domain"
	orderHandler "github.com/aifedorov/gophermart/internal/order/handler"
	"github.com/aifedorov/gophermart/internal/pkg/config"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"github.com/aifedorov/gophermart/internal/pkg/middleware"
	userDomain "github.com/aifedorov/gophermart/internal/user/domain"
	userHandler "github.com/aifedorov/gophermart/internal/user/handler"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	router       *chi.Mux
	config       config.Config
	userService  userDomain.Service
	orderService orderDomain.Service
}

func NewServer(
	cfg config.Config,
	userService userDomain.Service,
	orderService orderDomain.Service,
) *Server {
	return &Server{
		router:       chi.NewRouter(),
		config:       cfg,
		userService:  userService,
		orderService: orderService,
	}
}

func (s *Server) Run() error {
	s.mountHandlers()

	logger.Log.Info("server: running on", zap.String("address", s.config.ListenAddress))
	return http.ListenAndServe(s.config.ListenAddress, s.router)
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
		r.Get("/api/user/balance", orderHandler.NewBalanceHandler(s.orderService))
		r.Post("/api/user/balance/withdraw", orderHandler.NewWithdrawHandler(s.orderService))
		r.Get("/api/user/withdrawals", orderHandler.NewWithdrawalsHandler(s.orderService))
	})
}
