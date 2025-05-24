package server

import (
	"log"
	"net/http"

	"github.com/odysseymorphey/quotes-service/internal/handlers"
	"github.com/odysseymorphey/quotes-service/internal/repository"
)

type Server struct {
	mux  *http.ServeMux
	repo repository.Repository
}

func New(r repository.Repository) *Server {
	m := http.NewServeMux()
	h := handlers.New(r)

	registerRoutes(m, h)

	return &Server{
		mux:  m,
		repo: r,
	}
}

func registerRoutes(mux *http.ServeMux, h *handlers.BaseHandler) {
	mux.HandleFunc("POST /quotes", h.AddQuote)

	mux.HandleFunc("GET /quotes", h.GetQuotes)
	mux.HandleFunc("GET /quotes/random", h.GetRandomQuote)

	mux.HandleFunc("DELETE /quotes/{id}", h.DeleteQuote)
}

func (s *Server) Run() {
	if err := http.ListenAndServe(":8080", s.mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *Server) Stop() {
	// TODO: Need graceful shutdown
}
