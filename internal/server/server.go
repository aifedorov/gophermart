package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/aifedorov/gophermart/internal/config"
	"github.com/aifedorov/gophermart/internal/logger"
	"github.com/aifedorov/gophermart/internal/repository"
	"github.com/aifedorov/gophermart/internal/server/handlers"
	"github.com/aifedorov/gophermart/internal/server/middleware"
	"github.com/aifedorov/gophermart/internal/server/middleware/auth"
)

type Server struct {
	router *chi.Mux
	config config.Config
	repo   repository.Repository
}

func NewServer(cfg config.Config, repo repository.Repository) *Server {
	return &Server{
		router: chi.NewRouter(),
		config: cfg,
		repo:   repo,
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

	jwtMiddleware := auth.NewJWTMiddleware(s.config.SecretKey)

	s.router.Use(chimiddleware.Compress(6, "application/json", "text/plain", "text/html"))
	s.router.Use(middleware.RequestLogger)
	s.router.Use(middleware.ResponseLogger)

	s.router.Post("/api/user/register", handlers.NewRegisterHandler(s.config, s.repo))
	s.router.Post("/api/user/login", handlers.NewLoginHandler(s.config, s.repo))

	s.router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.CheckJWT)
		r.Post("/api/user/orders", handlers.NewCreateOrdersHandler(s.repo))
		r.Get("/api/user/orders", handlers.NewGetOrdersHandler(s.repo))
	})
}
