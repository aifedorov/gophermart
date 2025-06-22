package http

import (
	"github.com/aifedorov/gophermart/internal/domain/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	orderDomain "github.com/aifedorov/gophermart/internal/domain/order"
	userDomain "github.com/aifedorov/gophermart/internal/domain/user"
	"github.com/aifedorov/gophermart/internal/http/handlers/v1/order"
	"github.com/aifedorov/gophermart/internal/http/handlers/v1/user"
	"github.com/aifedorov/gophermart/internal/http/middleware"
)

type Server struct {
	router       *chi.Mux
	listenAddr   string
	userService  userDomain.Service
	orderService orderDomain.Service
	authService  auth.Service
}

func NewServer(
	listenAddr string,
	userService userDomain.Service,
	orderService orderDomain.Service,
	authService auth.Service,
) *Server {
	return &Server{
		router:       chi.NewRouter(),
		listenAddr:   listenAddr,
		userService:  userService,
		orderService: orderService,
		authService:  authService,
	}
}

func (s *Server) Start() error {
	s.setupMiddleware()
	s.setupRoutes()

	return http.ListenAndServe(s.listenAddr, s.router)
}

func (s *Server) setupMiddleware() {
	s.router.Use(chimiddleware.Logger)
	s.router.Use(chimiddleware.Recoverer)
	s.router.Use(chimiddleware.Compress(6, "application/json", "text/plain"))
}

func (s *Server) setupRoutes() {
	s.router.Post("/api/user/register", user.NewRegisterHandler(s.userService, s.authService).Handle)
	s.router.Post("/api/user/login", user.NewLoginHandler(s.userService, s.authService).Handle)

	s.router.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(s.authService).Handle)

		r.Post("/api/user/orders", order.NewCreateHandler(s.orderService, s.authService).Handle)
		r.Get("/api/user/orders", order.NewListHandler(s.orderService, s.authService).Handle)
	})
}
