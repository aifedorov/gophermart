package server

import (
	"github.com/aifedorov/gophermart/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
}

func NewServer() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

func (s *Server) mountHandlers() {
	s.router.Post("/api/user/register", handlers.NewRegisterHandler())
}
